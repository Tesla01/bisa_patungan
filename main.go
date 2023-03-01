package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"tesla01/bisa_patungan/handler"
	"tesla01/bisa_patungan/user"
)

func main() {

	dsn := "root:example@tcp(127.0.0.1:3306)/golang_crowdfunding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	userService.SaveAvatar(1, "images/1-profile.png")

	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	router.GET("api/check-health", healthCheck)

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checker", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)

	router.Run(":9001")

}

type Response struct {
	Message string
}

func healthCheck(c *gin.Context) {
	respons := Response{
		Message: "OK",
	}
	c.JSON(http.StatusOK, respons)
}
