package model

type User struct {
	Id int				`json:"id"`
	Username string		`json:"name"`
	Password string		`json:"password"`
	Email string		`jon:"email"`
}