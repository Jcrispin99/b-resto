package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupPOSRoutes configura las rutas para POS terminals
func SetupPOSRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/pos", controllers.GetPOSTerminals)
		api.GET("/pos/:id", controllers.GetPOSTerminal)
		api.POST("/pos", controllers.CreatePOSTerminal)
		api.PUT("/pos/:id", controllers.UpdatePOSTerminal)
		api.PATCH("/pos/:id/toggle", controllers.TogglePOSStatus)
	}
}
