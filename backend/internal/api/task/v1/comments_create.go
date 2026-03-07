package task_v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/task/v1/dto"
	taskconverter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/util"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/validation"
)

// CreateComment создаёт комментарий к задаче. POST /api/v1/tasks/{id}/comments.
func (api *API) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, err := metadata.UserID(r.Context())
	if err != nil {
		mapError(w, r, err)
		return
	}
	taskID, err := util.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		mapError(w, r, model.ErrTaskNotFound)
		return
	}

	var req dto.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapError(w, r, err)
		return
	}

	if err := validation.Validator.Struct(req); err != nil {
		mapError(w, r, err)
		return
	}

	comment, err := api.commentService.Create(r.Context(), taskID, userID, req.Content)
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.WriteJSON(r.Context(), w, http.StatusCreated, taskconverter.CommentToResponse(comment))
}
