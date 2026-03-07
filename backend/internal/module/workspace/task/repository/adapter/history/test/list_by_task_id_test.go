package history_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

func (s *AdapterSuite) TestListHistoryByTaskID_Success() {
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")
	changedBy := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	entries := []*model.TaskHistory{
		{
			ID:        uuid.New(),
			TaskID:    taskID,
			ChangedBy: changedBy,
			FieldName: "title",
			OldValue:  "Old",
			NewValue:  "New",
			ChangedAt: time.Now(),
		},
	}

	s.historyReader.On("ListByTaskID", mock.Anything, (*sqlx.Tx)(nil), taskID).
		Return(entries, nil).Once()

	got, err := s.repo.ListHistoryByTaskID(s.ctx, nil, taskID)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), got, 1)
	assert.Equal(s.T(), "title", got[0].FieldName)
	assert.Equal(s.T(), "Old", got[0].OldValue)
	assert.Equal(s.T(), "New", got[0].NewValue)
	s.historyReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestListHistoryByTaskID_Empty() {
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")

	s.historyReader.On("ListByTaskID", mock.Anything, mock.Anything, taskID).
		Return([]*model.TaskHistory{}, nil).Once()

	got, err := s.repo.ListHistoryByTaskID(s.ctx, nil, taskID)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), got)
	s.historyReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestListHistoryByTaskID_WithTx_ReaderError() {
	tx := &sqlx.Tx{}
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")

	s.historyReader.On("ListByTaskID", mock.Anything, tx, taskID).
		Return(([]*model.TaskHistory)(nil), assert.AnError).Once()

	got, err := s.repo.ListHistoryByTaskID(s.ctx, tx, taskID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.historyReader.AssertExpectations(s.T())
}
