package handler

import (
	"context"
	"log"

	"order-service/gen/order"
	"order-service/internal/domain"
	"order-service/internal/service"
)

type OrderGRPCHandler struct {
	order.UnimplementedOrderServiceServer
	service *service.OrderService
}

func NewOrderGRPCHandler(svc *service.OrderService) *OrderGRPCHandler {
	return &OrderGRPCHandler{
		service: svc,
	}
}

func (h *OrderGRPCHandler) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.OrderResponse, error) {
	// Convert request to domain objects
	var items []domain.OrderItem
	for _, item := range req.Items {
		items = append(items, domain.OrderItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
			Price:     item.Price,
		})
	}

	// Call service
	order, err := h.service.CreateOrder(ctx, req.UserId, items)
	if err != nil {
		log.Printf("CreateOrder failed: %v", err)
		return nil, err
	}

	// Convert response
	return &order.OrderResponse{
		OrderId: order.ID,
		Status:  string(order.Status),
		Total:   order.Total,
	}, nil
}

func (h *OrderGRPCHandler) ProcessPayment(ctx context.Context, req *order.PaymentRequest) (*order.PaymentResponse, error) {
	// Call service
	order, err := h.service.ProcessPayment(ctx, req.OrderId, req.PaymentMethod)
	if err != nil {
		log.Printf("ProcessPayment failed: %v", err)
		return nil, err
	}

	// Convert response
	return &order.PaymentResponse{
		PaymentId:  order.PaymentID,
		Status:     string(order.Status),
		PaymentUrl: order.PaymentURL,
	}, nil
}
