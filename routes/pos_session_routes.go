package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupPOSSessionRoutes configura las rutas para pos sessions
func SetupPOSSessionRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/pos-sessions", controllers.GetPOSSessions)
		api.POST("/pos-sessions/open", controllers.OpenPOSSession)
		api.GET("/pos-sessions/active", controllers.GetActivePOSSessions)

		// Rutas dinámicas
		api.GET("/pos-sessions/:id", controllers.GetPOSSession)
		api.PATCH("/pos-sessions/:id/close", controllers.ClosePOSSession)
		api.GET("/pos-sessions/:id/movements", controllers.GetSessionMovements)
		api.GET("/pos-sessions/:id/cash-movements", controllers.GetCashMovements)
		api.POST("/pos-sessions/:id/cash-movements", controllers.CreateCashMovement)
		api.DELETE("/pos-sessions/:id/cash-movements/:movement_id", controllers.DeleteCashMovement)
	}
}
