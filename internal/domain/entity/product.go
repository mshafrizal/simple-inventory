package entity

import "time"

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	SKU         string    `json:"sku" gorm:"unique;not null;index"`
	Barcode     string    `json:"barcode" gorm:"unique;index"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity" gorm:"default:0"`
	MinQuantity int       `json:"min_quantity" gorm:"default:0"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2)"`
	LocationID  *uint     `json:"location_id" gorm:"index"`
	Location    *Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Product) IsLowStock() bool {
	return p.Quantity <= p.MinQuantity
}

func (p *Product) UpdateQuantity(delta int) {
	p.Quantity += delta
	if p.Quantity < 0 {
		p.Quantity = 0
	}
}
