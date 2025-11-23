package usecase

import (
	"context"
	"errors"
	"simple-inventory/internal/domain/entity"
	"simple-inventory/internal/domain/repository"
)

type InventoryUseCase struct {
	productRepo     repository.ProductRepository
	transactionRepo repository.InventoryTransactionRepository
}

func NewInventoryUseCase(productRepo repository.ProductRepository, transactionRepo repository.InventoryTransactionRepository) *InventoryUseCase {
	return &InventoryUseCase{
		productRepo:     productRepo,
		transactionRepo: transactionRepo,
	}
}

func (uc *InventoryUseCase) ReceiveInventory(ctx context.Context, productID uint, quantity int, locationID *uint, userID uint, notes string) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	product.UpdateQuantity(quantity)
	if locationID != nil {
		product.LocationID = locationID
	}

	if err := uc.productRepo.Update(ctx, product); err != nil {
		return err
	}

	transaction := &entity.InventoryTransaction{
		ProductID:    productID,
		Type:         entity.TransactionTypeIn,
		Quantity:     quantity,
		ToLocationID: locationID,
		UserID:       userID,
		Notes:        notes,
	}

	return uc.transactionRepo.Create(ctx, transaction)
}

func (uc *InventoryUseCase) IssueInventory(ctx context.Context, productID uint, quantity int, userID uint, notes string) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	if product.Quantity < quantity {
		return errors.New("insufficient stock")
	}

	product.UpdateQuantity(-quantity)

	if err := uc.productRepo.Update(ctx, product); err != nil {
		return err
	}

	transaction := &entity.InventoryTransaction{
		ProductID:      productID,
		Type:           entity.TransactionTypeOut,
		Quantity:       quantity,
		FromLocationID: product.LocationID,
		UserID:         userID,
		Notes:          notes,
	}

	return uc.transactionRepo.Create(ctx, transaction)
}

func (uc *InventoryUseCase) AdjustInventory(ctx context.Context, productID uint, newQuantity int, userID uint, notes string) error {
	if newQuantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	delta := newQuantity - product.Quantity
	product.Quantity = newQuantity

	if err := uc.productRepo.Update(ctx, product); err != nil {
		return err
	}

	transaction := &entity.InventoryTransaction{
		ProductID: productID,
		Type:      entity.TransactionTypeAdjust,
		Quantity:  delta,
		UserID:    userID,
		Notes:     notes,
	}

	return uc.transactionRepo.Create(ctx, transaction)
}

func (uc *InventoryUseCase) TransferInventory(ctx context.Context, productID uint, quantity int, fromLocationID, toLocationID uint, userID uint, notes string) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	if product.LocationID == nil || *product.LocationID != fromLocationID {
		return errors.New("product is not at the source location")
	}

	if product.Quantity < quantity {
		return errors.New("insufficient stock at source location")
	}

	product.LocationID = &toLocationID

	if err := uc.productRepo.Update(ctx, product); err != nil {
		return err
	}

	transaction := &entity.InventoryTransaction{
		ProductID:      productID,
		Type:           entity.TransactionTypeTransfer,
		Quantity:       quantity,
		FromLocationID: &fromLocationID,
		ToLocationID:   &toLocationID,
		UserID:         userID,
		Notes:          notes,
	}

	return uc.transactionRepo.Create(ctx, transaction)
}

func (uc *InventoryUseCase) GetProductTransactions(ctx context.Context, productID uint, limit, offset int) ([]*entity.InventoryTransaction, error) {
	return uc.transactionRepo.GetByProductID(ctx, productID, limit, offset)
}

func (uc *InventoryUseCase) GetUserTransactions(ctx context.Context, userID uint, limit, offset int) ([]*entity.InventoryTransaction, error) {
	return uc.transactionRepo.GetByUserID(ctx, userID, limit, offset)
}
