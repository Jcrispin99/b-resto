package controllers

import (
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProfile obtiene el perfil del usuario autenticado
func GetProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"message":  "Welcome " + username.(string),
		"username": username,
		"role":     role,
	})
}

// GetUsers obtiene todos los usuarios (solo para admins	)
func GetUsers(c *gin.Context) {
	// En una aplicación real, obtendrías los usuarios de la BD
	users := []models.User{
		{Username: "admin", Role: models.AdminRole},
		{Username: "testuser", Role: models.UserRole},
		{Username: "guest", Role: models.GuestRole},
	}

	// Remover passwords de la respuesta
	var response []gin.H
	for _, user := range users {
		response = append(response, gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": response,
	})
}

// CreateResource crea un nuevo recurso (solo para admins)
func CreateResource(c *gin.Context) {
	var resource map[string]interface{}

	if err := c.ShouldBindJSON(&resource); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Resource created successfully",
		"resource": resource,
	})
}

// UpdateResource actualiza un recurso existente
func UpdateResource(c *gin.Context) {
	resourceID := c.Param("id")

	var resource map[string]interface{}
	if err := c.ShouldBindJSON(&resource); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Resource updated successfully",
		"resource_id": resourceID,
		"data":        resource,
	})
}

// DeleteResource elimina un recurso (solo para admins)
func DeleteResource(c *gin.Context) {
	resourceID := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"message":     "Resource deleted successfully",
		"resource_id": resourceID,
	})
}

// GetPublicData obtiene datos públicos (accesible para todos)
func GetPublicData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "This is public data, accessible to everyone",
		"data": []string{
			"Public item 1",
			"Public item 2",
			"Public item 3",
		},
	})
}
