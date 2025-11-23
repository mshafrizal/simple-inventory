package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"simple-inventory/internal/domain/entity"
	"simple-inventory/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type AuthUseCase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	jwtSecret   string
	jwtExpHours int
}

func NewAuthUseCase(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, jwtSecret string, jwtExpHours int) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
		jwtExpHours: jwtExpHours,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, username, email, password string) (*entity.User, error) {
	existingUser, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingUser, err = uc.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	user := &entity.User{
		Username: username,
		Email:    email,
		Password: password,
		Role:     "user",
		IsActive: true,
	}

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, username, password string) (*entity.Session, error) {
	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	session := &entity.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(uc.jwtExpHours) * time.Hour),
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	session.User = *user
	return session, nil
}

func (uc *AuthUseCase) Logout(ctx context.Context, token string) error {
	return uc.sessionRepo.Delete(ctx, token)
}

func (uc *AuthUseCase) ValidateSession(ctx context.Context, token string) (*entity.Session, error) {
	session, err := uc.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		uc.sessionRepo.Delete(ctx, token)
		return nil, errors.New("session expired")
	}

	return session, nil
}

func (uc *AuthUseCase) CleanupExpiredSessions(ctx context.Context) error {
	return uc.sessionRepo.DeleteExpired(ctx)
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
