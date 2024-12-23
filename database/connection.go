package database

import (
	"QuickPicsAuth/models"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_TIMEZONE"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	DB = db

	db.AutoMigrate(&models.User{})

	return db, nil

}

// func createDatabase() {
// 	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable TimeZone=%s", Config("DB_HOST"), Config("DB_PORT"), Config("DB_USER"), Config("DB_PASSWORD"), Config("DB_TIMEZONE"))
// 	DB, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// 	createDatabaseCommand := fmt.Sprintf("CREATE DATABASE %s", Config("DB_NAME"))
// 	DB.Exec(createDatabaseCommand)
// }

// func main() {
// 	createDatabase()
// }
