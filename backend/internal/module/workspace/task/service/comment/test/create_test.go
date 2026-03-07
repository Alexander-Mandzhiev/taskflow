package comment_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestCreate_Success() {
	taskID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	content := "New comment"
	task := &model.Task{ID: taskID, Title: "Task", TeamID: teamID}
	member := &teamModel.TeamMember{UserID: userID, TeamID: teamID}
	created := &model.TaskComment{
		ID:        uuid.New(),
		TaskID:    taskID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return(member, nil).Once()
	s.commentRepo.On("CreateComment", mock.Anything, mock.Anything, taskID, userID, content).
		Return(created, nil).Once()

	got, err := s.svc.Create(s.ctx, taskID, userID, content)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Equal(s.T(), content, got.Content)
	assert.Equal(s.T(), taskID, got.TaskID)
	assert.Equal(s.T(), userID, got.UserID)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertExpectations(s.T())
	s.commentRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreate_TaskNotFound() {
	taskID := uuid.New()
	userID := uuid.New()
	content := "Comment"

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).
		Return((*model.Task)(nil), model.ErrTaskNotFound).Once()

	got, err := s.svc.Create(s.ctx, taskID, userID, content)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertNotCalled(s.T(), "GetMember")
	s.commentRepo.AssertNotCalled(s.T(), "CreateComment")
}

func (s *ServiceSuite) TestCreate_NotMember() {
	taskID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	task := &model.Task{ID: taskID, Title: "Task", TeamID: teamID}
	content := "Comment"

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).
		Return((*teamModel.TeamMember)(nil), teamModel.ErrMemberNotFound).Once()

	got, err := s.svc.Create(s.ctx, taskID, userID, content)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertExpectations(s.T())
	s.commentRepo.AssertNotCalled(s.T(), "CreateComment")
}

func (s *ServiceSuite) TestCreate_CommentRepoError() {
	taskID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	content := "Comment"
	task := &model.Task{ID: taskID, Title: "Task", TeamID: teamID}
	member := &teamModel.TeamMember{UserID: userID, TeamID: teamID}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return(member, nil).Once()
	s.commentRepo.On("CreateComment", mock.Anything, mock.Anything, taskID, userID, content).
		Return((*model.TaskComment)(nil), assert.AnError).Once()

	got, err := s.svc.Create(s.ctx, taskID, userID, content)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertExpectations(s.T())
	s.commentRepo.AssertExpectations(s.T())
}
