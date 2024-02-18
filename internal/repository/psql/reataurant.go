package psql

import (
	"context"
	"database/sql"
	"fmt"
	"gitlab.com/aitalina/nocoin/internal/domain"
	"strings"
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

func (b *Restaurant) GetByID(ctx context.Context, id string) (domain.Restaurant, error) {
	var book domain.Restaurant
	err := b.db.QueryRow("SELECT id, name FROM restaurant WHERE id=$1", id).Scan(&book.ID, &book.Name)
	if err == sql.ErrNoRows {
		return book, domain.ErrBookNotFound
	}

	return book, err
}

func (b *Restaurant) GetAll(ctx context.Context) ([]domain.Restaurant, error) {
	rows, err := b.db.Query("SELECT id, name FROM restaurant")
	if err != nil {
		return nil, err
	}

	books := make([]domain.Restaurant, 0)
	for rows.Next() {
		var book domain.Restaurant
		if err := rows.Scan(&book.ID, &book.Name); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, rows.Err()
}

func (b *Restaurant) Update(ctx context.Context, id string, inp domain.UpdateRestaurantInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if inp.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *inp.Name)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE restaurant SET %s WHERE id=$%d", setQuery, argId)
	args = append(args, id)

	_, err := b.db.Exec(query, args...)
	return err
}
