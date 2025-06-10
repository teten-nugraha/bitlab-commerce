package handler

import (
	"context"
	"log"

	"product-service/gen/product"
	"product-service/internal/domain"
	"product-service/internal/service"
)

type ProductGRPCHandler struct {
	product.UnimplementedProductServiceServer
	service *service.ProductService
}

func NewProductGRPCHandler(svc *service.ProductService) *ProductGRPCHandler {
	return &ProductGRPCHandler{
		service: svc,
	}
}

func (h *ProductGRPCHandler) ValidateProducts(ctx context.Context, req *product.ValidateProductsRequest) (*product.ValidateProductsResponse, error) {
	// Convert request to domain objects
	var items []domain.ProductStock
	for _, item := range req.Items {
		items = append(items, domain.ProductStock{
			ID:    item.ProductId,
			Stock: int(item.Quantity),
		})
	}

	// Call service
	validation, err := h.service.ValidateProducts(ctx, items)
	if err != nil {
		log.Printf("ValidateProducts failed: %v", err)
		return nil, err
	}

	// Convert response
	resp := &product.ValidateProductsResponse{
		Valid:   validation.Valid,
		Message: validation.Message,
	}

	for _, item := range validation.UnavailableItems {
		resp.UnavailableItems = append(resp.UnavailableItems, &product.ProductItem{
			ProductId: item.ID,
			Quantity:  int32(item.Stock),
		})
	}

	return resp, nil
}

func (h *ProductGRPCHandler) GetProductDetails(ctx context.Context, req *product.GetProductDetailsRequest) (*product.GetProductDetailsResponse, error) {
	// Call service
	products, err := h.service.GetProducts(ctx, req.ProductIds)
	if err != nil {
		log.Printf("GetProductDetails failed: %v", err)
		return nil, err
	}

	// Convert response
	resp := &product.GetProductDetailsResponse{}
	for _, p := range products {
		resp.Products = append(resp.Products, &product.ProductDetail{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       int32(p.Stock),
		})
	}

	return resp, nil
}
