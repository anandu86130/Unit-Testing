package controller

import (
	"log"
	"net/http"
	"strconv"
	"userPage/database"
	"userPage/models"

	"github.com/gin-gonic/gin"
)

func UserList(c *gin.Context) {
	var users []models.Users
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
func EditUser(c *gin.Context) {
	var user models.Users

	// Parse the user ID from the URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	// Bind JSON request body to user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// Ensure we only update the fields that are provided
	updates := map[string]interface{}{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}

	// Execute the update query
	if err := database.DB.Model(&models.Users{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Printf("Failed to update user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully updated user",
	})
}
