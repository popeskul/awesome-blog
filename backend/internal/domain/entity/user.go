package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created"`
	UpdatedAt    time.Time `json:"updated"`
}

type UpdateUser struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email,omitempty" validate:"omitempty,email"`
	Password string    `json:"password,omitempty" validate:"omitempty,min=6"`
	Username string    `json:"username,omitempty" validate:"omitempty"`
}

type NewUser struct {
	Email        string `json:"email" validate:"required,email"`
	PasswordHash string `json:"password" validate:"required,min=6"`
	Username     string `json:"username" validate:"required"`
}
