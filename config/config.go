package config

import (
	"strconv"
	"sync"

	"os"

	"github.com/didiyudha/sse-redis/internal/platform/redis"
)

// Cfg is a configuration variable for the app.
var Cfg Config
var once sync.Once

// Config is a general configuration.
type Config struct {
	Port  int
	Redis redis.Config
}

// LoadEnv loads configuration from env variables.
func LoadEnv() {
	once.Do(func() {
		Cfg = Config{
			Port:  port(),
			Redis: redis.LoadEnv(),
		}
	})
}

func port() int {
	p, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return 8080
	}
	return p
}
