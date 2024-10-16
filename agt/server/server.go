package server

import (
	"IA04-hotel/agt"
	"IA04-hotel/agt/cleaner"
	"IA04-hotel/agt/client"
	"IA04-hotel/agt/employee"
	"IA04-hotel/agt/hotel"
	"IA04-hotel/agt/receptionnist"
	"IA04-hotel/agt/room"
	"IA04-hotel/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var day, hour, old_day int
var date agt.Time
var id_resa int
var id_empl int
var id_client int
var hotelsMap = make(map[int]*hotel.Hotel)                    // hotel.id : hotel
var roomsMap = make(map[int]map[int]*room.Room)               // hotel.id : room.number : room
var employeesMap = make(map[string]*employee.Employee)        // employee.id : employee
var recptsMap = make(map[string]*receptionnist.Receptionnist) // recpt.id : recpt
var cleanersMap = make(map[string]*cleaner.Cleaner)           // cl.id : cl
var clientsMap = make(map[string]agt.DataClient)              // client.id : DataClient (toutes les données relatives au client)
var currWorkingEmployee = make(map[string]*employee.Employee) // employee.id : employee

var requestRoomMap = make(map[string][]agt.BookRoomRequest)  // rcpt.id : []BookRoomRequest (les demandes (en cours) de réservation pour un réceptionniste)
var choiceRoomMap = make(map[string][]agt.ChooseRoomRequest) // rcpt.id : []ChooseRoomRequest
var reservationMap = make(map[int]map[int]agt.Reservation)   // room.number : Reservation (Réservation d'une chambre) | pas encore utilisé
var cleaningMap = make(map[int]map[int]*room.Room)           // hotel.id : room.number : room

var muRoomRequest, muChoice, muRoomMap, muReservation, muCleaning sync.RWMutex
var muRcpt, muClient sync.Mutex

func viewFront(w http.ResponseWriter, r *http.Request) {
	p, _ := ioutil.ReadFile("test_front/index.html")
	w.Write(p)
}

func loadScript(w http.ResponseWriter, r *http.Request) {
	p, _ := ioutil.ReadFile("test_front/sketch.js")
	w.Write(p)
}

func loadStyle(w http.ResponseWriter, r *http.Request) {
	p, _ := ioutil.ReadFile("test_front/style.css")
	w.Write(p)
}

func timeManage(w http.ResponseWriter, r *http.Request) {
	var res agt.Time
	res.Day = day
	res.Hour = hour
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func createHotel(w http.ResponseWriter, r *http.Request) {
	// décodage de la requête
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var req agt.CreateHotelRequest
	err := json.Unmarshal(buf.Bytes(), &req)
	log.Println("Demande de création d'un hôtel")

	// vérification que les données soient bonnes et qu'il n'y a pas d'erreurs
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	if req.NbRooms == 0 || req.NbEmployees == 0 || req.NbFloors == 0 { // pas nécessaire d'ajouter des chambres
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New("missing required field"))
		return
	}
	// création de l'hôtel
	hotel := hotel.NewHotel(len(hotelsMap)+1, req.NbRooms, req.NbRooms, req.NbEmployees, req.Money, req.RoomList, req.EmployeeList)
	hotelsMap[hotel.Id()] = hotel
	log.Println("  ↪  ", "L'hôtel "+fmt.Sprint(hotel.Id()), "a été créé avec", hotel.NbRooms(), hotel.NbFloors(), hotel.NbEmployees(), hotel.RoomList(), hotel.EmployeeList())

	if len(req.RoomList) > 0 {
		for _, j := range req.RoomList {
			// création de la chambre
			if _, ok := roomsMap[hotel.Id()][j.Number]; !ok {
				ro := room.NewRoom(j.Number, j.Capacity, j.Price, room.State(j.State))
				hotelsMap[hotel.Id()].AddRoom(*ro)
				if roomsMap[hotel.Id()] == nil {
					roomsMap[hotel.Id()] = make(map[int]*room.Room)
				}
				roomsMap[hotel.Id()][ro.GetNumber()] = ro
				muReservation.Lock()
				reservationMap[ro.GetNumber()] = make(map[int]agt.Reservation)
				muReservation.Unlock()
				log.Println("  ↪  ", "La chambre "+fmt.Sprint(ro.GetNumber()), "a été créé avec", ro.GetCapacity(), ro.GetState())
			}
		}
	}

	if len(req.EmployeeList) > 0 { // Pour l'instant ne fonctionne pas
		for _, j := range req.EmployeeList {
			id_empl++
			if j.Fct == 0 {
				// création du réceptionniste
				recpt := receptionnist.NewReceptionnist(j.Employee.GetSalary(), j.Employee.GetIdHotel(), j.Employee.GetSchedule(), j.Employee.GetShift(), "empl"+fmt.Sprint(id_empl), j.Employee.GetFirstName(), j.Employee.GetLastName(), receptionnist.Pref(0))
				employeesMap[recpt.GetId()] = &recpt.Employee
				log.Println("  ↪  ", "Le réceptionniste "+recpt.GetId(), "a été créé avec un salaire de", recpt.GetSalary(), "euros, et avec l'emploi du temps ", recpt.GetSchedule(), "et avec le shift ", recpt.GetShift())
				recpt.Start()
			} else if j.Fct == 1 {
				// création du cleaner
				id_empl++
				cl := cleaner.NewCleaner(j.Employee.GetSalary(), j.Employee.GetIdHotel(), j.Employee.GetSchedule(), j.Employee.GetShift(), "empl"+fmt.Sprint(id_empl), j.Employee.GetFirstName(), j.Employee.GetLastName())
				employeesMap[cl.GetId()] = &cl.Employee
				log.Println("  ↪  ", "L'employé de ménage "+cl.GetId(), "a été créé avec un salaire de", cl.GetSalary(), "euros, et avec l'emploi du temps ", cl.GetSchedule(), "et avec le shift ", cl.GetShift())
			}
		}
	}
	var res agt.CreateHotelResponse
	res.IdHotel = hotel.Id()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func createRoom(w http.ResponseWriter, r *http.Request) {
	// décodage de la requête
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var req agt.CreateRoomRequest
	err := json.Unmarshal(buf.Bytes(), &req)
	log.Println("Demande de création d'une chambre")

	// vérification que les données soient bonnes et qu'il n'y a pas d'erreurs
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	if req.IdHotel == 0 || req.Number == 0 || req.Capacity == 0 { // pas nécessaire d'ajouter l'état (par défaut état libre)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New("missing required field"))
		return
	}

	// création de la chambre
	if _, ok := roomsMap[req.IdHotel][req.Number]; !ok {
		ro := room.NewRoom(req.Number, req.Capacity, req.Price, room.State(req.State))
		hotelsMap[req.IdHotel].AddRoom(*ro)
		if roomsMap[req.IdHotel] == nil {
			roomsMap[req.IdHotel] = make(map[int]*room.Room)
		}
		roomsMap[req.IdHotel][ro.GetNumber()] = ro
		reservationMap[ro.GetNumber()] = make(map[int]agt.Reservation)
		log.Println("  ↪  ", "La chambre "+fmt.Sprint(ro.GetNumber()), "a été créé avec", ro.GetCapacity(), ro.GetState())

		var res agt.CreateRoomResponse
		res.Number = ro.GetNumber()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func createReceptionnist(w http.ResponseWriter, r *http.Request) {
	// décodage de la requête
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var req agt.CreateReceptionnistRequest
	err := json.Unmarshal(buf.Bytes(), &req)
	log.Println("Demande de création d'un réceptionniste")

	// vérification que les données soient bonnes et qu'il n'y a pas d'erreurs
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	if req.FirstName == "" || req.LastName == "" || req.Salary <= 0 || req.IdHotel < 0 || req.Schedule == nil || req.Shift < 0 || req.Pref < 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New("missing or invalid required field"))
		return
	}

	// création du réceptionniste
	id_empl++
	recpt := receptionnist.NewReceptionnist(req.Salary, req.IdHotel, req.Schedule, req.Shift, "empl"+fmt.Sprint(id_empl), req.FirstName, req.LastName, receptionnist.Pref(req.Pref))
	employeesMap[recpt.GetId()] = &recpt.Employee
	recptsMap[recpt.GetId()] = recpt
	currWorkingEmployee[recpt.GetId()] = &recpt.Employee
	log.Println("  ↪  ", "Le réceptionniste "+recpt.GetId(), "a été créé avec un salaire de", recpt.GetSalary(), "euros, et avec l'emploi du temps ", recpt.GetSchedule(), "et avec le shift, ", recpt.GetShift(), recpt.GetJob())
	recpt.Start()
	// envoi de l'id du réceptionniste
	var res agt.CreateReceptionnistResponse
	res.Id = recpt.GetId()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func createCleaner(w http.ResponseWriter, r *http.Request) {
	// décodage de la requête
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var req agt.CreateCleanerRequest
	err := json.Unmarshal(buf.Bytes(), &req)
	log.Println("Demande de création d'un employé de ménage")

	// vérification que les données soient bonnes et qu'il n'y a pas d'erreurs
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	if req.FirstName == "" || req.LastName == "" || req.Salary <= 0 || req.Schedule == nil || req.Shift < 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New("missing or invalid required field"))
		return
	}

	// création du cleaner
	id_empl++
	cl := cleaner.NewCleaner(req.Salary, req.IdHotel, req.Schedule, req.Shift, "empl"+fmt.Sprint(id_empl), req.FirstName, req.LastName)
	employeesMap[cl.GetId()] = &cl.Employee
	cleanersMap[cl.GetId()] = cl
	currWorkingEmployee[cl.GetId()] = &cl.Employee
	log.Println("  ↪  ", "L'employé de ménage "+cl.GetId(), "a été créé avec un salaire de", cl.GetSalary(), "euros, et avec l'emploi du temps ", cl.GetSchedule(), "et avec le shift ", cl.GetShift(), cl.GetJob())
	cl.Start()
	// envoi de l'id du cleaner
	var res agt.CreateCleanerResponse
	res.Id = cl.GetId()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func createClient(w http.ResponseWriter, r *http.Request) {
	// décodage de la requête
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var req agt.CreateClientRequest
	err := json.Unmarshal(buf.Bytes(), &req)
	//log.Println("Demande de création d'un client")

	// vérification que les données soient bonnes et qu'il n'y a pas d'erreurs
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	if req.FirstName == "" || req.LastName == "" || req.BeginDate < day || req.EndDate < day || req.BeginDate > req.EndDate || req.PrixMax <= 0 || req.Nb <= 0 || req.Pref < 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New("missing or invalid required field"))
		return
	}

	// création du client
	id_client++
	client := client.NewClient("clt"+fmt.Sprint(id_client), req.FirstName, req.LastName, req.Hotel, req.BeginDate, req.EndDate, req.Nb, req.PrixMax, client.Pref(req.Pref))
	//log.Println("  ↪  ", "Le client "+fmt.Sprint(client.Id()), "a été créé")
	defer client.Start()
	// envoi de l'id du client
	var res agt.CreateClientResponse
	res.Id = client.GetId()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func bookRoomRequest(w http.ResponseWriter, r *http.Request) {
	// Décoder la demande de réservation envoyée par le client
	var req agt.BookRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// Vérifier que la demande de réservation est valide
	if req.IdClient == "" || req.IdHotel <= 0 || req.BeginDate < 0 || req.EndDate < req.BeginDate || req.NbPpl == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New("missing or invalid required field"))
		return
	}

	// Sélectionner le réceptionniste qui s'occupera de la demande de réservation
	var recptID string
	muRcpt.Lock()
	for id, j := range recptsMap {
		if j.GetState() == employee.Idle && j.GetIsWorking() {
			recptID = id
			break
		}
	}
	muRcpt.Unlock()

	// S'il n'y a pas de réceptionniste disponible, sélectionner celui qui a le moins de demandes en cours
	if recptID == "" {
		min := 1000000
		muRcpt.Lock()
		for id, j := range recptsMap {
			if !j.GetIsWorking() && j.GetState() == employee.Working {
				j.SetState(employee.Idle)
			}
			if recptsMap[id].GetIsWorking() {
				muRoomRequest.Lock()
				length := len(requestRoomMap[id])
				muRoomRequest.Unlock()
				if length < min {
					min = length
					recptID = id
				}
			}
		}
		muRcpt.Unlock()
	}

	// S'il n'y a toujours pas de réceptionniste disponible, sélectionner aléatoirement un réceptionniste
	if recptID == "" {
		muRcpt.Lock()
		for id := range recptsMap {
			if recptsMap[id].GetIsWorking() {
				recptID = id
				break
			}

		}
		muRcpt.Unlock()
	}
	// Mettre à jour l'état du réceptionniste et ajouter la demande de réservation à sa liste de demandes en cours
	muRcpt.Lock()
	if recptID != "" {
		if recptsMap[recptID].GetIsWorking() {
			recptsMap[recptID].SetState(employee.Working)
			muRoomRequest.Lock()
			requestRoomMap[recptID] = append(requestRoomMap[recptID], req)
			muRoomRequest.Unlock()
		} else {
			recptsMap[recptID].SetState(employee.Idle)
			recptID = ""
			muRcpt.Lock()
			for id := range recptsMap {
				if recptsMap[id].GetIsWorking() {
					recptID = id
					break
				}
			}
			muRcpt.Unlock()
			if recptID != "" {
				recptsMap[recptID].SetState(employee.Working)
				muRoomRequest.Lock()
				requestRoomMap[recptID] = append(requestRoomMap[recptID], req)
				muRoomRequest.Unlock()
			}
		}
	}
	muRcpt.Unlock()
	// Envoyer une réponse au client indiquant que la demande de réservation a été acceptée
	w.WriteHeader(http.StatusAccepted)
}

