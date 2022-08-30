package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-pokerchips/config"
	"go-pokerchips/controllers"
	"go-pokerchips/hub"
	"go-pokerchips/routers"
	"go-pokerchips/services"
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
	r           *gin.Engine
	mongoClient *mongo.Client
	redisClient *redis.Client

	roomCollection      *mongo.Collection
	roomService         services.RoomService
	roomController      controllers.RoomController
	roomRouteController routers.RoomRouteController
)

func main() {

	cfg, err := config.LoadConfig(WORKDIR)
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT*time.Second)
	defer cancel()

	// Get redis and mongodb connection
	redisClient = config.InitRedis(cfg, ctx)
	mongoClient = config.InitMongo(cfg, ctx)
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Register all routes, controllers and services
	db := mongoClient.Database("poker-chips")
	roomCollection = db.Collection("rooms")
	roomService = services.NewRoomService(roomCollection)
	roomController = controllers.NewRoomController(roomService)
	roomRouteController = routers.NewRoomRouteController(roomController)

	// Start the websocket hub
	r = gin.Default()
	h := hub.NewHub()

	// Serve local file
	r.Use(static.Serve("/", static.LocalFile("./public", false)))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "pong"})
	})

	apiRouter := r.Group("/api")
	{
		roomRouteController.RoomRoute(apiRouter)
	}

	r.GET("/ws", func(c *gin.Context) {
		name, _ := c.Request.URL.Query()["name"]
		roomId, _ := c.Request.URL.Query()["uri"]

		ro := h.FindRoomByUri(roomId[0])

		fmt.Println(ro)

		if ro == nil {
			fmt.Println("room is nil in memory")

			room, _ := roomService.FindRoomByUri(roomId[0])
			ro = h.CreateRoom(room)
		}

		hub.ServeWS(h, ro, name[0], c)
	})

	log.Fatal(r.Run("localhost:" + cfg.Port))
}
