package cleaner

import (
	"IA04-hotel/agt"
	employee "IA04-hotel/agt/employee"
	"IA04-hotel/agt/room"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Cleaner struct {
	employee.Employee
}

func NewCleaner(salary int, idHotel int, schedule []employee.Day, shift employee.Shift, id string, firstName string, lastName string) *Cleaner {
	employe := employee.NewEmployee(1, salary, employee.Idle, true, idHotel, schedule, shift, id, firstName, lastName)
	return &Cleaner{*employe}
}

func (c *Cleaner) Start() {
	api_url := "http://localhost:8080/"
	go func() {
		for c.GetIsWorking() {
			time.Sleep(100 * time.Millisecond)
			data := agt.DataCleaner{}
			res, _ := http.Get(api_url + "dataCleaner")
			json.NewDecoder(res.Body).Decode(&data)
			var toClean room.Room
			if len(data.CleaningList[c.GetIdHotel()]) != 0 {
				log.Println(c.Agent.GetFirstName(), c.Agent.GetLastName(), "nettoie une chambre !")
				for _, room := range data.CleaningList[c.GetIdHotel()] {
					toClean = *room //récupère aléatoirement une chambre à nettoyer dans la map
					break
				}
				if toClean.GetState() == 2 { //vérifier malgré tout que la chambre est bien dans l'état cleaning, afin de bloquer tout cleaner qui aurait récupéré la donnée au mauvais moment
					c.SetState(employee.Working)
					var res agt.CleanRoomResponse
					res.IdHotel = c.GetIdHotel()
					res.IdRoom = toClean.Number
					json_data, err := json.Marshal(res)
					resp, err := http.Post(api_url+"cleanRoomResponse", "application/json", bytes.NewBuffer(json_data))
					if err != nil {
						log.Fatal(err, resp)
					}
				}
			} else {
				c.SetState(employee.Idle)
			}
		}
	}()
}
