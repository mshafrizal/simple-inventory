package repository

import (
	"context"
	"simple-inventory/internal/domain/entity"
)

type LocationRepository interface {
	Create(ctx context.Context, location *entity.Location) error
	GetByID(ctx context.Context, id uint) (*entity.Location, error)
	GetByCode(ctx context.Context, code string) (*entity.Location, error)
	Update(ctx context.Context, location *entity.Location) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*entity.Location, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*entity.Location, error)
}
