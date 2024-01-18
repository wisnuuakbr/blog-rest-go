package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/wisnuuakbr/blog-rest-go/db/initializers"
	"github.com/wisnuuakbr/blog-rest-go/internal/format_errors"
	"github.com/wisnuuakbr/blog-rest-go/internal/models"
	"github.com/wisnuuakbr/blog-rest-go/internal/pagination"
	"github.com/wisnuuakbr/blog-rest-go/internal/validations"
	"golang.org/x/crypto/bcrypt"
)

// Register User
func Register(c *gin.Context) {
	var userInput struct {
		Name     string `json:"name" binding:"required,min=2,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
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

	// Email validation
	if validations.IsUniqueValue("users", "email", userInput.Email) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"validations": map[string]interface{}{
				"Email": "The email is already exist!",
			},
		})
		return
	}

	// Hashing the password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{
		Name: userInput.Name,
		Email: userInput.Email,
		Password: string(hashPassword),
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		format_errors.InternalServerError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Login User
func Login(c *gin.Context) {
	var userInput struct {
		Email string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	
	if c.ShouldBindJSON(&userInput) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Find user by email
	var user models.User
	initializers.DB.First(&user, "email = ?", userInput.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"error": "Invalid email or password",
		})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign in and get the complete encoded token as a string using the .env secret_key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Set expired
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome " + user.Name + "!",
	})
}

// Logout
func Logout(c *gin.Context) {
	// Clear the cookie
	c.SetCookie("Authorization", "", 0, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"successMessage": "Logout successful",
	})
}

// Get all users
func GetUsers(c *gin.Context) {
	
	var users []models.User

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

	result, err := pagination.Paginate(initializers.DB, page, perPage, nil, &users)
	if err != nil {
		format_errors.InternalServerError(c)
		return
	}

	// Return the users
	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

// Edit User
func  Edit(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	result := initializers.DB.First(&user, id)

	if err := result.Error; err != nil {
		format_errors.RecordNotFound(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": user,
	})
}

// Update User
func Update(c *gin.Context) {
	id := c.Param("id")

	var userInput struct {
		Name  string `json:"name" binding:"required,min=2,max=50"`
		Email string `json:"email" binding:"required,email"`
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

	// find user by id
	var user models.User
	result := initializers.DB.First(&user, id)

	if err := result.Error; err != nil {
		format_errors.RecordNotFound(c, err)
		return
	}

	// Email validation
	if user.Email != userInput.Email && validations.IsUniqueValue("users", "email", userInput.Email) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"validaitons": map[string]interface{}{
				"Email": "The email is already exist!",
			},
		})
		return
	}

	updateUser := models.User{
		Name:  userInput.Name,
		Email: userInput.Email,
	}
	
	result = initializers.DB.Model(&user).Updates(&updateUser)

	if result.Error != nil {
		format_errors.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Temporary Delete User
func Delete(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	result := initializers.DB.First(&user, id)
	if err := result.Error; err != nil {
		format_errors.RecordNotFound(c, err)
		return
	}

	initializers.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "The user has been deleted successfully",
	})
}

// Trashed User
func GetTrashedUsers(c *gin.Context) {
	var users []models.User

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

	result, err := pagination.Paginate(initializers.DB.Unscoped().Where("deleted_at IS NOT NULL"), page, perPage, nil, &users)
	if err != nil {
		format_errors.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

// Permanent Delete
func PermanentDelete(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := initializers.DB.Unscoped().First(&user, id).Error; err != nil {
		format_errors.RecordNotFound(c, err)
		return
	}

	initializers.DB.Unscoped().Delete(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "The user has been deleted permanently",
	})
}
