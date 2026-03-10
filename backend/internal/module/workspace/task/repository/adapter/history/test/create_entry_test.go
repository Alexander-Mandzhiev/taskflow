package history_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

func (s *AdapterSuite) TestCreateHistoryEntry_Success() {
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")
	changedBy := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	entry := model.TaskHistory{
		TaskID:    taskID,
		ChangedBy: changedBy,
		FieldName: "title",
		OldValue:  "Old",
		NewValue:  "New",
		ChangedAt: time.Now(),
	}

	s.historyWriter.On("Create", mock.Anything, (*sqlx.Tx)(nil), entry).
		Return(nil).Once()

	err := s.repo.CreateHistoryEntry(s.ctx, nil, entry)

	assert.NoError(s.T(), err)
	s.historyWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreateHistoryEntry_WithTx() {
	tx := &sqlx.Tx{}
	entry := model.TaskHistory{
		TaskID:    uuid.New(),
		ChangedBy: uuid.New(),
		FieldName: "status",
		OldValue:  "todo",
		NewValue:  "done",
		ChangedAt: time.Now(),
	}

	s.historyWriter.On("Create", mock.Anything, tx, entry).
		Return(nil).Once()

	err := s.repo.CreateHistoryEntry(s.ctx, tx, entry)

	assert.NoError(s.T(), err)
	s.historyWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreateHistoryEntry_WriterError() {
	entry := model.TaskHistory{
		TaskID: uuid.New(), ChangedBy: uuid.New(), FieldName: "description",
		OldValue: "", NewValue: "Updated", ChangedAt: time.Now(),
	}

	s.historyWriter.On("Create", mock.Anything, mock.Anything, entry).
		Return(assert.AnError).Once()

	err := s.repo.CreateHistoryEntry(s.ctx, nil, entry)

	assert.Error(s.T(), err)
	s.historyWriter.AssertExpectations(s.T())
}
