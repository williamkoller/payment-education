package config

import (
	"fmt"
	"os"
	"strconv"
)

type DatabaseConfiguration struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
}

type AppConfiguration struct {
	Port    int
	AppName string
	Env     string
}

type Config struct {
	Database DatabaseConfiguration
	App      AppConfiguration
	Resend   string
}

func LoadConfig() (*Config, error) {
	dbCfg, err := loadDatabaseConfiguration()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar configuração do banco: %w", err)
	}

	appCfg, err := loadAppConfiguration()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar configuração da aplicação: %w", err)
	}

	resend := loadResend()

	return &Config{
		Database: *dbCfg,
		App:      *appCfg,
		Resend:   resend,
	}, nil
}

func loadDatabaseConfiguration() (*DatabaseConfiguration, error) {
	portStr := os.Getenv("DB_PORT")
	if portStr == "" {
		portStr = "5432"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("DB_PORT inválida: %v", err)
	}

	cfg := &DatabaseConfiguration{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		Username: getEnv("DB_USERNAME", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DbName:   getEnv("DB_DATABASE", "postgres"),
	}

	return cfg, nil
}

func loadAppConfiguration() (*AppConfiguration, error) {
	portStr := getEnv("APP_PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("APP_PORT inválida: %v", err)
	}

	cfg := &AppConfiguration{
		Port:    port,
		AppName: getEnv("APP_NAME", "myapp"),
		Env:     getEnv("APP_ENV", "development"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func loadResend() string {
	return getEnv("RESEND_API_KEY", "")
}
