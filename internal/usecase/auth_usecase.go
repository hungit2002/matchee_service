package usecase

import (
	"errors"

	"matchee/services/internal/entity"
	"matchee/services/internal/repository"
	"matchee/services/internal/service"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(req *entity.RegisterRequest) (*entity.AuthResponse, error)
	Login(req *entity.LoginRequest) (*entity.AuthResponse, error)
	RefreshToken(req *entity.RefreshTokenRequest) (*entity.AuthResponse, error)
	Logout(refreshToken string) error
	GetCurrentUser(userID uint64) (*entity.User, error)
	UpdateProfile(userID uint64, req *entity.UpdateProfileRequest) (*entity.User, error)
	ChangePassword(userID uint64, req *entity.ChangePasswordRequest) error
}

type authUsecase struct {
	authRepo   repository.AuthRepository
	userRepo   repository.UserRepository
	jwtService service.JWTService
}

func NewAuthUsecase(authRepo repository.AuthRepository, userRepo repository.UserRepository, jwtService service.JWTService) AuthUsecase {
	return &authUsecase{
		authRepo:   authRepo,
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (a *authUsecase) Register(req *entity.RegisterRequest) (*entity.AuthResponse, error) {
	// Check if phone already exists
	existingUser, _ := a.authRepo.GetUserByPhone(req.Phone)
	if existingUser != nil {
		return nil, errors.New("phone number already exists")
	}

	// Check if email already exists (if provided)
	if req.Email != "" {
		existingUser, _ := a.authRepo.GetUserByEmail(req.Email)
		if existingUser != nil {
			return nil, errors.New("email already exists")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &entity.User{
		FullName:     req.FullName,
		Phone:        req.Phone,
		Email:        &req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := a.authRepo.CreateUser(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Assign default player role
	playerRole, err := a.authRepo.GetRoleByName(entity.RolePlayer)
	if err == nil && playerRole != nil {
		// Assign player role to user
		userRole := &entity.UserRole{
			UserID: user.ID,
			RoleID: playerRole.ID,
		}
		if err := a.authRepo.CreateUserRole(userRole); err != nil {
			// If role assignment fails, we'll continue without role assignment
			// This ensures user registration doesn't fail due to role issues
		}
	}

	// Generate token pair
	return a.jwtService.GenerateTokenPair(user)
}

func (a *authUsecase) Login(req *entity.LoginRequest) (*entity.AuthResponse, error) {
	// Get user by phone
	user, err := a.authRepo.GetUserByPhone(req.Phone)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token pair
	return a.jwtService.GenerateTokenPair(user)
}

func (a *authUsecase) RefreshToken(req *entity.RefreshTokenRequest) (*entity.AuthResponse, error) {
	return a.jwtService.RefreshToken(req.RefreshToken)
}

func (a *authUsecase) Logout(refreshToken string) error {
	return a.jwtService.RevokeToken(refreshToken)
}

func (a *authUsecase) GetCurrentUser(userID uint64) (*entity.User, error) {
	return a.authRepo.GetUserByID(userID)
}

func (a *authUsecase) UpdateProfile(userID uint64, req *entity.UpdateProfileRequest) (*entity.User, error) {
	// Get current user
	user, err := a.authRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.Phone != "" {
		// Check if phone already exists for another user
		existingUser, _ := a.authRepo.GetUserByPhone(req.Phone)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("phone number already exists")
		}
		user.Phone = req.Phone
	}

	if req.Email != nil {
		if *req.Email != "" {
			// Check if email already exists for another user
			existingUser, _ := a.authRepo.GetUserByEmail(*req.Email)
			if existingUser != nil && existingUser.ID != userID {
				return nil, errors.New("email already exists")
			}
		}
		user.Email = req.Email
	}

	// Save updated user
	if err := a.authRepo.UpdateUser(user); err != nil {
		return nil, errors.New("failed to update user")
	}

	return user, nil
}

func (a *authUsecase) ChangePassword(userID uint64, req *entity.ChangePasswordRequest) error {
	// Get current user
	user, err := a.authRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	if err := a.authRepo.UpdateUser(user); err != nil {
		return errors.New("failed to update password")
	}

	// Revoke all refresh tokens for security
	return a.authRepo.DeleteUserRefreshTokens(userID)
}
