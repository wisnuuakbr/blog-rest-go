package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wisnuuakbr/blog-rest-go/api/controllers"
	"github.com/wisnuuakbr/blog-rest-go/api/middleware"
)

func GetRouter(r *gin.Engine) {
	// User routes
	r.POST("/api/register", controllers.Register)
	r.POST("/api/login", controllers.Login)
	r.Use(middleware.RequireAuth)
	r.POST("/api/logout", controllers.Logout)
	userRouter := r.Group("/api/users")
	{
		userRouter.GET("/", controllers.GetUsers)
		userRouter.GET("/:id/edit", controllers.Edit)
		userRouter.PUT("/:id/update", controllers.Update)
		userRouter.DELETE("/:id/delete", controllers.Delete)
		userRouter.GET("/all-trash", controllers.GetTrashedUsers)
		userRouter.DELETE("/delete-permanent/:id", controllers.PermanentDelete)
	}
}