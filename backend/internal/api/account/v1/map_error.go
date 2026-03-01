package account_v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	accountmodel "mkk/internal/module/identity/account/model"
	"mkk/internal/module/identity/user/model"
	pkghttp "mkk/pkg/http"
	"mkk/pkg/http/middleware"
	"mkk/pkg/logger"
	"mkk/pkg/metadata"
)

var validate = validator.New()

// mapError маппит все ошибки ручки (JSON-декодирование, валидация, домен, сессия) в HTTP-ответ.
// Использование: mapError(w, r, err); return
func mapError(w http.ResponseWriter, r *http.Request, err error) {
	if middleware.IsBodyError(w, err) {
		return
	}

	// Ошибки декодирования body (JSON) → 400
	var jsonErr *json.SyntaxError
	var jsonTypeErr *json.UnmarshalTypeError
	if errors.As(err, &jsonErr) || errors.As(err, &jsonTypeErr) {
		pkghttp.WriteJSON(w, http.StatusBadRequest, pkghttp.ErrorBody{
			Code:    http.StatusBadRequest,
			Message: "Некорректное тело запроса",
		})
		return
	}

	// Ошибки валидации запроса (validator) → 400 с деталями по полям
	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		items := make([]pkghttp.ValidationErrorItem, 0, len(valErrs))
		for _, e := range valErrs {
			items = append(items, pkghttp.ValidationErrorItem{
				Field:   e.Field(),
				Message: validationMessage(e),
			})
		}
		pkghttp.WriteJSON(w, http.StatusBadRequest, pkghttp.ValidationErrorBody{
			Code:    http.StatusBadRequest,
			Message: "Ошибка валидации запроса",
			Errors:  items,
		})
		return
	}

	code, message := mapDomainError(err)
	if code == http.StatusInternalServerError {
		logger.Error(r.Context(), "[Account API] unhandled error",
			zap.String("path", r.URL.Path),
			zap.Error(err),
		)
	}
	pkghttp.WriteJSON(w, code, pkghttp.ErrorBody{Code: code, Message: message})
}

// mapDomainError возвращает HTTP-код и сообщение для доменных ошибок.
func mapDomainError(err error) (int, string) {
	switch {
	case errors.Is(err, accountmodel.ErrInvalidCredentials):
		return http.StatusUnauthorized, "Неверные учётные данные"
	case errors.Is(err, model.ErrUserNotFound):
		return http.StatusNotFound, "Пользователь не найден"
	case errors.Is(err, model.ErrEmailDuplicate):
		return http.StatusConflict, "Пользователь с таким email уже существует"
	case errors.Is(err, model.ErrNilInput):
		return http.StatusBadRequest, "Некорректные данные запроса"
	case errors.Is(err, accountmodel.ErrSessionNotFound), errors.Is(err, metadata.ErrNotFound):
		return http.StatusUnauthorized, "Сессия не найдена или истекла"
	default:
		return http.StatusInternalServerError, "Внутренняя ошибка сервера"
	}
}

// validationMessage возвращает человекочитаемое сообщение для правила валидации.
func validationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "поле обязательно"
	case "email":
		return "некорректный формат email"
	case "min":
		return "значение слишком короткое (минимум " + e.Param() + ")"
	case "max":
		return "значение слишком длинное (максимум " + e.Param() + ")"
	default:
		return "некорректное значение"
	}
}
