package repository

import (
	"context"
	"matchee/services/internal/entity"
)

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var u entity.User
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
