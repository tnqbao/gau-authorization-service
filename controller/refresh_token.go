package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-authorization-service/entity"
	"github.com/tnqbao/gau-authorization-service/utils"
)

func (ctrl *Controller) CreateNewToken(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Create new token request received")

	var request CreateNewTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Token] Failed to bind JSON request")
		utils.JSON400(c, "Invalid request body")
		return
	}

	if request.UserID == uuid.Nil {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Token] User ID is required but not provided")
		utils.JSON400(c, "User ID is required")
		return
	}

	if request.Permission == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Token] Permission is required but not provided for user: %s", request.UserID.String())
		utils.JSON400(c, "Permission is required")
		return
	}

	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Token] X-Device-ID header is required but not provided for user: %s", request.UserID.String())
		utils.JSON400(c, "X-Device-ID header is required")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Creating new token for user: %s, device: %s, permission: %s", request.UserID.String(), deviceID, request.Permission)

	// Generate UUID for refresh token instead of allocating from bitmap
	refreshTokenID, err := ctrl.Repository.GenerateRefreshTokenID(c.Request.Context())
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Token] Failed to generate refresh token ID for user: %s", request.UserID.String())
		utils.JSON500(c, "Could not generate refresh token ID")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Generated refresh token ID: %s for user: %s", refreshTokenID.String(), request.UserID.String())

	refreshTokenPlain := ctrl.GenerateToken()
	refreshTokenHashed := ctrl.hashToken(refreshTokenPlain)
	refreshTokenExpiry := time.Now().Add(30 * 24 * time.Hour)

	refreshTokenModel := &entity.RefreshToken{
		ID:        refreshTokenID,
		UserID:    request.UserID,
		Token:     refreshTokenHashed,
		DeviceID:  deviceID,
		ExpiresAt: refreshTokenExpiry,
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Creating refresh token record for user: %s, device: %s, expires_at: %s",
		request.UserID.String(), deviceID, refreshTokenExpiry.Format(time.RFC3339))

	if err := ctrl.Repository.CreateRefreshToken(refreshTokenModel); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Token] Failed to save refresh token for user: %s", request.UserID.String())
		utils.JSON500(c, "Could not store refresh token")
		return
	}

	accessTokenDuration := time.Duration(ctrl.Config.EnvConfig.JWT.Expire) * time.Minute
	if accessTokenDuration <= 0 {
		accessTokenDuration = 15 * time.Minute
	}
	accessTokenExpiry := time.Now().Add(accessTokenDuration)

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Generating access token for user: %s, duration: %s",
		request.UserID.String(), accessTokenDuration.String())

	claims := &ClaimsToken{
		JID:        refreshTokenModel.ID,
		UserID:     request.UserID,
		Permission: request.Permission,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := ctrl.CreateAccessTokenModel(*claims)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Token] Failed to create access token for user: %s", request.UserID.String())
		utils.JSON500(c, "Could not create access token")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Token creation completed successfully for user: %s, device: %s, expires_in: %d",
		request.UserID.String(), deviceID, int(accessTokenDuration.Seconds()))

	utils.JSON200(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshTokenPlain,
		"expires_in":    int(accessTokenDuration.Seconds()),
	})
}

