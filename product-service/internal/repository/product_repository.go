package repository

import (
	"context"
	"product-service/internal/domain"
)

type ProductRepository interface {
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	FindMultipleByID(ctx context.Context, ids []string) ([]domain.Product, error)
	CheckStocks(ctx context.Context, items []domain.ProductStock) (domain.ProductValidation, error)
	UpdateStocks(ctx context.Context, items []domain.ProductStock) error
}
