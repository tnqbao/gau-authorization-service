package controller

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ClaimsToken struct {
	JID        uuid.UUID `json:"jid"`
	UserID     uuid.UUID `json:"user_id"`
	Permission string    `json:"permission"`
	FullName   string    `json:"fullname"`
	jwt.RegisteredClaims
}

type CreateNewTokenRequest struct {
	UserID     uuid.UUID `json:"user_id"`
	Permission string    `json:"permission"`
}
