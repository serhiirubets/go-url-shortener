package main

import (
	"api/internal/link"
	"api/internal/stat"
	"api/internal/user"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Migrator().AutoMigrate(&link.Link{}, &user.User{}, stat.Stat{})
}
