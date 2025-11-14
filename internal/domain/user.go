package domain

type User struct {
	Id       string
	Name     string
	IsActive bool
	TeamId   *int
}

func NewUser(id, name string, isActive bool, teamId *int) *User {
	return &User{
		Id:       id,
		Name:     name,
		IsActive: isActive,
		TeamId:   teamId,
	}
}
