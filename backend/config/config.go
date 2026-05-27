package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBServer   string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort int
	ServerHost string
}

func LoadConfig() *Config {
	// Load .env file if exists
	_ = godotenv.Load()

	port, _ := strconv.Atoi(getEnv("DB_PORT", "1433"))
	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "3000"))

	return &Config{
		DBServer:   getEnv("DB_SERVER", "localhost"),
		DBPort:     port,
		DBUser:     getEnv("DB_USER", "sa"),
		DBPassword: getEnv("DB_PASSWORD", "mypassword!1234"),
		DBName:     getEnv("DB_NAME", "master"),
		ServerPort: serverPort,
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s",
		c.DBServer,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
	)
}
