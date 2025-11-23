package persistence

import (
	"context"
	"simple-inventory/internal/domain/entity"
	"simple-inventory/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type inventoryTransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewInventoryTransactionRepository(db *gorm.DB) repository.InventoryTransactionRepository {
	return &inventoryTransactionRepositoryImpl{db: db}
}

func (r *inventoryTransactionRepositoryImpl) Create(ctx context.Context, transaction *entity.InventoryTransaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *inventoryTransactionRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.InventoryTransaction, error) {
	var transaction entity.InventoryTransaction
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("User").
		First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *inventoryTransactionRepositoryImpl) GetByProductID(ctx context.Context, productID uint, limit, offset int) ([]*entity.InventoryTransaction, error) {
	var transactions []*entity.InventoryTransaction
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("User").
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}

func (r *inventoryTransactionRepositoryImpl) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]*entity.InventoryTransaction, error) {
	var transactions []*entity.InventoryTransaction
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}

func (r *inventoryTransactionRepositoryImpl) GetByDateRange(ctx context.Context, start, end time.Time, limit, offset int) ([]*entity.InventoryTransaction, error) {
	var transactions []*entity.InventoryTransaction
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("User").
		Where("created_at BETWEEN ? AND ?", start, end).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}

func (r *inventoryTransactionRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.InventoryTransaction, error) {
	var transactions []*entity.InventoryTransaction
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}