func bookRoomResponse(w http.ResponseWriter, r *http.Request) {
	//le client veut réserver une chambre pour nbPpl personnes  du beginDate au endDate.
	var req agt.BookRoomResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// Mise à jour de l'état du réceptionniste et de la liste de demandes de chambres
	muRcpt.Lock()
	recptsMap[req.Recpt].SetState(employee.Idle)
	muRcpt.Unlock()
	muRoomRequest.Lock()
	requestRoomMap[req.Recpt] = requestRoomMap[req.Recpt][1:]
	muRoomRequest.Unlock()

	// Mise à jour des informations du client
	if len(req.RoomsList) > 0 {
		var client agt.DataClient
		client.RoomsList = req.RoomsList
		muClient.Lock()
		clientsMap[req.Client] = client
		muClient.Unlock()
	} else {
		muClient.Lock()
		delete(clientsMap, req.Client)
		muClient.Unlock()
	}

	//log.Println("Réponse de réservation d'une chambre")
	w.WriteHeader(http.StatusOK)
}

func chooseRoomRequest(w http.ResponseWriter, r *http.Request) {
	// Le client choisit sa chambre
	var req agt.ChooseRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// Vérification des paramètres de la requête
	if len(req.Rooms) == 0 || req.Book.IdClient == "" || req.Book.IdHotel <= 0 || req.Book.BeginDate < 0 || req.Book.EndDate < req.Book.BeginDate || req.Book.NbPpl == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New("missing or invalid required field"))
		return
	}

	// Récupération du réceptionniste disponible ou ayant le moins de demandes en cours
	recpt := recptsMap[""] // réceptionniste par défaut (au cas où aucun réceptionniste n'est disponible)
	minRequests := 1000000
	for id, j := range recptsMap {
		muRcpt.Lock()
		state := j.GetState()
		muRcpt.Unlock()
		if state == employee.Idle {
			recpt = j
			break
		} else {
			muRoomRequest.Lock()
			numRequests := len(choiceRoomMap[id])
			muRoomRequest.Unlock()
			if numRequests < minRequests {
				minRequests = numRequests
				recpt = j
			}
		}
	}

	recpt.SetState(employee.Working)

	muChoice.Lock()
	choiceRoomMap[recpt.GetId()] = append(choiceRoomMap[recpt.GetId()], req)
	muChoice.Unlock()
	w.WriteHeader(http.StatusAccepted)
}

