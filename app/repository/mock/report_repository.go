package mock

import (
	"context"
	"prestasi_mhs/app/model"
)

type ReportRepoMock struct {
	MockStatistics    *model.StatisticResponse
	MockStudentReport *model.StudentReportResponse
}

func NewReportRepoMock() *ReportRepoMock {
	return &ReportRepoMock{}
}

func (m *ReportRepoMock) GetStatistics(ctx context.Context, filter string) (*model.StatisticResponse, error) {
	return m.MockStatistics, nil
}

func (m *ReportRepoMock) GetStudentReport(ctx context.Context, studentID string) (*model.StudentReportResponse, error) {
	return m.MockStudentReport, nil
}
