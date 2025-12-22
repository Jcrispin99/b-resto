package controllers

import (
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProfile godoc
// @Summary      Obtener perfil del usuario
// @Description  Obtiene el perfil del usuario autenticado desde el token JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "user: datos del usuario"
// @Failure      400  {object}  map[string]string       "error: invalid token claims"
// @Router       /profile [get]
// @Security     Bearer
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

// GetUsers godoc
// @Summary      Listar usuarios
// @Description  Obtiene lista de todos los usuarios (solo admin)
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "users: array de usuarios"
// @Failure      500  {object}  map[string]string       "error: Failed to fetch users"
// @Router       /admin/users [get]
// @Security     Bearer
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

// CreateResource godoc
// @Summary      Crear recurso
// @Description  Crea un nuevo recurso (solo admin)
// @Tags         resources
// @Accept       json
// @Produce      json
// @Success      201  {object}  map[string]string  "message: Resource created by admin"
// @Router       /admin/resources [post]
// @Security     Bearer
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

// UpdateResource godoc
// @Summary      Actualizar recurso
// @Description  Actualiza un recurso (endpoint de ejemplo)
// @Tags         resources
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "ID del recurso"
// @Success      200  {object}  map[string]string  "message: Resource updated"
// @Router       /resources/{id} [put]
// @Security     Bearer
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

// DeleteResource godoc
// @Summary      Eliminar recurso
// @Description  Elimina un recurso (solo admin)
// @Tags         resources
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "ID del recurso"
// @Success      200  {object}  map[string]string  "message: Resource deleted by admin"
// @Router       /admin/resources/{id} [delete]
// @Security     Bearer
func DeleteResource(c *gin.Context) {
	resourceID := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"message":     "Resource deleted successfully",
		"resource_id": resourceID,
	})
}

// GetPublicData godoc
// @Summary      Obtener datos públicos
// @Description  Obtiene datos públicos sin autenticación
// @Tags         public
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string  "message: Public data accessible to everyone"
// @Router       /public/data [get]
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
