package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/wisnuuakbr/blog-rest-go/db/initializers"
	"github.com/wisnuuakbr/blog-rest-go/internal/models"
	"github.com/wisnuuakbr/blog-rest-go/internal/validations"
)

func CreateCategory(c *gin.Context) {
	var userInput struct {
		Name string `json:"name" binding:"required,min=2"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"validations": validations.FormatValidationErrors(errs),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
		return
	}

	// name unique validation
	if validations.IsUniqueValue("categories", "name", userInput.Name) || validations.IsUniqueValue("categories", "slug", slug.Make(userInput.Name)) {
		c.JSON(http.StatusConflict, gin.H{
			"validations": map[string]interface{}{
				"Name": "Name is already exists!",
			},
		})
		return
	}

	// Create category
	category := models.Category{
		Name: userInput.Name,
	}
	result := initializers.DB.Create(&category)

	if result.Error!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Can't create category",
		})
		return
	}

	// return category
	c.JSON(http.StatusOK, gin.H{
		"category": category,
	})
}