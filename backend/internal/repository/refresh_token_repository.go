package repository

import (
	"servixa/internal/model"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	*BaseRepository
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *RefreshTokenRepository) Create(
	refreshToken *model.RefreshToken,
) error {
	return r.db.Create(refreshToken).Error
}

func (r *RefreshTokenRepository) FindByToken(
	token string,
) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken

	err := r.db.
		Where("token = ?", token).
		First(&refreshToken).Error

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *RefreshTokenRepository) Revoke(
	refreshToken *model.RefreshToken,
) error {
	return r.db.
		Model(refreshToken).
		Update("revoked", true).
		Error
}
