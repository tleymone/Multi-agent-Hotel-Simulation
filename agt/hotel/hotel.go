package hotel

import (
	"IA04-hotel/agt/employee"
	"IA04-hotel/agt/room"
)

type Hotel struct {
	idHotel      int
	nbRooms      int
	nbFloors     int
	nbEmployees  int
	Money        int
	roomList     []room.Room
	employeeList []struct {
		Fct int
		employee.Employee
	}
}

func NewHotel(id int, nbRooms int, nbFloors int, nbEmployees int, monney int, roomList []room.Room, employeeList []struct {
	Fct int
	employee.Employee
}) *Hotel {
	return &Hotel{id, nbRooms, nbFloors, nbEmployees, monney, roomList, employeeList}
}

func (h Hotel) Id() int {
	return h.idHotel
}

func (h Hotel) NbRooms() int {
	return h.nbRooms
}

func (h *Hotel) SetNbRooms(nbRooms int) {
	h.nbRooms = nbRooms
}

func (h Hotel) NbFloors() int {
	return h.nbFloors
}

func (h *Hotel) SetNbFloors(nbFloors int) {
	h.nbFloors = nbFloors
}

func (h Hotel) NbEmployees() int {
	return h.nbEmployees
}

func (h *Hotel) SetMoney(money int) {
	h.Money = money
}

func (h Hotel) GetMoney() int {
	return h.Money
}

func (h *Hotel) SetNbEmployees(nbEmployees int) {
	h.nbEmployees = nbEmployees
}

func (h Hotel) RoomList() []room.Room {
	return h.roomList
}

func (h Hotel) EmployeeList() []struct {
	Fct int
	employee.Employee
} {
	return h.employeeList
}

func (h *Hotel) AddRoom(r room.Room) {
	h.roomList = append(h.roomList, r)
}
