package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tesla01/bisa_patungan/internal/model"
	"tesla01/bisa_patungan/internal/service"
	"tesla01/bisa_patungan/util"
)

type CampaignHandler struct {
	service service.CampaignService
}

func NewCampaignHandler(service service.CampaignService) *CampaignHandler {
	return &CampaignHandler{service}
}

func (h *CampaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := util.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("List of campaigns", http.StatusOK, "success", model.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	var input model.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := util.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	detailCampaign, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := util.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Campaign detail", http.StatusOK, "success", model.FormatCampaignDetail(detailCampaign))
	c.JSON(http.StatusOK, response)

}

func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var input model.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := util.FormatValidationError(err)
		// Gin Mapping
		errorMessage := gin.H{"errors": errors}

		response := util.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Get from JWT
	currentUser := c.MustGet("currentUser").(model.User)

	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)

	if err != nil {
		response := util.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success create campaign", http.StatusOK, "success", model.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)

}

func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID model.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := util.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData model.CreateCampaignInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := util.FormatValidationError(err)
		// Gin Mapping
		errorMessage := gin.H{"errors": errors}

		response := util.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(model.User)

	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := util.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success update campaign", http.StatusOK, "success", model.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *CampaignHandler) UploadImage(c *gin.Context) {
	var input model.CreateCampaignImageInput
	err := c.ShouldBind(&input)
	if err != nil {
		errors := util.FormatValidationError(err)
		// Gin Mapping
		errorMessage := gin.H{"errors": errors}
		response := util.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := util.APIResponse("Failed to upload campaign image, file error", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get from JWT
	currentUser := c.MustGet("currentUser").(model.User)
	input.User = currentUser
	userID := currentUser.ID
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := util.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.service.SaveCampaignImage(input, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := util.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	data := gin.H{"is_uploaded": true}
	response := util.APIResponse("Campaign image successfully uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)

}
