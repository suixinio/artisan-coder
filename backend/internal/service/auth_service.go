package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/fx"

	"artisan-coder/internal/models"
	"artisan-coder/internal/repository"
	"artisan-coder/pkg/jwt"
	"artisan-coder/pkg/password"
)

type AuthService interface {
	Register(ctx context.Context, username, email, userPassword string) (*models.User, string, string, error)
	Login(ctx context.Context, email, userPassword string) (*models.User, string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.User, string, string, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.Manager
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.Manager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *authService) Register(ctx context.Context, username, email, userPassword string) (*models.User, string, string, error) {
	// 检查用户是否已存在
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, "", "", repository.ErrUserAlreadyExists
	} else if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, "", "", err
	}

	// 加密密码
	hashedPassword, err := password.Hash(userPassword)
	if err != nil {
		return nil, "", "", err
	}

	// 创建用户
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", "", err
	}

	// 生成 Token
	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *authService) Login(ctx context.Context, email, userPassword string) (*models.User, string, string, error) {
	// 查找用户
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, "", "", errors.New("invalid credentials")
		}
		return nil, "", "", err
	}

	// 验证密码
	if !password.Verify(user.PasswordHash, userPassword) {
		return nil, "", "", errors.New("invalid credentials")
	}

	// 生成 Token
	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*models.User, string, string, error) {
	// 验证并刷新 Token
	accessToken, newRefreshToken, err := s.jwtManager.RefreshToken(refreshToken)
	if err != nil {
		return nil, "", "", errors.New("invalid or expired refresh token")
	}

	// 从 Access Token 中获取用户信息
	claims, err := s.jwtManager.ValidateToken(accessToken)
	if err != nil {
		return nil, "", "", err
	}

	// 获取用户信息
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, newRefreshToken, nil
}

func (s *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

// Module 返回 Service 模块的 FX 选项
func Module() fx.Option {
	return fx.Provide(
		NewAuthService,
	)
}
