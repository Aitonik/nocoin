package psql

import (
	"context"
	"database/sql"
	"gitlab.com/aitalina/nocoin/internal/domain"
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
