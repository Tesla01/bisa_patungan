package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tesla01/bisa_patungan/auth"
	"tesla01/bisa_patungan/internal/model"
	"tesla01/bisa_patungan/internal/service"
	"tesla01/bisa_patungan/util"
)

type UserHandler struct {
	userService service.UserService
	authService auth.Service
}

func NewUserHandler(userService service.UserService, authService auth.Service) *UserHandler {
	return &UserHandler{userService, authService}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var input model.RegisterUserInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := util.FormatValidationError(err)
		// Gin Mapping
		errorMessage := gin.H{"errors": errors}

		response := util.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)

	if err != nil {

		response := util.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {

		response := util.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	user := model.FormatUser(newUser, token)

	response := util.APIResponse("Account has been registered", http.StatusOK, "success", user)

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Login(c *gin.Context) {

	var input model.LoginInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := util.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := util.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedUser, err := h.userService.LoginUser(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := util.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedUser.ID)
	if err != nil {

		response := util.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := model.FormatUser(loggedUser, token)

	response := util.APIResponse("Success Login", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) CheckEmailAvailability(c *gin.Context) {
	var input model.CheckEmailInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := util.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := util.APIResponse("Email Checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)

	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}

		response := util.APIResponse("Email Checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{"is_available": isEmailAvailable}

	metaMessage := "Email has been registered"

	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := util.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := util.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get from JWT
	currentUser := c.MustGet("currentUser").(model.User)
	userID := currentUser.ID
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := util.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := util.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	data := gin.H{"is_uploaded": true}
	response := util.APIResponse("Avatar image successfully uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) FetchUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(model.User)

	formatter := model.FormatUser(currentUser, "")

	response := util.APIResponse("Successfuly fetch user data", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)

}
