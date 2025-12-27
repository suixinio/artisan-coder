package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

type ServerConfig struct {
	Port            string        `mapstructure:"port"`
	Mode            string        `mapstructure:"mode"`             // debug, release, test
	ReadTimeout     time.Duration `mapstructure:"readTimeout"`
	WriteTimeout    time.Duration `mapstructure:"writeTimeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbName"`
	SSLMode  string `mapstructure:"sslMode"`
}

type JWTConfig struct {
	Secret          string        `mapstructure:"secret"`
	AccessDuration  time.Duration `mapstructure:"accessDuration"`
	RefreshDuration time.Duration `mapstructure:"refreshDuration"`
	Issuer          string        `mapstructure:"issuer"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowedOrigins"`
	AllowedMethods []string `mapstructure:"allowedMethods"`
	AllowedHeaders []string `mapstructure:"allowedHeaders"`
}

// Load 加载配置
func Load() (*Config, error) {
	v := viper.New()

	// 1. 设置默认值
	setDefaults(v)

	// 2. 配置文件路径和名称
	// 根据环境变量 APP_ENV 选择配置文件（development, production）
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	v.SetConfigName("config." + env)
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	// 3. 读取配置文件（可选，不存在不报错）
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		log.Println("No config file found, using defaults and environment variables")
	} else {
		log.Printf("Using config file: %s", v.ConfigFileUsed())
	}

	// 4. 绑定环境变量
	// 环境变量可以覆盖配置文件
	// 例如：SERVER_PORT -> server.port
	//       DB_HOST -> database.host
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 5. 解析到结构体
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// setDefaults 设置默认配置值
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.readTimeout", "15s")
	v.SetDefault("server.writeTimeout", "15s")
	v.SetDefault("server.shutdownTimeout", "10s")

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.user", "artisan")
	v.SetDefault("database.password", "artisan123")
	v.SetDefault("database.dbName", "artisan_coder")
	v.SetDefault("database.sslMode", "disable")

	// JWT defaults
	v.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	v.SetDefault("jwt.accessDuration", "1h")
	v.SetDefault("jwt.refreshDuration", "168h") // 7 days
	v.SetDefault("jwt.issuer", "artisan-coder")

	// CORS defaults
	v.SetDefault("cors.allowedOrigins", []string{"http://localhost:5173"})
	v.SetDefault("cors.allowedMethods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allowedHeaders", []string{"Origin", "Content-Type", "Authorization"})
}

// Module 返回配置模块的 FX 选项
func Module() fx.Option {
	return fx.Provide(
		Load,
	)
}
