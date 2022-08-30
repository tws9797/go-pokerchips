package controllers

import (
	"github.com/gin-gonic/gin"
	"go-pokerchips/models"
	"go-pokerchips/services"
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

	var room *models.RoomInput

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

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": newRoom})
}
