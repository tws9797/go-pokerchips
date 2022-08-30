package routers

import (
	"github.com/gin-gonic/gin"
	"go-pokerchips/controllers"
)

type RoomRouteController struct {
	roomController controllers.RoomController
}

func NewRoomRouteController(roomController controllers.RoomController) RoomRouteController {
	return RoomRouteController{roomController}
}

func (rc *RoomRouteController) RoomRoute(rg *gin.RouterGroup) {
	router := rg.Group("/room")
	router.POST("/", rc.roomController.CreateRoom)
}