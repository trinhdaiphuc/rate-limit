package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/trinhdaiphuc/rate-limit-redis/middleware"
	"github.com/trinhdaiphuc/rate-limit-redis/redis"
)

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
	})
	redis.NewRedisClient()
	redisCli := redis.GetRedisCli()
	app.Use(recover.New())

	app.Get("/redis", middleware.RedisRateLimiter(), func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})

	app.Get("/fiber", middleware.FiberLimiter(), func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	// Close redis connection
	redisCli.Close()
	fmt.Println("Fiber was successful shutdown.")
}
