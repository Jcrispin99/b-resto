package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"b-resto/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var loginUser models.User

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Buscar usuario por email
	var storedUser models.User
	if err := config.DB.Where("email = ?", loginUser.Email).First(&storedUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !storedUser.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User account is inactive"})
		return
	}

	if err := utils.VerifyPassword(storedUser.Password, loginUser.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(storedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       storedUser.ID,
			"username": storedUser.Username,
			"email":    storedUser.Email,
			"role":     storedUser.Role,
		},
	})
}

func Register(c *gin.Context) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validaciones
	if newUser.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	if newUser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}
	if newUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	// Verificar si el email ya existe
	var existingUser models.User
	if err := config.DB.Where("email = ?", newUser.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Verificar si el username ya existe
	if err := config.DB.Where("username = ?", newUser.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Hashear password
	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	newUser.Password = hashedPassword

	// Asignar role por defecto si no se especifica
	if newUser.Role == "" {
		newUser.Role = models.UserRole
	}

	// Guardar en BD
	if err := config.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	// Sincronizar con Casbin
	roleName := string(newUser.Role) + "_role"
	config.Enforcer.AddRoleForUser(newUser.Username, roleName)
	config.Enforcer.SavePolicy()

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":        newUser.ID,
			"username":  newUser.Username,
			"email":     newUser.Email,
			"role":      newUser.Role,
			"is_active": newUser.IsActive,
		},
	})
}
