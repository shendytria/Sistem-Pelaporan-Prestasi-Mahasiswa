package model

import "time"

type AchievementReference struct {
	ID                 string    `db:"id"`
	StudentID          string    `db:"student_id"`
	MongoAchievementID string    `db:"mongo_achievement_id"`
	Status             string    `db:"status"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}
