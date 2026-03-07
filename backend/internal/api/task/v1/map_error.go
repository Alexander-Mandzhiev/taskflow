package task_v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
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
		logger.Error(r.Context(), "[Task API] unhandled error",
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
	case errors.Is(err, strconv.ErrSyntax):
		return http.StatusBadRequest, "Некорректное числовое значение параметра"
	case errors.Is(err, model.ErrInvalidAssigneeID):
		return http.StatusBadRequest, "Некорректный формат assignee_id (ожидается UUID)"
	case errors.Is(err, model.ErrTaskNotFound):
		return http.StatusNotFound, "Задача не найдена"
	case errors.Is(err, model.ErrForbidden):
		return http.StatusForbidden, "Недостаточно прав"
	case errors.Is(err, model.ErrNilInput):
		return http.StatusBadRequest, "Некорректные данные запроса"
	case errors.Is(err, model.ErrInvalidStatus):
		return http.StatusBadRequest, "Недопустимый статус задачи (допустимы: todo, in_progress, done)"
	case errors.Is(err, model.ErrAssigneeNotInTeam):
		return http.StatusBadRequest, "Исполнитель не является участником команды задачи"
	case errors.Is(err, model.ErrTeamIDRequired):
		return http.StatusBadRequest, "Укажите team_id в параметрах запроса"
	case errors.Is(err, model.ErrPaginationRequired):
		return http.StatusBadRequest, "Укажите limit > 0 для пагинации"
	case errors.Is(err, model.ErrInvalidLimit):
		return http.StatusBadRequest, "Параметр limit должен быть положительным"
	case errors.Is(err, model.ErrTemporaryFailure):
		return http.StatusServiceUnavailable, "Временная ошибка, попробуйте позже"
	case errors.Is(err, model.ErrCommentNotImplemented):
		return http.StatusNotImplemented, "Сервис комментариев пока не реализован"
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
	case "oneof":
		return "допустимые значения: " + e.Param()
	case "max":
		return "значение слишком длинное (максимум " + e.Param() + ")"
	default:
		return "некорректное значение"
	}
}
