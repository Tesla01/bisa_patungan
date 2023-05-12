package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"tesla01/bisa_patungan/auth"
	handler2 "tesla01/bisa_patungan/internal/handler"
	"tesla01/bisa_patungan/internal/repository"
	"tesla01/bisa_patungan/internal/service"
	pkg "tesla01/bisa_patungan/pkg/database"
	"tesla01/bisa_patungan/util"
)

func main() {

	db, err := pkg.NewDB()

	if err != nil {
		log.Fatal(err.Error())
	}

	// Repository
	userRepository := repository.NewUserRepository(db)
	campaignRepository := repository.NewCampaignRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)

	// Service
	userService := service.NewUserService(userRepository)
	campaignService := service.NewCampaignService(campaignRepository)
	paymentService := service.NewPaymentService()
	transactionService := service.NewTransactionService(transactionRepository, campaignRepository, paymentService)
	authService := auth.NewService()

	//Handler
	userHandler := handler2.NewUserHandler(userService, authService)
	campaignHandler := handler2.NewCampaignHandler(campaignService)
	transactionHandler := handler2.NewTransactionHandler(transactionService)
	checkHealthHandler := handler2.NewHealthHandler()

	router := gin.Default()
	router.Use(cors.Default())
	router.Static("/images", "./images")

	api := router.Group("/api/v1")

	//Place Middleware after path
	// Util
	api.GET("/health", checkHealthHandler.CheckServerHealth)
	// User
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checker", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)
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

	port := util.GetEnvVariable("APP_PORT", "9000")

	err = router.Run(":" + port)
	if err != nil {
		fmt.Println("Error Start Server")
		return
	}

}

func authMiddleware(authService auth.Service, userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := util.APIResponse("unauthorized", http.StatusUnauthorized, "error", nil)
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
			response := util.APIResponse("unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := util.APIResponse("unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		currentUser, err := userService.GetUserByID(userID)
		if err != nil {
			response := util.APIResponse("unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", currentUser)

	}
}
