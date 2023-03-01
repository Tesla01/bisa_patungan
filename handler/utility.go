package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tesla01/bisa_patungan/helper"
	"tesla01/bisa_patungan/utility"
)

type utilityHandler struct {
	utilityService utility.Service
}

func NewUtilityHandler(utilityService utility.Service) *utilityHandler {
	return &utilityHandler{utilityService}
}

func (h *utilityHandler) CheckHealth(c *gin.Context) {
	respons := helper.APIResponse("Server OK", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, respons)
	return
}
