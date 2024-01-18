package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/wisnuuakbr/blog-rest-go/db/initializers"
	"github.com/wisnuuakbr/blog-rest-go/internal/format_errors"
	"github.com/wisnuuakbr/blog-rest-go/internal/helpers"
	"github.com/wisnuuakbr/blog-rest-go/internal/models"
	"github.com/wisnuuakbr/blog-rest-go/internal/pagination"
	"github.com/wisnuuakbr/blog-rest-go/internal/validations"
	"gorm.io/gorm"
)

// Create Post
func CreatePost(c *gin.Context) {
	// get user input
	var userInput struct{
		Title      string `json:"title" binding:"required,min=2,max=200"`
		Body       string `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validations": validations.FormatValidationErrors(errs),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	// Create a post
	authID := helpers.GetAuthUser(c).ID

	post := models.Post{
		Title:      userInput.Title,
		Body:       userInput.Body,
		UserID:     authID,
	}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		format_errors.InternalServerError(c)
		return
	}

	// Return the post
	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// Get Post
func GetPost(c *gin.Context) {
	// get all
	var posts []models.Post

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		// Handle invalid page parameter
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	perPageStr := c.DefaultQuery("perPage", "5")
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 {
		// Handle invalid perPage parameter
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid perPage parameter"})
		return
	}

	preloadFunc := func(query *gorm.DB) *gorm.DB {
		return query.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email")
		})
	}

	result, err := pagination.Paginate(initializers.DB, page, perPage, preloadFunc, &posts)

	if err != nil {
		format_errors.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": result,
	})
}

// Show Post by ID
func ShowPost(c *gin.Context) {

	// get id
	id := c.Param("id")

	// find post
	var post models.Post

	result := initializers.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&post, id)

	if err := result.Error; err != nil {
		format_errors.RecordNotFound(c, err)
		return
	}

	// Return the post
	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})

}

// Update Post
func UpdatePost(c *gin.Context) {
	// Get the id from url
	id := c.Param("id")

	// Get the data from request body
	var userInput struct {
		Title      string `json:"title" binding:"required,min=2,max=200"`
		Body       string `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validations": validations.FormatValidationErrors(errs),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Find the post by id
	var post models.Post
	result := initializers.DB.First(&post, id)

	if err := result.Error; err != nil {
		format_errors.RecordNotFound(c, err)
		return
	}

	// Prepare data to update
	authID := helpers.GetAuthUser(c).ID
	updatePost := models.Post{
		Title:      userInput.Title,
		Body:       userInput.Body,
		UserID:     authID,
	}

	// Update the post
	result = initializers.DB.Model(&post).Updates(&updatePost)

	if result.Error != nil {
		format_errors.InternalServerError(c)
		return
	}

	// Return the post
	c.JSON(http.StatusOK, gin.H{
		"post": updatePost,
	})
}

// Delete Post
func DeletePost(c *gin.Context) {
	// get the id
	id  := c.Param("id")
	var post models.Post

	// Find the post
	if err := initializers.DB.Unscoped().First(&post, id).Error; err != nil {
		format_errors.RecordNotFound(c, err)
		return
	}

	// Delete the post
	initializers.DB.Unscoped().Delete(&post)

	c.JSON(http.StatusOK, gin.H{
		"message": "The post has been deleted successfully",
	})
}
