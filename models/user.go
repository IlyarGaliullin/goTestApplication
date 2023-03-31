package models

type User struct {
	Id       int
	Email    string
	Name     string
	Password string
	Roles    []Role
}
