package controllers

import (
	"github.com/gin-gonic/gin"
	"go-pokerchips/models"
	"go-pokerchips/services"
	"net/http"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{userService}
}

func (uc *UserController) GetMe(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*models.DBUser)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": currentUser}})
}
