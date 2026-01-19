package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           int64
	UUID         uuid.UUID
	Username     string
	PasswordHash string
	IsActive     bool
	CreatedAt    time.Time
}

type UserRepository struct {
	pool *pgxpool.Pool
}

func CreateUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
	query := `
		insert into users (uuid, username, password_hash, is_active)
		values ($1, $2, $3, $4)
		returning id, created_at
	`

	err := r.pool.QueryRow(ctx, query, user.UUID, user.Username, user.PasswordHash, user.IsActive).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {

	query := `
		select id, uuid, username, password_hash, is_active
		from users
		where username = $1
	`

	var user User

	err := r.pool.QueryRow(ctx, query, username).Scan(&user.ID, &user.UUID, &user.Username, &user.PasswordHash, &user.IsActive)
	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUUID(ctx context.Context, UUID uuid.UUID) (*User, error) {

	query := `
		select id, uuid, username, password_hash, is_active
		from users
		where uuid = $1
	`

	var user User

	err := r.pool.QueryRow(ctx, query, UUID).Scan(&user.ID, &user.UUID, &user.Username, &user.PasswordHash, &user.IsActive)
	if err != nil {
		return &User{}, err
	}

	return &user, nil
}
