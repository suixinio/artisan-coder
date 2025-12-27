package app

import (
	"go.uber.org/fx"

	"artisan-coder/internal/config"
	"artisan-coder/internal/database"
	"artisan-coder/internal/handler"
	"artisan-coder/internal/repository"
	"artisan-coder/internal/router"
	"artisan-coder/internal/server"
	"artisan-coder/internal/service"
	"artisan-coder/pkg/jwt"
)

// Module 返回 FX 应用模块
func Module() fx.Option {
	return fx.Options(
		// 基础模块
		config.Module(),

		// 数据层
		database.Module(),
		jwt.Module(),

		// 业务层
		repository.Module(),
		service.Module(),

		// HTTP 层
		handler.Module(),
		router.Module(),

		// 服务器层
		server.Module(),
	)
}
