package team_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestAcceptInvitation_NotImplemented() {
	token := "some-token"
	userID := uuid.New()

	got, err := s.svc.AcceptInvitation(s.ctx, token, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrNotImplemented)
	assert.Equal(s.T(), model.TeamMember{}, got)
}