func (ctrl *Controller) RenewAccessToken(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Renew access token request received")

	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		refreshToken, _ = c.Cookie("refresh_token")
	}

	deviceID := c.GetHeader("X-Device-ID")
	if refreshToken == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Token] Refresh token is required but not provided")
		utils.JSON400(c, "Refresh token is required")
		return
	}
	if deviceID == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Token] Device ID is required but not provided")
		utils.JSON400(c, "Device ID is required")
		return
	}

	oldAccessToken := c.GetHeader("X-Old-Access-Token")
	if oldAccessToken == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Token] Old access token is required but not provided for device: %s", deviceID)
		utils.JSON400(c, "Old access token is required")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Renewing access token for device: %s", deviceID)

	hashedRefreshToken := ctrl.hashToken(refreshToken)
	refreshTokenModel, err := ctrl.Repository.GetRefreshTokenByTokenAndDevice(hashedRefreshToken, deviceID)
	if err != nil || refreshTokenModel == nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Token] Failed to get refresh token for device: %s", deviceID)
		handleTokenError(c, err)
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Found refresh token for user: %s, device: %s",
		refreshTokenModel.UserID.String(), deviceID)

	// Kiểm tra hạn sử dụng của refresh token
	if time.Now().After(refreshTokenModel.ExpiresAt) {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Token] Refresh token expired for user: %s, device: %s, expired_at: %s",
			refreshTokenModel.UserID.String(), deviceID, refreshTokenModel.ExpiresAt.Format(time.RFC3339))
		utils.JSON401(c, "Refresh token expired")
		return
	}

	claims, err := ctrl.DecodeAccessToken(oldAccessToken)
	if err != nil || claims == nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Token] Invalid old access token for user: %s, device: %s",
			refreshTokenModel.UserID.String(), deviceID)
		utils.JSON401(c, "Invalid old access token")
		return
	}

	duration := time.Duration(ctrl.Config.EnvConfig.JWT.Expire) * time.Minute
	if duration <= 0 {
		duration = 15 * time.Minute
	}
	accessTokenExpiry := time.Now().Add(duration)

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Generating new access token for user: %s, device: %s, duration: %s",
		claims.UserID.String(), deviceID, duration.String())

	newClaims := ClaimsToken{
		JID:        refreshTokenModel.ID,
		UserID:     claims.UserID,
		Permission: claims.Permission,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := ctrl.CreateAccessTokenModel(newClaims)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Token] Failed to create new access token for user: %s", claims.UserID.String())
		utils.JSON500(c, "Could not create access token")
		return
	}

	ctrl.SetAccessCookie(c, accessToken, int(duration.Seconds()))

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Access token renewed successfully for user: %s, device: %s, expires_in: %d",
		claims.UserID.String(), deviceID, int(duration.Seconds()))

	utils.JSON200(c, gin.H{
		"access_token": accessToken,
		"expires_in":   int(duration.Seconds()),
	})
}

func (ctrl *Controller) CheckAccessToken(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Check access token request received")

	token := c.Query("token")
	if token == "" {
		token = c.GetHeader("Authorization")
	}
	if token == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Token] Access token is required but not provided")
		utils.JSON400(c, "Access token is required")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Validating access token")

	claims, err := utils.ValidateToken(c.Request.Context(), token, ctrl.Config.EnvConfig, ctrl.Repository)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Token] Access token validation failed")
		utils.JSON401(c, err.Error())
		return
	}

	if claims == nil {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Token] Invalid access token - claims are nil")
		utils.JSON401(c, "Invalid access token")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Access token validated successfully for user")

	utils.JSON200(c, gin.H{
		"message": "Access token is valid",
	})
}

func (ctrl *Controller) RevokeToken(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Revoke token request received")

	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		refreshToken, _ = c.Cookie("refresh_token")
	}

	if refreshToken == "" {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Token] No refresh token provided in header or cookie")
		utils.JSON400(c, "No refresh token provided")
		c.Abort()
		return
	}

	hashedToken := ctrl.hashToken(refreshToken)
	deviceID := c.GetHeader("X-Device-ID")

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Revoking token for device: %s", deviceID)

	refreshTokenRecord, err := ctrl.Repository.GetRefreshTokenByTokenAndDevice(hashedToken, deviceID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Token] Error fetching refresh token for device: %s", deviceID)
		utils.JSON500(c, "Internal server error")
		c.Abort()
		return
	}

	if refreshTokenRecord != nil {
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Found refresh token record for user: %s, device: %s",
			refreshTokenRecord.UserID.String(), deviceID)

		rowsAffected, err := ctrl.Repository.DeleteRefreshTokenByTokenAndDevice(hashedToken, deviceID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Token] Error deleting refresh token for user: %s, device: %s",
				refreshTokenRecord.UserID.String(), deviceID)
			utils.JSON500(c, "Internal server error")
			c.Abort()
			return
		}

		if rowsAffected > 0 {
			ttl := time.Until(refreshTokenRecord.ExpiresAt)
			if ttl > 0 {
				// Use UUID-based blacklist instead of bitmap
				if err := ctrl.Repository.BlacklistRefreshTokenIDWithTTL(
					c.Request.Context(),
					refreshTokenRecord.ID,
					ttl,
				); err != nil {
					ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Token] Failed to blacklist refresh token ID: %s for user: %s",
						refreshTokenRecord.ID.String(), refreshTokenRecord.UserID.String())
				} else {
					ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Refresh token ID %s blacklisted for %s for user: %s",
						refreshTokenRecord.ID.String(), ttl.String(), refreshTokenRecord.UserID.String())
				}
			}
		}

		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Token] Refresh token revoked successfully for user: %s, device: %s",
			refreshTokenRecord.UserID.String(), deviceID)
	} else {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Token] No refresh token record found for device: %s", deviceID)
	}

	utils.JSON200(c, gin.H{"message": "Refresh token revoked successfully"})
}
