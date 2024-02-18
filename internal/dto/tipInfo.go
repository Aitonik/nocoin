package dto

type TipInfo struct {
	Tips         []Tip  `json:"tips"`
	RestaurantId string `json:"restaurantId"`
}
