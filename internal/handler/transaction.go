package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tesla01/bisa_patungan/internal/model"
	"tesla01/bisa_patungan/internal/service"
	"tesla01/bisa_patungan/util"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service}
}

func (h *TransactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input model.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := util.APIResponse("Failed to get campaign transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(model.User)
	input.User = currentUser

	transactions, err := h.service.GetTransactionByCampaignID(input)

	if err != nil {
		response := util.APIResponse("Failed to get campaign transactions, not authorized user", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Campaign Transactions", http.StatusOK, "success", model.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(model.User)
	userID := currentUser.ID

	transactions, err := h.service.GetTransactionByUserID(userID)

	if err != nil {
		response := util.APIResponse("Failed to get user transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("User's transaction", http.StatusOK, "success", model.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var input model.CreateTransactionInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := util.FormatValidationError(err)
		// Gin Mapping
		errorMessage := gin.H{"errors": errors}

		response := util.APIResponse("Failed to create transaction", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get from JWT
	currentUser := c.MustGet("currentUser").(model.User)

	input.User = currentUser

	newTransaction, err := h.service.CreateTransaction(input)

	if err != nil {
		response := util.APIResponse("Failed to create transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success to create transaction", http.StatusOK, "success", model.FormatTransaction(newTransaction))
	c.JSON(http.StatusOK, response)
}

func (h *TransactionHandler) GetNotification(c *gin.Context) {
	var input model.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := util.APIResponse("Failed to process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = h.service.ProcessPayment(input)

	if err != nil {
		response := util.APIResponse("Failed to process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, input)
}
