package rbac

import (
	"encoding/json"
	"net/http"
	"strconv"

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

// --- Role Handlers ---

func (h Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, statusCode, err := h.service.CreateRole(r.Context(), req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, role)
}

func (h Handler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	roles, statusCode, err := h.service.GetAllRoles(r.Context())
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, roles)
}

func (h Handler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, statusCode, err := h.service.GetRoleByID(r.Context(), id)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, role)
}

func (h Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, statusCode, err := h.service.UpdateRole(r.Context(), id, req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, role)
}

func (h Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	statusCode, err := h.service.DeleteRole(r.Context(), id)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, map[string]string{"message": "Role deleted successfully"})
}

// --- Permission Handlers ---

func (h Handler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var req CreatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	permission, statusCode, err := h.service.CreatePermission(r.Context(), req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, permission)
}

func (h Handler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	permissions, statusCode, err := h.service.GetAllPermissions(r.Context())
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, permissions)
}

func (h Handler) GetPermissionsByModule(w http.ResponseWriter, r *http.Request) {
	permissions, statusCode, err := h.service.GetPermissionsByModule(r.Context())
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, permissions)
}

func (h Handler) GetPermissionByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	permission, statusCode, err := h.service.GetPermissionByID(r.Context(), id)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, permission)
}

func (h Handler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	var req UpdatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	permission, statusCode, err := h.service.UpdatePermission(r.Context(), id, req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, permission)
}

func (h Handler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	statusCode, err := h.service.DeletePermission(r.Context(), id)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, map[string]string{"message": "Permission deleted successfully"})
}

// --- Role-Permission Handlers ---

func (h Handler) AssignPermissionsToRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var req AssignPermissionsToRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	statusCode, err := h.service.AssignPermissionsToRole(r.Context(), id, req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, map[string]string{"message": "Permissions assigned successfully"})
}

func (h Handler) GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	permissions, statusCode, err := h.service.GetRolePermissions(r.Context(), id)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, permissions)
}

// --- User-Role Handlers ---

func (h Handler) AssignRolesToUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req AssignRolesToUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	statusCode, err := h.service.AssignRolesToUser(r.Context(), userID, req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, map[string]string{"message": "Roles assigned successfully"})
}

func (h Handler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	roles, statusCode, err := h.service.GetUserRoles(r.Context(), userID)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, roles)
}

// --- Module Access Handlers ---

func (h Handler) UpdateModuleAccess(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var req UpdateModuleAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	moduleAccess, statusCode, err := h.service.UpdateModuleAccess(r.Context(), id, req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, moduleAccess)
}

func (h Handler) GetModuleAccessByRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	moduleAccess, statusCode, err := h.service.GetModuleAccessByRole(r.Context(), id)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, moduleAccess)
}

// --- Permission Checking Handlers ---

func (h Handler) CheckPermission(w http.ResponseWriter, r *http.Request) {
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

	var req CheckPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	hasPermission, statusCode, err := h.service.CheckPermission(r.Context(), userID, req.Permission)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, CheckPermissionResponse{HasPermission: hasPermission})
}

func (h Handler) CheckModuleAccess(w http.ResponseWriter, r *http.Request) {
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

	var req CheckModuleAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	hasAccess, statusCode, err := h.service.CheckModuleAccess(r.Context(), userID, req.Module, req.Action)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, CheckModuleAccessResponse{HasAccess: hasAccess})
}
