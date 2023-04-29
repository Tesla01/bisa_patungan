package helper

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func NewDB() *gorm.DB {
	dbHost := GetEnvVariable("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := GetEnvVariable("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	dbUsername := GetEnvVariable("DB_USERNAME")
	dbPassword := GetEnvVariable("DB_PASSWORD")
	dbDatabase := GetEnvVariable("DB_DATABASE")

	dsn := dbUsername + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbDatabase + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}