func chooseRoomResponse(w http.ResponseWriter, r *http.Request) {
	//le client veut réserver une chambre pour nbPpl personnes  du beginDate au endDate.
	var req agt.ChooseRoomResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("err", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// Mise à jour de l'état du réceptionniste et de l'index de la demande de chambre
	muRcpt.Lock()
	recptsMap[req.Recpt].SetState(employee.Idle)
	muRcpt.Unlock()
	var idx int
	muChoice.Lock()
	for i, j := range choiceRoomMap[req.Recpt] {
		if j.Book.IdClient == req.Client {
			idx = i
			break
		}
	}
	muChoice.Unlock()

	// Vérification de la disponibilité de la chambre
	roomIsAccepted := true
	muReservation.RLock()
	reservation, ok := reservationMap[req.Room.GetNumber()]
	muReservation.RUnlock()
	if ok {
		for _, j := range reservation {
			if !agt.ResaIsOk(j, req.Reservation, day) {
				roomIsAccepted = false
				var values agt.ChooseRoomRequest
				muChoice.Lock()
				if len(choiceRoomMap[req.Recpt]) > 0 {
					if len(choiceRoomMap[req.Recpt][idx].Rooms) > 0 {
						values = agt.ChooseRoomRequest{
							Book:  choiceRoomMap[req.Recpt][idx].Book,
							Rooms: choiceRoomMap[req.Recpt][idx].Rooms[1:],
						}
						choiceRoomMap[req.Recpt][idx] = values
					}
				}
				muChoice.Unlock()
				break
			}
		}
	} else if req.Reservation.BeginDate >= day {
		roomIsAccepted = true
	} else {
		roomIsAccepted = false
	}

	// Réservation de la chambre
	if roomIsAccepted {
		id_resa++
		muReservation.Lock()
		reservationMap[req.Room.GetNumber()][id_resa] = req.Reservation
		muReservation.Unlock()
		log.Println("La chambre", req.Room, "a bien été réservé pour", req.Client)
		w.WriteHeader(http.StatusOK)
	}
}

func cleanRoomResponse(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var req agt.CleanRoomResponse
	err := json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	ro := roomsMap[req.IdHotel][req.IdRoom]
	log.Println("Je nettoie ", ro)                                  // rajouter l'id du cleaner
	time.Sleep(time.Duration(ro.Capacity) * time.Millisecond * 150) //le temps de nettoyage est set à 150 * la capacité de la chambre, en millisecondes (~= 20 min pour une chambre de 2 pers)
	muRoomMap.Lock()
	roomsMap[req.IdHotel][req.IdRoom].SetState(0)
	muRoomMap.Unlock()
	muCleaning.Lock()
	delete(cleaningMap[req.IdHotel], req.IdRoom)
	muCleaning.Unlock()
	w.WriteHeader(http.StatusOK)
}

func getDataReceptionist(w http.ResponseWriter, r *http.Request) {
	// Traitement de la requête
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var req agt.DataReceptionistRequest
	err := json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error(), req)
		return
	}
	// Envoi des données
	muReservation.Lock()
	reservations := make(map[int]map[int]agt.Reservation)
	for k, v := range reservationMap {
		reservations[k] = make(map[int]agt.Reservation)
		for i, j := range v {
			reservations[k][i] = j
		}
	}
	muReservation.Unlock()
	muRoomMap.Lock()
	rooms := make(map[int]*room.Room)
	for k, v := range roomsMap[req.Hotel] {
		rooms[k] = v
	}
	muRoomMap.Unlock()
	muChoice.Lock()
	choices := make([]agt.ChooseRoomRequest, 0)
	choices = choiceRoomMap[req.ID]
	choiceRoomMap[req.ID] = nil
	muChoice.Unlock()
	muRoomRequest.Lock()
	res := agt.DataReceptionist{
		Time:         date,
		Reservations: reservations,
		Rooms:        rooms,
		Requests:     requestRoomMap[req.ID],
		Choice:       choices,
	}
	muRoomRequest.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(res)

	if err != nil {
		fmt.Println(err)
	}

}

