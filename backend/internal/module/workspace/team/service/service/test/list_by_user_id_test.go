package team_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestListByUserID_Success() {
	userID := uuid.New()
	want := []model.TeamWithRole{
		{Team: model.Team{ID: uuid.New(), Name: "Team A"}, Role: model.RoleOwner},
		{Team: model.Team{ID: uuid.New(), Name: "Team B"}, Role: model.RoleMember},
	}

	s.teamRepo.On("ListByUserID", mock.Anything, mock.Anything, userID).Return(want, nil).Once()

	got, err := s.svc.ListByUserID(s.ctx, userID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), want, got)
	assert.Len(s.T(), got, 2)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListByUserID_Empty() {
	userID := uuid.New()

	s.teamRepo.On("ListByUserID", mock.Anything, mock.Anything, userID).Return([]model.TeamWithRole{}, nil).Once()

	got, err := s.svc.ListByUserID(s.ctx, userID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Empty(s.T(), got)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListByUserID_RepoError() {
	userID := uuid.New()

	s.teamRepo.On("ListByUserID", mock.Anything, mock.Anything, userID).Return(nil, assert.AnError).Once()

	got, err := s.svc.ListByUserID(s.ctx, userID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamRepo.AssertExpectations(s.T())
}
