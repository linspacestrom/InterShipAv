package domain

type UserShort struct {
	Id       string
	IsActive bool
}

type User struct {
	Id       string
	Username string
	TeamName string
	IsActive bool
}
