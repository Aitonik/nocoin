package domain

import (
	"errors"
	"time"
)

var (
	ErrTipNotFound = errors.New("tip not found")
)

type Tip struct {
	ID           string    `json:"id"`
	Count        int       `json:"count"`
	Transaction  string    `json:"transaction"`
	Status       string    `json:"status"`
	RestaurantId string    `json:"restaurantId"`
	CreateDate   time.Time `json:"createDate"`
}
