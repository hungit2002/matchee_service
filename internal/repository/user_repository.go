package repository

import (
	"context"

	"gorm.io/gorm"

	"matchee/services/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id uint64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *entity.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var u entity.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *userRepository) GetByID(ctx context.Context, id uint64) (*entity.User, error) {
	var u entity.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var u entity.User
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
