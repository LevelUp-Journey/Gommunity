package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// KafkaConfig holds Kafka-related configuration
type KafkaConfig struct {
	BootstrapServers string
	GroupID          string
	Topics           []string
}

type Config struct {
	Port                 string
	MongoURI             string
	MongoDatabase        string
	MongoTimeout         time.Duration
	Kafka                KafkaConfig
	JWTSecret            string
	ServiceDiscoveryURL  string
	ServerIP             string
	ServiceName          string
	CORSAllowedOrigins   []string
	CORSAllowedMethods   []string
	CORSAllowedHeaders   []string
	CORSAllowCredentials bool
	CORSMaxAge           time.Duration
}

func Load() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default values")
	}

	config := &Config{
		Port:          getEnv("PORT", "8080"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "gommunity"),
		MongoTimeout:  getEnvDuration("MONGO_TIMEOUT", 10*time.Second),
		Kafka: KafkaConfig{
			BootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
			GroupID:          getEnv("KAFKA_GROUP_ID", "gommunity-group"),
			Topics:           getEnvSlice("KAFKA_TOPICS", []string{}),
		},
		JWTSecret:            getEnv("JWT_SECRET", ""),
		ServiceDiscoveryURL:  strings.TrimSuffix(getEnv("SERVICE_DISCOVERY_URL", "http://127.0.0.1:8761/eureka"), "/"),
		ServerIP:             getEnv("SERVER_IP", "127.0.0.1"),
		ServiceName:          getEnv("SERVICE_NAME", "gommunity-service"),
		CORSAllowedOrigins:   getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		CORSAllowedMethods:   getEnvSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}),
		CORSAllowedHeaders:   getEnvSlice("CORS_ALLOWED_HEADERS", []string{"*"}),
		CORSAllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", true),
		CORSMaxAge:           getEnvDuration("CORS_MAX_AGE", 12*time.Hour),
	}

	return config, nil
}

// Helper functions for environment variables

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.Split(value, ",")
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true"
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Invalid duration for %s: %v, using default", key, err)
		return defaultValue
	}
	return duration
}
