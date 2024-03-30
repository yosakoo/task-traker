package models

type User struct {
	ID       int    `json:"id"`
	Name string 	`json:"username"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}
