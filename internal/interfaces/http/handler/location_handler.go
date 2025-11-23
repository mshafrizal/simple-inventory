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

type LocationHandler struct {
	locationUseCase *usecase.LocationUseCase
}

func NewLocationHandler(locationUseCase *usecase.LocationUseCase) *LocationHandler {
	return &LocationHandler{locationUseCase: locationUseCase}
}

func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var req dto.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	location := &entity.Location{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Building:    req.Building,
		Floor:       req.Floor,
		Aisle:       req.Aisle,
		Shelf:       req.Shelf,
		IsActive:    true,
	}

	if err := h.locationUseCase.CreateLocation(c.Request.Context(), location); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "creation_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{Message: "Location created successfully", Data: h.toLocationResponse(location)})
}

func (h *LocationHandler) GetLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_id"})
		return
	}

	location, err := h.locationUseCase.GetLocation(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "location_not_found"})
		return
	}

	c.JSON(http.StatusOK, h.toLocationResponse(location))
}

func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_id"})
		return
	}

	var req dto.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	location, err := h.locationUseCase.GetLocation(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "location_not_found"})
		return
	}

	if req.Name != "" {
		location.Name = req.Name
	}
	if req.Code != "" {
		location.Code = req.Code
	}
	if req.Description != "" {
		location.Description = req.Description
	}
	if req.Building != "" {
		location.Building = req.Building
	}
	if req.Floor != "" {
		location.Floor = req.Floor
	}
	if req.Aisle != "" {
		location.Aisle = req.Aisle
	}
	if req.Shelf != "" {
		location.Shelf = req.Shelf
	}
	if req.IsActive != nil {
		location.IsActive = *req.IsActive
	}

	if err := h.locationUseCase.UpdateLocation(c.Request.Context(), location); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "update_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Location updated successfully", Data: h.toLocationResponse(location)})
}

func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_id"})
		return
	}

	if err := h.locationUseCase.DeleteLocation(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "deletion_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Location deleted successfully"})
}

func (h *LocationHandler) ListLocations(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	locations, err := h.locationUseCase.ListLocations(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "list_failed", Message: err.Error()})
		return
	}

	response := make([]dto.LocationResponse, len(locations))
	for i, l := range locations {
		response[i] = h.toLocationResponse(l)
	}

	c.JSON(http.StatusOK, response)
}

func (h *LocationHandler) SearchLocations(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "query_required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	locations, err := h.locationUseCase.SearchLocations(c.Request.Context(), query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "search_failed", Message: err.Error()})
		return
	}

	response := make([]dto.LocationResponse, len(locations))
	for i, l := range locations {
		response[i] = h.toLocationResponse(l)
	}

	c.JSON(http.StatusOK, response)
}

func (h *LocationHandler) toLocationResponse(l *entity.Location) dto.LocationResponse {
	return dto.LocationResponse{
		ID:          l.ID,
		Name:        l.Name,
		Code:        l.Code,
		Description: l.Description,
		Building:    l.Building,
		Floor:       l.Floor,
		Aisle:       l.Aisle,
		Shelf:       l.Shelf,
		FullPath:    l.GetFullPath(),
		IsActive:    l.IsActive,
		CreatedAt:   util.FormatTimeUTC(l.CreatedAt),
		UpdatedAt:   util.FormatTimeUTC(l.UpdatedAt),
	}
}
