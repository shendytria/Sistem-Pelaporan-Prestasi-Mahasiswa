package repository

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/database"
)

type LecturerRepository struct{}

func NewLecturerRepository() *LecturerRepository { return &LecturerRepository{} }

func (r *LecturerRepository) FindAll(ctx context.Context) ([]model.Lecturer, error) {
	const q = `SELECT id, user_id, lecturer_id, department, created_at FROM lecturers`
	rows, err := database.PG.Query(ctx, q)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []model.Lecturer
	for rows.Next() {
		var l model.Lecturer
		rows.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt)
		list = append(list, l)
	}
	return list, nil
}

func (r *LecturerRepository) FindAdvisees(ctx context.Context, lecturerID string) ([]map[string]interface{}, error) {
	const q = `
		SELECT s.id, s.student_id, s.program_study, s.academic_year
		FROM students s
		WHERE s.advisor_id = $1
	`

	rows, err := database.PG.Query(ctx, q, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, sid, ps, ay string
		rows.Scan(&id, &sid, &ps, &ay)
		list = append(list, map[string]interface{}{
			"id":            id,
			"student_id":    sid,
			"program_study": ps,
			"academic_year": ay,
		})
	}
	return list, nil
}

func (r *LecturerRepository) FindByUserID(ctx context.Context, userID string) (*model.Lecturer, error) {
	const q = `SELECT id, user_id, lecturer_id, department, created_at FROM lecturers WHERE user_id=$1`
	row := database.PG.QueryRow(ctx, q, userID)

	var l model.Lecturer
	if err := row.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt); err != nil {
		return nil, err
	}
	return &l, nil
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

func (r *LecturerRepository) FindByID(ctx context.Context, lecturerID string) (*model.Lecturer, error) {
	const q = `SELECT id, user_id, lecturer_id, department, created_at FROM lecturers WHERE id=$1 LIMIT 1`
	row := database.PG.QueryRow(ctx, q, lecturerID)

	var l model.Lecturer
	if err := row.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt); err != nil {
		return nil, err
	}
	return &l, nil
}