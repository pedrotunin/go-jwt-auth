package models

type UserID = int
type UserEmail = string
type UserPassword = string

type User struct {
	ID       UserID
	Email    UserEmail
	Password UserPassword
}
