package config

import (
	"log"
	"os"
)

// Config is the main configuration for the app
type Config struct {
	RedisHost string
	RedisPass string
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
}
