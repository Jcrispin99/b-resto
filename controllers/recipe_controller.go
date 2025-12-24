package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetRecipes godoc
// @Summary      Listar recetas
// @Description  Obtiene lista de todas las recetas de productos
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        product_id  query  int  false  "Filtrar por producto"
// @Success      200  {object}  map[string]interface{}  "data: array de recipes"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /recipes [get]
// @Security     Bearer
func GetRecipes(c *gin.Context) {
	var recipes []models.Recipe

	query := config.DB
	if productID := c.Query("product_id"); productID != "" {
		query = query.Where("product_id = ?", productID)
	}

	if err := query.Preload("Product").Preload("Items").Find(&recipes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": recipes})
}

// GetRecipe godoc
// @Summary      Obtener receta
// @Description  Obtiene una receta por ID con todos sus ingredientes
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la receta"
// @Success      200  {object}  map[string]interface{}  "data: recipe con items"
// @Failure      404  {object}  map[string]string       "error: Recipe not found"
// @Router       /recipes/{id} [get]
// @Security     Bearer
func GetRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe

	if err := config.DB.Preload("Product").Preload("Items").First(&recipe, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": recipe})
}

// CreateRecipe godoc
// @Summary      Crear receta
// @Description  Crea una nueva receta de producto con sus ingredientes
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        recipe  body  models.Recipe  true  "Datos de la receta"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /recipes [post]
// @Security     Bearer
func CreateRecipe(c *gin.Context) {
	var recipe models.Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Recipe created successfully",
		"data":    recipe,
	})
}

// UpdateRecipe godoc
// @Summary      Actualizar receta
// @Description  Actualiza los datos de una receta existente
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        id      path  int            true  "ID de la receta"
// @Param        recipe  body  models.Recipe  true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Recipe not found"
// @Router       /recipes/{id} [put]
// @Security     Bearer
func UpdateRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe

	if err := config.DB.First(&recipe, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	var updateData models.Recipe
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&recipe).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe updated successfully",
		"data":    recipe,
	})
}

// DeleteRecipe godoc
// @Summary      Eliminar receta
// @Description  Elimina una receta (soft delete)
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la receta"
// @Success      200  {object}  map[string]string  "message: Recipe deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Recipe not found"
// @Router       /recipes/{id} [delete]
// @Security     Bearer
func DeleteRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe

	if err := config.DB.First(&recipe, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	if err := config.DB.Delete(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}
