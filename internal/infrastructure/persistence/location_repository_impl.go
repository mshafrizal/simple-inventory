package persistence

import (
	"context"
	"simple-inventory/internal/domain/entity"
	"simple-inventory/internal/domain/repository"

	"gorm.io/gorm"
)

type locationRepositoryImpl struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) repository.LocationRepository {
	return &locationRepositoryImpl{db: db}
}

func (r *locationRepositoryImpl) Create(ctx context.Context, location *entity.Location) error {
	return r.db.WithContext(ctx).Create(location).Error
}

func (r *locationRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Location, error) {
	var location entity.Location
	err := r.db.WithContext(ctx).First(&location, id).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *locationRepositoryImpl) GetByCode(ctx context.Context, code string) (*entity.Location, error) {
	var location entity.Location
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&location).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *locationRepositoryImpl) Update(ctx context.Context, location *entity.Location) error {
	return r.db.WithContext(ctx).Save(location).Error
}

func (r *locationRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Location{}, id).Error
}

func (r *locationRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.Location, error) {
	var locations []*entity.Location
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&locations).Error
	return locations, err
}

func (r *locationRepositoryImpl) Search(ctx context.Context, query string, limit, offset int) ([]*entity.Location, error) {
	var locations []*entity.Location
	searchPattern := "%" + query + "%"
	err := r.db.WithContext(ctx).
		Where("name ILIKE ? OR code ILIKE ? OR building ILIKE ? OR description ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern).
		Limit(limit).
		Offset(offset).
		Find(&locations).Error
	return locations, err
}
