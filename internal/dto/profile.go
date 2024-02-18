package dto

type Profile struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	RestaurantName string `json:"restaurantName"`
}
