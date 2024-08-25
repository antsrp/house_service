package models

import "github.com/google/uuid"

type UserType = string

const (
	Client    UserType = "client"
	Moderator UserType = "moderator"
)

type User struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	UserType UserType  `json:"user_type"`
}
