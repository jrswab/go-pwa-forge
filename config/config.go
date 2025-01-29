package config

import (
	"log"
	"os"
)

// Config is the main configuration for the app
type Config struct {
	RedisHost string
	RedisPass string

	PostgresUser string
	PostgresPass string
	PostgresHost string
	PostgresPort string
	PostgresDB   string
}

// Load loads the config and should only happen in main.go
func (c *Config) Load() {
	redisHost := os.Getenv("REDIS_HOST")
	c.RedisHost = redisHost
	if c.RedisHost == "" {
		log.Println("REDIS_HOST env var is blank", c)
	}

	redisPass := os.Getenv("REDIS_PASS")
	c.RedisPass = redisPass
	if c.RedisPass == "" {
		log.Println("REDIS_PASS env var is blank", c)
	}

	c.PostgresUser = os.Getenv("DB_USER")
	if c.PostgresUser == "" {
		log.Println("DB_USER env var is blank", c)
	}

	c.PostgresPass = os.Getenv("DB_PASS")
	if c.PostgresPass == "" {
		log.Println("DB_PASS env var is blank", c)
	}

	c.PostgresHost = os.Getenv("DB_HOST")
	if c.PostgresHost == "" {
		log.Println("DB_HOST env var is blank", c)
	}

	c.PostgresPort = os.Getenv("DB_PORT")
	if c.PostgresPort == "" {
		log.Println("DB_PORT env var is blank", c)
	}

	c.PostgresDB = os.Getenv("DB_NAME")
	if c.PostgresDB == "" {
		log.Println("DB_NAME env var is blank", c)
	}
}
