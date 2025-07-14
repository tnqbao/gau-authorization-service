package controller

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (ctrl *Controller) CreateAccessTokenModel(claims ClaimsToken) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jid":        claims.JID,
		"user_id":    claims.UserID,
		"permission": claims.Permission,
		"fullname":   claims.FullName,
		"exp":        claims.ExpiresAt.Unix(),
		"iat":        time.Now().Unix(),
	})

	return token.SignedString([]byte(ctrl.Config.EnvConfig.JWT.SecretKey))
}

func (ctrl *Controller) SetAccessCookie(c *gin.Context, token string, timeExpired int) {
	globalDomain := ctrl.Config.EnvConfig.CORS.GlobalDomain
	c.SetCookie("access_token", token, timeExpired, "/", globalDomain, false, true)
}

func (ctrl *Controller) SetRefreshCookie(c *gin.Context, token string, timeExpired int) {
	globalDomain := ctrl.Config.EnvConfig.CORS.GlobalDomain
	c.SetCookie("refresh_token", token, timeExpired, "/", globalDomain, false, true)
}

func (ctrl *Controller) GenerateToken() string {
	return uuid.NewString() + uuid.NewString()
}

func (ctrl *Controller) hashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (ctrl *Controller) CheckNullString(str *string) string {
	if str == nil || *str == "" {
		return ""
	}
	return *str
}

func (ctrl *Controller) IsValidEmail(email string) bool {
	// Simple regex for email validation
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	at := 0
	for i, char := range email {
		if char == '@' {
			at++
			if at > 1 || i == 0 || i == len(email)-1 {
				return false
			}
		} else if char == '.' && (i == 0 || i == len(email)-1 || email[i-1] == '@') {
			return false
		}
	}
	return at == 1
}

func (ctrl *Controller) IsValidPhone(phone string) bool {
	if len(phone) < 10 || len(phone) > 15 {
		return false
	}
	for _, char := range phone {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func handleTokenError(c *gin.Context, err error) {
	if err == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Refresh token not found"})
		return
	}

	switch err.Error() {
	case "record not found":
		c.JSON(http.StatusNotFound, gin.H{"error": "Refresh token not found or revoked"})
	case "refresh token expired":
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}

func (ctrl *Controller) DecodeAccessToken(tokenString string) (*ClaimsToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ClaimsToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ctrl.Config.EnvConfig.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ClaimsToken); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
