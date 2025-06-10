package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"user-service/internal/config"
	"user-service/internal/controllers"
	"user-service/internal/infrastructure/repositories"
	"user-service/internal/middleware"
	"user-service/internal/services"
	"user-service/pkg/eventbus"
	"user-service/pkg/jwt"
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

	// Initialize JWT Manager
	jwtManager := jwt.NewManager(cfg.JWTSecret, 24*time.Hour)

	// Initialize Repository
	userRepo := repositories.NewMongoUserRepository(mongoClient.Database(cfg.MongoDB), 5*time.Second)

	// Initialize Services
	userService := services.NewUserService(userRepo, jwtManager, eventBus, 5*time.Second)

	// Initialize Controllers
	userController := controllers.NewUserController(userService, 5*time.Second)

	// Initialize Gin Router
	router := gin.Default()

	// Routes
	api := router.Group("/api/v1")
	{
		api.POST("/register", userController.Register)
		api.POST("/login", userController.Login)

		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware(jwtManager))
		{
			auth.GET("/profile", userController.GetProfile)
			auth.PUT("/profile", userController.UpdateProfile)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}

	log.Println("server exiting")
}
