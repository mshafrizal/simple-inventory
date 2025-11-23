package usecase

import (
	"context"
	"errors"
	"simple-inventory/internal/domain/entity"
	"simple-inventory/internal/domain/repository"

	"gorm.io/gorm"
)

type LocationUseCase struct {
	locationRepo repository.LocationRepository
}

func NewLocationUseCase(locationRepo repository.LocationRepository) *LocationUseCase {
	return &LocationUseCase{
		locationRepo: locationRepo,
	}
}

func (uc *LocationUseCase) CreateLocation(ctx context.Context, location *entity.Location) error {
	existing, err := uc.locationRepo.GetByCode(ctx, location.Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		return errors.New("location with this code already exists")
	}

	return uc.locationRepo.Create(ctx, location)
}

func (uc *LocationUseCase) GetLocation(ctx context.Context, id uint) (*entity.Location, error) {
	return uc.locationRepo.GetByID(ctx, id)
}

func (uc *LocationUseCase) GetLocationByCode(ctx context.Context, code string) (*entity.Location, error) {
	return uc.locationRepo.GetByCode(ctx, code)
}

func (uc *LocationUseCase) UpdateLocation(ctx context.Context, location *entity.Location) error {
	existing, err := uc.locationRepo.GetByID(ctx, location.ID)
	if err != nil {
		return err
	}

	if existing.Code != location.Code {
		codeCheck, err := uc.locationRepo.GetByCode(ctx, location.Code)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if codeCheck != nil {
			return errors.New("location with this code already exists")
		}
	}

	return uc.locationRepo.Update(ctx, location)
}

func (uc *LocationUseCase) DeleteLocation(ctx context.Context, id uint) error {
	return uc.locationRepo.Delete(ctx, id)
}

func (uc *LocationUseCase) ListLocations(ctx context.Context, limit, offset int) ([]*entity.Location, error) {
	return uc.locationRepo.List(ctx, limit, offset)
}

func (uc *LocationUseCase) SearchLocations(ctx context.Context, query string, limit, offset int) ([]*entity.Location, error) {
	return uc.locationRepo.Search(ctx, query, limit, offset)
}
