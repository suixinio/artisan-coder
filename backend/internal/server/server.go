package server

import (
	"context"
	"log"
	"net/http"

	"go.uber.org/fx"
	"github.com/gin-gonic/gin"

	"artisan-coder/internal/config"
)

// Module 返回服务器模块的 FX 选项
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewServer),
		fx.Invoke(RegisterHooks),
	)
}

// ServerIn 服务器模块的依赖组
type ServerIn struct {
	fx.In
	Router *gin.Engine
	Config *config.Config
}

// NewServer 创建 HTTP 服务器
func NewServer(in ServerIn) *http.Server {
	return &http.Server{
		Addr:         ":" + in.Config.Server.Port,
		Handler:      in.Router,
		ReadTimeout:  in.Config.Server.ReadTimeout,
		WriteTimeout: in.Config.Server.WriteTimeout,
	}
}

// RegisterHooks 注册服务器生命周期钩子
func RegisterHooks(lc fx.Lifecycle, server *http.Server, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Starting HTTP server on port %s", cfg.Server.Port)

			// 在 goroutine 中启动服务器
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("HTTP server failed: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down HTTP server...")
			return server.Shutdown(ctx)
		},
	})
}
