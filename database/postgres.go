package database

import (
	"fmt"
	"log"

	"github.com/visitha2001/go-jwt-auth/configs"
	"github.com/visitha2001/go-jwt-auth/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the exported database connection
var DB *gorm.DB

// ConnectDB initializes the database connection and runs migrations
func ConnectDB() {
	var err error

	// Load config from .env
	host := configs.EnvConfig("DB_HOST")
	user := configs.EnvConfig("DB_USER")
	password := configs.EnvConfig("DB_PASSWORD")
	dbname := configs.EnvConfig("DB_NAME")
	port := configs.EnvConfig("DB_PORT")
	sslmode := configs.EnvConfig("DB_SSLMODE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	fmt.Println("Database connection successful")

	// Run migrations
	err = models.MigrateUser(DB)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	fmt.Println("Database migration successful")
}
