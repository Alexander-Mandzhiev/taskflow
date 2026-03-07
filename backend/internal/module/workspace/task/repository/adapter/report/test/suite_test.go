package report_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	reportAdapter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/adapter/report"
	reportMocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/report/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

type AdapterSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	reportReader *reportMocks.ReportReaderRepository
	repo         repository.ReportRepository
}

func (s *AdapterSuite) SetupTest() {
	s.ctx = context.Background()
	if err := logger.InitDefault(); err != nil {
		panic(err)
	}
	s.reportReader = reportMocks.NewReportReaderRepository(s.T())
	s.repo = reportAdapter.NewAdapter(s.reportReader)
}

func TestAdapter(t *testing.T) {
	suite.Run(t, new(AdapterSuite))
}
