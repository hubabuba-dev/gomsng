package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	SecretKey  []byte
	Issuer     string
	TTLAccess  time.Duration
	TTLRefresh time.Duration
}

type UserClaims struct {
	jwt.RegisteredClaims
}

func NewJWTManager(secret_key []byte, issuer string, ttl_access time.Duration, ttl_refresh time.Duration) *JWTManager {
	return &JWTManager{
		SecretKey:  secret_key,
		Issuer:     issuer,
		TTLAccess:  ttl_access,
		TTLRefresh: ttl_refresh,
	}
}

func (j *JWTManager) GenerateAccessToken(userID uuid.UUID) (signedString string, err error) {
	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Issuer:    j.Issuer,
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(j.TTLAccess)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err = token.SignedString(j.SecretKey)
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func (j *JWTManager) GenerateRefreshToken(userID uuid.UUID) (jti uuid.UUID, signedString string, err error) {
	jti = uuid.New()
	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		ID:        jti.String(),
		Issuer:    j.Issuer,
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(j.TTLRefresh)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err = token.SignedString(j.SecretKey)
	if err != nil {
		return uuid.UUID{}, "", err
	}

	return jti, signedString, nil
}

func (j *JWTManager) ValidateAccessToken(tokenString string) (userId uuid.UUID, err error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unexpected signing method")
			}
			return j.SecretKey, nil
		},
		jwt.WithIssuer(j.Issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		return userId, err
	}

	claims, ok := token.Claims.(*UserClaims)

	if !ok || !token.Valid {
		return uuid.Nil, errors.New("Token not valid")
	}

	if claims.ID != "" {
		return uuid.Nil, errors.New("Token should not be refresh")
	}

	userId, err = uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.New("Token not valid")
	}

	return userId, nil
}

func (j *JWTManager) ValidateRefreshToken(tokenString string) (userId uuid.UUID, jti uuid.UUID, err error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unexpected signing method")
			}
			return j.SecretKey, nil
		},
		jwt.WithIssuer(j.Issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	claims, ok := token.Claims.(*UserClaims)

	if !ok || !token.Valid {
		return uuid.Nil, uuid.Nil, errors.New("Token not valid")
	}

	userId, err = uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("Token not valid")
	}

	if claims.ID == "" {
		return uuid.Nil, uuid.Nil, errors.New("JTI field is empty")
	}
	jti, err = uuid.Parse(claims.ID)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("Token not valid")
	}

	return userId, jti, nil
}
