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

type InventoryHandler struct {
	inventoryUseCase *usecase.InventoryUseCase
}

func NewInventoryHandler(inventoryUseCase *usecase.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{inventoryUseCase: inventoryUseCase}
}

func (h *InventoryHandler) ReceiveInventory(c *gin.Context) {
	var req dto.ReceiveInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	userID := h.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	err := h.inventoryUseCase.ReceiveInventory(c.Request.Context(), req.ProductID, req.Quantity, req.LocationID, userID, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "receive_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Inventory received successfully"})
}

func (h *InventoryHandler) IssueInventory(c *gin.Context) {
	var req dto.IssueInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	userID := h.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	err := h.inventoryUseCase.IssueInventory(c.Request.Context(), req.ProductID, req.Quantity, userID, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "issue_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Inventory issued successfully"})
}

func (h *InventoryHandler) AdjustInventory(c *gin.Context) {
	var req dto.AdjustInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	userID := h.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	err := h.inventoryUseCase.AdjustInventory(c.Request.Context(), req.ProductID, req.NewQuantity, userID, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "adjust_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Inventory adjusted successfully"})
}

func (h *InventoryHandler) TransferInventory(c *gin.Context) {
	var req dto.TransferInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	userID := h.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	err := h.inventoryUseCase.TransferInventory(c.Request.Context(), req.ProductID, req.Quantity, req.FromLocationID, req.ToLocationID, userID, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "transfer_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Inventory transferred successfully"})
}

func (h *InventoryHandler) GetProductTransactions(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_product_id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	transactions, err := h.inventoryUseCase.GetProductTransactions(c.Request.Context(), uint(productID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "list_failed", Message: err.Error()})
		return
	}

	response := make([]dto.TransactionResponse, len(transactions))
	for i, t := range transactions {
		response[i] = h.toTransactionResponse(t)
	}

	c.JSON(http.StatusOK, response)
}

func (h *InventoryHandler) getUserID(c *gin.Context) uint {
	userValue, exists := c.Get("user")
	if !exists {
		return 0
	}

	userDTO, ok := userValue.(dto.UserDTO)
	if !ok {
		return 0
	}

	return userDTO.ID
}

func (h *InventoryHandler) toTransactionResponse(t *entity.InventoryTransaction) dto.TransactionResponse {
	return dto.TransactionResponse{
		ID:             t.ID,
		ProductID:      t.ProductID,
		Type:           string(t.Type),
		Quantity:       t.Quantity,
		FromLocationID: t.FromLocationID,
		ToLocationID:   t.ToLocationID,
		UserID:         t.UserID,
		Notes:          t.Notes,
		CreatedAt:      util.FormatTimeUTC(t.CreatedAt),
	}
}
