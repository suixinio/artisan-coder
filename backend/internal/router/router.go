package router

import (
	"go.uber.org/fx"
	"github.com/gin-gonic/gin"

	"artisan-coder/internal/config"
	"artisan-coder/internal/handler"
	"artisan-coder/internal/middleware"
	"artisan-coder/pkg/jwt"
)

// Module 返回路由模块的 FX 选项
func Module() fx.Option {
	return fx.Provide(
		NewRouter,
	)
}

// RouterIn 路由模块的依赖组
type RouterIn struct {
	fx.In
	AuthHandler *handler.AuthHandler
	JWTManager  *jwt.Manager
	Config      *config.Config
}

// NewRouter 创建 Gin 路由
func NewRouter(in RouterIn) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(in.Config.Server.Mode)

	router := gin.New()

	// 全局中间件
	router.Use(middleware.CORS(in.Config.CORS.AllowedOrigins))
	router.Use(middleware.Logger())
	router.Use(gin.Recovery())

	// 注册路由
	setupRoutes(router, in.AuthHandler, in.JWTManager)

	return router
}

// setupRoutes 配置所有路由
func setupRoutes(router *gin.Engine, authHandler *handler.AuthHandler, jwtManager *jwt.Manager) {
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.RefreshToken)

			// 需要认证的路由
			auth.GET("/me", middleware.Auth(jwtManager), authHandler.GetCurrentUser)
		}
	}
}
