package config

import (
	"log"
	"os"
)

type Config struct {
	HttpPort    string
	DatabaseDSN string
	RedisHost   string
	RedisPort   string
	RedisPass   string
}

func LoadConfigs() *Config {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	dbDSN := os.Getenv("DATABASE_DSN")
	if dbDSN == "" {
		log.Println("Warning: DATABASE_DSN is not set")
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisPass := os.Getenv("REDIS_PASS")

	return &Config{
		HttpPort:    httpPort,
		DatabaseDSN: dbDSN,
		RedisHost:   redisHost,
		RedisPort:   redisPort,
		RedisPass:   redisPass,
	}
}
