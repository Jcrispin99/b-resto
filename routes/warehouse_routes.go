package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupWarehouseRoutes configura las rutas para warehouses
func SetupWarehouseRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/warehouses", controllers.GetWarehouses)
		api.GET("/warehouses/:id", controllers.GetWarehouse)
		api.POST("/warehouses", controllers.CreateWarehouse)
		api.PUT("/warehouses/:id", controllers.UpdateWarehouse)
		api.DELETE("/warehouses/:id", controllers.DeleteWarehouse)
		api.PATCH("/warehouses/:id/toggle", controllers.ToggleWarehouseStatus)
	}
}
