package employee

import (
	"IA04-hotel/agt/agent"
)

type State int
type Day int
type Shift int

const (
	Idle State = iota
	Working
)

const (
	Sunday Day = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

const (
	Noon Shift = iota
	Morning
	Night
)

type Employee struct {
	Job         int   `json:"job"` // 0: rcpt, 1: cleaner
	Salary      int   `json:"salary"`
	State       State `json:"state"`
	IsWorking   bool  `json:"is-working"`
	IdHotel     int   `json:"id-hotel"`
	Schedule    []Day `json:"schedule"` // 0:Dimanche, 1:Lundi, 2:Mardi,....
	Shift       Shift
	agent.Agent `json:"agent"`
}

func NewEmployee(job int, salary int, state State, isWorking bool, idHotel int, schedule []Day, shift Shift, id string, firstName string, lastName string) *Employee {
	agent := agent.NewAgent(id, firstName, lastName)
	return &Employee{job, salary, state, isWorking, idHotel, schedule, shift, *agent}
}

func (e Employee) GetJob() int {
	return e.Job
}

func (e *Employee) SetJob(job int) {
	e.Job = job
}

func (e Employee) GetSalary() int {
	return e.Salary
}

func (e *Employee) SetSalary(salary int) {
	e.Salary = salary
}

func (e Employee) GetState() State {
	return e.State
}

func (e *Employee) SetState(state State) {
	e.State = state
}

func (e Employee) GetIsWorking() bool {
	return e.IsWorking
}

func (e *Employee) SetIsWorking(isWorking bool) {
	e.IsWorking = isWorking
}

func (e Employee) GetIdHotel() int {
	return e.IdHotel
}

func (e *Employee) SetIdHotel(idHotel int) {
	e.IdHotel = idHotel
}

func (e Employee) GetSchedule() []Day {
	return e.Schedule
}

func (e *Employee) SetSchedule(schedule []Day) {
	e.Schedule = schedule
}

func (e Employee) GetShift() Shift {
	return e.Shift
}

func (e *Employee) SetShift(shift Shift) {
	e.Shift = shift
}
