package entity

import "time"

type TransactionType string

const (
	TransactionTypeIn       TransactionType = "IN"
	TransactionTypeOut      TransactionType = "OUT"
	TransactionTypeAdjust   TransactionType = "ADJUST"
	TransactionTypeTransfer TransactionType = "TRANSFER"
)

type InventoryTransaction struct {
	ID              uint            `json:"id" gorm:"primaryKey"`
	ProductID       uint            `json:"product_id" gorm:"not null;index"`
	Product         Product         `json:"product" gorm:"foreignKey:ProductID"`
	Type            TransactionType `json:"type" gorm:"not null"`
	Quantity        int             `json:"quantity" gorm:"not null"`
	FromLocationID  *uint           `json:"from_location_id" gorm:"index"`
	FromLocation    *Location       `json:"from_location,omitempty" gorm:"foreignKey:FromLocationID"`
	ToLocationID    *uint           `json:"to_location_id" gorm:"index"`
	ToLocation      *Location       `json:"to_location,omitempty" gorm:"foreignKey:ToLocationID"`
	UserID          uint            `json:"user_id" gorm:"not null;index"`
	User            User            `json:"user" gorm:"foreignKey:UserID"`
	Notes           string          `json:"notes"`
	CreatedAt       time.Time       `json:"created_at"`
}
