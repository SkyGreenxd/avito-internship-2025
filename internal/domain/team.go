package domain

type Team struct {
	Id   int
	Name string
}

func NewTeam(name string) Team {
	return Team{
		Name: name,
	}
}
