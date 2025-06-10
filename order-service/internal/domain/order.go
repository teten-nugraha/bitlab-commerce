package domain

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusFailed    OrderStatus = "failed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID         string      `json:"id" bson:"_id"`
	UserID     string      `json:"user_id" bson:"user_id"`
	Items      []OrderItem `json:"items" bson:"items"`
	Total      float64     `json:"total" bson:"total"`
	Status     OrderStatus `json:"status" bson:"status"`
	PaymentID  string      `json:"payment_id,omitempty" bson:"payment_id,omitempty"`
	PaymentURL string      `json:"payment_url,omitempty" bson:"payment_url,omitempty"`
	CreatedAt  time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at" bson:"updated_at"`
}

type OrderItem struct {
	ProductID string  `json:"product_id" bson:"product_id"`
	Quantity  int     `json:"quantity" bson:"quantity"`
	Price     float64 `json:"price" bson:"price"`
}

// Events
type OrderCreatedEvent struct {
	OrderID   string      `json:"order_id"`
	UserID    string      `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	CreatedAt time.Time   `json:"created_at"`
}

type PaymentProcessedEvent struct {
	OrderID    string    `json:"order_id"`
	PaymentID  string    `json:"payment_id"`
	Status     string    `json:"status"`
	Amount     float64   `json:"amount"`
	OccurredAt time.Time `json:"occurred_at"`
}
