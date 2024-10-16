package agent

type Agent struct {
	Id        string `json:"id"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
}

func NewAgent(id string, firstName string, lastName string) *Agent {
	return &Agent{id, firstName, lastName}
}

func (a Agent) GetId() string {
	return a.Id
}

func (a *Agent) SetId(id string) {
	a.Id = id
}

func (a Agent) GetFirstName() string {
	return a.FirstName
}

func (a *Agent) SetFirstName(firstName string) {
	a.FirstName = firstName
}

func (a Agent) GetLastName() string {
	return a.LastName
}

func (a *Agent) SetLastName(lastName string) {
	a.LastName = lastName
}
