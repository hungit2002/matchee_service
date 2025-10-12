package usecase

import (
	"context"
	"matchee/services/internal/entity"
)

func (u *userUsecase) Get(ctx context.Context, id uint64) (*entity.User, error) {
	return u.users.GetByID(ctx, id)
}
