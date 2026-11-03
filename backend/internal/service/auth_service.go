package service

import (
	"errors"
	"strings"

	"servixa/internal/auth"
	"servixa/internal/config"
	"servixa/internal/model"
	"servixa/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepository      *repository.UserRepository
	config              config.Config
	refreshTokenService *RefreshTokenService
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	User         *model.User
}

func NewAuthService(
	userRepository *repository.UserRepository,
	refreshTokenService *RefreshTokenService,
	cfg config.Config,
) *AuthService {
	return &AuthService{
		userRepository:      userRepository,
		refreshTokenService: refreshTokenService,
		config:              cfg,
	}
}

func (s *AuthService) Register(
	name string,
	email string,
	password string,
) (*model.User, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" {
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	existingUser, err := s.userRepository.FindByEmail(email)

	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     model.RoleCustomer,
		IsActive: true,
	}

	if err := s.userRepository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(
	email string,
	password string,
) (*LoginResult, error) {

	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	user, err := s.userRepository.FindByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}

		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := auth.GenerateAccessToken(
		user.ID,
		user.Role,
		s.config.JWTSecret,
		s.config.JWTAccessTokenDuration,
	)

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.refreshTokenService.Create(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		User:         user,
	}, nil
}

func (s *AuthService) RefreshAccessToken(
	refreshToken string,
) (string, string, error) {
	token, err := s.refreshTokenService.Validate(refreshToken)
	if err != nil {
		return "", "", err
	}

	user, err := s.userRepository.FindByID(token.UserID)
	if err != nil {
		return "", "", err
	}

	if !user.IsActive {
		return "", "", errors.New("user is inactive")
	}

	newAccessToken, err := auth.GenerateAccessToken(
		user.ID,
		user.Role,
		s.config.JWTSecret,
		s.config.JWTAccessTokenDuration,
	)

	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.refreshTokenService.Rotate(
		token,
		user.ID,
	)

	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken.Token, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	token, err := s.refreshTokenService.Validate(refreshToken)
	if err != nil {
		return err
	}

	return s.refreshTokenService.Revoke(token)
}
