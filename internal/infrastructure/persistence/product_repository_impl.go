package persistence

import (
	"context"
	"simple-inventory/internal/domain/entity"
	"simple-inventory/internal/domain/repository"

	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).Preload("Location").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) GetBySKU(ctx context.Context, sku string) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).Preload("Location").Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) GetByBarcode(ctx context.Context, barcode string) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).Preload("Location").Where("barcode = ?", barcode).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) Update(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Product{}, id).Error
}

func (r *productRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	var products []*entity.Product
	err := r.db.WithContext(ctx).Preload("Location").Limit(limit).Offset(offset).Find(&products).Error
	return products, err
}

func (r *productRepositoryImpl) Search(ctx context.Context, query string, limit, offset int) ([]*entity.Product, error) {
	var products []*entity.Product
	searchPattern := "%" + query + "%"
	err := r.db.WithContext(ctx).
		Preload("Location").
		Where("name ILIKE ? OR sku ILIKE ? OR barcode ILIKE ? OR description ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern).
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}

func (r *productRepositoryImpl) GetLowStock(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	var products []*entity.Product
	err := r.db.WithContext(ctx).
		Preload("Location").
		Where("quantity <= min_quantity").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}

func (r *productRepositoryImpl) GetByLocation(ctx context.Context, locationID uint, limit, offset int) ([]*entity.Product, error) {
	var products []*entity.Product
	err := r.db.WithContext(ctx).
		Preload("Location").
		Where("location_id = ?", locationID).
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}
