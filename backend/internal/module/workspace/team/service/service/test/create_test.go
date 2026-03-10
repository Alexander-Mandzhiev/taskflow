package team_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestCreate_NilInput() {
	ownerID := uuid.New()

	got, err := s.svc.Create(s.ctx, model.TeamInput{}, ownerID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrNilInput)
	assert.Equal(s.T(), model.Team{}, got)
	s.teamRepo.AssertNotCalled(s.T(), "Create")
	s.memberRepo.AssertNotCalled(s.T(), "AddMember")
}

func (s *ServiceSuite) TestCreate_Success() {
	input := model.TeamInput{Name: "New Team"}
	ownerID := uuid.New()
	teamID := uuid.New()
	created := model.Team{
		ID:        teamID,
		Name:      input.Name,
		CreatedBy: ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	member := model.TeamMember{
		ID: uuid.New(), UserID: ownerID, TeamID: teamID, Role: model.RoleOwner, CreatedAt: time.Now(),
	}

	s.teamRepo.On("Create", mock.Anything, mock.Anything, input, ownerID).Return(created, nil).Once()
	s.memberRepo.On("AddMember", mock.Anything, mock.Anything, teamID, ownerID, model.RoleOwner).Return(member, nil).Once()

	got, err := s.svc.Create(s.ctx, input, ownerID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), teamID, got.ID)
	assert.Equal(s.T(), input.Name, got.Name)
	s.teamRepo.AssertExpectations(s.T())
	s.memberRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreate_RepoCreateError() {
	input := model.TeamInput{Name: "New Team"}
	ownerID := uuid.New()

	s.teamRepo.On("Create", mock.Anything, mock.Anything, input, ownerID).Return(model.Team{}, model.ErrInternal).Once()
	s.memberRepo.AssertNotCalled(s.T(), "AddMember")

	got, err := s.svc.Create(s.ctx, input, ownerID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrInternal)
	assert.Equal(s.T(), model.Team{}, got)
	s.teamRepo.AssertExpectations(s.T())
	s.memberRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreate_AddMemberError() {
	input := model.TeamInput{Name: "New Team"}
	ownerID := uuid.New()
	teamID := uuid.New()
	created := model.Team{ID: teamID, Name: input.Name, CreatedBy: ownerID, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	s.teamRepo.On("Create", mock.Anything, mock.Anything, input, ownerID).Return(created, nil).Once()
	s.memberRepo.On("AddMember", mock.Anything, mock.Anything, teamID, ownerID, model.RoleOwner).Return(model.TeamMember{}, assert.AnError).Once()

	got, err := s.svc.Create(s.ctx, input, ownerID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.Team{}, got)
	s.teamRepo.AssertExpectations(s.T())
	s.memberRepo.AssertExpectations(s.T())
}
