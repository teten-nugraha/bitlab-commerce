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
	"google.golang.org/grpc/reflection"

	"product-service/internal/config"
	"product-service/internal/handler"
	"product-service/internal/repository"
	"product-service/internal/service"
	"product-service/pkg/grpcutil"
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

	// Initialize Repository
	productRepo := repository.NewMongoProductRepository(mongoClient.Database(cfg.MongoDB), 5*time.Second)

	// Initialize Services
	productService := service.NewProductService(productRepo, 5*time.Second)

	// Initialize gRPC Server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcutil.LoggingInterceptor),
	)

	// Register Services
	productHandler := handler.NewProductGRPCHandler(productService)
	product.RegisterProductServiceServer(grpcServer, productHandler)
	reflection.Register(grpcServer)

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
