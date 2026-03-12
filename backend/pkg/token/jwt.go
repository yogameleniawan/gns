package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/base-go/backend/pkg/config"
	"github.com/base-go/backend/pkg/constants"
)

type JWTCustomClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	MitraType string `json:"mitra_type,omitempty"`
	jwt.RegisteredClaims
}

func GenerateTokenPair(userID string, email, role, mitraType string) (accessToken, refreshToken string, err error) {
	cfg := config.GetConfig()
	now := time.Now()

	accessExp := now.Add(time.Duration(cfg.Auth.TokenExpiration) * time.Second)
	accessClaims := JWTCustomClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		MitraType: mitraType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "backend",
		},
	}

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessJWT.SignedString([]byte(cfg.Auth.JwtSecret))
	if err != nil {
		return constants.EMPTY, constants.EMPTY, err
	}

	refreshExp := now.Add(time.Duration(cfg.Auth.RefreshTokenExpiration) * time.Second)
	refreshClaims := JWTCustomClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		MitraType: mitraType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "backend",
		},
	}

	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshJWT.SignedString([]byte(cfg.Auth.RefreshTokenSecret))
	if err != nil {
		return constants.EMPTY, constants.EMPTY, err
	}

	return accessString, refreshString, nil
}
