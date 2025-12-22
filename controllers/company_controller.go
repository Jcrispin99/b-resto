package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCompanies godoc
// @Summary      Listar compañías
// @Description  Obtiene lista de todas las compañías y sucursales
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        is_active      query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Param        with_parent    query  string  false  "Incluir compañía matriz"    Enums(true, false)
// @Param        with_branches  query  string  false  "Incluir sucursales"         Enums(true, false)
// @Success      200  {object}  map[string]interface{}  "data: array de companies"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /companies [get]
// @Security     Bearer
func GetCompanies(c *gin.Context) {
	var companies []models.Company

	// Filtro opcional por is_active
	query := config.DB
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	// Preload de relaciones si se solicita
	if c.Query("with_parent") == "true" {
		query = query.Preload("Parent")
	}
	if c.Query("with_branches") == "true" {
		query = query.Preload("Branches")
	}

	if err := query.Find(&companies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch companies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": companies,
	})
}

// GetCompany godoc
// @Summary      Obtener compañía
// @Description  Obtiene una compañía por ID con sus relaciones
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la compañía"
// @Success      200  {object}  map[string]interface{}  "data: company con parent y branches"
// @Failure      404  {object}  map[string]string       "error: Company not found"
// @Router       /companies/{id} [get]
// @Security     Bearer
func GetCompany(c *gin.Context) {
	id := c.Param("id")
	var company models.Company

	query := config.DB
	// Cargar relaciones automáticamente
	query = query.Preload("Parent").Preload("Branches")

	if err := query.First(&company, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": company,
	})
}

// CreateCompany godoc
// @Summary      Crear compañía
// @Description  Crea una nueva compañía (casa matriz) o sucursal
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        company  body  models.Company  true  "Datos de la compañía (incluir parent_id para sucursal)"
// @Success      201  {object}  map[string]interface{}  "message y data con la compañía creada"
// @Failure      400  {object}  map[string]string       "error: validación o parent no encontrado"
// @Failure      409  {object}  map[string]string       "error: nombre ya existe"
// @Router       /companies [post]
// @Security     Bearer
func CreateCompany(c *gin.Context) {
	var company models.Company

	// Gin valida automáticamente según los binding tags del modelo
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Solo validaciones de negocio
	// Validar que parent_id existe (si es sucursal)
	if company.ParentID != nil {
		var parent models.Company
		if err := config.DB.First(&parent, *company.ParentID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent company not found"})
			return
		}
	}

	// Verificar unicidad del nombre
	var existing models.Company
	if err := config.DB.Where("name = ?", company.Name).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Company with this name already exists"})
		return
	}

	if err := config.DB.Create(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Company created successfully",
		"data":    company,
	})
}

// UpdateCompany godoc
// @Summary      Actualizar compañía
// @Description  Actualiza los datos de una compañía existente
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        id       path  int              true  "ID de la compañía"
// @Param        company  body  models.Company   true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data actualizada"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Company not found"
// @Router       /companies/{id} [put]
// @Security     Bearer
func UpdateCompany(c *gin.Context) {
	id := c.Param("id")
	var company models.Company

	// Buscar la compañía
	if err := config.DB.First(&company, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	var updateData models.Company
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar parent_id si se está actualizando
	if updateData.ParentID != nil {
		// No puede ser su propio padre
		if *updateData.ParentID == company.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Company cannot be its own parent"})
			return
		}

		// Verificar que el padre existe
		var parent models.Company
		if err := config.DB.First(&parent, *updateData.ParentID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent company not found"})
			return
		}
	}

	// Actualizar solo campos no vacíos
	if err := config.DB.Model(&company).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Company updated successfully",
		"data":    company,
	})
}

// DeleteCompany godoc
// @Summary      Eliminar compañía
// @Description  Elimina una compañía (soft delete). No permite eliminar si tiene sucursales activas.
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la compañía"
// @Success      200  {object}  map[string]string  "message: Company deleted successfully"
// @Failure      400  {object}  map[string]string  "error: tiene sucursales activas"
// @Failure      404  {object}  map[string]string  "error: Company not found"
// @Router       /companies/{id} [delete]
// @Security     Bearer
func DeleteCompany(c *gin.Context) {
	id := c.Param("id")
	var company models.Company

	if err := config.DB.First(&company, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	// Verificar si tiene sucursales activas
	var branchCount int64
	config.DB.Model(&models.Company{}).Where("parent_id = ? AND deleted_at IS NULL", id).Count(&branchCount)
	if branchCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete company with active branches"})
		return
	}

	if err := config.DB.Delete(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Company deleted successfully",
	})
}

// ToggleCompanyStatus godoc
// @Summary      Activar/Desactivar compañía
// @Description  Cambia el estado is_active de una compañía
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la compañía"
// @Success      200  {object}  map[string]interface{}  "message y data con nuevo estado"
// @Failure      404  {object}  map[string]string       "error: Company not found"
// @Router       /companies/{id}/toggle [patch]
// @Security     Bearer
func ToggleCompanyStatus(c *gin.Context) {
	id := c.Param("id")
	var company models.Company

	if err := config.DB.First(&company, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	company.IsActive = !company.IsActive

	if err := config.DB.Save(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status toggled successfully",
		"data":    company,
	})
}
