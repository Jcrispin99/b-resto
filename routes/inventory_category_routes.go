package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupInventoryCategoryRoutes configura las rutas para inventory categories
func SetupInventoryCategoryRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/inventory-categories", controllers.GetInventoryCategories)
		api.GET("/inventory-categories/:id", controllers.GetInventoryCategory)
		api.POST("/inventory-categories", controllers.CreateInventoryCategory)
		api.PUT("/inventory-categories/:id", controllers.UpdateInventoryCategory)
		api.DELETE("/inventory-categories/:id", controllers.DeleteInventoryCategory)
		api.PATCH("/inventory-categories/:id/toggle", controllers.ToggleInventoryCategoryStatus)
	}
}
