package models

import (
	"strings"

	"github.com/pedrotunin/go-jwt-auth/internal/validators"
)

type UserID = int
type UserEmail = string
type UserPassword = string
type UserStatus = string

type User struct {
	ID       UserID
	Email    UserEmail
	Password UserPassword
	Status   UserStatus
}

func NewUser(email, password string) (*User, error) {
	password = strings.TrimSpace(password)

	if err := validators.IsValidEmail(email); err != nil {
		return nil, err
	}

	if err := validators.IsValidPassword(password); err != nil {
		return nil, err
	}

	return &User{
		Email:    email,
		Password: password,
	}, nil
}
