package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupProductRoutes configura todas las rutas de productos
func SetupProductRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Product Attributes
		api.GET("/product-attributes", controllers.GetProductAttributes)
		api.GET("/product-attributes/:id", controllers.GetProductAttribute)
		api.POST("/product-attributes", controllers.CreateProductAttribute)
		api.PUT("/product-attributes/:id", controllers.UpdateProductAttribute)
		api.DELETE("/product-attributes/:id", controllers.DeleteProductAttribute)

		// Product Attribute Values
		api.GET("/product-attribute-values", controllers.GetProductAttributeValues)
		api.GET("/product-attribute-values/:id", controllers.GetProductAttributeValue)
		api.POST("/product-attribute-values", controllers.CreateProductAttributeValue)
		api.PUT("/product-attribute-values/:id", controllers.UpdateProductAttributeValue)
		api.DELETE("/product-attribute-values/:id", controllers.DeleteProductAttributeValue)

		// Product Templates
		api.GET("/product-templates", controllers.GetProductTemplates)
		api.GET("/product-templates/:id", controllers.GetProductTemplate)
		api.POST("/product-templates", controllers.CreateProductTemplate)
		api.PUT("/product-templates/:id", controllers.UpdateProductTemplate)
		api.DELETE("/product-templates/:id", controllers.DeleteProductTemplate)
		api.PATCH("/product-templates/:id/toggle", controllers.ToggleProductTemplateStatus)

		// Product Variants
		api.GET("/product-variants", controllers.GetProductVariants)
		api.GET("/product-variants/:id", controllers.GetProductVariant)
		api.POST("/product-variants", controllers.CreateProductVariant)
		api.PUT("/product-variants/:id", controllers.UpdateProductVariant)
		api.DELETE("/product-variants/:id", controllers.DeleteProductVariant)
		api.PATCH("/product-variants/:id/toggle", controllers.ToggleProductVariantStatus)

		// Combos
		api.GET("/combos", controllers.GetCombos)
		api.GET("/combos/:id", controllers.GetCombo)
		api.POST("/combos", controllers.CreateCombo)
		api.PUT("/combos/:id", controllers.UpdateCombo)
		api.DELETE("/combos/:id", controllers.DeleteCombo)
		api.PATCH("/combos/:id/toggle", controllers.ToggleComboStatus)

		// Recipes
		api.GET("/recipes", controllers.GetRecipes)
		api.GET("/recipes/:id", controllers.GetRecipe)
		api.POST("/recipes", controllers.CreateRecipe)
		api.PUT("/recipes/:id", controllers.UpdateRecipe)
		api.DELETE("/recipes/:id", controllers.DeleteRecipe)
	}
}
