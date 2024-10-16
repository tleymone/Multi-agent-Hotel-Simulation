package room

type State int

const (
	Free State = iota // permet de faire un enum en go
	Reserved
	Cleaning
	Closed
)

type Room struct {
	Number   int
	Capacity int
	Price    int
	State    State
}

func NewRoom(number int, capacity int, price int, state State) *Room {
	return &Room{number, capacity, price, state}
}

func (r Room) GetNumber() int {
	return r.Number
}

func (r *Room) SetNumber(number int) {
	r.Number = number
}

func (r Room) GetCapacity() int {
	return r.Capacity
}

func (r *Room) SetCapacity(capacity int) {
	r.Capacity = capacity
}

func (r Room) GetPrice() int {
	return r.Price
}

func (r *Room) SetPrce(price int) {
	r.Price = price
}

func (r Room) GetState() State {
	return r.State
}

func (r *Room) SetState(state State) {
	r.State = state
}
