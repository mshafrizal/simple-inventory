package persistence

import (
	"context"
	"simple-inventory/internal/domain/entity"
	"simple-inventory/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type sessionRepositoryImpl struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) repository.SessionRepository {
	return &sessionRepositoryImpl{db: db}
}

func (r *sessionRepositoryImpl) Create(ctx context.Context, session *entity.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepositoryImpl) GetByToken(ctx context.Context, token string) (*entity.Session, error) {
	var session entity.Session
	err := r.db.WithContext(ctx).Preload("User").Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepositoryImpl) GetByUserID(ctx context.Context, userID uint) ([]*entity.Session, error) {
	var sessions []*entity.Session
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&sessions).Error
	return sessions, err
}

func (r *sessionRepositoryImpl) Delete(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&entity.Session{}).Error
}

func (r *sessionRepositoryImpl) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&entity.Session{}).Error
}

func (r *sessionRepositoryImpl) DeleteByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&entity.Session{}).Error
}
