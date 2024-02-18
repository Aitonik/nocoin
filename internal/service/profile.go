package service

import (
	"context"
	"gitlab.com/aitalina/nocoin/internal/domain"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile domain.Profile) error
	FindProfileByEmail(ctx context.Context, id string) (domain.Profile, error)
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

func (b *Profile) FindProfileByEmail(ctx context.Context, email string) (domain.Profile, error) {
	return b.repo.FindProfileByEmail(ctx, email)
}
