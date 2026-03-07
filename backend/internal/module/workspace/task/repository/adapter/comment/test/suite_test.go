package comment_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	commentAdapter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/adapter/comment"
	commentMocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/comment/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

type AdapterSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	commentReader *commentMocks.TaskCommentReaderRepository
	commentWriter *commentMocks.TaskCommentWriterRepository
	repo          repository.TaskCommentRepository
}

func (s *AdapterSuite) SetupTest() {
	s.ctx = context.Background()
	if err := logger.InitDefault(); err != nil {
		panic(err)
	}
	s.commentReader = commentMocks.NewTaskCommentReaderRepository(s.T())
	s.commentWriter = commentMocks.NewTaskCommentWriterRepository(s.T())
	s.repo = commentAdapter.NewAdapter(s.commentReader, s.commentWriter)
}

func TestAdapter(t *testing.T) {
	suite.Run(t, new(AdapterSuite))
}
