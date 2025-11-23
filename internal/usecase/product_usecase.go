package usecase

import (
	"context"
	"errors"
	"simple-inventory/internal/domain/entity"
	"simple-inventory/internal/domain/repository"

	"gorm.io/gorm"
)

type ProductUseCase struct {
	productRepo repository.ProductRepository
}

func NewProductUseCase(productRepo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
	}
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, product *entity.Product) error {
	existing, err := uc.productRepo.GetBySKU(ctx, product.SKU)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		return errors.New("product with this SKU already exists")
	}

	if product.Barcode != "" {
		existing, err = uc.productRepo.GetByBarcode(ctx, product.Barcode)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if existing != nil {
			return errors.New("product with this barcode already exists")
		}
	}

	return uc.productRepo.Create(ctx, product)
}

func (uc *ProductUseCase) GetProduct(ctx context.Context, id uint) (*entity.Product, error) {
	return uc.productRepo.GetByID(ctx, id)
}

func (uc *ProductUseCase) GetProductBySKU(ctx context.Context, sku string) (*entity.Product, error) {
	return uc.productRepo.GetBySKU(ctx, sku)
}

func (uc *ProductUseCase) GetProductByBarcode(ctx context.Context, barcode string) (*entity.Product, error) {
	return uc.productRepo.GetByBarcode(ctx, barcode)
}

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, product *entity.Product) error {
	existing, err := uc.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return err
	}

	if existing.SKU != product.SKU {
		skuCheck, err := uc.productRepo.GetBySKU(ctx, product.SKU)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if skuCheck != nil {
			return errors.New("product with this SKU already exists")
		}
	}

	return uc.productRepo.Update(ctx, product)
}

func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id uint) error {
	return uc.productRepo.Delete(ctx, id)
}

func (uc *ProductUseCase) ListProducts(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	return uc.productRepo.List(ctx, limit, offset)
}

func (uc *ProductUseCase) SearchProducts(ctx context.Context, query string, limit, offset int) ([]*entity.Product, error) {
	return uc.productRepo.Search(ctx, query, limit, offset)
}

func (uc *ProductUseCase) GetLowStockProducts(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	return uc.productRepo.GetLowStock(ctx, limit, offset)
}

func (uc *ProductUseCase) GetProductsByLocation(ctx context.Context, locationID uint, limit, offset int) ([]*entity.Product, error) {
	return uc.productRepo.GetByLocation(ctx, locationID, limit, offset)
}
