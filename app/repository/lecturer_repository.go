package repository

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/database"

	"github.com/jackc/pgx/v5"
)

type LecturerRepository struct{}

func NewLecturerRepository() *LecturerRepository {
	return &LecturerRepository{}
}

func (r *LecturerRepository) FindByUserID(ctx context.Context, userID string) (*model.Lecturer, error) {

	const q = `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE user_id = $1
	`

	row := database.PG.QueryRow(ctx, q, userID)

	var lec model.Lecturer
	err := row.Scan(
		&lec.ID,
		&lec.UserID,
		&lec.LecturerID,
		&lec.Department,
		&lec.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &lec, nil
}
