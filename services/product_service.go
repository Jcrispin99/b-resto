package services

import (
	"b-resto/config"
	"b-resto/models"
	"errors"
	"fmt"
)

// ProductService maneja la lógica de negocio de productos
type ProductService struct{}

// NewProductService crea una nueva instancia del servicio
func NewProductService() *ProductService {
	return &ProductService{}
}

// CreateTemplateWithDefaultVariant crea un template y automáticamente genera una variante default
func (s *ProductService) CreateTemplateWithDefaultVariant(template *models.ProductTemplate) (*models.ProductTemplate, *models.ProductProduct, error) {
	// Comenzar transacción
	tx := config.DB.Begin()
	if tx.Error != nil {
		return nil, nil, tx.Error
	}

	// Rollback si hay error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Crear el template
	if err := tx.Create(template).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to create template: %w", err)
	}

	// 2. Crear variante default
	defaultSKU := template.InternalReference
	if defaultSKU == "" {
		defaultSKU = fmt.Sprintf("PROD-%d", template.ID)
	}

	defaultVariant := &models.ProductProduct{
		TemplateID: template.ID,
		SKU:        defaultSKU,
		Barcode:    template.Barcode,
		SalePrice:  &template.SalePrice,
		IsActive:   template.IsActive,
	}

	if err := tx.Create(defaultVariant).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to create default variant: %w", err)
	}

	// Commit transacción
	if err := tx.Commit().Error; err != nil {
		return nil, nil, err
	}

	return template, defaultVariant, nil
}

// GenerateVariantsFromAttributes genera variantes basadas en combinaciones de atributos
func (s *ProductService) GenerateVariantsFromAttributes(templateID uint, attributeValueIDs []uint) ([]models.ProductProduct, error) {
	// Validar que el template existe
	var template models.ProductTemplate
	if err := config.DB.First(&template, templateID).Error; err != nil {
		return nil, errors.New("template not found")
	}

	// Obtener los valores de atributos
	var attributeValues []models.ProductAttributeValue
	if err := config.DB.Preload("Attribute").Find(&attributeValues, attributeValueIDs).Error; err != nil {
		return nil, err
	}

	// Comenzar transacción
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Eliminar variante default si existe (SKU tipo PROD-X)
	defaultSKU := fmt.Sprintf("PROD-%d", templateID)
	tx.Where("template_id = ? AND sku = ?", templateID, defaultSKU).Delete(&models.ProductProduct{})

	// Generar variantes (si es solo 1 atributo, crear una variante por cada valor)
	var variants []models.ProductProduct
	for _, attrValue := range attributeValues {
		sku := fmt.Sprintf("%s-%s", template.InternalReference, attrValue.Value)
		if template.InternalReference == "" {
			sku = fmt.Sprintf("PROD-%d-%s", templateID, attrValue.Value)
		}

		variant := models.ProductProduct{
			TemplateID: templateID,
			SKU:        sku,
			SalePrice:  &template.SalePrice,
			IsActive:   true,
		}

		if err := tx.Create(&variant).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Asociar el valor de atributo con la variante
		if err := tx.Model(&variant).Association("AttributeValues").Append(&attrValue); err != nil {
			tx.Rollback()
			return nil, err
		}

		variants = append(variants, variant)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return variants, nil
}
