package routes

import (
	"github.com/gin-gonic/gin"
	"go-pokerchips/controllers"
	"go-pokerchips/middlewares"
	"go-pokerchips/services"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup, userService services.UserService) {

	router := rg.Group("users")
	router.Use(middlewares.DeserializeUser(userService))
	router.GET("/me", uc.userController.GetMe)
}
