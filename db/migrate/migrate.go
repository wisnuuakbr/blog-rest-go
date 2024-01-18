package main

import (
	"fmt"
	"log"

	"github.com/wisnuuakbr/blog-rest-go/config"
	"github.com/wisnuuakbr/blog-rest-go/db/initializers"
	"github.com/wisnuuakbr/blog-rest-go/internal/models"
)

func init() {
	config.LoadEnv()
	initializers.ConnectDB()
}

func main() {
	err := initializers.DB.Migrator().DropTable(models.User{}, models.Post{})
	if err != nil {
		log.Fatal("Table dropping failed", err)
	}

	err = initializers.DB.AutoMigrate(models.User{}, models.Post{})
	if err != nil {
		log.Fatal("Migration failed", err)
	}

	// Don't forget to close the database connection when done
	sqlDB, err := initializers.DB.DB()
	if err != nil {
		log.Fatal("Failed to get the database connection", err)
	}
	defer sqlDB.Close()

	// Print a success message
	fmt.Println("Migration successful!")
}