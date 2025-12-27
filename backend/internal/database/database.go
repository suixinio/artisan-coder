package database

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"artisan-coder/internal/config"
	"artisan-coder/internal/models"
)

// Module 返回数据库模块的 FX 选项
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewDB),
		fx.Invoke(RegisterHooks),
	)
}

// NewDB 创建数据库连接
func NewDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("Database connected successfully")

	// 启用 uuid-ossp 扩展
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return nil, fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	// 自动迁移（开发环境）
	// 生产环境应使用 golang-migrate
	if cfg.Server.Mode == "debug" {
		if err := db.AutoMigrate(&models.User{}); err != nil {
			return nil, fmt.Errorf("failed to auto migrate: %w", err)
		}
		log.Println("Database auto migration completed")
	}

	return db, nil
}

// RegisterHooks 注册数据库生命周期钩子
func RegisterHooks(lc fx.Lifecycle, db *gorm.DB) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting database connection...")
			// GORM v2 自动连接，无需额外操作
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Closing database connection...")
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		},
	})
}
