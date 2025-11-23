package repository

import (
	"context"
	"simple-inventory/internal/domain/entity"
	"time"
)

type InventoryTransactionRepository interface {
	Create(ctx context.Context, transaction *entity.InventoryTransaction) error
	GetByID(ctx context.Context, id uint) (*entity.InventoryTransaction, error)
	GetByProductID(ctx context.Context, productID uint, limit, offset int) ([]*entity.InventoryTransaction, error)
	GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]*entity.InventoryTransaction, error)
	GetByDateRange(ctx context.Context, start, end time.Time, limit, offset int) ([]*entity.InventoryTransaction, error)
	List(ctx context.Context, limit, offset int) ([]*entity.InventoryTransaction, error)
}
