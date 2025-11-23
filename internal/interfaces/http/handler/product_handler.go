package handler

import (
	"net/http"
	"simple-inventory/internal/domain/entity"
	"simple-inventory/internal/interfaces/http/dto"
	"simple-inventory/internal/interfaces/http/util"
	"simple-inventory/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
}

func NewProductHandler(productUseCase *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{productUseCase: productUseCase}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	product := &entity.Product{
		Name:        req.Name,
		SKU:         req.SKU,
		Barcode:     req.Barcode,
		Description: req.Description,
		Quantity:    req.Quantity,
		MinQuantity: req.MinQuantity,
		Price:       req.Price,
		LocationID:  req.LocationID,
	}

	if err := h.productUseCase.CreateProduct(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "creation_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{Message: "Product created successfully", Data: h.toProductResponse(product)})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_id"})
		return
	}

	product, err := h.productUseCase.GetProduct(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product_not_found"})
		return
	}

	c.JSON(http.StatusOK, h.toProductResponse(product))
}

func (h *ProductHandler) GetProductByBarcode(c *gin.Context) {
	barcode := c.Query("barcode")
	if barcode == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "barcode_required"})
		return
	}

	product, err := h.productUseCase.GetProductByBarcode(c.Request.Context(), barcode)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product_not_found"})
		return
	}

	c.JSON(http.StatusOK, h.toProductResponse(product))
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_id"})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	product, err := h.productUseCase.GetProduct(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product_not_found"})
		return
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.SKU != "" {
		product.SKU = req.SKU
	}
	if req.Barcode != "" {
		product.Barcode = req.Barcode
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	product.Quantity = req.Quantity
	product.MinQuantity = req.MinQuantity
	product.Price = req.Price
	product.LocationID = req.LocationID

	if err := h.productUseCase.UpdateProduct(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "update_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Product updated successfully", Data: h.toProductResponse(product)})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_id"})
		return
	}

	if err := h.productUseCase.DeleteProduct(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "deletion_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Product deleted successfully"})
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	products, err := h.productUseCase.ListProducts(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "list_failed", Message: err.Error()})
		return
	}

	response := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		response[i] = h.toProductResponse(p)
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "query_required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	products, err := h.productUseCase.SearchProducts(c.Request.Context(), query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "search_failed", Message: err.Error()})
		return
	}

	response := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		response[i] = h.toProductResponse(p)
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetLowStockProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	products, err := h.productUseCase.GetLowStockProducts(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "list_failed", Message: err.Error()})
		return
	}

	response := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		response[i] = h.toProductResponse(p)
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) toProductResponse(p *entity.Product) dto.ProductResponse {
	response := dto.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		SKU:         p.SKU,
		Barcode:     p.Barcode,
		Description: p.Description,
		Quantity:    p.Quantity,
		MinQuantity: p.MinQuantity,
		Price:       p.Price,
		LocationID:  p.LocationID,
		IsLowStock:  p.IsLowStock(),
		CreatedAt:   util.FormatTimeUTC(p.CreatedAt),
		UpdatedAt:   util.FormatTimeUTC(p.UpdatedAt),
	}

	if p.Location != nil {
		response.Location = &dto.LocationResponse{
			ID:       p.Location.ID,
			Name:     p.Location.Name,
			Code:     p.Location.Code,
			FullPath: p.Location.GetFullPath(),
		}
	}

	return response
}
