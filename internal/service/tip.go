package service

import (
	"context"
	"gitlab.com/aitalina/nocoin/internal/domain"
	"gitlab.com/aitalina/nocoin/internal/dto"
)

type TipRepository interface {
	Create(ctx context.Context, tip domain.Tip) error
	FindAllByOwnerId(ctx context.Context, ownerId string) ([]dto.Tip, error)
}

type Tip struct {
	repo TipRepository
}

func NewTips(repo TipRepository) *Tip {
	return &Tip{
		repo: repo,
	}
}

func (b *Tip) Create(ctx context.Context, tip domain.Tip) error {
	return b.repo.Create(ctx, tip)
}

func (b *Tip) FindAllByOwnerId(ctx context.Context, ownerId string) ([]dto.Tip, error) {
	return b.repo.FindAllByOwnerId(ctx, ownerId)
}
