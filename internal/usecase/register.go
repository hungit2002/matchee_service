package usecase

import (
	"context"
	"errors"
	"matchee/services/internal/entity"

	"gorm.io/gorm"
)

func (u *userUsecase) Register(ctx context.Context, email, name string, phone string, password, confirmPassword string, role int) (*entity.User, error) {
	// STEP 1: Check password and confirm password
	if password != confirmPassword {
		return nil, errors.New("password and confirm password do not match")
	}

	// STEP 2: Check phone or email is exist
	user, err := u.users.GetByEmail(ctx, email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("email is already exist")
	}
	// STEP 3: Check phone is exist
	user, err = u.users.GetByPhone(ctx, phone)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("phone is already exist")
	}
	// STEP 4: Create user
	user = &entity.User{Email: &email, FullName: name, Phone: phone, PasswordHash: password}
	if err := u.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil

}
