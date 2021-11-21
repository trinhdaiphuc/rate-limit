package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis"
)

func FiberLimiter() fiber.Handler {
	store := redis.New(redis.Config{
		Host:     "localhost",
		Port:     6379,
		Username: "",
		Password: "password",
		Database: 1,
		Reset:    false,
	})
	return limiter.New(limiter.Config{
		Max:        maxRequest,
		Expiration: ttl * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).SendString("Too many requests")
		},
		Storage: store,
	})
}
