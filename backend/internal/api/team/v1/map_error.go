package team_v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http/middleware"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

// mapError маппит ошибки ручки в HTTP-ответ.
func mapError(w http.ResponseWriter, r *http.Request, err error) {
	if middleware.IsBodyError(r.Context(), w, err) {
		return
	}

	var jsonErr *json.SyntaxError
	var jsonTypeErr *json.UnmarshalTypeError
	if errors.As(err, &jsonErr) || errors.As(err, &jsonTypeErr) {
		pkghttp.WriteJSON(r.Context(), w, http.StatusBadRequest, pkghttp.ErrorBody{
			Code:    http.StatusBadRequest,
			Message: "Некорректное тело запроса",
		})
		return
	}

	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		items := make([]pkghttp.ValidationErrorItem, 0, len(valErrs))
		for _, e := range valErrs {
			items = append(items, pkghttp.ValidationErrorItem{
				Field:   e.Field(),
				Message: validationMessage(e),
			})
		}
		pkghttp.WriteJSON(r.Context(), w, http.StatusBadRequest, pkghttp.ValidationErrorBody{
			Code:    http.StatusBadRequest,
			Message: "Ошибка валидации запроса",
			Errors:  items,
		})
		return
	}

	code, message := mapDomainError(err)
	if code == http.StatusInternalServerError {
		logger.Error(r.Context(), "[Team API] unhandled error",
			zap.String("path", r.URL.Path),
			zap.Error(err),
		)
	}

	pkghttp.WriteJSON(r.Context(), w, code, pkghttp.ErrorBody{Code: code, Message: message})
}

func mapDomainError(err error) (int, string) {
	switch {
	case errors.Is(err, metadata.ErrNotFound):
		return http.StatusUnauthorized, "Необходима аутентификация"
	case errors.Is(err, model.ErrTeamNotFound):
		return http.StatusNotFound, "Команда не найдена"
	case errors.Is(err, model.ErrMemberNotFound):
		return http.StatusNotFound, "Участник не найден"
	case errors.Is(err, model.ErrForbidden):
		return http.StatusForbidden, "Недостаточно прав"
	case errors.Is(err, model.ErrAlreadyMember):
		return http.StatusConflict, "Пользователь уже в команде"
	case errors.Is(err, model.ErrAlreadyInvited):
		return http.StatusConflict, "Приглашение на этот email уже отправлено"
	case errors.Is(err, model.ErrUserNotFound):
		return http.StatusNotFound, "Пользователь с указанным email не найден"
	case errors.Is(err, model.ErrNilInput):
		return http.StatusBadRequest, "Некорректные данные запроса"
	case errors.Is(err, model.ErrInvalidID):
		return http.StatusBadRequest, "Некорректный идентификатор команды"
	case errors.Is(err, model.ErrInvitationNotFound):
		return http.StatusNotFound, "Приглашение не найдено"
	case errors.Is(err, model.ErrInvalidRole):
		return http.StatusBadRequest, "Недопустимая роль приглашения (допустимы: member, admin)"
	case errors.Is(err, model.ErrTemporaryFailure):
		return http.StatusServiceUnavailable, "Временная ошибка, попробуйте позже"
	default:
		return http.StatusInternalServerError, "Внутренняя ошибка сервера"
	}
}

func validationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "поле обязательно"
	case "uuid":
		return "некорректный формат UUID"
	case "email":
		return "некорректный формат email"
	case "oneof":
		return "допустимые значения: " + e.Param()
	case "max":
		return "значение слишком длинное (максимум " + e.Param() + ")"
	default:
		return "некорректное значение"
	}
}
