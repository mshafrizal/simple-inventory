package dto

type ReceiveInventoryRequest struct {
	ProductID  uint   `json:"product_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	LocationID *uint  `json:"location_id"`
	Notes      string `json:"notes"`
}

type IssueInventoryRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	Notes     string `json:"notes"`
}

type AdjustInventoryRequest struct {
	ProductID   uint   `json:"product_id" binding:"required"`
	NewQuantity int    `json:"new_quantity" binding:"required,min=0"`
	Notes       string `json:"notes"`
}

type TransferInventoryRequest struct {
	ProductID      uint   `json:"product_id" binding:"required"`
	Quantity       int    `json:"quantity" binding:"required,min=1"`
	FromLocationID uint   `json:"from_location_id" binding:"required"`
	ToLocationID   uint   `json:"to_location_id" binding:"required"`
	Notes          string `json:"notes"`
}

type TransactionResponse struct {
	ID             uint              `json:"id"`
	ProductID      uint              `json:"product_id"`
	Product        *ProductResponse  `json:"product,omitempty"`
	Type           string            `json:"type"`
	Quantity       int               `json:"quantity"`
	FromLocationID *uint             `json:"from_location_id"`
	FromLocation   *LocationResponse `json:"from_location,omitempty"`
	ToLocationID   *uint             `json:"to_location_id"`
	ToLocation     *LocationResponse `json:"to_location,omitempty"`
	UserID         uint              `json:"user_id"`
	Notes          string            `json:"notes"`
	CreatedAt      string            `json:"created_at"`
}
