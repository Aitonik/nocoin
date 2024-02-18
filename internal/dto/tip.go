package dto

import (
	"errors"
)

var (
	ErrTipNotFound = errors.New("tip not found")
)

type Tip struct {
	Count        int    `json:"count"`
	RestaurantId string `json:"restaurantId"`
}
