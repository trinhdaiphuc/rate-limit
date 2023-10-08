package ratelimit

import (
	"context"
	"time"
)

// Limiter is the limiter instance.
type Limiter struct {
	cfg *Config
}

type Algorithm interface {
	Process(store Store, ctx context.Context, key string) (*Context, error)
}

type Rate struct {
	Period time.Duration
	Limit  int64
}

// Context is the limit context.
type Context struct {
	Limit     int64
	Remaining int64
	Reset     int64
	Reached   bool
}

const (
	// X-RateLimit-* headers
	XRateLimitLimit     = "X-RateLimit-Limit"
	XRateLimitRemaining = "X-RateLimit-Remaining"
	XRateLimitReset     = "X-RateLimit-Reset"
)

func New(options ...Option) *Limiter {
	config := configDefault()
	for _, o := range options {
		o.apply(config)
	}

	return &Limiter{
		cfg: config,
	}
}
func (limiter *Limiter) Process(ctx context.Context, key string) (*Context, error) {
	return limiter.cfg.Algorithm.Process(limiter.cfg.Store, ctx, key)
}
