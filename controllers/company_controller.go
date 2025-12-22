package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCompanies obtiene todas las compañías
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

// GetCompany obtiene una compañía por ID
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

// CreateCompany crea una nueva compañía
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

// UpdateCompany actualiza una compañía existente
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

// DeleteCompany elimina una compañía (soft delete)
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

// ToggleCompanyStatus activa/desactiva una compañía
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
