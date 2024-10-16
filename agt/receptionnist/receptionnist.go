package receptionnist

import (
	"IA04-hotel/agt"
	"IA04-hotel/agt/employee"
	"IA04-hotel/agt/room"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"
)

type Pref int

const (
	Default Pref = iota // permet de faire un enum en go
	Expensive
	Smallest
)

type Receptionnist struct {
	employee.Employee
	pref Pref
}

func NewReceptionnist(salary int, idHotel int, schedule []employee.Day, shift employee.Shift, id string, firstName string, lastName string, pref Pref) *Receptionnist {
	employe := employee.NewEmployee(0, salary, employee.Idle, true, idHotel, schedule, shift, id, firstName, lastName)
	return &Receptionnist{*employe, pref}
}

func (r Receptionnist) Pref() Pref {
	return r.pref
}

func (r *Receptionnist) SetPref(pref Pref) {
	r.pref = pref
}

func (r *Receptionnist) Start() {
	var data agt.DataReceptionist
	api_url := "http://localhost:8080/"
	dataRequest := agt.DataReceptionistRequest{
		ID:    r.GetId(),
		Hotel: r.GetIdHotel(),
	}
	json_data, err := json.Marshal(dataRequest)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for r.GetIsWorking() {
			time.Sleep(250 * time.Millisecond)
			data = agt.DataReceptionist{}
			res, err := http.Post(api_url+"dataReceptionist", "application/json", bytes.NewBuffer(json_data))
			if err == nil {
				if res != nil {
					json.NewDecoder(res.Body).Decode(&data)
				}
				for _, j := range data.Choice {
					r.validBook(data, j)
				}
				if len(data.Requests) > 0 {
					r.work(data)
				}
			}
		}
	}()
}

func (r *Receptionnist) work(data agt.DataReceptionist) {
	request := data.Requests[0]
	// Il faut faire une update de la base de données
	log.Println()
	log.Println("Je suis", r.GetId(), r.GetFirstName(), r.GetLastName(), "et je m'occupe de", request.IdClient)
	var roomsList = make([]room.Room, 0)
	if data.Time.Day <= request.BeginDate {
		for _, ro := range data.Rooms {
			roomIsAccepted := true
			if len(data.Reservations[ro.GetNumber()]) > 0 {
				for _, j := range data.Reservations[ro.GetNumber()] {
					if !agt.RoomIsFree(*data.Rooms[j.IdRoom], j, request) {
						roomIsAccepted = false
						break
					}
				}
			}
			if roomIsAccepted { // Différentes sélections par rapport aux prefs du réceptionniste
				roomsList = append(roomsList, *ro)
			}
		}
	}

	if r.Pref() == Smallest {
		sort.Slice(roomsList, func(i, j int) bool {
			return roomsList[i].Capacity < roomsList[j].Capacity
		})
	} else if r.Pref() == Expensive {
		sort.Slice(roomsList, func(i, j int) bool {
			return roomsList[i].Price > roomsList[j].Price
		})
	}

	if len(roomsList) > 5 {
		roomsList = roomsList[:5]
	}

	if len(roomsList) <= 0 {
		log.Println("  ↪  ", "Il n'y a pas de chambre disponible", roomsList)
	} else {
		log.Println("  ↪  ", "Je propose", roomsList)
	}
	r.bookRoomResponse(request, roomsList)
}

func (r *Receptionnist) bookRoomResponse(request agt.BookRoomRequest, rooms []room.Room) {
	values := agt.BookRoomResponse{
		Recpt:     r.GetId(),
		Client:    request.IdClient,
		RoomsList: rooms,
	}

	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://localhost:8080/bookRoomResponse", "application/json", bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err, resp)
	}
}

func (r *Receptionnist) validBook(data agt.DataReceptionist, c agt.ChooseRoomRequest) {
	log.Println()
	log.Println("J'ai reçu la validation pour la chambre :", c)
	// Traitement de la validation
	if data.Time.Day > c.Book.BeginDate {
		if len(data.Reservations[c.Rooms[0]]) > 0 {
			for _, j := range data.Reservations[c.Rooms[0]] {
				if agt.RoomIsFree(*data.Rooms[j.IdRoom], j, c.Book) ||
					data.Time.Day > c.Book.BeginDate { // condition inverse de celle L75 dans work()
					values := c.Book
					json_data, err := json.Marshal(values)
					resp, err := http.Post("http://localhost:8080/bookRoomRequest", "application/json", bytes.NewBuffer(json_data))

					if err != nil {
						log.Fatal(err, resp)
					}
					return
				}
			}
		}
	}
	if len(c.Rooms) > 0 {
		reservation := agt.Reservation{
			IdHotel:   c.Book.IdHotel,
			IdClient:  c.Book.IdClient,
			IdRoom:    c.Rooms[0],
			NbPpl:     c.Book.NbPpl,
			BeginDate: c.Book.BeginDate,
			EndDate:   c.Book.EndDate,
		}

		values := agt.ChooseRoomResponse{
			Recpt:       r.GetId(),
			Client:      c.Book.IdClient,
			Room:        *data.Rooms[c.Rooms[0]],
			Reservation: reservation,
		}
		json_data, err := json.Marshal(values)

		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.Post("http://localhost:8080/chooseRoomResponse", "application/json", bytes.NewBuffer(json_data))

		if err != nil {
			log.Fatal(err, resp)
		}
	}
}