func getDataClient(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var req agt.DataClientRequest
	err := json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	muClient.Lock()
	res := &agt.DataClient{
		RoomsList: clientsMap[req.Client].RoomsList,
	}
	delete(clientsMap, req.Client)
	muClient.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func getDataCleaner(w http.ResponseWriter, r *http.Request) {
	res := &agt.DataCleaner{
		CleaningList: cleaningMap,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func getData(w http.ResponseWriter, r *http.Request) {
	// Envoi des données
	//log.Println("Demande d'envoi de données")
	res := &agt.Data{
		Time:         agt.Time{Day: day, Hour: hour},
		Hotel:        *hotelsMap[1],
		Reservations: reservationMap,
		Rooms:        roomsMap,
		Agents:       employeesMap,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func update(day, hour int) {
	// update de l'état des employés
	for _, employee := range employeesMap {
		if day%7 == 0 && hour == 23 {
			hotelsMap[employee.IdHotel].SetMoney(hotelsMap[employee.IdHotel].GetMoney() - employee.GetSalary()/4)
		}
		if utils.CheckWorkingSchedule(employee, day, hour) {
			employee.SetIsWorking(true)
			//relancer les agents s'ils se mettent à travailler alors qu'ils ne travaillaient pas avant
			if _, ok := currWorkingEmployee[employee.GetId()]; !ok {
				//vérifier s'il s'agit d'un cleaner ou d'un réceptionniste
				if cl, ok := cleanersMap[employee.GetId()]; ok {
					cl.Start()
				}
				if recpt, ok := recptsMap[employee.GetId()]; ok {
					recpt.Start()
				}
				//ajouter l'employé dans la liste des employés travaillant
				currWorkingEmployee[employee.GetId()] = employee
			}
		} else {
			employee.SetIsWorking(false)
			//enlever l'employé de la liste des employés travaillant
			delete(currWorkingEmployee, employee.GetId())
		}
	}

	// update de l'état de l'hotel
	date = agt.Time{Day: day, Hour: hour}
	for r, reservations := range reservationMap {
		for id_resa, resa := range reservations {
			//log.Println("Pour la chambre", r, ",il a", resa, id_resa)

			if resa.EndDate == day && hour == 11 && roomsMap[resa.IdHotel][resa.IdRoom].GetState() == room.Reserved {
				//oldResList = append(oldResList, reservationMap[r][id_resa])
				muReservation.Lock()
				delete(reservationMap[r], id_resa)
				muReservation.Unlock()
				muRoomMap.Lock()
				roomsMap[resa.IdHotel][resa.IdRoom].SetState(2)
				muRoomMap.Unlock()
				// ajouter la chambre dans la liste des chambres à nettoyer
				ro := roomsMap[resa.IdHotel][resa.IdRoom]
				muCleaning.Lock()
				if cleaningMap[resa.IdHotel] == nil { //initialiser la map interne si elle n'existe pas encore (premier passage)
					cleaningMap[resa.IdHotel] = make(map[int]*room.Room)
				}
				cleaningMap[resa.IdHotel][resa.IdRoom] = ro
				muCleaning.Unlock()
			}
			if resa.EndDate < day {
				muReservation.Lock()
				delete(reservationMap[r], id_resa)
				muReservation.Unlock()
			}
			if resa.BeginDate <= day && hour >= 18 {
				if resa.BeginDate == day && hour == 18 {
					log.Println("Cette réservation commence", resa)
				}
				if hour == 18 {
					hotelsMap[resa.IdHotel].SetMoney(hotelsMap[resa.IdHotel].GetMoney() + roomsMap[resa.IdHotel][resa.IdRoom].Price)
				}
				if roomsMap[resa.IdHotel][resa.IdRoom].GetState() == room.Free {
					muRoomMap.Lock()
					roomsMap[resa.IdHotel][resa.IdRoom].SetState(1)
					muRoomMap.Unlock()
					break
				} else if roomsMap[resa.IdHotel][resa.IdRoom].GetState() != room.Reserved {
					//log.Println("Impossible la chambre n'est pas libre")
				}
			}
		}
	}
	if old_day != day {
		log.Println("----- Jour", day, "-----")
		log.Println("Etat des chambres :")
		for _, r := range roomsMap[1] {
			hotelsMap[1].SetMoney(hotelsMap[1].GetMoney() - 20)
			log.Println("", r.GetNumber(), " capacity ", r.GetCapacity(), " price ", r.GetPrice(), ":", r.GetState())
		}
		log.Println("Liste des réservation :", reservationMap)
		//log.Println("Old res :", oldResList)
		log.Println()
		old_day = day
	}
}

var url string // Port du serveur
func LaunchServ(port string) {
	url = ":" + port
	day = 0
	// Création du multiplexeur
	mux := http.NewServeMux()
	mux.HandleFunc("/time", timeManage)
	mux.HandleFunc("/createHotel", createHotel)
	mux.HandleFunc("/createRoom", createRoom)
	mux.HandleFunc("/createClient", createClient)
	mux.HandleFunc("/createReceptionnist", createReceptionnist)
	mux.HandleFunc("/createCleaner", createCleaner)
	mux.HandleFunc("/bookRoomRequest", bookRoomRequest)
	mux.HandleFunc("/bookRoomResponse", bookRoomResponse)
	mux.HandleFunc("/chooseRoomRequest", chooseRoomRequest)
	mux.HandleFunc("/chooseRoomResponse", chooseRoomResponse)
	mux.HandleFunc("/dataReceptionist", getDataReceptionist)
	mux.HandleFunc("/dataClient", getDataClient)
	mux.HandleFunc("/dataCleaner", getDataCleaner)
	mux.HandleFunc("/cleanRoomResponse", cleanRoomResponse)
	mux.HandleFunc("/data", getData)
	mux.HandleFunc("/view", viewFront)
	mux.HandleFunc("/script", loadScript)
	mux.HandleFunc("/style", loadStyle)

	//ajouter les routes d'appels à l'API

	// Démarage du serveur
	s := &http.Server{
		Addr:           url,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	log.Println("démarrage du serveur...")

	// Attente de requêtes
	go s.ListenAndServe()
	log.Println("----- Jour", day, "-----")
	for {
		time.Sleep(500 * time.Millisecond)
		hour++
		update(day, hour)
		if hour == 24 {
			hour = 0
			day++
		}
	}
	// Ctrl + C pour arrêter le serveur
}
