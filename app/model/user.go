package model

import "time"

type User struct {
	ID           string    `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FullName     string    `db:"full_name"`
	RoleID       string    `db:"role_id"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type CreateUserReq struct {
    Username     string `json:"username"`
    Email        string `json:"email"`
    Password     string `json:"password_hash"`
    FullName     string `json:"full_name"`
    RoleID       string `json:"role_id"`
}
