package repository

import (
	"context"
	"simple-inventory/internal/domain/entity"
)

type SessionRepository interface {
	Create(ctx context.Context, session *entity.Session) error
	GetByToken(ctx context.Context, token string) (*entity.Session, error)
	GetByUserID(ctx context.Context, userID uint) ([]*entity.Session, error)
	Delete(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID uint) error
}
