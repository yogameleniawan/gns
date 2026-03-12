package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/base-go/backend/pkg/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{
		service: service,
	}
}

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	authResp, statusCode, err := h.service.Register(r.Context(), req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, authResp)
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	authResp, statusCode, err := h.service.Login(r.Context(), req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, authResp)
}

func (h Handler) LoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	var req GoogleOAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	authResp, statusCode, err := h.service.LoginWithGoogle(r.Context(), req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, authResp)
}

func (h Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	authResp, statusCode, err := h.service.RefreshToken(r.Context(), req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, authResp)
}

func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get user context
	userCtx, ok := r.Context().Value("user_context").(response.UserContext)
	if !ok {
		response.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	statusCode, err := h.service.Logout(r.Context(), userID)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, map[string]string{"message": "Logged out successfully"})
}

func (h Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user context
	userCtx, ok := r.Context().Value("user_context").(response.UserContext)
	if !ok {
		response.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	profile, statusCode, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, profile)
}

func (h Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Get user context
	userCtx, ok := r.Context().Value("user_context").(response.UserContext)
	if !ok {
		response.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, statusCode, err := h.service.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, user)
}

func (h Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get user context
	userCtx, ok := r.Context().Value("user_context").(response.UserContext)
	if !ok {
		response.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	statusCode, err := h.service.ChangePassword(r.Context(), userID, req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, map[string]string{"message": "Password changed successfully"})
}

// --- User Management (Admin) ---

func (h Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	var req UserListRequest

	// Parse query parameters
	req.Page = parseIntParam(r.URL.Query().Get("page"), 1)
	req.PageSize = parseIntParam(r.URL.Query().Get("page_size"), 10)
	req.Search = r.URL.Query().Get("search")
	req.SortBy = r.URL.Query().Get("sort_by")
	req.SortDir = r.URL.Query().Get("sort_dir")

	if isActive := r.URL.Query().Get("is_active"); isActive != "" {
		val := isActive == "true"
		req.IsActive = &val
	}

	if roleID := r.URL.Query().Get("role_id"); roleID != "" {
		val := parseIntParam(roleID, 0)
		if val > 0 {
			req.RoleID = &val
		}
	}

	result, statusCode, err := h.service.ListUsers(r.Context(), req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, result)
}

func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, statusCode, err := h.service.CreateUserByAdmin(r.Context(), req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, user)
}

func (h Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, statusCode, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, user)
}

func (h Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req AdminUpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, statusCode, err := h.service.UpdateUserByAdmin(r.Context(), id, req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, user)
}

func (h Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	statusCode, err := h.service.DeleteUserByAdmin(r.Context(), id)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, map[string]string{"message": "User deleted successfully"})
}

func (h Handler) ToggleUserStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, statusCode, err := h.service.ToggleUserStatus(r.Context(), id)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, user)
}

// Helper function
func parseIntParam(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			val = val*10 + int(c-'0')
		} else {
			return defaultVal
		}
	}
	return val
}

// --- Deleted Users (Restore) ---

func (h Handler) ListDeletedUsers(w http.ResponseWriter, r *http.Request) {
	var req UserListRequest

	// Parse query parameters
	req.Page = parseIntParam(r.URL.Query().Get("page"), 1)
	req.PageSize = parseIntParam(r.URL.Query().Get("page_size"), 10)
	req.Search = r.URL.Query().Get("search")
	req.SortBy = r.URL.Query().Get("sort_by")
	req.SortDir = r.URL.Query().Get("sort_dir")

	result, statusCode, err := h.service.ListDeletedUsers(r.Context(), req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, result)
}

func (h Handler) RestoreUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, statusCode, err := h.service.RestoreUserByAdmin(r.Context(), id)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, user)
}
