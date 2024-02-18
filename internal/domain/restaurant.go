package domain

import (
	"errors"
)

var (
	ErrBookNotFound = errors.New("book not found")
)

type Restaurant struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerId string `json:"ownerId"`
}

type UpdateRestaurantInput struct {
	Name *string `json:"name"`
}
