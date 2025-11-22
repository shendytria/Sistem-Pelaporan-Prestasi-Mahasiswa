package model

import "time"

type Lecturer struct {
	ID         string    `db:"id"`
	UserID     string    `db:"user_id"`
	LecturerID string    `db:"lecturer_id"`
	Department string    `db:"department"`
	CreatedAt  time.Time `db:"created_at"`
}
