package controller

import (
	"net/http"
	"userPage/database"
	"userPage/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignIn(c *gin.Context) {
	var singin SignInInput
	if err := c.ShouldBind(&singin); err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to bind data",
		})
	}
	var user models.Users
	if err := database.DB.Where("email=?", singin.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(singin.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
