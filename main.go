package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/wisnuuakbr/blog-rest-go/api/router"
	"github.com/wisnuuakbr/blog-rest-go/config"
	"github.com/wisnuuakbr/blog-rest-go/db/initializers"
)

func init() {
	config.LoadEnv()
	initializers.ConnectDB()
}

func main() {
	fmt.Println("Hello Bro!")
	// set releaseMode for production
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	router.GetRouter(r)

	r.Run()
}