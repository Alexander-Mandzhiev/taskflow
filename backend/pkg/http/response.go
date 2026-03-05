package http

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// WriteJSON пишет ответ с кодом statusCode и телом body в формате JSON.
func WriteJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if body != nil {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			logger.Error(context.Background(), "[WriteJSON] failed to encode response", zap.Error(err))
		}
	}
}

// ErrorBody — тело ответа с ошибкой в формате JSON (code + message).
type ErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ValidationErrorItem — ошибка валидации по одному полю (для 400 при ошибках валидации запроса).
type ValidationErrorItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorBody — тело ответа при ошибках валидации (400).
type ValidationErrorBody struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Errors  []ValidationErrorItem `json:"errors"`
}
