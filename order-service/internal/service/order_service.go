package service

import (
	"context"
	"errors"
	"time"

	"order-service/gen/payment"
	"order-service/gen/product"
	"order-service/internal/domain"
	"order-service/internal/repository"
	"order-service/pkg/eventbus"

	"github.com/google/uuid"
)

var (
	ErrInvalidOrder      = errors.New("invalid order")
	ErrProductValidation = errors.New("product validation failed")
	ErrPaymentProcessing = errors.New("payment processing failed")
	ErrOrderNotFound     = errors.New("order not found")
)

type OrderService struct {
	orderRepo  repository.OrderRepository
	productCli *client.ProductClient
	paymentCli *client.PaymentClient
	eventBus   eventbus.EventBus
	timeout    time.Duration
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	productCli *client.ProductClient,
	paymentCli *client.PaymentClient,
	eventBus eventbus.EventBus,
	timeout time.Duration,
) *OrderService {
	return &OrderService{
		orderRepo:  orderRepo,
		productCli: productCli,
		paymentCli: paymentCli,
		eventBus:   eventBus,
		timeout:    timeout,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, items []domain.OrderItem) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Validate input
	if userID == "" || len(items) == 0 {
		return nil, ErrInvalidOrder
	}

	// Validate product stock
	if err := s.validateProducts(ctx, items); err != nil {
		return nil, err
	}

	// Calculate total
	total := calculateTotal(items)

	// Create order
	order := &domain.Order{
		ID:        generateID(),
		UserID:    userID,
		Items:     items,
		Total:     total,
		Status:    domain.OrderStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Publish OrderCreated event
	go s.eventBus.Publish("order.created", domain.OrderCreatedEvent{
		OrderID:   order.ID,
		UserID:    order.UserID,
		Items:     order.Items,
		Total:     order.Total,
		CreatedAt: order.CreatedAt,
	})

	return order, nil
}

func (s *OrderService) ProcessPayment(ctx context.Context, orderID, paymentMethod string) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Get order
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// Process payment via Payment Service
	paymentResp, err := s.paymentCli.CreatePayment(ctx, &payment.PaymentRequest{
		OrderId:       order.ID,
		UserId:        order.UserID,
		Amount:        order.Total,
		Currency:      "USD",
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		// Update order status to failed
		order.Status = domain.OrderStatusFailed
		s.orderRepo.Update(ctx, order)
		return nil, ErrPaymentProcessing
	}

	// Update order with payment details
	order.PaymentID = paymentResp.PaymentId
	order.PaymentURL = paymentResp.PaymentUrl

	if paymentResp.Status == "success" {
		order.Status = domain.OrderStatusPaid
	} else {
		order.Status = domain.OrderStatusFailed
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	// Publish PaymentProcessed event
	go s.eventBus.Publish("payment.processed", domain.PaymentProcessedEvent{
		OrderID:    order.ID,
		PaymentID:  paymentResp.PaymentId,
		Status:     paymentResp.Status,
		Amount:     order.Total,
		OccurredAt: time.Now(),
	})

	return order, nil
}

// Helper functions
func (s *OrderService) validateProducts(ctx context.Context, items []domain.OrderItem) error {
	var productItems []*product.ProductItem
	for _, item := range items {
		productItems = append(productItems, &product.ProductItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
	}

	resp, err := s.productCli.ValidateProducts(ctx, productItems)
	if err != nil {
		return ErrProductValidation
	}
	if !resp.Valid {
		return ErrProductValidation
	}
	return nil
}

func calculateTotal(items []domain.OrderItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

func generateID() string {
	return uuid.New().String()
}
