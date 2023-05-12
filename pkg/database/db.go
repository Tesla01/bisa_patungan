package pkg

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tesla01/bisa_patungan/util"
)

func NewDB() (*gorm.DB, error) {

	dbHost := util.GetEnvVariable("DB_HOST", "localhost")
	dbPort := util.GetEnvVariable("DB_PORT", "3306")
	dbUsername := util.GetEnvVariable("DB_USERNAME", "")
	dbPassword := util.GetEnvVariable("DB_PASSWORD", "")
	dbDatabase := util.GetEnvVariable("DB_DATABASE", "")

	dsn := dbUsername + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbDatabase + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return db, err
	}

	return db, nil
}
