package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-pokerchips/models"
	"go-pokerchips/services"
	"log"
	"net/http"
	"strings"
)

const DefaultChipsPerUser = 1000

type RoomController struct {
	roomService services.RoomService
}

func NewRoomController(roomService services.RoomService) RoomController {
	return RoomController{roomService}
}

func (rc *RoomController) CreateRoom(c *gin.Context) {

	var room *models.CreateRoomInput

	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	room.Record = make(map[string]int)
	room.Record[room.Creator] = DefaultChipsPerUser

	newRoom, err := rc.roomService.CreateRoom(room)

	if err != nil {
		if strings.Contains(err.Error(), "room already exists") {
			c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
		} else {
			c.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		}
		return
	}

	userSession := map[string]string{
		"uri":  room.Uri,
		"name": room.Creator,
	}

	encodedStr, err := json.Marshal(userSession)

	if err != nil {
		log.Println(err)
	}

	c.SetCookie("session", string(encodedStr), 60*60*3600, "/", "localhost", false, false)
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": newRoom})
}

func (rc *RoomController) JoinRoom(c *gin.Context) {

	var joinRoom *models.JoinRoomInput

	if err := c.ShouldBindJSON(&joinRoom); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
}
