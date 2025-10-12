package repository

import (
	"context"
	"matchee/services/internal/entity"
)

func (r *userRepository) GetByID(ctx context.Context, id uint64) (*entity.User, error) {
	var u entity.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
