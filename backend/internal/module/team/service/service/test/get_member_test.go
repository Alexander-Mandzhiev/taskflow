package service_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

func (s *ServiceSuite) TestGetMember_Success() {
	teamID := uuid.New()
	userID := uuid.New()
	want := &model.TeamMember{
		ID:        uuid.New(),
		UserID:    userID,
		TeamID:    teamID,
		Role:      model.RoleAdmin,
		CreatedAt: time.Now(),
	}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return(want, nil).Once()

	got, err := s.svc.GetMember(s.ctx, teamID, userID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), want, got)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetMember_NotFound() {
	teamID := uuid.New()
	userID := uuid.New()

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return((*model.TeamMember)(nil), model.ErrMemberNotFound).Once()

	got, err := s.svc.GetMember(s.ctx, teamID, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrMemberNotFound)
	assert.Nil(s.T(), got)
	s.teamRepo.AssertExpectations(s.T())
}
