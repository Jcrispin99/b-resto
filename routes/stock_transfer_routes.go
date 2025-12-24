package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupStockTransferRoutes configura las rutas para stock transfers
func SetupStockTransferRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/stock-transfers", controllers.GetStockTransfers)
		api.GET("/stock-transfers/:id", controllers.GetStockTransfer)
		api.POST("/stock-transfers", controllers.CreateStockTransfer)
		api.PUT("/stock-transfers/:id", controllers.UpdateStockTransfer)
		api.PATCH("/stock-transfers/:id/send", controllers.SendStockTransfer)
		api.PATCH("/stock-transfers/:id/receive", controllers.ReceiveStockTransfer)
		api.PATCH("/stock-transfers/:id/cancel", controllers.CancelStockTransfer)
	}
}
