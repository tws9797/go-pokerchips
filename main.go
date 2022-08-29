package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-pokerchips/config"
	"go-pokerchips/hub"
	"go.mongodb.org/mongo-driver/mongo"
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
)

func main() {

	cfg, err := config.LoadConfig(WORKDIR)
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT*time.Second)
	defer cancel()

	redisClient = config.InitRedis(cfg, ctx)
	mongoClient = config.InitMongo(cfg, ctx)

	server = gin.Default()
	fmt.Println("run HUb")
	h := hub.NewHub()
	go h.Run()

	// Serve local file
	server.Use(static.Serve("/", static.LocalFile("./public", false)))

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "pong"})
	})

	server.GET("/ws", func(c *gin.Context) {
		fmt.Println("start serving websocket")
		hub.ServeWS(h, c)
	})

	log.Fatal(server.Run("localhost:" + cfg.Port))
}
