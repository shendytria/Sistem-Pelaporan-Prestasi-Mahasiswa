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

func (r *UserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	const q = `
		SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`

	rows, err := database.PG.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var u model.User
		err := rows.Scan(
			&u.ID, &u.Username, &u.Email,
			&u.FullName, &u.RoleID,
			&u.IsActive, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	const q = `
		SELECT id, username, email, password_hash, full_name, role_id,
		       is_active, created_at, updated_at
		FROM users WHERE id = $1
	`

	row := database.PG.QueryRow(ctx, q, id)

	var u model.User
	err := row.Scan(
		&u.ID, &u.Username, &u.Email,
		&u.PasswordHash, &u.FullName,
		&u.RoleID, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	const q = `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := database.PG.Exec(ctx, q,
		u.ID, u.Username, u.Email,
		u.PasswordHash, u.FullName, u.RoleID,
	)

	return err
}

func (r *UserRepository) Update(ctx context.Context, u *model.User) error {
	const q = `
		UPDATE users
		SET username=$1, email=$2, full_name=$3, is_active=$4, updated_at=NOW()
		WHERE id=$5
	`

	_, err := database.PG.Exec(ctx, q,
		u.Username, u.Email, u.FullName,
		u.IsActive, u.ID,
	)

	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	const q = `
		DELETE FROM users WHERE id = $1
	`

	_, err := database.PG.Exec(ctx, q, id)
	return err
}

func (r *UserRepository) UpdateRole(ctx context.Context, id string, roleID string) error {
	const q = `
		UPDATE users SET role_id = $1, updated_at = NOW() WHERE id = $2
	`

	_, err := database.PG.Exec(ctx, q, roleID, id)
	return err
}