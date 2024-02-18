package psql

import (
	"context"
	"database/sql"
	"gitlab.com/aitalina/nocoin/internal/domain"
	"gitlab.com/aitalina/nocoin/internal/dto"
)

type Tip struct {
	db *sql.DB
}

func NewTip(db *sql.DB) *Tip {
	return &Tip{db}
}

func (b *Tip) Create(ctx context.Context, tip domain.Tip) error {
	_, err := b.db.Exec("INSERT INTO tip (count, transaction, status, restaurant_id) values ($1, $2, $3, $4)", tip.Count, tip.Transaction, tip.Status, tip.RestaurantId)

	return err
}

func (b *Tip) FindAllByOwnerId(ctx context.Context, ownerId string) ([]dto.Tip, error) {

	var data []dto.Tip
	// Выполняем запрос к базе данных
	rows, err := b.db.Query("SELECT t.count FROM tip as t INNER JOIN restaurant as r ON t.restaurant_id=r.id WHERE r.owner_id=$1", ownerId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Обрабатываем результаты запроса
	for rows.Next() {
		var tip dto.Tip
		if err := rows.Scan(&tip.Count); err != nil {
			panic(err)
		}
		data = append(data, tip)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	return data, err
}
