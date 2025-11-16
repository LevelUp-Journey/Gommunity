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
	community_acl "Gommunity/platform/community/application/outboundservices/acl"
	community_commandservices "Gommunity/platform/community/application/commandservices"
	community_queryservices "Gommunity/platform/community/application/queryservices"
	community_repositories "Gommunity/platform/community/infrastructure/persistence/repositories"
	community_controllers "Gommunity/platform/community/interfaces/rest/controllers"
	posts_acl_impl "Gommunity/platform/posts/application/acl"
	posts_commandservices "Gommunity/platform/posts/application/commandservices"
	posts_acl "Gommunity/platform/posts/application/outboundservices/acl"
	posts_queryservices "Gommunity/platform/posts/application/queryservices"
	posts_repositories "Gommunity/platform/posts/infrastructure/persistence/repositories"
	posts_controllers "Gommunity/platform/posts/interfaces/rest/controllers"
	reactions_commandservices "Gommunity/platform/reactions/application/commandservices"
	reactions_acl "Gommunity/platform/reactions/application/outboundservices/acl"
	reactions_queryservices "Gommunity/platform/reactions/application/queryservices"
	reactions_repositories "Gommunity/platform/reactions/infrastructure/persistence/repositories"
	reactions_controllers "Gommunity/platform/reactions/interfaces/rest/controllers"
	"Gommunity/platform/users/application/commandservices"
	"Gommunity/platform/users/application/eventhandlers"
	"Gommunity/platform/users/application/queryservices"
	"Gommunity/platform/users/infrastructure/messaging"
	"Gommunity/platform/users/infrastructure/persistence/repositories"
	"Gommunity/platform/users/interfaces/rest/controllers"
	"Gommunity/shared/config"
	"Gommunity/shared/infrastructure/discovery"
	"Gommunity/shared/infrastructure/messaging/kafka"
	"Gommunity/shared/infrastructure/middleware"
	"Gommunity/shared/infrastructure/persistence/mongodb"

	// Subscriptions BC imports
	communities_acl "Gommunity/platform/community/application/acl"
	subscriptions_acl "Gommunity/platform/subscriptions/application/acl"
	subscription_commandservices "Gommunity/platform/subscriptions/application/commandservices"
	subscriptions_outbound_acl "Gommunity/platform/subscriptions/application/outboundservices/acl"
	subscription_queryservices "Gommunity/platform/subscriptions/application/queryservices"
	subscription_repositories "Gommunity/platform/subscriptions/infrastructure/persistence/repositories"
	subscription_controllers "Gommunity/platform/subscriptions/interfaces/rest/controllers"
	users_acl "Gommunity/platform/users/application/acl"

	// Feed BC imports
	feed_acl "Gommunity/platform/feed/application/outboundservices/acl"
	feed_queryservices "Gommunity/platform/feed/application/queryservices"
	feed_controllers "Gommunity/platform/feed/interfaces/rest/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Gommunity API
