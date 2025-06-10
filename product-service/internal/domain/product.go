package domain

import (
	"time"
)

type Product struct {
	ID          string    `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Price       float64   `json:"price" bson:"price"`
	Stock       int       `json:"stock" bson:"stock"`
	Category    string    `json:"category" bson:"category"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type ProductStock struct {
	ID    string `json:"id" bson:"_id"`
	Stock int    `json:"stock" bson:"stock"`
}

type ProductValidation struct {
	Valid            bool
	UnavailableItems []ProductStock
	Message          string
}
