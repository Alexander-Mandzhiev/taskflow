package comment_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

func (s *AdapterSuite) TestListCommentsByTaskID_Success() {
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	comments := []*model.TaskComment{
		{
			ID:        uuid.New(),
			TaskID:    taskID,
			UserID:    userID,
			Content:   "First comment",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	s.commentReader.On("ListByTaskID", mock.Anything, (*sqlx.Tx)(nil), taskID).
		Return(comments, nil).Once()

	got, err := s.repo.ListCommentsByTaskID(s.ctx, nil, taskID)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), got, 1)
	assert.Equal(s.T(), "First comment", got[0].Content)
	assert.Equal(s.T(), taskID, got[0].TaskID)
	assert.Equal(s.T(), userID, got[0].UserID)
	s.commentReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestListCommentsByTaskID_Empty() {
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")

	s.commentReader.On("ListByTaskID", mock.Anything, mock.Anything, taskID).
		Return([]*model.TaskComment{}, nil).Once()

	got, err := s.repo.ListCommentsByTaskID(s.ctx, nil, taskID)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), got)
	s.commentReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestListCommentsByTaskID_WithTx_ReaderError() {
	tx := &sqlx.Tx{}
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")

	s.commentReader.On("ListByTaskID", mock.Anything, tx, taskID).
		Return(([]*model.TaskComment)(nil), assert.AnError).Once()

	got, err := s.repo.ListCommentsByTaskID(s.ctx, tx, taskID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.commentReader.AssertExpectations(s.T())
}
