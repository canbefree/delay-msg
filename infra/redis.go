package infra

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

func NewRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		MaxActive:   20,
		IdleTimeout: 30 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "redis:6379")
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
