package agt

import (
	"IA04-hotel/agt/employee"
	"IA04-hotel/agt/hotel"
	"IA04-hotel/agt/room"
)

type Time struct {
	Day  int `json:"day"`
	Hour int `json:"hour"`
}

func RoomIsFree(room room.Room, res Reservation, request BookRoomRequest) bool {
	if !((res.BeginDate <= request.BeginDate && request.BeginDate < res.EndDate) ||
		(res.BeginDate < request.EndDate && request.EndDate <= res.EndDate)) {
		if room.GetCapacity() >= request.NbPpl {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func ResaIsOk(res Reservation, request Reservation, day int) bool {
	if !((res.BeginDate <= request.BeginDate && request.BeginDate < res.EndDate) ||
		(res.BeginDate < request.EndDate && request.EndDate <= res.EndDate)) &&
		day <= request.BeginDate {
		return true
	} else {
		return false
	}
}

type CreateHotelRequest struct {
	NbRooms      int         `json:"nb-rooms"`
	NbFloors     int         `json:"nb-floors"`
	NbEmployees  int         `json:"nb-employees"`
	Money        int         `json:"money"`
	RoomList     []room.Room `json:"room-list"`
	EmployeeList []struct {
		Fct int
		employee.Employee
	} `json:"employees-list"`
}

type CreateHotelResponse struct {
	IdHotel int `json:"id-hotel"`
}

type CreateRoomRequest struct {
	IdHotel  int `json:"id-hotel"`
	Number   int `json:"number"`
	Capacity int `json:"capacity"`
	Price    int `json:"price"`
	State    int `json:"state"` // 0: Free, 1: Reserved, 2: Cleaning, 3: Closed
}

type CreateRoomResponse struct {
	Number int `json:"number"`
}

type CreateClientRequest struct {
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Hotel     int    `json:"hotel"`
	BeginDate int    `json:"begin-date"`
	EndDate   int    `json:"end-date"`
	Nb        int    `json:"#pers"`
	PrixMax   int    `json:"prix-max"`
	Pref      int    `json:"pref"`
}

type CreateClientResponse struct {
	Id string
}

type CreateCleanerRequest struct {
	FirstName string         `json:"first-name"`
	LastName  string         `json:"last-name"`
	State     employee.State `json:"state"`
	Salary    int
	IdHotel   int `json:"id-hotel"`
	Schedule  []employee.Day
	Shift     employee.Shift
}

type CreateCleanerResponse struct {
	Id string
}

type CreateReceptionnistRequest struct {
	FirstName string         `json:"first-name"`
	LastName  string         `json:"last-name"`
	State     employee.State `json:"state"`
	Salary    int
	IdHotel   int `json:"id-hotel"`
	Schedule  []employee.Day
	Shift     employee.Shift
	Pref      int
}

type CreateReceptionnistResponse struct {
	Id string
}

type Reservation struct {
	IdClient  string `json:"id-client"`
	IdHotel   int    `json:"id-hotel"`
	IdRoom    int    `json:"id-room"`
	NbPpl     int    `json:"nb-ppl"`
	BeginDate int    `json:"date-start"`
	EndDate   int    `json:"date-end"`
}

type BookRoomRequest struct {
	IdClient  string `json:"id-client"`
	IdHotel   int    `json:"id-hotel"`
	BeginDate int    `json:"date-start"`
	EndDate   int    `json:"date-end"`
	NbPpl     int    `json:"#pers"`
	Pref      int    `json:"pref"`
}
type BookRoomResponse struct {
	Recpt     string      `json:"id-recpt"`
	Client    string      `json:"id-client"`
	RoomsList []room.Room `json:"rooms-list"`
}

type ChooseRoomRequest struct {
	Rooms []int           `json:"rooms"`
	Book  BookRoomRequest `json:"book"`
}

type ChooseRoomResponse struct {
	Recpt       string      `json:"id-recpt"`
	Client      string      `json:"id-client"`
	Room        room.Room   `json:"room"`
	Reservation Reservation `json:"reservation"`
}

type DataReceptionistRequest struct {
	ID    string `json:"id-employee"`
	Hotel int    `json:"id-hotel"`
}

type DataReceptionist struct {
	Time         Time                        `json:"time"`
	Rooms        map[int]*room.Room          `json:"rooms"`
	Reservations map[int]map[int]Reservation `json:"reservations"`
	Requests     []BookRoomRequest           `json:"requests"`
	Choice       []ChooseRoomRequest         `json:"choice"`
}

type DataClientRequest struct {
	Client string `json:"id-client"`
}

type DataClient struct {
	RoomsList []room.Room `json:"rooms-list"`
}

type DataCleaner struct {
	CleaningList map[int]map[int]*room.Room `json:"cleaning-list"`
}

type CleanRoomRequest struct {
	IdHotel   int    `json:"id-hotel"`
	IdRoom    int    `json:"id-room"`
	IdCleaner string `json:"id-cleaner"`
}

type CleanRoomResponse struct {
	IdHotel int `json:"id-hotel"`
	IdRoom  int `json:"id-room"`
}
type Data struct {
	Time         Time                          `json:"time"`
	Hotel        hotel.Hotel                   `json:"hotel"`
	Rooms        map[int]map[int]*room.Room    `json:"rooms"`
	Reservations map[int]map[int]Reservation   `json:"reservations"`
	Agents       map[string]*employee.Employee `json:"agents"`
}
