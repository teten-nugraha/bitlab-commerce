package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	MongoURI     string
	MongoDB      string
	JWTSecret    string
	KafkaBrokers []string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		MongoURI:     getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:      getEnv("MONGO_DB", "user_service"),
		JWTSecret:    getEnv("JWT_SECRET", "secret"),
		KafkaBrokers: getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}, ","),
	}

	return cfg, nil
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
