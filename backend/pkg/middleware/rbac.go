package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/base-go/backend/pkg/response"
)

// RequireAuth ensures the user is authenticated
// This should be used after JWTAuthMiddleware
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userCtx, ok := r.Context().Value("user_context").(response.UserContext)
		if !ok || userCtx.UserID == "" {
			response.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireRole returns a middleware that checks if user has any of the specified roles
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userCtx, ok := r.Context().Value("user_context").(response.UserContext)
			if !ok {
				response.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			// Check if user has any of the required roles
			hasRole := false
			for _, role := range roles {
				if userCtx.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				response.ResponseError(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RBACChecker interface for checking permissions and module access
type RBACChecker interface {
	UserHasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
	UserCanAccessModule(ctx context.Context, userID uuid.UUID, module string, action string) (bool, error)
}

// RequirePermission returns a middleware that checks if user has the specified permission
func RequirePermission(checker RBACChecker, permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userCtx, ok := r.Context().Value("user_context").(response.UserContext)
			if !ok {
				response.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			userID, err := uuid.Parse(userCtx.UserID)
			if err != nil {
				response.ResponseError(w, http.StatusUnauthorized, "Invalid user ID")
				return
			}

			hasPermission, err := checker.UserHasPermission(r.Context(), userID, permission)
			if err != nil {
				response.ResponseError(w, http.StatusInternalServerError, "Failed to check permission")
				return
			}

			if !hasPermission {
				response.ResponseError(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireModuleAccess returns a middleware that checks if user can access a module with specified action
func RequireModuleAccess(checker RBACChecker, module string, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userCtx, ok := r.Context().Value("user_context").(response.UserContext)
			if !ok {
				response.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			userID, err := uuid.Parse(userCtx.UserID)
			if err != nil {
				response.ResponseError(w, http.StatusUnauthorized, "Invalid user ID")
				return
			}

			hasAccess, err := checker.UserCanAccessModule(r.Context(), userID, module, action)
			if err != nil {
				response.ResponseError(w, http.StatusInternalServerError, "Failed to check module access")
				return
			}

			if !hasAccess {
				response.ResponseError(w, http.StatusForbidden, "Insufficient permissions for this module")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
