package grpcutil

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"product-service/gen/product"
)

type ProductClient struct {
	conn    *grpc.ClientConn
	client  product.ProductServiceClient
	timeout time.Duration
}

func NewProductClient(addr string, timeout time.Duration) (*ProductClient, error) {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(timeout),
	)
	if err != nil {
		return nil, err
	}

	return &ProductClient{
		conn:    conn,
		client:  product.NewProductServiceClient(conn),
		timeout: timeout,
	}, nil
}

func (c *ProductClient) Close() error {
	return c.conn.Close()
}

func (c *ProductClient) ValidateProducts(ctx context.Context, items []*product.ProductItem) (*product.ValidateProductsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.client.ValidateProducts(ctx, &product.ValidateProductsRequest{
		Items: items,
	})
}

func (c *ProductClient) GetProductDetails(ctx context.Context, productIDs []string) (*product.GetProductDetailsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.client.GetProductDetails(ctx, &product.GetProductDetailsRequest{
		ProductIds: productIDs,
	})
}
