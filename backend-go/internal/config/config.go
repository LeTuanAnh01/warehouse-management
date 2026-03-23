package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	App      AppConfig
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type AppConfig struct {
	DefaultPageSize int
	MaxPageSize     int
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// Try loading from parent directory (for cmd/api)
		_ = godotenv.Load("../../.env")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "anhlt"),
			Password:        getEnv("DB_PASSWORD", ""),
			DBName:          getEnv("DB_NAME", "warehouse_db"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key"),
			Expiration: getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
		},
		App: AppConfig{
			DefaultPageSize: getEnvAsInt("DEFAULT_PAGE_SIZE", 20),
			MaxPageSize:     getEnvAsInt("MAX_PAGE_SIZE", 100),
		},
	}
}

func (c *Config) GetDSN() string {
	password := c.Database.Password
	if password == "" {
		// For empty password, use simple format without password parameter
		return fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=%s",
			c.Database.Host,
			c.Database.Port,
			c.Database.User,
			c.Database.DBName,
			c.Database.SSLMode,
		)
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
