package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	"order-service/internal/client"
	"order-service/internal/config"
	"order-service/internal/handler"
	"order-service/internal/repository"
	"order-service/internal/service"
	"order-service/pkg/eventbus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize MongoDB
	mongoCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("failed to disconnect mongodb: %v", err)
		}
	}()

	// Initialize Kafka Event Bus
	eventBus := eventbus.NewKafkaEventBus(cfg.KafkaBrokers)
	defer eventBus.Close()

	// Initialize Clients
	productCli, err := client.NewProductClient(cfg.ProductServiceAddr, 5*time.Second)
	if err != nil {
		log.Fatalf("failed to create product client: %v", err)
	}
	defer productCli.Close()

	paymentCli, err := client.NewPaymentClient(cfg.PaymentServiceAddr, 5*time.Second)
	if err != nil {
		log.Fatalf("failed to create payment client: %v", err)
	}
	defer paymentCli.Close()

	// Initialize Repository
	orderRepo := repository.NewMongoOrderRepository(mongoClient.Database(cfg.MongoDB), 5*time.Second)

	// Initialize Services
	orderService := service.NewOrderService(orderRepo, productCli, paymentCli, eventBus, 5*time.Second)

	// Initialize gRPC Server
	grpcServer := grpc.NewServer()
	orderHandler := handler.NewOrderGRPCHandler(orderService)
	order.RegisterOrderServiceServer(grpcServer, orderHandler)

	// Start gRPC Server
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Printf("gRPC server listening on %s", listener.Addr())
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down gRPC server...")

	grpcServer.GracefulStop()
	log.Println("server exited")
}
