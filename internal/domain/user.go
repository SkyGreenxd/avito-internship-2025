package domain

// TODO: изменить id на стринг и использовать UUID.
type User struct {
	Id       string
	Name     string
	IsActive bool
	TeamId   *string
}

func NewUser(id, name string, isActive bool, teamId *string) *User {
	return &User{
		Id:       id,
		Name:     name,
		IsActive: isActive,
		TeamId:   teamId,
	}
}
