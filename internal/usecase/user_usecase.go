package usecase

import (
	"context"

	"matchee/services/internal/entity"
	"matchee/services/internal/repository"
)

type UserUsecase interface {
	Register(ctx context.Context, email, name string, phone string, password, confirmPassword string, role int) (*entity.User, error)
	Get(ctx context.Context, id uint64) (*entity.User, error)
}

type userUsecase struct {
	users repository.UserRepository
}

func NewUserUsecase(users repository.UserRepository) UserUsecase {
	return &userUsecase{users: users}
}
