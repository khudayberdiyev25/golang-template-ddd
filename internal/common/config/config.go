package config

import (
	"os"
)

type Config struct {
	APP         string
	Environment string
	LogLevel    string
	RPCPort     string
	Context     struct {
		Timeout string
	}
	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		Sslmode  string
	}
	ALIF struct {
		BaseURL string
		Token   string
	}
	PAYME struct {
		BaseURL  string
		Timeout  string
		Username string
		Password string
	}
	OTLPCollector struct {
		Host string
		Port string
	}
}

func New() *Config {
	var config Config

	config.APP = getEnv("APP", "iman_merchant_service")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.RPCPort = getEnv("RPC_PORT", ":5001")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	// initialization db
	config.DB.Host = getEnv("POSTGRES_HOST", "localhost")
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "payment_service")
	config.DB.User = getEnv("POSTGRES_USER", "iman_develop")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "iman_develop")
	config.DB.Sslmode = getEnv("POSTGRES_SSLMODE", "disable")

	config.ALIF.BaseURL = getEnv("ALIF_PAY_BASE_URL", "https://panda-dev.alifpay.uz")
	config.ALIF.Token = getEnv("ALIF_PAY_TOKEN", "c07f8954-0b1e-438a-a51e-e69c5ff34277")

	return &config
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
