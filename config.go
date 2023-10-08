package ratelimit

import "time"

type Config struct {
	Store     Store
	Algorithm Algorithm
	Rate      Rate
}

// Option is used to define configuration.
type Option interface {
	apply(config *Config)
}

type option func(*Config)

func (o option) apply(config *Config) {
	o(config)
}

// WithStore will configure to use the given Store.
func WithStore(store Store) Option {
	return option(func(c *Config) {
		c.Store = store
	})
}

// WithAlgorithm will configure to use the given Algorithm.
func WithAlgorithm(algorithm Algorithm) Option {
	return option(func(c *Config) {
		c.Algorithm = algorithm
	})
}

func configDefault() *Config {
	return &Config{
		Store: NewRedisStore(),
		Algorithm: FixedWindow{
			Limit:      10,
			Expiration: time.Minute,
		},
	}
}
