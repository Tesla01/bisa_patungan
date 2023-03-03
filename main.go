package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"tesla01/bisa_patungan/handler"
	"tesla01/bisa_patungan/user"
	"tesla01/bisa_patungan/utility"
)

func main() {

	dsn := "root:@tcp(127.0.0.1:3306)/golang_crowdfunding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	utilityRepository := utility.NewRepository()

	userService := user.NewService(userRepository)
	utilityService := utility.NewService(utilityRepository)

	//Test
	userService.SaveAvatar(1, "images/1-profile.png")

	userHandler := handler.NewUserHandler(userService)
	utilityHandler := handler.NewUtilityHandler(utilityService)

	router := gin.Default()

	api := router.Group("/api/v1")

	api.GET("/check", utilityHandler.CheckHealth)
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checker", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)

	err = router.Run(":8080")
	if err != nil {
		fmt.Println("Error Start Server")
		return
	}

}
