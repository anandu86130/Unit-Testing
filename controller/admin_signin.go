package controller

import (
	"net/http"
	"userPage/database"
	"userPage/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminLoginInput struct {
	Email    string `json:"email" `
	Password string `json:"password" `
}

func AdminLogin(c *gin.Context) {
	var input AdminLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin models.Admin
	if err := database.DB.Where("email = ?", input.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successfully"})
}
