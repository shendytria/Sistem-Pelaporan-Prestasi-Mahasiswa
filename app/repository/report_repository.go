package repository

import (
	"context"
	"fmt"
	"sort"

	"go.mongodb.org/mongo-driver/bson"
	"prestasi_mhs/app/model"
	"prestasi_mhs/database"
)

type ReportRepository struct{}

func NewReportRepository() *ReportRepository { return &ReportRepository{} }

func (r *ReportRepository) GetStatistics(ctx context.Context, filterQuery string, args ...interface{}) (*model.StatisticResponse, error) {
	var res model.StatisticResponse

	const qStatus = `
		SELECT status, COUNT(*)
		FROM achievement_references
		%s
		GROUP BY status
	`
	rows1, err := database.PG.Query(ctx, fmt.Sprintf(qStatus, filterQuery), args...)
	if err != nil {
		return nil, err
	}
	defer rows1.Close()

	for rows1.Next() {
		var status string
		var count int
		rows1.Scan(&status, &count)
		res.TotalAchievements += count
		if status == "verified" {
			res.Verified = count
		} else if status == "draft" {
			res.Draft = count
		} else if status == "submitted" {
			res.Submitted = count
		} else if status == "rejected" {
			res.Rejected = count
		}
	}

	studentIDs := []string{}
	if filterQuery == "" {
		rows, _ := database.PG.Query(ctx, `SELECT DISTINCT student_id FROM achievement_references`)
		for rows.Next() {
			var sid string
			rows.Scan(&sid)
			studentIDs = append(studentIDs, sid)
		}
	} else {
		rows, _ := database.PG.Query(ctx, `SELECT DISTINCT student_id FROM achievement_references `+filterQuery, args...)
		for rows.Next() {
			var sid string
			rows.Scan(&sid)
			studentIDs = append(studentIDs, sid)
		}
	}

	collection := database.Mongo.Collection("achievements")

	filter := bson.M{"studentId": bson.M{"$in": studentIDs}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type periodCount struct {
		period string
		count  int
	}
	periodMap := map[string]int{}
	typeMap := map[string]int{}
	levelMap := map[string]int{}
	scoreMap := map[string]struct {
		points int
		count  int
	}{}

	for cursor.Next(ctx) {
		var a model.Achievement
		cursor.Decode(&a)

		typeMap[a.AchievementType]++

		period := a.Details.EventDate.Format("2006-01")
		periodMap[period]++

		levelMap[a.Details.CompetitionLevel]++

		st := scoreMap[a.StudentID]
		st.points += int(a.Points)
		st.count++
		scoreMap[a.StudentID] = st
	}

	for t, c := range typeMap {
		res.ByType = append(res.ByType, model.StatisticByType{Type: t, Count: c})
	}
	for p, c := range periodMap {
		res.ByPeriod = append(res.ByPeriod, model.StatisticByPeriod{Period: p, Count: c})
	}
	for lv, c := range levelMap {
		res.ByLevel = append(res.ByLevel, model.StatisticByLevel{Level: lv, Count: c})
	}
	for sid, val := range scoreMap {
		res.TopStudents = append(res.TopStudents, model.TopStudent{
			StudentID: sid,
			Count:     val.count,
			Points:    val.points,
		})
	}

	sort.Slice(res.TopStudents, func(i, j int) bool {
		return res.TopStudents[i].Points > res.TopStudents[j].Points
	})
	if len(res.TopStudents) > 5 {
		res.TopStudents = res.TopStudents[:5]
	}

	return &res, nil
}

func (r *ReportRepository) GetStudentReport(ctx context.Context, studentID string) (*model.StudentReportResponse, error) {
	var res model.StudentReportResponse
	res.StudentID = studentID

	const q = `
		SELECT status, COUNT(*)
		FROM achievement_references
		WHERE student_id = $1
		GROUP BY status
	`

	rows, err := database.PG.Query(ctx, q, studentID)
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
