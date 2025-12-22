package middlewares

import (
	"b-resto/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CasbinMiddleware verifica permisos usando Casbin
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el username del contexto (del JWT)
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Obtener ruta y método
		obj := c.Request.URL.Path
		act := c.Request.Method

		// Verificar con Casbin
		allowed, err := config.Enforcer.Enforce(username.(string), obj, act)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authorization error"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Access denied",
				"path":   obj,
				"method": act,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
