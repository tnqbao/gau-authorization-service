package repository

import (
	"fmt"
	"github.com/tnqbao/gau-authorization-service/entity"
	"time"
)

func (r *Repository) CreateRefreshToken(token *entity.RefreshToken) error {
	if err := r.DB.Create(token).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteRefreshTokenByTokenAndDevice(token string, deviceID string) (int64, error) {
	result := r.DB.Where("token = ? AND device_id = ?", token, deviceID).Delete(&entity.RefreshToken{})
	return result.RowsAffected, result.Error
}

func (r *Repository) GetRefreshTokenByTokenAndDevice(token string, deviceID string) (*entity.RefreshToken, error) {
	var refreshToken entity.RefreshToken
	if err := r.DB.Where("token = ? AND device_id = ?", token, deviceID).First(&refreshToken).Error; err != nil {
		return nil, err
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		if err := r.DB.Delete(&refreshToken).Error; err != nil {
			return nil, fmt.Errorf("failed to delete expired refresh token: %w", err)
		}
		return nil, fmt.Errorf("refresh token expired")
	}

	return &refreshToken, nil
}
