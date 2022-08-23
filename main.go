package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-pokerchips/config"
	"go-pokerchips/controllers"
	"go-pokerchips/routes"
	"go-pokerchips/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

const (
	WORKDIR = "."
	TIMEOUT = 20
)

var (
	server      *gin.Engine
	mongoClient *mongo.Client
	redisClient *redis.Client

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController
)

func initRedis(cfg config.Config, ctx context.Context) *redis.Client {
	// Create a new Redis client
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisUri,
	})

	_, err := client.Ping(ctx).Result()

	// Test the connection
	if err != nil {
		panic(err)
	}

	// Test Redis save
	err = client.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB", 0).Err()
	if err != nil {
		panic(err)
	}

	// Test Redis read
	_, err = client.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	fmt.Println("Redis successfully connected...")

	return client
}

func initMongo(cfg config.Config, ctx context.Context) *mongo.Client {
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DBUri))

	if err != nil {
		panic(err)
	}

	// Ping the primary
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	return client
}

func main() {
	// Load the .env variables
	cfg, err := config.LoadConfig(WORKDIR)
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	// Create a non-nil, empty Context
	ctx, cancel := context.WithTimeout(context.TODO(), TIMEOUT*1000*time.Millisecond)

	redisClient = initRedis(cfg, ctx)
	mongoClient = initMongo(cfg, ctx)

	authCollection = mongoClient.Database("poker-chips").Collection("users")
	userService = services.NewUserService(authCollection)
	authService = services.NewAuthService(authCollection)
	AuthController = controllers.NewAuthController(authService, userService)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)

	// Create the Gin Engine instance
	server = gin.Default()

	defer cancel()

	// Disconnect mongoDB
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			fmt.Println("MongoDB disconnected")
		}
	}()

	router := server.Group("/api")
	router.GET("/health-checker", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Welcome to Gin!"})
	})

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)

	log.Fatal(server.Run(":" + cfg.Port))
}
