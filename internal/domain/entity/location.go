package entity

import "time"

type Location struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Code        string    `json:"code" gorm:"unique;not null;index"`
	Description string    `json:"description"`
	Building    string    `json:"building"`
	Floor       string    `json:"floor"`
	Aisle       string    `json:"aisle"`
	Shelf       string    `json:"shelf"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (l *Location) GetFullPath() string {
	path := l.Building
	if l.Floor != "" {
		path += "/" + l.Floor
	}
	if l.Aisle != "" {
		path += "/" + l.Aisle
	}
	if l.Shelf != "" {
		path += "/" + l.Shelf
	}
	return path
}
