package controllers

import (
	"net/http"
	"ssr-metaverse/internal/auth/services"

	"github.com/gin-gonic/gin"
)

// LoginRequest represents the payload for user login.
// It contains the username and password required for authentication.
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

// LoginResponse represents the response returned after a successful login.
// It contains the JWT token generated for the authenticated user.
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ErrorResponse represents an error response.
// It provides a message describing the error.
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

type AuthController struct {
	Service *services.UserService
}

func NewAuthController(service *services.UserService) *AuthController {
	return &AuthController{Service: service}
}

// Login godoc
// @Summary Authenticate user and generate JWT token
// @Description Authenticates a user using the provided username and password, then returns a JWT token if the credentials are valid.
// @Tags authentication
// @Accept json
// @Produce json
// @Param login body LoginRequest true "User login credentials"
// @Success 200 {object} LoginResponse "Successful authentication"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 401 {object} ErrorResponse "Unauthorized - invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var input LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := ctrl.Service.Authenticate(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	roles, err := ctrl.Service.GetUserRoles(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	token, err := services.GenerateToken(user.ID, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}
