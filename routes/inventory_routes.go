package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupInventoryRoutes configura las rutas para inventories
func SetupInventoryRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/inventories", controllers.GetInventories)
		api.GET("/inventories/:id", controllers.GetInventory)
		api.GET("/inventories/warehouse/:warehouse_id/product/:product_id", controllers.GetInventoryByWarehouseAndProduct)
		api.POST("/inventories/adjust", controllers.AdjustInventory)
		api.GET("/inventories/low-stock", controllers.GetLowStockProducts)
	}
}
