package main

import (
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
	err := initializers.DB.Migrator().DropTable(models.User{})
	if err != nil {
		log.Fatal("Table dropping failed", err)
	}

	err = initializers.DB.AutoMigrate(models.User{})
	if err != nil {
		log.Fatal("Migration failed", err)
	}

	// Check if the table exists after migration
	if !initializers.DB.Migrator().HasTable(&models.User{}) {
		log.Fatal("Table not found after migration")
	}

	log.Println("Migration completed successfully")

	// Don't forget to close the database connection when done
	sqlDB, _ := initializers.DB.DB()
	defer sqlDB.Close()
}