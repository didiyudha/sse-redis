package redis

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

// Config redis.
type Config struct {
	Host string
	Port int
}

// LoadEnv - load configuration from env.
func LoadEnv() Config {
	return Config{
		Host: os.Getenv("REDIS_HOST"),
		Port: port(),
	}
}

// NewRedis creates new connection to redis and return the connection
func NewRedis(cfg Config) (redis.Conn, error) {
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	pong, err := conn.Do("PING")
	if err != nil {
		return nil, err
	}
	_, err = redis.String(pong, err)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func port() int {
	p, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		return 6379
	}
	return p
}

func db() int {
	d, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return 1
	}
	return d
}
