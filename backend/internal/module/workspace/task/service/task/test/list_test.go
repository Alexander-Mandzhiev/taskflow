package task_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestList_NilFilter() {
	userID := uuid.New()

	got, total, err := s.svc.List(s.ctx, userID, nil)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrPaginationRequired)
	assert.Nil(s.T(), got)
	assert.Equal(s.T(), 0, total)
	s.teamSvc.AssertNotCalled(s.T(), "GetMember")
	s.taskRepo.AssertNotCalled(s.T(), "List")
}

func (s *ServiceSuite) TestList_NoPagination() {
	userID := uuid.New()
	teamID := uuid.New()
	filter := &model.TaskListFilter{TeamID: &teamID, Limit: 0, Offset: 0}

	got, total, err := s.svc.List(s.ctx, userID, filter)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrPaginationRequired)
	assert.Nil(s.T(), got)
	assert.Equal(s.T(), 0, total)
	s.teamSvc.AssertNotCalled(s.T(), "GetMember")
	s.taskRepo.AssertNotCalled(s.T(), "List")
}

func (s *ServiceSuite) TestList_NoTeamID() {
	userID := uuid.New()
	filter := &model.TaskListFilter{Limit: 10, Offset: 0}

	got, total, err := s.svc.List(s.ctx, userID, filter)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrForbidden)
	assert.Nil(s.T(), got)
	assert.Equal(s.T(), 0, total)
	s.teamSvc.AssertNotCalled(s.T(), "GetMember")
	s.taskRepo.AssertNotCalled(s.T(), "List")
}

func (s *ServiceSuite) TestList_NotMember() {
	userID := uuid.New()
	teamID := uuid.New()
	filter := &model.TaskListFilter{TeamID: &teamID, Limit: 10, Offset: 0}

	s.teamSvc.On("GetMember", mock.Anything, teamID, userID).
		Return((*teamModel.TeamMember)(nil), teamModel.ErrMemberNotFound).Once()

	got, total, err := s.svc.List(s.ctx, userID, filter)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Nil(s.T(), got)
	assert.Equal(s.T(), 0, total)
	s.teamSvc.AssertExpectations(s.T())
	s.taskRepo.AssertNotCalled(s.T(), "List")
}

func (s *ServiceSuite) TestList_Success() {
	userID := uuid.New()
	teamID := uuid.New()
	filter := &model.TaskListFilter{TeamID: &teamID, Limit: 10, Offset: 0}
	member := &teamModel.TeamMember{UserID: userID, TeamID: teamID}
	tasks := []*model.Task{
		{ID: uuid.New(), Title: "Task 1", TeamID: teamID},
		{ID: uuid.New(), Title: "Task 2", TeamID: teamID},
	}

	s.teamSvc.On("GetMember", mock.Anything, teamID, userID).Return(member, nil).Once()
	s.taskRepo.On("List", mock.Anything, mock.Anything, filter).Return(tasks, 2, nil).Once()

	got, total, err := s.svc.List(s.ctx, userID, filter)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), got, 2)
	assert.Equal(s.T(), 2, total)
	s.teamSvc.AssertExpectations(s.T())
	s.taskRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestList_RepoError() {
	userID := uuid.New()
	teamID := uuid.New()
	filter := &model.TaskListFilter{TeamID: &teamID, Limit: 10, Offset: 0}
	member := &teamModel.TeamMember{UserID: userID, TeamID: teamID}

	s.teamSvc.On("GetMember", mock.Anything, teamID, userID).Return(member, nil).Once()
	s.taskRepo.On("List", mock.Anything, mock.Anything, filter).
		Return(([]*model.Task)(nil), 0, assert.AnError).Once()

	got, total, err := s.svc.List(s.ctx, userID, filter)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	assert.Equal(s.T(), 0, total)
	s.teamSvc.AssertExpectations(s.T())
	s.taskRepo.AssertExpectations(s.T())
}
