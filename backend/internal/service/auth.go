package service

import (
	"context"
	"errors"
	"log"
	"time"

	"msng/internal/jwtutil"
	"msng/internal/repository"
	tk "msng/internal/token"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	jwt       *jwtutil.JWTManager
	hashToken *tk.HashToken
	users     *repository.UserRepository
	token     *repository.TokenRepository
}

type AuthTokenVerifier interface {
	VerifyAccessToken(token string) (userID uuid.UUID, err error)
}

func NewAuthService(jwt *jwtutil.JWTManager, users *repository.UserRepository, hashToken *tk.HashToken, token *repository.TokenRepository) *AuthService {
	return &AuthService{
		jwt:       jwt,
		hashToken: hashToken,
		users:     users,
		token:     token,
	}
}

func (s *AuthService) createRefreshToken(ctx context.Context, tx pgx.Tx, uuID uuid.UUID, userID int64) (string, uuid.UUID, error) {
	expiresAt := time.Now().Add(s.jwt.TTLRefresh)

	jti, refresh_token, err := s.jwt.GenerateRefreshToken(uuID)
	if err != nil {
		return "", uuid.Nil, err
	}

	//Hash jwttoken
	hash_token := s.hashToken.HashRefreshToken(refresh_token)

	//Insert jwt to db
	err = s.token.AddRefreshToken(ctx, tx, jti, userID, hash_token, expiresAt)
	if err != nil {
		return "", uuid.Nil, err
	}

	return refresh_token, jti, err
}

func (s *AuthService) LoginUser(ctx context.Context, username string, password string) (access string, refresh string, err error) {

	user, err := s.users.GetUserByUsername(ctx, username)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return access, refresh, ErrWrongCreds
		}
		return access, refresh, err
	}

	if !user.IsActive {
		return access, refresh, ErrBlocked
	}

	valid_password, err := checkPasswordHash(user.PasswordHash, password)
	if err != nil {
		return access, refresh, ErrWrongCreds
	}

	if !valid_password {
		return access, refresh, ErrWrongCreds
	} else {
		access, err = s.jwt.GenerateAccessToken(user.UUID)
		if err != nil {
			return access, refresh, err
		}
		tx, err := s.token.BeginTransaction(ctx)
		if err != nil {
			return access, refresh, err
		}
		defer tx.Rollback(ctx)

		refresh, _, err = s.createRefreshToken(ctx, tx, user.UUID, user.ID)
		if err != nil {
			return access, refresh, err
		}
		err = tx.Commit(ctx)
		if err != nil {
			return access, refresh, err
		}
	}

	return access, refresh, nil
}

func (s *AuthService) RefreshUserToken(ctx context.Context, refresh_token string) (access string, refresh string, err error) {
	log.Println("Refreshing token")
	userId, jti, err := s.jwt.ValidateRefreshToken(refresh_token)
	if err != nil {
		log.Println("Refreshing token")
		return access, refresh, err
	}

	db_token, err := s.token.GetRefreshToken(ctx, jti, userId)
	if err != nil {
		return access, refresh, err
	}
	user, err := s.users.GetUserByUUID(ctx, userId)
	if err != nil {
		return access, refresh, err
	}

	if !user.IsActive {
		return access, refresh, ErrBlocked
	}

	valid_token := s.hashToken.CompareToken(db_token.TokenHash, refresh_token)

	if !valid_token {
		log.Println("err1")
		return access, refresh, ErrTokenNotValid
	}

	if db_token.ExpiresAt.Before(time.Now()) {
		log.Println("err")
		return access, refresh, ErrTokenNotValid
	}

	if db_token.Revoked {
		log.Println("err322345")
		return access, refresh, ErrTokenNotValid
	}

	tx, err := s.token.BeginTransaction(ctx)
	if err != nil {
		return access, refresh, err
	}
	defer tx.Rollback(ctx)

	err = s.token.RevokeRefreshTokenByID(ctx, tx, jti)
	if err != nil {
		return access, refresh, err
	}

	access, err = s.jwt.GenerateAccessToken(userId)
	if err != nil {
		return "", "", err
	}

	refresh, _, err = s.createRefreshToken(ctx, tx, user.UUID, user.ID)
	if err != nil {
		log.Println("Errr")
		return "", "", err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", "", err
	}
	log.Println(access, refresh)

	return access, refresh, nil
}

func (s *AuthService) VerifyAccessToken(access_token string) (userID uuid.UUID, err error) {
	userID, err = s.jwt.ValidateAccessToken(access_token)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (s *AuthService) LogoutUser(ctx context.Context, refresh_token string) error {

	log.Println("start jopa")
	tx, err := s.token.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, jti, err := s.jwt.ValidateRefreshToken(refresh_token)
	if err != nil {
		return err
	}
	err = s.token.RevokeRefreshTokenByID(ctx, tx, jti)
	if err != nil {
		return err
	}
	tx.Commit(ctx)

	return nil
}
