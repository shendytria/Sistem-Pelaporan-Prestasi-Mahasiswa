package repository

import (
	"context"
	"fmt"
	"prestasi_mhs/app/model"
	"prestasi_mhs/database"
)

type ReportRepository struct{}

func NewReportRepository() *ReportRepository {
	return &ReportRepository{}
}

func (r *ReportRepository) GetStatistics(ctx context.Context, filterQuery string, args ...interface{}) (*model.StatisticResponse, error) {
	var res model.StatisticResponse

	const qStatus = `
		SELECT status, COUNT(*)
		FROM achievement_references
		%s
		GROUP BY status
	`

	rows, err := database.PG.Query(ctx,
		fmt.Sprintf(qStatus, filterQuery),
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}

		if status == "verified" {
			res.Verified = count
		} else if status == "draft" {
			res.Draft = count
		} else if status == "submitted" {
			res.Submitted = count
		} else if status == "rejected" {
			res.Rejected = count
		}
		res.TotalAchievements += count
	}

	return &res, nil
}

func (r *ReportRepository) GetStudentReport(ctx context.Context, studentID string) (*model.StudentReportResponse, error) {
	var res model.StudentReportResponse
	res.StudentID = studentID

	const qStatus = `
		SELECT status, COUNT(*)
		FROM achievement_references
		WHERE student_id = $1
		GROUP BY status
	`

	rows, err := database.PG.Query(ctx, qStatus, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}

		if status == "verified" {
			res.Verified = count
		} else if status == "draft" {
			res.Draft = count
		} else if status == "submitted" {
			res.Submitted = count
		} else if status == "rejected" {
			res.Rejected = count
		}
		res.TotalAchievements += count
	}

	return &res, nil
}
