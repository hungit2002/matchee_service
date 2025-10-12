package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"matchee/services/internal/entity"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateTokenPair(user *entity.User) (*entity.AuthResponse, error)
	ValidateToken(tokenString string) (*entity.JWTClaims, error)
	RefreshToken(refreshToken string) (*entity.AuthResponse, error)
	RevokeToken(tokenString string) error
}

type jwtService struct {
	secret        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	authRepo      AuthRepository
}

type AuthRepository interface {
	GetRefreshToken(token string) (*entity.AuthToken, error)
	DeleteRefreshToken(token string) error
	CreateRefreshToken(token *entity.AuthToken) error
	GetUserByID(id uint64) (*entity.User, error)
}

func NewJWTService(secret string, accessExpiry, refreshExpiry time.Duration, authRepo AuthRepository) JWTService {
	return &jwtService{
		secret:        secret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		authRepo:      authRepo,
	}
}

func (j *jwtService) GenerateTokenPair(user *entity.User) (*entity.AuthResponse, error) {
	// Generate access token
	accessToken, err := j.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := j.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Store refresh token in database
	authToken := &entity.AuthToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(j.refreshExpiry),
	}

	if err := j.authRepo.CreateRefreshToken(authToken); err != nil {
		return nil, err
	}

	return &entity.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (j *jwtService) ValidateToken(tokenString string) (*entity.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*entity.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (j *jwtService) RefreshToken(refreshToken string) (*entity.AuthResponse, error) {
	// Validate refresh token
	authToken, err := j.authRepo.GetRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	user, err := j.authRepo.GetUserByID(authToken.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Delete old refresh token
	if err := j.authRepo.DeleteRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	// Generate new token pair
	return j.GenerateTokenPair(user)
}

func (j *jwtService) RevokeToken(tokenString string) error {
	return j.authRepo.DeleteRefreshToken(tokenString)
}

func (j *jwtService) generateAccessToken(user *entity.User) (string, error) {
	claims := &entity.JWTClaims{
		UserID:   user.ID,
		Email:    getStringValue(user.Email),
		Phone:    user.Phone,
		FullName: user.FullName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *jwtService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
