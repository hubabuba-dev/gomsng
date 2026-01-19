package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshToken struct {
	JTI       uuid.UUID
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	Revoked   bool
}

type RefreshTokenValidate struct {
	JTI       uuid.UUID
	UUID      uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	Revoked   bool
}

type TokenRepository struct {
	pool *pgxpool.Pool
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{
		pool: pool,
	}
}

func (r *TokenRepository) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	return r.pool.BeginTx(ctx, pgx.TxOptions{})
}

func (r *TokenRepository) AddRefreshToken(ctx context.Context, tx pgx.Tx, jti uuid.UUID, userId int64, hashToken string, expiresAt time.Time) error {
	query := `
		insert into refresh_tokens (id, user_id, token_hash, expires_at)
		values ($1, $2, $3, $4)
	`

	_, err := tx.Exec(ctx, query, jti, userId, hashToken, expiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *TokenRepository) GetRefreshToken(ctx context.Context, jti uuid.UUID, uuser_id uuid.UUID) (*RefreshTokenValidate, error) {
	var refresh_token RefreshTokenValidate

	log.Println(jti, uuser_id)
	query := `
		SELECT refresh_tokens.id, users.uuid, refresh_tokens.token_hash, refresh_tokens.expires_at, refresh_tokens.revoked
		FROM refresh_tokens 
		INNER JOIN users ON refresh_tokens.user_id = users.id
		WHERE refresh_tokens.id = $1 AND users.uuid = $2
	`

	err := r.pool.QueryRow(ctx, query, jti, uuser_id).Scan(&refresh_token.JTI, &refresh_token.UUID, &refresh_token.TokenHash, &refresh_token.ExpiresAt, &refresh_token.Revoked)
	if err != nil {

		return &RefreshTokenValidate{}, err
	}

	return &refresh_token, nil
}

func (r *TokenRepository) RevokeRefreshTokenByID(ctx context.Context, tx pgx.Tx, jti uuid.UUID) error {

	query := `
		update refresh_tokens
		set revoked = true
		where id = $1 and revoked = false
	`

	cmdTag, err := tx.Exec(ctx, query, jti)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("Tokne already revoked")
	}

	return nil
}

func (r *TokenRepository) RevokeRefreshTokenBy(ctx context.Context, tx pgx.Tx, jti uuid.UUID) error {

	query := `
		update refresh_tokens
		set revoked = true
		where id = $1 and revoked = false
	`

	cmdTag, err := tx.Exec(ctx, query, jti)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("Tokne already revoked")
	}

	return nil
}
