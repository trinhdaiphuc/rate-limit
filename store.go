package ratelimit

import (
	"context"
	"time"
)

// Store is the common interface for limiter stores.
type Store interface {
	// Increment increments the limit by given count & gives back the new limit for given identifier
	Increment(ctx context.Context, key string, expired time.Duration) (*Result, error)
}

// StoreOptions are options for store.
type StoreOptions struct {
	// Prefix is the prefix to use for the key.
	Prefix string

	// MaxRetry is the maximum number of retry under race conditions on redis store.
	// Deprecated: this option is no longer required since all operations are atomic now.
	MaxRetry int

	// CleanUpInterval is the interval for cleanup (run garbage collection) on stale entries on memory store.
	// Setting this to a low value will optimize memory consumption, but will likely
	// reduce performance and increase lock contention.
	// Setting this to a high value will maximum throughput, but will increase the memory footprint.
	CleanUpInterval time.Duration
}

type Result struct {
	Count int64
	TTL   int64
}
