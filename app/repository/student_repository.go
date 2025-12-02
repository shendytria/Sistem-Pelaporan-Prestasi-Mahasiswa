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

func (r *StudentRepository) IsMyStudent(ctx context.Context, dosenUserID, studentID string) (bool, error) {

    const q = `
        SELECT COUNT(*)
        FROM students
        WHERE id = $1
        AND advisor_id = (
            SELECT id FROM lecturers WHERE user_id = $2
        )
    `

    var count int
    err := database.PG.QueryRow(ctx, q, studentID, dosenUserID).Scan(&count)
    if err != nil {
        return false, err
    }

    return count > 0, nil
}

func (r *StudentRepository) FindAll(ctx context.Context) ([]model.Student, error) {
	const q = `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at FROM students`
	rows, err := database.PG.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Student
	for rows.Next() {
		var s model.Student
		rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt)
		list = append(list, s)
	}
	return list, nil
}

func (r *StudentRepository) FindByID(ctx context.Context, id string) (*model.Student, error) {
	const q = `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at FROM students WHERE id=$1`
	row := database.PG.QueryRow(ctx, q, id)

	var s model.Student
	err := row.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StudentRepository) UpdateAdvisor(ctx context.Context, studentID string, lecturerID string) error {
	const q = `UPDATE students SET advisor_id = $1 WHERE id = $2`
	_, err := database.PG.Exec(ctx, q, lecturerID, studentID)
	return err
}

func (r *StudentRepository) CheckAdvisor(ctx context.Context, lecturerID, studentID string) (bool, error) {
	const q = `SELECT COUNT(*) FROM students WHERE id = $1 AND advisor_id = $2`
	var count int
	if err := database.PG.QueryRow(ctx, q, studentID, lecturerID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *LecturerRepository) FindStudentByUserID(ctx context.Context, userID string) (*model.Student, error) {
	const q = `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id
		FROM students
		WHERE user_id = $1
		LIMIT 1
	`
	row := database.PG.QueryRow(ctx, q, userID)

	var s model.Student
	if err := row.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID); err != nil {
		return nil, err
	}
	return &s, nil
}