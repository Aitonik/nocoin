package psql

import (
	"context"
	"database/sql"
	"gitlab.com/aitalina/nocoin/internal/domain"
)

type Restaurant struct {
	db *sql.DB
}

func NewRestaurant(db *sql.DB) *Restaurant {
	return &Restaurant{db}
}

func (b *Restaurant) Create(ctx context.Context, restaurant domain.Restaurant) error {
	_, err := b.db.Exec("INSERT INTO restaurant (name, owner_id) values ($1, $2)", restaurant.Name, restaurant.OwnerId)

	return err
}

func (b *Restaurant) GetByOwnerId(ctx context.Context, ownerId string) (domain.Restaurant, error) {
	var restaurant domain.Restaurant
	err := b.db.QueryRow("SELECT id, name FROM restaurant WHERE owner_id=$1", ownerId).Scan(&restaurant.ID, &restaurant.Name)
	if err == sql.ErrNoRows {
		return restaurant, domain.ErrRestaurantNotFound
	}

	return restaurant, err
}
