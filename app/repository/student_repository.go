package repository

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/database"

	"github.com/jackc/pgx/v5"
)

type StudentRepository struct{}

func NewStudentRepository() *StudentRepository { return &StudentRepository{} }

func (r *StudentRepository) FindByUserID(ctx context.Context, userID string) (*model.Student, error) {

	const q = `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students WHERE user_id = $1
	`

	row := database.PG.QueryRow(ctx, q, userID)

	var s model.Student
	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.StudentID,
		&s.ProgramStudy,
		&s.AcademicYear,
		&s.AdvisorID,
		&s.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}
