package controllers

import (
	"net/http"
	"ssr-metaverse/internal/auth/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Service *services.UserService
}

func NewAuthController(service *services.UserService) *AuthController {
	return &AuthController{Service: service}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.Service.Authenticate(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	roles, err := ctrl.Service.GetUserRoles(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := services.GenerateToken(user.ID, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}