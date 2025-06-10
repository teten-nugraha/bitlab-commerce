package client

import (
	"context"
	"time"

	"order-service/gen/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentClient struct {
	client  payment.PaymentServiceClient
	conn    *grpc.ClientConn
	timeout time.Duration
}

func NewPaymentClient(addr string, timeout time.Duration) (*PaymentClient, error) {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithBlock(),
		grpc.WithTimeout(timeout),
	)
	if err != nil {
		return nil, err
	}

	return &PaymentClient{
		client:  payment.NewPaymentServiceClient(conn),
		conn:    conn,
		timeout: timeout,
	}, nil
}

func (c *PaymentClient) Close() error {
	return c.conn.Close()
}

func (c *PaymentClient) CreatePayment(ctx context.Context, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.client.CreatePayment(ctx, req)
}
