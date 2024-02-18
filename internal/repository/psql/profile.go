package psql

import (
	"context"
	"database/sql"
	"gitlab.com/aitalina/nocoin/internal/domain"
)

type Profile struct {
	db *sql.DB
}

func NewProfile(db *sql.DB) *Profile {
	return &Profile{db}
}

func (b *Profile) Create(ctx context.Context, profile domain.Profile) error {
	_, err := b.db.Exec("INSERT INTO profile (name, email, role, password) values ($1, $2, $3, $4)", profile.Name, profile.Email, profile.Role, profile.Password)

	return err
}

func (b *Profile) FindProfileByEmail(ctx context.Context, email string) (domain.Profile, error) {
	var profile domain.Profile
	err := b.db.QueryRow("SELECT id, name, email, password, role FROM profile WHERE email=$1", email).Scan(&profile.ID, &profile.Name, &profile.Email, &profile.Password, &profile.Role)
	if err == sql.ErrNoRows {
		return profile, domain.ErrProfileNotFound
	}

	return profile, err
}
