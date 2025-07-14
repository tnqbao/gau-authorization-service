package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-authorization-service/entity"
	"github.com/tnqbao/gau-authorization-service/utils"
	"log"
	"time"
)

func (ctrl *Controller) CreateNewToken(c *gin.Context) {
	var request CreateNewTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.JSON400(c, "Invalid request body")
		return
	}

	if request.UserID == uuid.Nil {
		utils.JSON400(c, "User ID is required")
		return
	}

	if request.Permission == "" {
		utils.JSON400(c, "Permission is required")
		return
	}

	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		utils.JSON400(c, "X-Device-ID header is required")
		return
	}

	// === Refresh Token ===
	refreshTokenID, err := ctrl.Repository.AllocateRefreshTokenID(c.Request.Context())
	if err != nil {
		log.Println("[CreateNewToken] Failed to allocate refresh token ID:", err)
		utils.JSON500(c, "Could not allocate refresh token ID")
		return
	}

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

	if err := ctrl.Repository.CreateRefreshToken(refreshTokenModel); err != nil {
		log.Println("[CreateNewToken] Failed to save refresh token:", err)
		_ = ctrl.Repository.ReleaseID(c.Request.Context(), refreshTokenID)
		utils.JSON500(c, "Could not store refresh token")
		return
	}

	// === Access Token ===
	accessTokenDuration := time.Duration(ctrl.Config.EnvConfig.JWT.Expire) * time.Minute
	if accessTokenDuration <= 0 {
		accessTokenDuration = 15 * time.Minute
	}
	accessTokenExpiry := time.Now().Add(accessTokenDuration)

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
		log.Println("[CreateNewToken] Failed to create access token:", err)
		utils.JSON500(c, "Could not create access token")
		return
	}

	// === Response ===
	utils.JSON200(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshTokenPlain,
		"expires_in":    int(accessTokenDuration.Seconds()),
	})
}

func (ctrl *Controller) RenewAccessToken(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		refreshToken, _ = c.Cookie("refresh_token")
	}

	deviceID := c.GetHeader("X-Device-ID")
	if refreshToken == "" {
		utils.JSON400(c, "Refresh token is required")
		return
	}
	if deviceID == "" {
		utils.JSON400(c, "Device ID is required")
		return
	}

	oldAccessToken := c.GetHeader("X-Old-Access-Token")
	if oldAccessToken == "" {
		utils.JSON400(c, "Old access token is required")
		return
	}

	hashedRefreshToken := ctrl.hashToken(refreshToken)
	refreshTokenModel, err := ctrl.Repository.GetRefreshTokenByTokenAndDevice(hashedRefreshToken, deviceID)
	if err != nil || refreshTokenModel == nil {
		handleTokenError(c, err)
		return
	}

	// Kiểm tra hạn sử dụng của refresh token
	if time.Now().After(refreshTokenModel.ExpiresAt) {
		utils.JSON401(c, "Refresh token expired")
		return
	}

	// Decode old access token
	claims, err := ctrl.DecodeAccessToken(oldAccessToken)
	if err != nil || claims == nil {
		utils.JSON401(c, "Invalid old access token")
		return
	}

	// === Access Token mới ===
	duration := time.Duration(ctrl.Config.EnvConfig.JWT.Expire) * time.Minute
	if duration <= 0 {
		duration = 15 * time.Minute
	}
	accessTokenExpiry := time.Now().Add(duration)

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
		utils.JSON500(c, "Could not create access token")
		return
	}

	ctrl.SetAccessCookie(c, accessToken, int(duration.Seconds()))

	utils.JSON200(c, gin.H{
		"access_token": accessToken,
		"expires_in":   int(duration.Seconds()),
	})
}

func (ctrl *Controller) CheckAccessToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		utils.JSON400(c, "Access token is required")
		return
	}

	claims, err := utils.ValidateToken(c.Request.Context(), token, ctrl.Config.EnvConfig, ctrl.Repository)
	if err != nil {
		utils.JSON401(c, err.Error())
		return
	}

	if claims == nil {
		utils.JSON401(c, "Invalid access token")
		return
	}

	utils.JSON200(c, gin.H{
		"message": "Access token is valid",
	})
}

func (ctrl *Controller) RevokeToken(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		refreshToken, _ = c.Cookie("refresh_token")
	}

	if refreshToken == "" {
		log.Println("No refresh token provided in header or cookie")
		utils.JSON400(c, "No refresh token provided")
		c.Abort()
		return
	}

	hashedToken := ctrl.hashToken(refreshToken)
	deviceID := c.GetHeader("X-Device-ID")

	refreshTokenRecord, err := ctrl.Repository.GetRefreshTokenByTokenAndDevice(hashedToken, deviceID)
	if err != nil {
		log.Println("Error fetching refresh token:", err)
		utils.JSON500(c, "Internal server error")
		c.Abort()
		return
	}

	if refreshTokenRecord != nil {
		rowsAffected, err := ctrl.Repository.DeleteRefreshTokenByTokenAndDevice(hashedToken, deviceID)
		if err != nil {
			log.Println("Error deleting refresh token:", err)
			utils.JSON500(c, "Internal server error")
			c.Abort()
			return
		}

		if rowsAffected > 0 {
			ttl := time.Until(refreshTokenRecord.ExpiresAt)
			if ttl > 0 {
				if err := ctrl.Repository.ReleaseAndBlacklistIDWithTTL(
					c.Request.Context(),
					refreshTokenRecord.ID,
					ttl,
				); err != nil {
					log.Println("Failed to blacklist refresh token ID with TTL:", err)
				} else {
					log.Printf("Refresh token ID %d blacklisted for %s\n", refreshTokenRecord.ID, ttl)
				}
			}
		}
	}
	utils.JSON200(c, gin.H{"message": "Refresh token revoked successfully"})
}
