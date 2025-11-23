package dto

type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	SKU         string   `json:"sku" binding:"required"`
	Barcode     string   `json:"barcode"`
	Description string   `json:"description"`
	Quantity    int      `json:"quantity"`
	MinQuantity int      `json:"min_quantity"`
	Price       float64  `json:"price" binding:"required,min=0"`
	LocationID  *uint    `json:"location_id"`
}

type UpdateProductRequest struct {
	Name        string   `json:"name"`
	SKU         string   `json:"sku"`
	Barcode     string   `json:"barcode"`
	Description string   `json:"description"`
	Quantity    int      `json:"quantity"`
	MinQuantity int      `json:"min_quantity"`
	Price       float64  `json:"price" binding:"min=0"`
	LocationID  *uint    `json:"location_id"`
}

type ProductResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	SKU         string           `json:"sku"`
	Barcode     string           `json:"barcode"`
	Description string           `json:"description"`
	Quantity    int              `json:"quantity"`
	MinQuantity int              `json:"min_quantity"`
	Price       float64          `json:"price"`
	LocationID  *uint            `json:"location_id"`
	Location    *LocationResponse `json:"location,omitempty"`
	IsLowStock  bool             `json:"is_low_stock"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
}
