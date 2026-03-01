package service

import (
	"context"

	"github.com/google/uuid"
)

// Whoami возвращает userID по sessionID. При отсутствии или истечении сессии — model.ErrSessionNotFound.
func (s *accountService) Whoami(ctx context.Context, sessionID uuid.UUID) (userID uuid.UUID, err error) {
	session, err := s.sessionRepo.Get(ctx, sessionID)
	if err != nil {
		return uuid.Nil, err
	}
	return session.UserID, nil
}
