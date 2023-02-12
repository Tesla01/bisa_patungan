package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"tesla01/bisa_patungan/user"
)

func main() {

	dsn := "root:@tcp(127.0.0.1:3306)/golang_crowdfunding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	user := user.User{
		Name: "Hana",
	}

	userRepository.Save(user)
}