package domain

import (
	"errors"
)

var (
	ErrProfileNotFound = errors.New("profile not found")
)

type Profile struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
