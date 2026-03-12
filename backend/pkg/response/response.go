package response

import (
	"encoding/json"
	"net/http"
)

type UserContext struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	MitraType string `json:"mitra_type,omitempty"`
}

type JSON struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Status  bool        `json:"status,omitempty"`
}

func ResponseJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to encode response"))
		}
	}
}

func ResponseError(w http.ResponseWriter, statusCode int, message string) {
	ResponseJSON(w, statusCode, map[string]string{"error": message})
}
