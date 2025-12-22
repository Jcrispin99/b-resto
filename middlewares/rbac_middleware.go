package middlewares

import (
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Permission string

const (
	CreatePermission Permission = "create"
	ReadPermission   Permission = "read"
	UpdatePermission Permission = "update"
	DeletePermission Permission = "delete"
)

var rolePermissions = map[models.Role][]Permission{
	models.AdminRole: {CreatePermission, ReadPermission, UpdatePermission, DeletePermission},
	models.UserRole:  {ReadPermission, UpdatePermission},
	models.GuestRole: {ReadPermission},
}

func RBACMiddleware(requiredPermission Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No claims found"})
			c.Abort()
			return
		}

		userClaims, ok := claims.(*models.Claims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
			c.Abort()
			return
		}

		userRole := userClaims.Role

		permissions, exists := rolePermissions[userRole]
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found"})
			c.Abort()
			return
		}

		hasPermission := false
		for _, permission := range permissions {
			if permission == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":    "Insufficient permissions",
				"required": requiredPermission,
				"role":     userRole,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireRole(allowedRoles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No claims found"})
			c.Abort()
			return
		}

		userClaims, ok := claims.(*models.Claims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
			c.Abort()
			return
		}

		userRole := userClaims.Role
		roleAllowed := false

		for _, role := range allowedRoles {
			if userRole == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error":          "Access denied",
				"required_roles": allowedRoles,
				"user_role":      userRole,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
