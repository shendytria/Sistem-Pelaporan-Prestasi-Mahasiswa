package repository

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/database"
)

type AchievementReferenceRepository struct{}

func NewAchievementReferenceRepository() *AchievementReferenceRepository {
	return &AchievementReferenceRepository{}
}

func (r *AchievementReferenceRepository) Insert(ctx context.Context, ref *model.AchievementReference) error {

	const q = `
		INSERT INTO achievement_references
		(id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`

	_, err := database.PG.Exec(ctx, q,
		ref.ID, ref.StudentID, ref.MongoAchievementID, ref.Status,
	)

	return err
}

func (r *AchievementReferenceRepository) FindMongoIDsByStudent(ctx context.Context, studentID string) ([]string, error) {

	const q = `
		SELECT mongo_achievement_id
		FROM achievement_references
		WHERE student_id::text = $1
	`

	rows, err := database.PG.Query(ctx, q, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			list = append(list, id)
		}
	}

	return list, nil
}