// @version 1.0
// @description Community management API with Kafka event processing
// @host localhost
// @BasePath /api/v1
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
	communityCollection := mongoConn.GetCollection("communities")
	subscriptionCollection := mongoConn.GetCollection("subscriptions")
	postCollection := mongoConn.GetCollection("posts")
	reactionCollection := mongoConn.GetCollection("reactions")

	// Create indexes
	indexCtx, indexCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer indexCancel()
	if err := mongodb.CreateUserIndexes(indexCtx, userCollection); err != nil {
		log.Printf("Warning: Failed to create indexes: %v", err)
	}

	userRepository := repositories.NewUserRepository(userCollection)
	communityRepository := community_repositories.NewCommunityRepository(communityCollection)
	subscriptionRepository := subscription_repositories.NewSubscriptionRepository(subscriptionCollection)
	postRepository := posts_repositories.NewPostRepository(postCollection)
	reactionRepository := reactions_repositories.NewReactionRepository(reactionCollection)

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

	// Initialize ACL facades
	usersFacade := users_acl.NewUsersFacade(userRepository)
	communitiesFacade := communities_acl.NewCommunitiesFacade(communityRepository)
	subscriptionsFacade := subscriptions_acl.NewSubscriptionsFacade(subscriptionRepository)

	// Initialize services
	userQueryService := queryservices.NewUserQueryService(userRepository)
	userCommandService := commandservices.NewUserCommandService(userRepository)
	communityQueryService := community_queryservices.NewCommunityQueryService(communityRepository)

	// Initialize Subscriptions BC services
	externalUsersService := subscriptions_outbound_acl.NewExternalUsersService(usersFacade)
	externalCommunitiesService := subscriptions_outbound_acl.NewExternalCommunitiesService(communitiesFacade)
	subscriptionCommandService := subscription_commandservices.NewSubscriptionCommandService(
		subscriptionRepository,
		externalUsersService,
		externalCommunitiesService,
	)
	subscriptionQueryService := subscription_queryservices.NewSubscriptionQueryService(
		subscriptionRepository,
	)

	// Initialize Community BC ACL service for subscriptions
	communityExternalSubscriptionsService := community_acl.NewExternalSubscriptionsService(subscriptionCommandService)

	// Initialize Community BC command service with subscription dependency
	communityCommandService := community_commandservices.NewCommunityCommandService(
		communityRepository,
		communityExternalSubscriptionsService,
	)

	postExternalUsersService := posts_acl.NewExternalUsersService(usersFacade)
	postExternalCommunitiesService := posts_acl.NewExternalCommunitiesService(communitiesFacade)
	postExternalSubscriptionsService := posts_acl.NewExternalSubscriptionsService(subscriptionsFacade)
	postCommandService := posts_commandservices.NewPostCommandService(
		postRepository,
		postExternalUsersService,
		postExternalCommunitiesService,
		postExternalSubscriptionsService,
	)
	postQueryService := posts_queryservices.NewPostQueryService(postRepository)

	// Initialize Posts ACL facade
	postsFacade := posts_acl_impl.NewPostsFacade(postQueryService, postRepository)

	// Initialize Reactions BC services
	reactionsExternalPostsService := reactions_acl.NewExternalPostsService(postsFacade)
	reactionsExternalUsersService := reactions_acl.NewExternalUsersService(usersFacade)
	reactionCommandService := reactions_commandservices.NewReactionCommandService(
		reactionRepository,
		reactionsExternalPostsService,
		reactionsExternalUsersService,
	)
	reactionQueryService := reactions_queryservices.NewReactionQueryService(reactionRepository)

	// Initialize Feed BC services
	feedExternalSubscriptionsService := feed_acl.NewExternalSubscriptionsService(subscriptionsFacade)
	feedExternalPostsService := feed_acl.NewExternalPostsService(postsFacade)
	feedQueryService := feed_queryservices.NewFeedQueryService(
		feedExternalSubscriptionsService,
		feedExternalPostsService,
	)

	// Initialize event handlers
	registrationHandler := eventhandlers.NewUserRegistrationHandler(userRepository)
	profileUpdateHandler := eventhandlers.NewProfileUpdatedHandler(userRepository)

	// Initialize controllers
	userController := controllers.NewUserController(userCommandService, userQueryService)
	communityController := community_controllers.NewCommunityController(communityCommandService, communityQueryService)
	subscriptionController := subscription_controllers.NewSubscriptionController(
		subscriptionCommandService,
		subscriptionQueryService,
		externalUsersService,
	)
	postController := posts_controllers.NewPostController(postCommandService, postQueryService)
	reactionController := reactions_controllers.NewReactionController(reactionCommandService, reactionQueryService)
	feedController := feed_controllers.NewFeedController(feedQueryService)

	// Initialize JWT middleware
	// Note: Roles (STUDENT, TEACHER, ADMIN) come directly from IAM service via JWT
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWTSecret)

	// Initialize Kafka event consumer
	kafkaEventConsumer := messaging.NewKafkaEventConsumer(registrationHandler, profileUpdateHandler)

	// Initialize Kafka consumer
	kafkaConsumer := kafka.NewKafkaConsumer(kafka.KafkaConfig{
		BootstrapServers: cfg.Kafka.BootstrapServers,
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

	// API routes with prefix
	api := r.Group(cfg.APIPrefix)

	// User routes (protected with JWT)
	userRoutes := api.Group("/users")
	userRoutes.Use(jwtMiddleware.AuthMiddleware())
	{
		userRoutes.GET("/:id", userController.GetUserByID)
		userRoutes.GET("/username/:username", userController.GetUserByUsername)
		userRoutes.PUT("/:id/banner", userController.UpdateBannerURL)
	}

	// Community routes (protected with JWT)
	communityRoutes := api.Group("/communities")
	communityRoutes.Use(jwtMiddleware.AuthMiddleware())
	{
		communityRoutes.POST("", communityController.CreateCommunity)
		communityRoutes.GET("", communityController.GetAllCommunities)
		communityRoutes.GET("/my-communities", communityController.GetMyCommunitiesAsOwner)
		communityRoutes.GET("/:community_id/posts", postController.GetPostsByCommunity)
		communityRoutes.POST("/:community_id/posts", postController.CreatePost)
		communityRoutes.GET("/:community_id/posts/:post_id", postController.GetPostByID)
		communityRoutes.DELETE("/:community_id/posts/:post_id", postController.DeletePost)
		communityRoutes.GET("/:community_id", communityController.GetCommunityByID)
		communityRoutes.PUT("/:community_id", communityController.UpdateCommunityInfo)
		communityRoutes.DELETE("/:community_id", communityController.DeleteCommunity)
		communityRoutes.PATCH("/:community_id/privacy", communityController.UpdateCommunityPrivacy)
	}

	// Subscription routes (protected with JWT)
	subscriptionRoutes := api.Group("/subscriptions")
	subscriptionRoutes.Use(jwtMiddleware.AuthMiddleware())
	{
		subscriptionRoutes.POST("", subscriptionController.SubscribeUser)
		subscriptionRoutes.DELETE("", subscriptionController.UnsubscribeUser)
		subscriptionRoutes.GET("/communities/:community_id/count", subscriptionController.GetSubscriptionCount)
		subscriptionRoutes.GET("/communities/:community_id", subscriptionController.GetAllSubscriptionsByCommunity)
		subscriptionRoutes.GET("/users/:user_id/communities/:community_id", subscriptionController.GetSubscriptionByUserAndCommunity)
	}

	// Reaction routes (protected with JWT)
	postRoutes := api.Group("/posts")
	postRoutes.Use(jwtMiddleware.AuthMiddleware())
	{
		postRoutes.POST("/:post_id/reactions", reactionController.AddReaction)
		postRoutes.DELETE("/:post_id/reactions", reactionController.RemoveReaction)
		postRoutes.GET("/:post_id/reactions/count", reactionController.GetReactionCountByPost)
		postRoutes.GET("/:post_id/reactions/me", reactionController.GetUserReactionOnPost)
	}

	// Feed routes (protected with JWT)
	feedRoutes := api.Group("/feed")
	feedRoutes.Use(jwtMiddleware.AuthMiddleware())
	{
		feedRoutes.GET("", feedController.GetUserFeed)
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
