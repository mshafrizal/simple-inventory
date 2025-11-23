package repository

import (
	"context"
	"simple-inventory/internal/domain/entity"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id uint) (*entity.Product, error)
	GetBySKU(ctx context.Context, sku string) (*entity.Product, error)
	GetByBarcode(ctx context.Context, barcode string) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*entity.Product, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*entity.Product, error)
	GetLowStock(ctx context.Context, limit, offset int) ([]*entity.Product, error)
	GetByLocation(ctx context.Context, locationID uint, limit, offset int) ([]*entity.Product, error)
}
