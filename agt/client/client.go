package client

import (
	"IA04-hotel/agt"
	"IA04-hotel/agt/agent"
	"IA04-hotel/agt/room"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/google/go-cmp/cmp"
)

type Pref int

const (
	Cheapest Pref = iota // permet de faire un enum en go
	Biggest
	Smallest
)

type Client struct {
	agent.Agent
	hotel     int
	beginDate int
	endDate   int
	nb        int
	prixMax   int
	pref      Pref
}

func NewClient(id string, firstName string, lastName string, hotel, beginDate, endDate, nb, prixMax int, pref Pref) *Client {
	agent := agent.NewAgent(id, firstName, lastName)
	return &Client{*agent, hotel, beginDate, endDate, nb, prixMax, pref}
}

func (c Client) Hotel() int {
	return c.hotel
}

func (c Client) BeginDate() int {
	return c.beginDate
}
func (c Client) EndDate() int {
	return c.endDate
}
func (c Client) Nb() int {
	return c.nb
}

func (c Client) PrixMax() int {
	return c.prixMax
}

func (c Client) Pref() Pref {
	return c.pref
}

func (c *Client) Start() {
	api_url := "http://localhost:8080/"
	go func() {
		// Demande d'une chambre
		//log.Println()
		//log.Println("Je suis", c.Id(), "un client qui veut une chambre")
		values := agt.BookRoomRequest{
			IdClient:  c.GetId(),
			IdHotel:   c.Hotel(),
			BeginDate: c.BeginDate(),
			EndDate:   c.EndDate(),
			NbPpl:     c.Nb(),
			Pref:      int(c.Pref()),
		}
		//log.Println("  ↪  ", "Je veux un chambre de", values.NbPpl, "de personnes du J", values.BeginDate, "au J", values.EndDate)
		json_data, err := json.Marshal(values)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.Post(api_url+"bookRoomRequest", "application/json", bytes.NewBuffer(json_data))
		if err != nil {
			log.Fatal(err, resp)
		}

		// Reception de la liste des chambres disponibles
		var data agt.DataClient
		dataRequest := agt.DataClientRequest{Client: c.GetId()}
		json_data, err = json.Marshal(dataRequest)
		old_data := data
		nb := 0
		for cmp.Equal(old_data, data) {
			time.Sleep(500 * time.Millisecond)
			nb++
			if nb > 20 {
				return
			}
			resp, err := http.Post(api_url+"dataClient", "application/json", bytes.NewBuffer(json_data))
			if err != nil {
				log.Fatal(err, resp)
			}
			json.NewDecoder(resp.Body).Decode(&data)
		}
		//log.Println()
		//log.Println("Je suis", c.Id(), "et j'ai reçu", data.RoomsList)

		var rooms []room.Room
		for _, j := range data.RoomsList {
			if j.GetPrice() <= c.PrixMax() {
				rooms = append(rooms, j)
			}
		}

		if len(rooms) <= 0 {
			return
		}

		// Choix de la chambre
		// Tri de la liste en fonction des préférences du client
		sort.Slice(rooms, func(i, j int) bool {
			switch c.Pref() {
			case 0: //Cheapest
				return rooms[i].Price < rooms[j].Price
			case 1: //Biggest
				return rooms[i].Capacity < rooms[j].Capacity
			case 2: //Smallest
				return rooms[i].Capacity > rooms[j].Capacity
			default:
				return true
			}
		})

		var roomsSorted []int
		for _, j := range rooms {
			roomsSorted = append(roomsSorted, j.GetNumber())
		}
		choice := agt.ChooseRoomRequest{
			Rooms: roomsSorted,
			Book:  values,
		}

		log.Println("Je suis", c.GetId(), "j'ai choisis", roomsSorted)

		// Envoi de la requête
		json_data, err = json.Marshal(choice)
		if err != nil {
			log.Fatal(err)
		}
		resp, err = http.Post(api_url+"chooseRoomRequest", "application/json", bytes.NewBuffer(json_data))
		if err != nil {
			log.Fatal(err, resp)
		}

		log.Println("Fin de la demande de réservation", c.GetId())
	}()
}
