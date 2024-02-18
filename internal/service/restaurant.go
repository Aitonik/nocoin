package service

import (
	"context"
	"gitlab.com/aitalina/nocoin/internal/domain"
)

type RestaurantRepository interface {
	Create(ctx context.Context, restaurant domain.Restaurant) error
	GetByOwnerId(ctx context.Context, id string) (domain.Restaurant, error)
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

func (b *Restaurant) GetByOwnerId(ctx context.Context, ownerId string) (domain.Restaurant, error) {
	return b.repo.GetByOwnerId(ctx, ownerId)
}
