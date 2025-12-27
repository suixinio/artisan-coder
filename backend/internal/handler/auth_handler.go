package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/fx"

	"artisan-coder/internal/models"
	"artisan-coder/internal/repository"
	"artisan-coder/internal/service"
	"artisan-coder/pkg/response"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=50"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=7"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"rememberMe"`
}

type AuthResponse struct {
	User         *UserResponse `json:"user"`
	Token        string        `json:"token"`
	RefreshToken string        `json:"refreshToken"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证密码一致性
	if req.Password != req.ConfirmPassword {
		response.BadRequest(c, "Passwords do not match")
		return
	}

	// 调用服务层
	user, accessToken, refreshToken, err := h.authService.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			response.Conflict(c, "User with this email already exists")
		} else {
			response.InternalError(c)
		}
		return
	}

	response.Created(c, &AuthResponse{
		User:         toUserResponse(user),
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用服务层
	user, accessToken, refreshToken, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Unauthorized(c, "Invalid email or password")
		return
	}

	response.Success(c, &AuthResponse{
		User:         toUserResponse(user),
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT 无状态，客户端删除 Token 即可
	// 如果需要实现 Token 黑名单，可以在此添加逻辑
	response.Success(c, nil)
}

// RefreshToken 刷新 Token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		response.Unauthorized(c, "Missing refresh token")
		return
	}

	// 移除 "Bearer " 前缀
	if len(refreshToken) > 7 && refreshToken[:7] == "Bearer " {
		refreshToken = refreshToken[7:]
	}

	// 调用服务层
	user, accessToken, newRefreshToken, err := h.authService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		response.Unauthorized(c, "Invalid or expired refresh token")
		return
	}

	response.Success(c, &AuthResponse{
		User:         toUserResponse(user),
		Token:        accessToken,
		RefreshToken: newRefreshToken,
	})
}

// GetCurrentUser 获取当前用户
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, toUserResponse(user))
}

func toUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// Module 返回 Handler 模块的 FX 选项
func Module() fx.Option {
	return fx.Provide(
		NewAuthHandler,
	)
}
