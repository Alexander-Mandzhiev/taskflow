package writer

import (
	"github.com/jmoiron/sqlx"

	commentDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/comment"
)

var _ commentDef.TaskCommentWriterRepository = (*repository)(nil)

type repository struct {
	writePool *sqlx.DB
}

// NewRepository создаёт репозиторий записи комментариев к задачам.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{writePool: db}
}
