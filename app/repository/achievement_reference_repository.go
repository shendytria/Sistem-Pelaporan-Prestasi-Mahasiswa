package repository

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/database"
	"time"
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

func (r *AchievementReferenceRepository) FindAll(ctx context.Context) ([]model.AchievementReference, error) {

	const q = `
        SELECT id, student_id, mongo_achievement_id, status, created_at, updated_at
        FROM achievement_references
    `

	rows, err := database.PG.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AchievementReference

	for rows.Next() {
		var ref model.AchievementReference
		if err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		); err == nil {
			list = append(list, ref)
		}
	}

	return list, nil
}

func (r *AchievementReferenceRepository) FindByID(ctx context.Context, id string) (*model.AchievementReference, error) {

	const q = `
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE id = $1
    `

	row := database.PG.QueryRow(ctx, q, id)

	var ref model.AchievementReference

	err := row.Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID,
		&ref.Status, &ref.SubmittedAt, &ref.VerifiedAt,
		&ref.VerifiedBy, &ref.RejectionNote,
		&ref.CreatedAt, &ref.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &ref, nil
}

func (r *AchievementReferenceRepository) FindByStudent(ctx context.Context, studentID string) ([]model.AchievementReference, error) {

	const q = `
        SELECT id, student_id, mongo_achievement_id, status, created_at, updated_at
        FROM achievement_references
        WHERE student_id::text = $1
        ORDER BY created_at DESC
    `

	rows, err := database.PG.Query(ctx, q, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AchievementReference

	for rows.Next() {
		var ref model.AchievementReference
		if err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		); err == nil {
			list = append(list, ref)
		}
	}

	return list, nil
}

func (r *AchievementReferenceRepository) UpdateStatus(ctx context.Context, id string, status string, submittedAt *time.Time, verifiedAt *time.Time, verifiedBy *string, rejectionNote *string, ) error {

	const q = `
        UPDATE achievement_references
        SET status = $1,
            submitted_at = $2,
            verified_at = $3,
            verified_by = $4,
            rejection_note = $5,
            updated_at = NOW()
        WHERE id = $6
    `

	_, err := database.PG.Exec(ctx, q,
		status,
		submittedAt,
		verifiedAt,
		verifiedBy,
		rejectionNote,
		id,
	)

	return err
}
