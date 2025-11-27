package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port      string
	JWTSecret string
}

func LoadConfig() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "super_secret_key_for_dev_only"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func (c *Config) GetPortInt() int {
	p, _ := strconv.Atoi(c.Port)
	return p
}
