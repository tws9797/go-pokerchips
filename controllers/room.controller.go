package controllers

import (
	"encoding/json"
	"fmt"
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

func (rc *RoomController) GetRoom(c *gin.Context) {

	var session string
	session, err := c.Cookie("session")

	if err != nil {
		fmt.Println(err)
	}

	var roomUser *models.JoinRoomInput
	if err = json.Unmarshal([]byte(session), &roomUser); err != nil {
		log.Println(err)
	}

	if c.Param("uri") != roomUser.Uri {
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "room uri does not match with session"})
		return
	}

	room, err := rc.roomService.FindRoomByUri(roomUser.Uri)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": room})
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

	var roomUser *models.JoinRoomInput

	if err := c.ShouldBindJSON(&roomUser); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	room, err := rc.roomService.FindRoomByUri(roomUser.Uri)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err = rc.roomService.RegisterUserInRoom(room.Id.Hex(), roomUser.User); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	userSession := map[string]string{
		"uri":  room.Uri,
		"name": roomUser.User,
	}

	encodedStr, err := json.Marshal(userSession)

	if err != nil {
		log.Println(err)
	}

	c.SetCookie("session", string(encodedStr), 60*60*3600, "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": room})
}
