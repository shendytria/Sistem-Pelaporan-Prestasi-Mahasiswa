package repository

import (
	"context"
	"prestasi_mhs/app/model"
	"prestasi_mhs/database"

	"github.com/jackc/pgx/v5"
)

type UserRepo interface {
    FindByUsername(ctx context.Context, username string) (*model.User, error)
    FindAll(ctx context.Context, limit, offset int) ([]model.User, int, error)
    FindByID(ctx context.Context, id string) (*model.User, error)
    Create(ctx context.Context, u *model.User) error
    Update(ctx context.Context, u *model.User) error
    Delete(ctx context.Context, id string) error
    UpdateRole(ctx context.Context, id string, roleID string) error
    GetPermissionsByRole(ctx context.Context, roleID string) ([]string, error)
    GetRoleName(ctx context.Context, roleID string) (string, error)
}

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

func (r *UserRepository) FindAll(ctx context.Context, limit, offset int) ([]model.User, int, error) {
	const qCount = `SELECT COUNT(*) FROM users`
	var total int
	if err := database.PG.QueryRow(ctx, qCount).Scan(&total); err != nil {
		return nil, 0, err
	}
	
	const q = `
		SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := database.PG.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		rows.Scan(
			&u.ID, &u.Username, &u.Email,
			&u.FullName, &u.RoleID,
			&u.IsActive, &u.CreatedAt, &u.UpdatedAt,
		)
		users = append(users, u)
	}

	return users, total, nil
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

func (r *UserRepository) GetPermissionsByRole(ctx context.Context, roleID string) ([]string, error) {
	const q = `
		SELECT p.name
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`
	rows, err := database.PG.Query(ctx, q, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		perms = append(perms, name)
	}

	return perms, nil
}

func (r *UserRepository) GetRoleName(ctx context.Context, roleID string) (string, error) {
	const q = `SELECT name FROM roles WHERE id = $1`
	var name string
	err := database.PG.QueryRow(ctx, q, roleID).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}