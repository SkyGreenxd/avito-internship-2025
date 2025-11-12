package domain

type Team struct {
	Id   string
	Name string
}

type TeamWithUsers struct {
	Team  *Team
	Users []User
}

func NewTeam(id, name string) *Team {
	return &Team{
		Id:   id,
		Name: name,
	}
}
