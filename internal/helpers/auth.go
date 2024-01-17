package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wisnuuakbr/blog-rest-go/api/middleware"
)

func GetAuthUser(c *gin.Context) *middleware.AuthUser {
	authUser, exists := c.Get("authUser")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get the user",
		})
		return  nil
	}
	
	if user, ok := authUser.(middleware.AuthUser); ok {
		return &user
	}
	return nil
}