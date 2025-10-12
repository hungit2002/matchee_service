package repository

import (
	"context"
	"matchee/services/internal/entity"
)

func (r *userRepository) Create(ctx context.Context, u *entity.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}
