package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-pokerchips/config"
	"go-pokerchips/controllers"
	"go-pokerchips/hub"
	"go-pokerchips/models"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	h := hub.NewHub(roomService)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "pong"})
	})

	apiRouter := r.Group("/api")
	{
		roomRouteController.RoomRoute(apiRouter)
	}

	r.GET("/ws", func(c *gin.Context) {

		var session string
		session, err = c.Cookie("session")

		fmt.Println(session)

		var roomUser *models.JoinRoomInput
		if err = json.Unmarshal([]byte(session), &roomUser); err != nil {
			log.Println(err)
		}

		fmt.Println(roomUser)

		foundRoom := h.FindRoomByUri(roomUser.Uri)

		if foundRoom == nil {

			// Get room from database
			var room *models.DBRoom
			room, err = roomService.FindRoomByUri(roomUser.Uri)
			if err != nil {
				fmt.Println(err)
			}

			// Create the room in the memory
			foundRoom = h.CreateRoom(room)
		}

		hub.ServeWS(h, foundRoom, roomUser.User, c)
	})

	log.Fatal(r.Run("localhost:" + cfg.Port))
}
