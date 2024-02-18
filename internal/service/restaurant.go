package service

import (
	"context"
	"gitlab.com/aitalina/nocoin/internal/domain"
)

type RestaurantRepository interface {
	Create(ctx context.Context, restaurant domain.Restaurant) error
	GetByID(ctx context.Context, id string) (domain.Restaurant, error)
	GetAll(ctx context.Context) ([]domain.Restaurant, error)
	Update(ctx context.Context, id string, inp domain.UpdateRestaurantInput) error
}

type Restaurant struct {
	repo RestaurantRepository
}

func NewRestaurants(repo RestaurantRepository) *Restaurant {
	return &Restaurant{
		repo: repo,
	}
}

func (b *Restaurant) Create(ctx context.Context, restaurant domain.Restaurant) error {
	return b.repo.Create(ctx, restaurant)
}

func (b *Restaurant) GetByID(ctx context.Context, id string) (domain.Restaurant, error) {
	return b.repo.GetByID(ctx, id)
}

func (b *Restaurant) GetAll(ctx context.Context) ([]domain.Restaurant, error) {
	return b.repo.GetAll(ctx)
}

func (b *Restaurant) Update(ctx context.Context, id string, inp domain.UpdateRestaurantInput) error {
	return b.repo.Update(ctx, id, inp)
}
