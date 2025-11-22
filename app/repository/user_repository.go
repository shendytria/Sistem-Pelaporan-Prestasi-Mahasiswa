package repository

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/database"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	const q = `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users WHERE username = $1
	`

	row := database.PG.QueryRow(ctx, q, username)

	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
