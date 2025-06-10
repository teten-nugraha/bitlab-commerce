package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort           string
	MongoURI           string
	MongoDB            string
	KafkaBrokers       []string
	ProductServiceAddr string
	PaymentServiceAddr string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	return &Config{
		GRPCPort:           getEnv("GRPC_PORT", "50052"),
		MongoURI:           getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:            getEnv("MONGO_DB", "order_service"),
		KafkaBrokers:       getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}, ","),
		ProductServiceAddr: getEnv("PRODUCT_SERVICE_ADDR", "product-service:50051"),
		PaymentServiceAddr: getEnv("PAYMENT_SERVICE_ADDR", "payment-service:50053"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValues []string, sep string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, sep)
	}
	return defaultValues
}
