package dto

type CreateLocationRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Building    string `json:"building"`
	Floor       string `json:"floor"`
	Aisle       string `json:"aisle"`
	Shelf       string `json:"shelf"`
}

type UpdateLocationRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Building    string `json:"building"`
	Floor       string `json:"floor"`
	Aisle       string `json:"aisle"`
	Shelf       string `json:"shelf"`
	IsActive    *bool  `json:"is_active"`
}

type LocationResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Building    string `json:"building"`
	Floor       string `json:"floor"`
	Aisle       string `json:"aisle"`
	Shelf       string `json:"shelf"`
	FullPath    string `json:"full_path"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
