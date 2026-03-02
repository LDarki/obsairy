package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUnauthorized       = errors.New("unauthorized")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, email, password string) (*User, error) {
	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hash),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password, userAgent, ip string) (*Session, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	expiresAt := time.Now().Add(24 * 7 * time.Hour)
	token, err := s.generateJWT(user.ID, expiresAt)
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     token,
		UserAgent: userAgent,
		IPAddress: ip,
		ExpiresAt: expiresAt,
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) generateJWT(userID string, expiresAt time.Time) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "very-secret-key" // fallback
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *Service) GetUserSessions(ctx context.Context, userID string) ([]Session, error) {
	return s.repo.GetSessionsByUserID(ctx, userID)
}

func (s *Service) RevokeSession(ctx context.Context, userID, sessionID string) error {
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return nil // already gone or doesn't exist
	}

	if session.UserID != userID {
		return ErrUnauthorized
	}

	return s.repo.DeleteSession(ctx, sessionID)
}

func (s *Service) RevokeAllUserSessions(ctx context.Context, userID string) error {
	return s.repo.DeleteUserSessions(ctx, userID)
}
