package service

import (
	"context"
	"errors"
	"time"

	"product-service/internal/domain"
	"product-service/internal/repository"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidStock    = errors.New("invalid stock quantity")
)

type ProductService struct {
	repo    repository.ProductRepository
	timeout time.Duration
}

func NewProductService(repo repository.ProductRepository, timeout time.Duration) *ProductService {
	return &ProductService{
		repo:    repo,
		timeout: timeout,
	}
}

func (s *ProductService) ValidateProducts(ctx context.Context, items []domain.ProductStock) (domain.ProductValidation, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Validate input
	for _, item := range items {
		if item.Quantity <= 0 {
			return domain.ProductValidation{}, ErrInvalidStock
		}
	}

	return s.repo.CheckStocks(ctx, items)
}

func (s *ProductService) GetProducts(ctx context.Context, ids []string) ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.repo.FindMultipleByID(ctx, ids)
}

func (s *ProductService) UpdateProductStocks(ctx context.Context, items []domain.ProductStock) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// First validate stocks
	validation, err := s.ValidateProducts(ctx, items)
	if err != nil {
		return err
	}
	if !validation.Valid {
		return ErrInvalidStock
	}

	return s.repo.UpdateStocks(ctx, items)
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	return product, nil
}
