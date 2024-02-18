package service

import (
	"context"
	"gitlab.com/aitalina/nocoin/internal/domain"
)

type TipRepository interface {
	Create(ctx context.Context, tip domain.Tip) error
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
