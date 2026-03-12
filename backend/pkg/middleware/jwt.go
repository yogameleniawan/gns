package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/base-go/backend/pkg/cache"
	"github.com/base-go/backend/pkg/config"
	"github.com/base-go/backend/pkg/response"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// this private type is required if we're parsing data claim to other
type contextKey string

var (
	ContextKey    contextKey = contextKey("jwt-claims")
	ContextUserID contextKey = contextKey("user-id")
	ContextEmail  contextKey = contextKey("email")
	ContextRole   contextKey = contextKey("role")
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	cfg := config.GetConfig()
	res := response.JSON{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			res.Code = http.StatusUnauthorized
			res.Message = "Missing Authorization header"
			response.ResponseJSON(w, res.Code, res)

			return
		}

		// split token from bearer
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			res.Code = http.StatusUnauthorized
			res.Message = "Invalid Authorization header format"
			response.ResponseJSON(w, res.Code, res)

			return
		}

		tokenStr := parts[1]
		claims := Claims{}

		if len(strings.Split(tokenStr, ".")) != 3 {
			res.Code = http.StatusUnauthorized
			res.Message = "Invalid token"
			response.ResponseJSON(w, res.Code, res)

			return
		}

		// Parse dan verifikasi token
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(cfg.Auth.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			res.Code = http.StatusUnauthorized
			res.Message = "Invalid token"
			response.ResponseJSON(w, res.Code, res)

			return
		}

		ch := cache.NewCache()
		accessTokenKey := fmt.Sprintf("%s_USER_ACCESS_TOKEN", claims.Email)

		accessTokenExist := ch.Exists(ctx, accessTokenKey).Val()
		if accessTokenExist < 1 {
			res.Code = http.StatusUnauthorized
			res.Message = "Expired token"
			response.ResponseJSON(w, res.Code, res)

			return
		}

		accessTokenValue, err := ch.Get(ctx, accessTokenKey)
		if accessTokenValue != tokenStr {
			res.Code = http.StatusUnauthorized
			res.Message = "Invalid token"
			response.ResponseJSON(w, res.Code, res)

			return
		}

		ctxRes := response.UserContext{}
		ctxRes.UserID = claims.UserID
		ctxRes.Email = claims.Email
		ctxRes.Role = claims.Role

		sharingCtx := context.WithValue(ctx, "user_context", ctxRes)
		next.ServeHTTP(w, r.WithContext(sharingCtx))
	})
}
