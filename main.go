package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"tesla01/bisa_patungan/auth"
	"tesla01/bisa_patungan/campaign"
	"tesla01/bisa_patungan/handler"
	"tesla01/bisa_patungan/helper"
	"tesla01/bisa_patungan/payment"
	"tesla01/bisa_patungan/transaction"
	"tesla01/bisa_patungan/user"
	"tesla01/bisa_patungan/utility"
)

func main() {

	dsn := "root:@tcp(127.0.0.1:3306)/golang_crowdfunding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	// Repository
	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)
	utilityRepository := utility.NewRepository()

	// Service
	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)
	utilityService := utility.NewService(utilityRepository)
	authService := auth.NewService()

	//Handler
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	utilityHandler := handler.NewUtilityHandler(utilityService)

	router := gin.Default()
	router.Static("/images", "./images")

	api := router.Group("/api/v1")

	//Place Middleware after path
	// Util
	api.GET("/check", utilityHandler.CheckHealth)
	// User
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checker", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	// Campaign
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaigns-images", authMiddleware(authService, userService), campaignHandler.UploadImage)
	// Transaction
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notifications", transactionHandler.GetNotification)

	err = router.Run(":9000")
	if err != nil {
		fmt.Println("Error Start Server")
		return
	}

}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//Get Token from "Bearer Token"
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		currentUser, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", currentUser)

	}
}
