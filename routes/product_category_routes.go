package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupProductCategoryRoutes configura las rutas para product categories
func SetupProductCategoryRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/product-categories", controllers.GetProductCategories)
		api.GET("/product-categories/:id", controllers.GetProductCategory)
		api.POST("/product-categories", controllers.CreateProductCategory)
		api.PUT("/product-categories/:id", controllers.UpdateProductCategory)
		api.DELETE("/product-categories/:id", controllers.DeleteProductCategory)
		api.PATCH("/product-categories/:id/toggle", controllers.ToggleProductCategoryStatus)
	}
}
