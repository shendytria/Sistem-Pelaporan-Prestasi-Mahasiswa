package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	StudentID       string             `bson:"studentId"`
	AchievementType string             `bson:"achievementType"`
	Title           string             `bson:"title"`
	Description     string             `bson:"description"`

	Details     AchievementDetails `bson:"details"`
	Attachments []AchievementFile  `bson:"attachments"`
	Tags        []string           `bson:"tags"`
	Points      float64            `bson:"points"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type AchievementDetails struct {
	CompetitionName  string                 `bson:"competitionName"`
	CompetitionLevel string                 `bson:"competitionLevel"`
	Rank             int                    `bson:"rank"`
	MedalType        string                 `bson:"medalType"`
	EventDate        time.Time              `bson:"eventDate"`
	Location         string                 `bson:"location"`
	Organizer        string                 `bson:"organizer"`
	CustomFields     map[string]interface{} `bson:"customFields"`
}

type AchievementFile struct {
	FileName   string    `bson:"fileName"`
	FileURL    string    `bson:"fileUrl"`
	FileType   string    `bson:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt"`
}

type AchievementReference struct {
	ID                 string     `db:"id"`
	StudentID          string     `db:"student_id"`
	MongoAchievementID string     `db:"mongo_achievement_id"`
	Status             string     `db:"status"`
	SubmittedAt        *time.Time `db:"submitted_at"`
	VerifiedAt         *time.Time `db:"verified_at"`
	VerifiedBy         *string    `db:"verified_by"`
	RejectionNote      *string    `db:"rejection_note"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
}
