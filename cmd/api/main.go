package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"Gommunity/docs"
	"Gommunity/internal/community/users/application/commandservices"
	"Gommunity/internal/community/users/application/eventhandlers"
	"Gommunity/internal/community/users/application/queryservices"
	"Gommunity/internal/community/users/domain/model/entities"
	"Gommunity/internal/community/users/domain/model/valueobjects"
	domain_repos "Gommunity/internal/community/users/domain/repositories"
	"Gommunity/internal/community/users/infrastructure/messaging"
	"Gommunity/internal/community/users/infrastructure/persistence/repositories"
	"Gommunity/internal/community/users/interfaces/rest/controllers"
	"Gommunity/shared/infrastructure/discovery"
	"Gommunity/shared/infrastructure/messaging/kafka"
	"Gommunity/shared/infrastructure/middleware"
	"Gommunity/shared/infrastructure/persistence/mongodb"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Gommunity API
// @version 1.0
// @description Community management API with Kafka event processing
// @host localhost
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default values")
	}

	// Get configuration from environment
	port := getEnv("PORT", "8080")
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDatabase := getEnv("MONGO_DATABASE", "gommunity")
	mongoTimeout := getEnvDuration("MONGO_TIMEOUT", 10*time.Second)
	kafkaBootstrapServers := getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
	jwtSecret := getEnv("JWT_SECRET", "")
	serviceDiscoveryURL := getEnv("SERVICE_DISCOVERY_URL", "http://127.0.0.1:8761/eureka/")
	serverIP := getEnv("SERVER_IP", "127.0.0.1")
	serviceName := getEnv("SERVICE_NAME", "gommunity-service")

	// Set Swagger host dynamically
	docs.SwaggerInfo.Host = "localhost:" + port

	// Initialize MongoDB connection
	mongoConn, err := mongodb.NewMongoConnection(mongodb.MongoConfig{
		URI:      mongoURI,
		Database: mongoDatabase,
		Timeout:  mongoTimeout,
	})
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mongoConn.Close(ctx); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}()

	// Initialize repositories
	userCollection := mongoConn.GetCollection("users")
	roleCollection := mongoConn.GetCollection("roles")

	// Create indexes
	indexCtx, indexCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer indexCancel()
	if err := mongodb.CreateUserIndexes(indexCtx, userCollection); err != nil {
		log.Printf("Warning: Failed to create indexes: %v", err)
	}

	userRepository := repositories.NewUserRepository(userCollection)
	roleRepository := repositories.NewRoleRepository(roleCollection)

	// Seed roles
	if err := seedRoles(context.Background(), roleRepository); err != nil {
		log.Printf("Warning: Failed to seed roles: %v", err)
	}

	// Initialize Eureka client
	var eurekaClient *discovery.EurekaClient
	eurekaClient, err = discovery.NewEurekaClient(discovery.EurekaConfig{
		ServiceName:     serviceName,
		ServerIP:        serverIP,
		Port:            port,
		DiscoveryURL:    serviceDiscoveryURL,
		HealthCheckURL:  fmt.Sprintf("http://%s:%s/health", serverIP, port),
		StatusPageURL:   fmt.Sprintf("http://%s:%s/swagger/index.html", serverIP, port),
		HomePageURL:     fmt.Sprintf("http://%s:%s/", serverIP, port),
		RenewalInterval: 30 * time.Second,
		DurationInSecs:  90,
	})
	if err != nil {
		log.Printf("Warning: Failed to create Eureka client: %v", err)
		eurekaClient = nil
	} else {
		// Register with Eureka
		if err := eurekaClient.Register(); err != nil {
			log.Printf("Warning: Failed to register with Eureka: %v", err)
			eurekaClient = nil
		} else {
			// Start heartbeat
			eurekaClient.StartHeartbeat()
			log.Println("Successfully registered with Eureka and started heartbeat")
		}
	}

	// Initialize services
	userQueryService := queryservices.NewUserQueryService(userRepository)
	userCommandService := commandservices.NewUserCommandService(userRepository)

	// Initialize event handlers
	registrationHandler := eventhandlers.NewUserRegistrationHandler(userRepository)
	profileUpdateHandler := eventhandlers.NewProfileUpdatedHandler(userRepository)

	// Initialize controllers
	userController := controllers.NewUserController(userCommandService, userQueryService)

	// Initialize JWT middleware
	jwtMiddleware := middleware.NewJWTMiddleware(jwtSecret)

	// Initialize Kafka event consumer
	kafkaEventConsumer := messaging.NewKafkaEventConsumer(registrationHandler, profileUpdateHandler)

	// Initialize Kafka consumer
	kafkaConsumer := kafka.NewKafkaConsumer(kafka.KafkaConfig{
		BootstrapServers: kafkaBootstrapServers,
		GroupID:          "gommunity-consumer-group",
		Topics: []string{
			messaging.TopicCommunityRegistration,
			messaging.TopicProfileUpdated,
		},
	})

	// Start Kafka consumer in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := kafkaConsumer.ConsumeMessages(ctx, kafkaEventConsumer.HandleMessage); err != nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	corsConfig := cors.Config{
		AllowOrigins:     getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		AllowMethods:     getEnvSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}),
		AllowHeaders:     getEnvSlice("CORS_ALLOWED_HEADERS", []string{"*"}),
		AllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", true),
		MaxAge:           getEnvDuration("CORS_MAX_AGE", 12*time.Hour),
	}
	r.Use(cors.New(corsConfig))

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	r.GET("/health", healthHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// User routes (protected with JWT)
	userRoutes := r.Group("/users")
	userRoutes.Use(jwtMiddleware.AuthMiddleware())
	{
		userRoutes.GET("/:id", userController.GetUserByID)
		userRoutes.GET("/username/:username", userController.GetUserByUsername)
		userRoutes.PUT("/:id/banner", userController.UpdateBannerURL)
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		log.Printf("Swagger UI available at: http://localhost:%s/swagger/index.html", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Deregister from Eureka
	if eurekaClient != nil {
		if err := eurekaClient.Deregister(); err != nil {
			log.Printf("Error deregistering from Eureka: %v", err)
		}
	}

	// Cancel Kafka consumer context
	cancel()

	log.Println("Server exited")
}

// healthHandler godoc
// @Summary Health check
// @Description Get health status of the service
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "gommunity",
	})
}

// seedRoles seeds the default roles if they don't exist
func seedRoles(ctx context.Context, roleRepo domain_repos.RoleRepository) error {
	roles := []struct {
		id   string
		name string
	}{
		{valueobjects.UserRoleIDStr, "user"},
		{valueobjects.MemberRoleIDStr, "member"},
		{valueobjects.AdminRoleIDStr, "admin"},
		{valueobjects.OwnerRoleIDStr, "owner"},
	}

	for _, r := range roles {
		roleID, err := valueobjects.NewRoleID(r.id)
		if err != nil {
			return err
		}

		// Check if role exists
		existing, err := roleRepo.FindByID(ctx, roleID)
		if err != nil {
			return err
		}

		if existing == nil {
			// Create role
			role, err := entities.NewRole(roleID, r.name)
			if err != nil {
				return err
			}

			if err := roleRepo.Save(ctx, role); err != nil {
				return err
			}

			log.Printf("Seeded role: %s", r.name)
		}
	}

	return nil
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
