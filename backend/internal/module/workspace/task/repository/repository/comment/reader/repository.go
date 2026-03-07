package reader

import (
	"github.com/jmoiron/sqlx"

	commentDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/comment"
)

var _ commentDef.TaskCommentReaderRepository = (*repository)(nil)

type repository struct {
	readPool *sqlx.DB
}

// NewRepository создаёт репозиторий чтения комментариев к задачам.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{readPool: db}
}
