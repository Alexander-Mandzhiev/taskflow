package comment_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

func (s *AdapterSuite) TestCreateComment_Success() {
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	content := "New comment"
	created := &model.TaskComment{
		ID:        uuid.New(),
		TaskID:    taskID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.commentWriter.On("Create", mock.Anything, (*sqlx.Tx)(nil), taskID, userID, content).
		Return(created, nil).Once()

	got, err := s.repo.CreateComment(s.ctx, nil, taskID, userID, content)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Equal(s.T(), content, got.Content)
	assert.Equal(s.T(), taskID, got.TaskID)
	assert.Equal(s.T(), userID, got.UserID)
	s.commentWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreateComment_WithTx() {
	tx := &sqlx.Tx{}
	taskID := uuid.New()
	userID := uuid.New()
	content := "Comment in tx"
	created := &model.TaskComment{
		ID:        uuid.New(),
		TaskID:    taskID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.commentWriter.On("Create", mock.Anything, tx, taskID, userID, content).
		Return(created, nil).Once()

	got, err := s.repo.CreateComment(s.ctx, tx, taskID, userID, content)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Equal(s.T(), content, got.Content)
	s.commentWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreateComment_WriterError() {
	taskID := uuid.New()
	userID := uuid.New()
	content := "Will fail"

	s.commentWriter.On("Create", mock.Anything, mock.Anything, taskID, userID, content).
		Return((*model.TaskComment)(nil), assert.AnError).Once()

	got, err := s.repo.CreateComment(s.ctx, nil, taskID, userID, content)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.commentWriter.AssertExpectations(s.T())
}
