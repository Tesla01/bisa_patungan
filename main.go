package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"tesla01/bisa_patungan/user"
)

func main() {
	//// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	//dsn := "root:@tcp(127.0.0.1:3306)/golang_crowdfunding?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//
	//fmt.Println("Connection Success")
	//
	//var users []user.User
	//length := len(users)
	//fmt.Println(length)
	//
	//db.Find(&users)
	//length = len(users)
	//fmt.Println(length)
	//
	//for _, user := range users {
	//	fmt.Printf("%s-%s\n", user.Name, user.Email)
	//	fmt.Println("========================")
	//}
	router := gin.Default()
	router.GET("/handler", handler)
	router.Run()
}

func handler(c *gin.Context) {
	dsn := "root:@tcp(127.0.0.1:3306)/golang_crowdfunding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connection Success")

	var users []user.User

	db.Find(&users)

	c.JSON(http.StatusOK, users)

	// Handler
	// Service
	// Repository
	// DB
}
