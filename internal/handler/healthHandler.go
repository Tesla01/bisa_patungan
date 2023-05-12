package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tesla01/bisa_patungan/util"
)

type healthHandler struct{}

func NewHealthHandler() *healthHandler {
	return &healthHandler{}
}

func (h *healthHandler) CheckServerHealth(c *gin.Context) {
	response := util.APIResponse("Server OK", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
	return
}
