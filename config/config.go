package config

import (
	"os"
	"strconv"
)

// Config представляет конфигурацию приложения
type Config struct {
	ServerAddress string
	StorageType   string
	GroupSize     int
	DBConfig      DBConfig
}

// DBConfig представляет параметры конфигурации для базы данных
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() Config {
	groupSize, _ := strconv.Atoi(getEnv("GROUP_SIZE", "4"))
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	return Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		StorageType:   getEnv("STORAGE_TYPE", "memory"),
		GroupSize:     groupSize,
		DBConfig: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "matchmaker"),
		},
	}
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
