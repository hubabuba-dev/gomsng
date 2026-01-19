package service

import (
	"context"
	"log"

	"msng/internal/repository"

	"github.com/google/uuid"
)

type UserService struct {
	repo *repository.UserRepository
}

func CreateUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, username string, password string) (uuid.UUID, string, error) {
	hash, err := hashPassword(password)
	if err != nil {
		return uuid.Nil, "", err
	}

	user := &repository.User{
		UUID:         uuid.New(),
		Username:     username,
		PasswordHash: hash,
		IsActive:     true,
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		log.Println(err)
		return uuid.Nil, "", err
	}

	return user.UUID, user.Username, nil
}
