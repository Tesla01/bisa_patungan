package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tesla01/bisa_patungan/helper"
	"tesla01/bisa_patungan/user"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	// Input from user
	// Map to user struct
	// Save to DB

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		// Gin Mapping
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)

	if err != nil {

		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	user := user.FormatUser(newUser, "tokentoken")

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", user)

	c.JSON(http.StatusOK, response)
}
