package main

import (
	"IA04-hotel/agt"
	"IA04-hotel/agt/employee"
	"IA04-hotel/agt/room"
	"IA04-hotel/agt/server"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func main() {
	logDebug := false
	defer func() {
		if r := recover(); r != nil {
			// Get the stack trace for the current goroutine
			stackTrace := debug.Stack()

			// Print the stack trace, but skip over goroutines from the net/http package
			for _, line := range strings.Split(string(stackTrace), "\n") {
				log.Println(line)
			}
		}
	}()
	f, err := os.OpenFile("log4.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	if logDebug {
		log.SetOutput(f)
	}

	nb_rcpt := 5
	nb_cleaner := 5
	//nb_client := 1000
	nb_room := 50
	money := 10000
	priceRoomTwoPpl := 60
	rand.Seed(time.Now().UnixNano())
	//lancement du serveur
	port := "8080"
	go server.LaunchServ(port)
	//time.Sleep(200 * time.Millisecond)
	// Création et envoi de la requête de l'hôtel
	nbFloors := 3            // nombre d'étages
	var roomList []room.Room // liste des chambres de l'hôtel
	for i := 0; i < nb_room; i++ {
		number := i
		capacity := rand.Intn(3)*2 + 2
		var price int
		switch capacity {
		case 2:
			price = priceRoomTwoPpl
		case 4:
			price = int(float64(priceRoomTwoPpl) * 1.5)
		case 6:
			price = priceRoomTwoPpl * 2
		default:
			price = 100
		}
		state := room.Free
		roomList = append(roomList, room.Room{Number: number, Capacity: capacity, Price: price, State: state})
	}
	var employeeList []struct {
		Fct int
		employee.Employee
	}
	hotel := agt.CreateHotelRequest{
		NbRooms:      nb_room,
		NbFloors:     nbFloors,
		NbEmployees:  nb_cleaner + nb_rcpt,
		Money:        money,
		RoomList:     roomList,
		EmployeeList: employeeList,
	}
	json_data, err := json.Marshal(hotel)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post("http://localhost:"+port+"/createHotel", "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err, resp)
	}
	var res map[string]int
	json.NewDecoder(resp.Body).Decode(&res)
	//log.Println("L'hôtel", res["id-hotel"], "a été créé.")

	for i := 0; i < nb_cleaner; i++ {
		// Création et envoi de la requête pour l'ajout d'un cleaner
		salary := 2000
		firstName := "Brigitte"
		lastName := "Lambert"
		var schedule = make([]employee.Day, 5)
		sl := rand.Perm(7)[:5] //prendre 5 jours au hasard dans la semaine
		for i, d := range sl {
			schedule[i] = employee.Day(d)
		}
		shift := employee.Shift(rand.Perm(3)[0])
		cleaner := agt.CreateCleanerRequest{
			FirstName: firstName,
			LastName:  lastName,
			State:     employee.Idle,
			Salary:    salary,
			Schedule:  schedule,
			Shift:     shift,
			IdHotel:   1,
		}
		json_data, err = json.Marshal(cleaner)
		if err != nil {
			log.Fatal(err)
		}
		resp, err = http.Post("http://localhost:"+port+"/createCleaner", "application/json",
			bytes.NewBuffer(json_data))
		if err != nil {
			log.Fatal(err, resp)
		}
	}
	// Création et envoi de la requête pour l'ajout d'un réceptionniste
	for i := 0; i < nb_rcpt; i++ {
		salary := 1750
		firstName := "Luke"
		lastName := "Skywalker"
		var schedule = make([]employee.Day, 5)
		sl := rand.Perm(7)[:5] //prendre 5 jours au hasard dans la semaine
		for i, d := range sl {
			schedule[i] = employee.Day(d)
		}
		shift := employee.Shift(rand.Perm(3)[0])
		receptionnist := agt.CreateReceptionnistRequest{
			FirstName: firstName,
			LastName:  lastName,
			State:     employee.Idle,
			Salary:    salary,
			Schedule:  schedule,
			Shift:     shift,
			IdHotel:   1,
			Pref:      rand.Intn(3),
		}
		json_data, err = json.Marshal(receptionnist)
		if err != nil {
			log.Fatal(err)
		}
		resp, err = http.Post("http://localhost:"+port+"/createReceptionnist", "application/json",
			bytes.NewBuffer(json_data))
		if err != nil {
			log.Fatal(err, resp)
		}
	}

	// Création des clients
	var begin int
	var i int
	for {
		i++
		time.Sleep(500 * time.Millisecond)
		firstName := "Michel"
		lastName := "Dupont"
		begin = rand.Intn(3) + i/18
		pref := rand.Perm(3)[0]
		client := agt.CreateClientRequest{
			FirstName: firstName,
			LastName:  lastName,
			Hotel:     1,
			BeginDate: begin,
			EndDate:   begin + rand.Intn(5) + 1,
			Nb:        rand.Intn(8) + 1,
			PrixMax:   (rand.Intn(5) + 4) * 20,
			Pref:      pref,
		}
		json_data, err = json.Marshal(client)
		if err != nil {
			log.Fatal(err)
		}
		resp, err = http.Post("http://localhost:"+port+"/createClient", "application/json",
			bytes.NewBuffer(json_data))
		if err != nil {
			log.Fatal(err, resp)
		}
		//log.Println("Le client", client.Id, "a été créé.")
	}
	fmt.Scanln()
}
