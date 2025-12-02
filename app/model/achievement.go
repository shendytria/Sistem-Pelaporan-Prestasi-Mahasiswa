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
	Details         AchievementDetails `bson:"details"`
	Attachments     []AchievementFile  `bson:"attachments"`
	Tags            []string           `bson:"tags"`
	Points          float64            `bson:"points"`
	CreatedAt       time.Time          `bson:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt"`
}

type CustomField struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

type AchievementDetails struct {
	CompetitionName  string        `bson:"competitionName"`
	CompetitionLevel string        `bson:"competitionLevel"`
	Rank             int           `bson:"rank"`
	MedalType        string        `bson:"medalType"`
	EventDate        time.Time     `bson:"eventDate"`
	Location         string        `bson:"location"`
	Organizer        string        `bson:"organizer"`
	CustomFields     []CustomField `bson:"customFields"`
}

type AchievementFile struct {
	FileName   string    `bson:"fileName"`
	FileURL    string    `bson:"fileUrl"`
	FileType   string    `bson:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt"`
}

type AchievementMongoUpdate struct {
	Title           *string             `bson:"title,omitempty"`
	Description     *string             `bson:"description,omitempty"`
	AchievementType *string             `bson:"achievementType,omitempty"`
	Details         *AchievementDetails `bson:"details,omitempty"`
	Tags            []string            `bson:"tags,omitempty"`
	Points          *float64            `bson:"points,omitempty"`
	UpdatedAt       time.Time           `bson:"updatedAt"`
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

type AchievementUpdate struct {
	Title           *string             `json:"title"`
	Description     *string             `json:"description"`
	AchievementType *string             `json:"achievementType"`
	Details         *AchievementDetails `json:"details"`
	Tags            []string            `json:"tags"`
	Points          *float64            `json:"points"`
}

type AchievementHistory struct {
	Status string     `json:"status"`
	At     *time.Time `json:"at"`
	By     *string    `json:"by,omitempty"`
	Note   *string    `json:"note,omitempty"`
}
