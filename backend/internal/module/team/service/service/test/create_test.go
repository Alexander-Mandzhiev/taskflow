package service_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

func (s *ServiceSuite) TestCreate_NilInput() {
	ownerID := uuid.New()

	got, err := s.svc.Create(s.ctx, nil, ownerID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrNilInput)
	assert.Nil(s.T(), got)
	s.teamRepo.AssertNotCalled(s.T(), "Create")
}

func (s *ServiceSuite) TestCreate_Success() {
	input := &model.TeamInput{Name: "New Team"}
	ownerID := uuid.New()
	teamID := uuid.New()
	created := &model.Team{
		ID:        teamID,
		Name:      input.Name,
		CreatedBy: ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	member := &model.TeamMember{
		ID: uuid.New(), UserID: ownerID, TeamID: teamID, Role: model.RoleOwner, CreatedAt: time.Now(),
	}

	s.teamRepo.On("Create", mock.Anything, mock.Anything, input, ownerID).Return(created, nil).Once()
	s.teamRepo.On("AddMember", mock.Anything, mock.Anything, teamID, ownerID, model.RoleOwner).Return(member, nil).Once()

	got, err := s.svc.Create(s.ctx, input, ownerID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Equal(s.T(), teamID, got.ID)
	assert.Equal(s.T(), input.Name, got.Name)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreate_RepoCreateError() {
	input := &model.TeamInput{Name: "New Team"}
	ownerID := uuid.New()

	s.teamRepo.On("Create", mock.Anything, mock.Anything, input, ownerID).Return((*model.Team)(nil), model.ErrInternal).Once()

	got, err := s.svc.Create(s.ctx, input, ownerID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrInternal)
	assert.Nil(s.T(), got)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreate_AddMemberError() {
	input := &model.TeamInput{Name: "New Team"}
	ownerID := uuid.New()
	teamID := uuid.New()
	created := &model.Team{ID: teamID, Name: input.Name, CreatedBy: ownerID, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	s.teamRepo.On("Create", mock.Anything, mock.Anything, input, ownerID).Return(created, nil).Once()
	s.teamRepo.On("AddMember", mock.Anything, mock.Anything, teamID, ownerID, model.RoleOwner).Return((*model.TeamMember)(nil), assert.AnError).Once()

	got, err := s.svc.Create(s.ctx, input, ownerID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamRepo.AssertExpectations(s.T())
}
