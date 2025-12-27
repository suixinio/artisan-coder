package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/fx"

	"artisan-coder/internal/config"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret          []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
	issuer          string
}

func NewManager(secret string, accessDuration, refreshDuration time.Duration, issuer string) *Manager {
	return &Manager{
		secret:          []byte(secret),
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
		issuer:          issuer,
	}
}

// GenerateTokenPair 生成访问令牌和刷新令牌
func (m *Manager) GenerateTokenPair(userID uuid.UUID, email string) (accessToken, refreshToken string, err error) {
	// 生成 Access Token
	accessToken, err = m.generateToken(userID, email, m.accessDuration)
	if err != nil {
		return "", "", err
	}

	// 生成 Refresh Token
	refreshToken, err = m.generateToken(userID, email, m.refreshDuration)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (m *Manager) generateToken(userID uuid.UUID, email string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ValidateToken 验证并解析 Token
func (m *Manager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshToken 刷新 Token
func (m *Manager) RefreshToken(refreshToken string) (accessToken, newRefreshToken string, err error) {
	claims, err := m.ValidateToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	return m.GenerateTokenPair(claims.UserID, claims.Email)
}

// Module 返回 JWT 模块的 FX 选项
func Module() fx.Option {
	return fx.Provide(
		NewManagerFromConfig,
	)
}

// NewManagerFromConfig 从配置创建 JWT Manager
func NewManagerFromConfig(cfg *config.Config) *Manager {
	return NewManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessDuration,
		cfg.JWT.RefreshDuration,
		cfg.JWT.Issuer,
	)
}
