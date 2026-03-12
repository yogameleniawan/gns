package router

import (
	"compress/zlib"
	"net/http"

	"github.com/base-go/backend/internal/auth"
	"github.com/base-go/backend/internal/rbac"
	"github.com/base-go/backend/pkg/middleware"
	"github.com/base-go/backend/pkg/response"
	"github.com/go-chi/chi/v5"
	cmiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/secure"
)

// SetupRoutes this function for centralize setup all route in this app.
// why wee need to centralize?, it's for easies debugging if any issue
//
// parameters: all interface handlers we need to expose with rest
func SetupRoutes(
	authHandler auth.Handler,
	rbacHandler rbac.Handler,
	rbacRepo rbac.Repository,
) *chi.Mux {
	mux := chi.NewRouter()

	// chi middleware
	mux.Use(cmiddleware.Logger)
	mux.Use(cmiddleware.Recoverer)
	mux.Use(cmiddleware.RealIP)
	mux.Use(cmiddleware.NoCache)
	mux.Use(cmiddleware.GetHead)
	mux.Use(cmiddleware.Compress(zlib.BestCompression))
	mux.Use(cmiddleware.AllowContentType("application/json"))
	mux.Use(secure.New(secure.Options{
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		STSSeconds:           900,
	}).Handler)

	mux.MethodNotAllowed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := response.JSON{Code: http.StatusMethodNotAllowed, Message: "Route method not allowed"}
		response.ResponseJSON(w, res.Code, res)
	}))

	mux.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := response.JSON{Code: http.StatusNotFound, Message: "Route not found"}
		response.ResponseJSON(w, res.Code, res)
	}))

	// set cors middleware
	mux.Use(middleware.Cors())
	// set middleware rate limiter
	mux.Use(middleware.RateLimit(1000, 10))

	// set prefix v1
	mux.Route("/v1", func(r chi.Router) {

		// Authentication routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/oauth/google", authHandler.LoginWithGoogle)
			r.Post("/refresh", authHandler.RefreshToken)

			// Protected auth routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.JWTAuthMiddleware)
				r.Post("/logout", authHandler.Logout)
				r.Get("/profile", authHandler.GetProfile)
				r.Put("/profile", authHandler.UpdateProfile)
				r.Post("/change-password", authHandler.ChangePassword)
			})
		})

		// RBAC routes (protected, admin only)
		r.Route("/rbac", func(r chi.Router) {
			r.Use(middleware.JWTAuthMiddleware)
			r.Use(middleware.RequireRole("Super Admin", "Admin"))

			// Roles
			r.Route("/roles", func(r chi.Router) {
				r.Post("/", rbacHandler.CreateRole)
				r.Get("/", rbacHandler.GetAllRoles)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", rbacHandler.GetRoleByID)
					r.Put("/", rbacHandler.UpdateRole)
					r.Delete("/", rbacHandler.DeleteRole)

					// Role permissions
					r.Post("/permissions", rbacHandler.AssignPermissionsToRole)
					r.Get("/permissions", rbacHandler.GetRolePermissions)

					// Module access
					r.Post("/module-access", rbacHandler.UpdateModuleAccess)
					r.Get("/module-access", rbacHandler.GetModuleAccessByRole)
				})
			})

			// Permissions
			r.Route("/permissions", func(r chi.Router) {
				r.Post("/", rbacHandler.CreatePermission)
				r.Get("/", rbacHandler.GetAllPermissions)
				r.Get("/by-module", rbacHandler.GetPermissionsByModule)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", rbacHandler.GetPermissionByID)
					r.Put("/", rbacHandler.UpdatePermission)
					r.Delete("/", rbacHandler.DeletePermission)
				})
			})

			// User roles
			r.Route("/users/{userId}/roles", func(r chi.Router) {
				r.Post("/", rbacHandler.AssignRolesToUser)
				r.Get("/", rbacHandler.GetUserRoles)
			})

			// Permission checking (available to all authenticated users)
			r.Group(func(r chi.Router) {
				r.Post("/check-permission", rbacHandler.CheckPermission)
				r.Post("/check-module-access", rbacHandler.CheckModuleAccess)
			})
		})

		// User management routes (protected - Admin only)
		r.Route("/users", func(r chi.Router) {
			r.Use(middleware.JWTAuthMiddleware)
			r.Use(middleware.RequireRole("Super Admin", "Admin"))

			r.Get("/", authHandler.ListUsers)
			r.Post("/", authHandler.CreateUser)
			r.Get("/deleted", authHandler.ListDeletedUsers)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", authHandler.GetUserByID)
				r.Put("/", authHandler.UpdateUser)
				r.Delete("/", authHandler.DeleteUser)
				r.Post("/toggle-status", authHandler.ToggleUserStatus)
				r.Post("/restore", authHandler.RestoreUser)
			})
		})
	})

	return mux
}
