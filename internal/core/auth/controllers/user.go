package controllers

import (
	"net/http"
	"strconv"

	"ssr-metaverse/internal/core/auth/services"
	"ssr-metaverse/internal/core/error"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{Service: service}
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		error.RespondWithError(c, error.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid Data: " + err.Error(),
		})
		return
	}

	user, err := ctrl.Service.CreateUser(input.Username, input.Email, input.Password)
	if err != nil {
		error.RespondWithError(c, *err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (ctrl *UserController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error.RespondWithError(c, error.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
		return
	}

	user, err := ctrl.Service.GetUserByID(id)
	if err != nil {
		if apiErr, ok := err.(*error.APIError); ok {
			error.RespondWithError(c, *apiErr)
		} else {
			error.RespondWithError(c, error.APIError{
				Code:    http.StatusInternalServerError,
				Message: "Unexpected error occurred looking for the user",
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error.RespondWithError(c, error.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
		return
	}

	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		error.RespondWithError(c, error.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid Data: " + err.Error(),
		})
		return
	}

	errAPI := ctrl.Service.UpdateUser(id, input.Username, input.Password)
	if errAPI != nil {
		error.RespondWithError(c, *errAPI)
		return
	}

	c.Status(http.StatusOK)
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error.RespondWithError(c, error.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
		return
	}

	errAPI := ctrl.Service.DeleteUser(id)
	if errAPI != nil {
		error.RespondWithError(c, *errAPI)
		return
	}

	c.Status(http.StatusOK)
}
