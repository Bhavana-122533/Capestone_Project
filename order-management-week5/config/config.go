package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort, Environment, DBHost, DBPort, DBName, DBUser, DBPassword, LogLevel, APIKey, ProductServiceURL string
	DBSSLMode                                                                                             string
	DBMaxConns                                                                                            int32
	ProductTimeout                                                                                        time.Duration
	ProductRetries                                                                                        int
}

func Load() (Config, error) {
	_ = godotenv.Load()
	c := Config{
		AppPort: env("APP_PORT", "8080"), Environment: env("APP_ENV", "development"),
		DBHost: os.Getenv("DB_HOST"), DBPort: env("DB_PORT", "5432"), DBName: env("DB_NAME", "orderdb"),
		DBUser: os.Getenv("DB_USER"), DBPassword: os.Getenv("DB_PASSWORD"), DBSSLMode: env("DB_SSLMODE", "disable"),
		LogLevel: env("LOG_LEVEL", "info"), APIKey: os.Getenv("API_KEY"),
		ProductServiceURL: os.Getenv("PRODUCT_SERVICE_URL"), DBMaxConns: 10,
		ProductTimeout: 2 * time.Second, ProductRetries: 2,
	}
	if n := os.Getenv("DB_MAX_CONNS"); n != "" {
		v, err := strconv.Atoi(n)
		if err != nil || v < 1 {
			return c, fmt.Errorf("DB_MAX_CONNS must be a positive integer")
		}
		c.DBMaxConns = int32(v)
	}
	if n := os.Getenv("PRODUCT_TIMEOUT_MS"); n != "" {
		v, err := strconv.Atoi(n)
		if err != nil || v < 1 {
			return c, fmt.Errorf("PRODUCT_TIMEOUT_MS must be positive")
		}
		c.ProductTimeout = time.Duration(v) * time.Millisecond
	}
	if n := os.Getenv("PRODUCT_RETRIES"); n != "" {
		v, err := strconv.Atoi(n)
		if err != nil || v < 0 || v > 5 {
			return c, fmt.Errorf("PRODUCT_RETRIES must be between 0 and 5")
		}
		c.ProductRetries = v
	}
	return c, nil
}
func env(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
