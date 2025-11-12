package domain

type User struct {
	Id       int
	Name     string
	IsActive bool
	TeamId   int
}

func NewUser(name string, isActive bool, teamId int) *User {
	return &User{
		Name:     name,
		IsActive: isActive,
		TeamId:   teamId,
	}
}
