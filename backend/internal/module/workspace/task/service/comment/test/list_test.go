package comment_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestListByTaskID_Success() {
	taskID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	task := model.Task{ID: taskID, Title: "Task", TeamID: teamID}
	member := teamModel.TeamMember{UserID: userID, TeamID: teamID}
	comments := []model.TaskComment{
		{
			ID:        uuid.New(),
			TaskID:    taskID,
			UserID:    userID,
			Content:   "First comment",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.memberRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return(member, nil).Once()
	s.commentRepo.On("ListCommentsByTaskID", mock.Anything, mock.Anything, taskID).Return(comments, nil).Once()

	got, err := s.svc.ListByTaskID(s.ctx, taskID, userID)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), got, 1)
	assert.Equal(s.T(), "First comment", got[0].Content)
	s.taskRepo.AssertExpectations(s.T())
	s.memberRepo.AssertExpectations(s.T())
	s.commentRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListByTaskID_TaskNotFound() {
	taskID := uuid.New()
	userID := uuid.New()

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).
		Return(model.Task{}, model.ErrTaskNotFound).Once()

	got, err := s.svc.ListByTaskID(s.ctx, taskID, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.memberRepo.AssertNotCalled(s.T(), "GetMember")
	s.commentRepo.AssertNotCalled(s.T(), "ListCommentsByTaskID")
}

func (s *ServiceSuite) TestListByTaskID_NotMember() {
	taskID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	task := model.Task{ID: taskID, Title: "Task", TeamID: teamID}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.memberRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).
		Return(teamModel.TeamMember{}, teamModel.ErrMemberNotFound).Once()

	got, err := s.svc.ListByTaskID(s.ctx, taskID, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.memberRepo.AssertExpectations(s.T())
	s.commentRepo.AssertNotCalled(s.T(), "ListCommentsByTaskID")
}

func (s *ServiceSuite) TestListByTaskID_CommentRepoError() {
	taskID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	task := model.Task{ID: taskID, Title: "Task", TeamID: teamID}
	member := teamModel.TeamMember{UserID: userID, TeamID: teamID}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.memberRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return(member, nil).Once()
	s.commentRepo.On("ListCommentsByTaskID", mock.Anything, mock.Anything, taskID).
		Return(nil, assert.AnError).Once()

	got, err := s.svc.ListByTaskID(s.ctx, taskID, userID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.memberRepo.AssertExpectations(s.T())
	s.commentRepo.AssertExpectations(s.T())
}
