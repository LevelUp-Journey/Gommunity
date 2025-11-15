package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Gommunity/docs"
	community_commandservices "Gommunity/internal/community/communities/application/commandservices"
	community_queryservices "Gommunity/internal/community/communities/application/queryservices"
	community_repositories "Gommunity/internal/community/communities/infrastructure/persistence/repositories"
	community_controllers "Gommunity/internal/community/communities/interfaces/rest/controllers"
	"Gommunity/internal/community/users/application/commandservices"
	"Gommunity/internal/community/users/application/eventhandlers"
	"Gommunity/internal/community/users/application/queryservices"
	user_services "Gommunity/internal/community/users/application/services"
	"Gommunity/internal/community/users/infrastructure/messaging"
	"Gommunity/internal/community/users/infrastructure/persistence/repositories"
	"Gommunity/internal/community/users/interfaces/rest/controllers"
	"Gommunity/shared/config"
	"Gommunity/shared/infrastructure/discovery"
	"Gommunity/shared/infrastructure/messaging/kafka"
	"Gommunity/shared/infrastructure/middleware"
	"Gommunity/shared/infrastructure/persistence/mongodb"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
// @description Enter your JWT token directly (Bearer prefix is optional).
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Swagger host dynamically
	docs.SwaggerInfo.Host = "localhost:" + cfg.Port

	// Initialize MongoDB connection
	mongoConn, err := mongodb.NewMongoConnection(mongodb.MongoConfig{
		URI:      cfg.MongoURI,
		Database: cfg.MongoDatabase,
		Timeout:  cfg.MongoTimeout,
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
	communityCollection := mongoConn.GetCollection("communities")

	// Create indexes
	indexCtx, indexCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer indexCancel()
	if err := mongodb.CreateUserIndexes(indexCtx, userCollection); err != nil {
		log.Printf("Warning: Failed to create indexes: %v", err)
	}

	userRepository := repositories.NewUserRepository(userCollection)
	roleRepository := repositories.NewRoleRepository(roleCollection)
	communityRepository := community_repositories.NewCommunityRepository(communityCollection)

	// Seed roles
	if err := eventhandlers.SeedRoles(context.Background(), roleRepository); err != nil {
		log.Printf("Warning: Failed to seed roles: %v", err)
	}

	// Initialize Eureka client
	var eurekaClient *discovery.EurekaClient
	eurekaClient, err = discovery.NewEurekaClient(discovery.EurekaConfig{
		ServiceName:     cfg.ServiceName,
		ServerIP:        cfg.ServerIP,
		Port:            cfg.Port,
		DiscoveryURL:    cfg.ServiceDiscoveryURL,
		HealthCheckURL:  fmt.Sprintf("http://%s:%s/", cfg.ServerIP, cfg.Port),
		StatusPageURL:   fmt.Sprintf("http://%s:%s/swagger/index.html", cfg.ServerIP, cfg.Port),
		HomePageURL:     fmt.Sprintf("http://%s:%s/", cfg.ServerIP, cfg.Port),
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
	communityQueryService := community_queryservices.NewCommunityQueryService(communityRepository)
	communityCommandService := community_commandservices.NewCommunityCommandService(communityRepository)

	// Initialize event handlers
	registrationHandler := eventhandlers.NewUserRegistrationHandler(userRepository)
	profileUpdateHandler := eventhandlers.NewProfileUpdatedHandler(userRepository)

	// Initialize controllers
	userController := controllers.NewUserController(userCommandService, userQueryService, roleRepository)
	communityController := community_controllers.NewCommunityController(communityCommandService, communityQueryService)

	// Initialize user role provider
	userRoleProvider := user_services.NewUserRoleProviderService(userRepository, roleRepository)

	// Initialize JWT middleware with role provider
	jwtMiddleware := middleware.NewJWTMiddlewareWithRoleProvider(cfg.JWTSecret, userRoleProvider)

	// Initialize Kafka event consumer
	kafkaEventConsumer := messaging.NewKafkaEventConsumer(registrationHandler, profileUpdateHandler)

	// Initialize Kafka consumer
	kafkaConsumer := kafka.NewKafkaConsumer(kafka.KafkaConfig{
		BootstrapServers: cfg.KafkaBootstrapServers,
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
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     cfg.CORSAllowedMethods,
		AllowHeaders:     cfg.CORSAllowedHeaders,
		AllowCredentials: cfg.CORSAllowCredentials,
		MaxAge:           cfg.CORSMaxAge,
	}
	r.Use(cors.New(corsConfig))

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// User routes (protected with JWT)
	userRoutes := r.Group("/users")
	userRoutes.Use(jwtMiddleware.AuthMiddleware())
	{
		userRoutes.GET("/:id", userController.GetUserByID)
		userRoutes.GET("/username/:username", userController.GetUserByUsername)
		userRoutes.PUT("/:id/banner", userController.UpdateBannerURL)
	}

	// Community routes (protected with JWT)
	communityRoutes := r.Group("/communities")
	communityRoutes.Use(jwtMiddleware.AuthMiddleware())
	{
		communityRoutes.POST("", communityController.CreateCommunity)
		communityRoutes.GET("", communityController.GetAllCommunities)
		communityRoutes.GET("/my-communities", communityController.GetMyCommunitiesAsOwner)
		communityRoutes.GET("/:id", communityController.GetCommunityByID)
		communityRoutes.PUT("/:id", communityController.UpdateCommunityInfo)
		communityRoutes.DELETE("/:id", communityController.DeleteCommunity)
		communityRoutes.PATCH("/:id/privacy", communityController.UpdateCommunityPrivacy)
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		log.Printf("Swagger UI available at: http://localhost:%s/swagger/index.html", cfg.Port)
		if err := r.Run(":" + cfg.Port); err != nil {
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
