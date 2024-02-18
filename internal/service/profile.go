package service

import (
	"context"
	"gitlab.com/aitalina/nocoin/internal/domain"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile domain.Profile) error
	GetByID(ctx context.Context, id string) (domain.Profile, error)
	GetPasswordByEmail(ctx context.Context, email string) (domain.Profile, error)
}

type Profile struct {
	repo ProfileRepository
}

func NewProfiles(repo ProfileRepository) *Profile {
	return &Profile{
		repo: repo,
	}
}

func (b *Profile) Create(ctx context.Context, profile domain.Profile) error {
	return b.repo.Create(ctx, profile)
}

func (b *Profile) GetByID(ctx context.Context, id string) (domain.Profile, error) {
	return b.repo.GetByID(ctx, id)
}

func (b *Profile) GetPasswordByEmail(ctx context.Context, email string) (domain.Profile, error) {
	return b.repo.GetPasswordByEmail(ctx, email)
}
