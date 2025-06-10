package repository

import (
	"context"
	"time"

	"product-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoProductRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

func NewMongoProductRepository(db *mongo.Database, timeout time.Duration) *MongoProductRepository {
	return &MongoProductRepository{
		collection: db.Collection("products"),
		timeout:    timeout,
	}
}

func (r *MongoProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var product domain.Product
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *MongoProductRepository) FindMultipleByID(ctx context.Context, ids []string) ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *MongoProductRepository) CheckStocks(ctx context.Context, items []domain.ProductStock) (domain.ProductValidation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// Get current stocks for all products
	productIDs := make([]string, len(items))
	for i, item := range items {
		productIDs[i] = item.ID
	}

	products, err := r.FindMultipleByID(ctx, productIDs)
	if err != nil {
		return domain.ProductValidation{}, err
	}

	// Validate stocks
	validation := domain.ProductValidation{Valid: true}
	productStockMap := make(map[string]int)
	for _, p := range products {
		productStockMap[p.ID] = p.Stock
	}

	for _, item := range items {
		if stock, exists := productStockMap[item.ID]; !exists || stock < item.Quantity {
			validation.Valid = false
			validation.UnavailableItems = append(validation.UnavailableItems, domain.ProductStock{
				ID:    item.ID,
				Stock: stock,
			})
		}
	}

	if !validation.Valid {
		validation.Message = "some products are unavailable or out of stock"
	}

	return validation, nil
}

func (r *MongoProductRepository) UpdateStocks(ctx context.Context, items []domain.ProductStock) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// Use bulk write for better performance
	var operations []mongo.WriteModel
	for _, item := range items {
		update := bson.M{
			"$inc": bson.M{"stock": -item.Quantity},
			"$set": bson.M{"updated_at": time.Now()},
		}
		operation := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": item.ID}).
			SetUpdate(update)
		operations = append(operations, operation)
	}

	_, err := r.collection.BulkWrite(ctx, operations)
	return err
}
