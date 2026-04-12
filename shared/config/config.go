package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Environment      string
	Port             string
	MongoDBURI       string
	JWTSecret        string
	JWTExpirationHrs int
	RefreshTokenHrs  int
}

func Load() *Config {
	loadEnvFile(".env")

	return &Config{
		Environment:      getEnv("ENVIRONMENT", "development"),
		Port:             getEnv("PORT", "9090"),
		MongoDBURI:       os.Getenv("MONGODB_URI"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTExpirationHrs: getEnvInt("JWT_EXPIRATION_HOURS", 24),
		RefreshTokenHrs:  getEnvInt("REFRESH_TOKEN_HOURS", 168),
	}
}

func loadEnvFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
