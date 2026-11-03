package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"servixa/internal/config"
	"servixa/internal/model"
	"servixa/internal/repository"
)

type RefreshTokenService struct {
	refreshTokenRepository *repository.RefreshTokenRepository
	config                 config.Config
}

func NewRefreshTokenService(
	refreshTokenRepository *repository.RefreshTokenRepository,
	cfg config.Config,
) *RefreshTokenService {
	return &RefreshTokenService{
		refreshTokenRepository: refreshTokenRepository,
		config:                 cfg,
	}
}

func (s *RefreshTokenService) Create(userID uint) (
	*model.RefreshToken,
	error,
) {
	tokenBytes := make([]byte, 32)

	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}

	token := base64.RawURLEncoding.EncodeToString(tokenBytes)

	refreshToken := &model.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(s.config.JWTRefreshTokenDuration),
		Revoked:   false,
	}

	if err := s.refreshTokenRepository.Create(refreshToken); err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (s *RefreshTokenService) Validate(
	token string,
) (*model.RefreshToken, error) {
	if token == "" {
		return nil, errors.New("refresh token is required")
	}

	refreshToken, err := s.refreshTokenRepository.FindByToken(token)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if refreshToken.Revoked {
		return nil, errors.New("refresh token has been revoked")
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, errors.New("refresh token has expired")
	}

	return refreshToken, nil
}

func (s *RefreshTokenService) Revoke(
	refreshToken *model.RefreshToken,
) error {
	if refreshToken.Revoked {
		return nil
	}

	return s.refreshTokenRepository.Revoke(refreshToken)
}

func (s *RefreshTokenService) Rotate(
	oldToken *model.RefreshToken,
	userID uint,
) (*model.RefreshToken, error) {
	if oldToken.Revoked {
		return nil, errors.New("refresh token has been revoked")
	}

	if time.Now().After(oldToken.ExpiresAt) {
		return nil, errors.New("refresh token has expired")
	}

	if oldToken.UserID != userID {
		return nil, errors.New("invalid refresh token")
	}

	if err := s.Revoke(oldToken); err != nil {
		return nil, err
	}

	return s.Create(userID)
}
