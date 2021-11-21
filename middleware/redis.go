package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/trinhdaiphuc/rate-limit-redis/redis"
)

const (
	maxRequest = 10
	ttl        = 20

	// X-RateLimit-* headers
	xRateLimitLimit     = "X-RateLimit-Limit"
	xRateLimitRemaining = "X-RateLimit-Remaining"
	xRateLimitReset     = "X-RateLimit-Reset"
)

func RedisRateLimiter() fiber.Handler {
	redisCli := redis.GetRedisCli()
	// Return new handler
	return func(c *fiber.Ctx) (err error) {
		var (
			ip  = c.IP()
			ctx = context.Background()
		)
		numReq := redisCli.Incr(ctx, ip).Val()
		if numReq == 1 {
			redisCli.Expire(ctx, ip, ttl*time.Second)
		}

		timeExpired := redisCli.TTL(ctx, ip).Val().Seconds()
		c.Set(xRateLimitLimit, fmt.Sprintf("%d", maxRequest))
		c.Set(xRateLimitRemaining, fmt.Sprintf("%d", numReq))
		c.Set(xRateLimitReset, fmt.Sprintf("%0.0f", timeExpired))
		if numReq > maxRequest {
			return c.Status(http.StatusTooManyRequests).SendString("Too many requests")
		}
		return c.Next()
	}
}
