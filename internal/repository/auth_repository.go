package repository

import (
	"time"

	"matchee/services/internal/entity"

	"gorm.io/gorm"
)

type AuthRepository interface {
	// User authentication methods
	CreateUser(user *entity.User) error
	GetUserByPhone(phone string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id uint64) (*entity.User, error)
	UpdateUser(user *entity.User) error

	// Role management methods
	GetRoleByName(name string) (*entity.Role, error)
	CreateRole(role *entity.Role) error
	CreateUserRole(userRole *entity.UserRole) error

	// Token management methods
	CreateRefreshToken(token *entity.AuthToken) error
	GetRefreshToken(token string) (*entity.AuthToken, error)
	DeleteRefreshToken(token string) error
	DeleteUserRefreshTokens(userID uint64) error
	CleanupExpiredTokens() error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// User authentication methods
func (r *authRepository) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) GetUserByPhone(phone string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("UserRoles.Role").Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("UserRoles.Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByID(id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("UserRoles.Role").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) UpdateUser(user *entity.User) error {
	return r.db.Save(user).Error
}

// Token management methods
func (r *authRepository) CreateRefreshToken(token *entity.AuthToken) error {
	return r.db.Create(token).Error
}

func (r *authRepository) GetRefreshToken(token string) (*entity.AuthToken, error) {
	var authToken entity.AuthToken
	err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&authToken).Error
	if err != nil {
		return nil, err
	}
	return &authToken, nil
}

func (r *authRepository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&entity.AuthToken{}).Error
}

func (r *authRepository) DeleteUserRefreshTokens(userID uint64) error {
	return r.db.Where("user_id = ?", userID).Delete(&entity.AuthToken{}).Error
}

func (r *authRepository) CleanupExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&entity.AuthToken{}).Error
}

// Role management methods
func (r *authRepository) GetRoleByName(name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *authRepository) CreateRole(role *entity.Role) error {
	return r.db.Create(role).Error
}

func (r *authRepository) CreateUserRole(userRole *entity.UserRole) error {
	return r.db.Create(userRole).Error
}
