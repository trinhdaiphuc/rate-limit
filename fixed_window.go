package ratelimit

import (
	"context"
	"time"
)

type FixedWindow struct {
	Limit      int64
	Expiration time.Duration
}

func NewFixedWindow() FixedWindow {
	return FixedWindow{}
}

func (algo FixedWindow) Process(store Store, ctx context.Context, key string) (*Context, error) {
	result, err := store.Increment(ctx, key, algo.Expiration)
	if err != nil {
		return nil, err
	}

	var (
		now        = time.Now()
		remaining  = int64(0)
		reached    = true
		expiration = now.Add(time.Duration(result.TTL) * time.Second)
	)

	if result.Count < algo.Limit {
		remaining = algo.Limit - result.Count
		reached = false
	}

	return &Context{
		Limit:     algo.Limit,
		Reset:     expiration.Unix(),
		Remaining: remaining,
		Reached:   reached,
	}, nil
}
