package main

import (
	"go.uber.org/fx"
	"artisan-coder/internal/app"
)

func main() {
	// 创建并启动 FX 应用
	app := fx.New(
		app.Module(),
	)

	// 运行应用（阻塞直到收到信号）
	app.Run()
}
