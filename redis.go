package ratelimit

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client redis.UniversalClient
}

const (
	increaseLuaScript = `
	local key = KEYS[1]
	local ttl = tonumber(ARGV[1])
	local value = redis.call('INCR', key)

	if  value == 1
	then 
		redis.call('EXPIRE', key, ARGV[1])
		return {value, ttl}
	end 

	ttl = redis.call('TTL', key)
	return {value, ttl}
	`
)

func NewRedisStore(client ...redis.UniversalClient) RedisStore {
	if len(client) < 1 {
		cli := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
		return RedisStore{
			client: cli,
		}
	}

	return RedisStore{
		client: client[0],
	}
}

// Increment increments the limit by given count & gives back the new limit for given identifier
func (store RedisStore) Increment(ctx context.Context, key string, expired time.Duration) (*Result, error) {
	result := store.client.Eval(ctx, increaseLuaScript, []string{key}, expired.Seconds())

	count, ttl, err := parseCountAndTTL(result)
	if err != nil {
		return nil, err
	}

	return &Result{
		Count: count,
		TTL:   ttl,
	}, nil
}

// parseCountAndTTL parse count and ttl from lua script output.
func parseCountAndTTL(cmd *redis.Cmd) (count int64, ttl int64, err error) {
	result, err := cmd.Result()
	if err != nil {
		return 0, 0, err
	}

	fields, ok := result.([]interface{})
	if !ok || len(fields) != 2 {
		return 0, 0, errors.New("two elements in result were expected")
	}

	count, ok1 := fields[0].(int64)
	ttl, ok2 := fields[1].(int64)
	if !ok1 || !ok2 {
		return 0, 0, errors.New("type of the count and/or ttl should be number")
	}

	return count, ttl, nil
}
