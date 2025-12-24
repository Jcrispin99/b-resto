package routes

import (
	"b-resto/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetupProductTemplateVariantRoutes configura rutas para generar variantes
func SetupProductTemplateVariantRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Endpoint para generar variantes con atributos
		api.POST("/product-templates/:id/generate-variants", func(c *gin.Context) {
			id := c.Param("id")

			var request struct {
				AttributeValueIDs []uint `json:"attribute_value_ids" binding:"required"`
			}

			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			parsedID, err := strconv.ParseUint(id, 10, 64)
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid template ID"})
				return
			}
			templateID := uint(parsedID)

			productService := services.NewProductService()
			variants, err := productService.GenerateVariantsFromAttributes(templateID, request.AttributeValueIDs)

			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, gin.H{
				"message":  "Variants generated successfully",
				"variants": variants,
			})
		})
	}
}
