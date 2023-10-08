package http

import (
	"net/http"

	limiter "github.com/trinhdaiphuc/rate-limit"
)

type Option func(*Middleware)

// ErrorHandler is an handler used to inform when an error has occurred.
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

func WithLimiter(limiter *limiter.Limiter) Option {
	return func(middleware *Middleware) {
		middleware.Limiter = limiter
	}
}

// WithErrorHandler will configure the Middleware to use the given ErrorHandler.
func WithErrorHandler(handler ErrorHandler) Option {
	return func(middleware *Middleware) {
		middleware.OnError = handler
	}
}

// DefaultErrorHandler is the default ErrorHandler used by a new Middleware.
func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	panic(err)
}

// LimitReachedHandler is an handler used to inform when the limit has exceeded.
type LimitReachedHandler func(w http.ResponseWriter, r *http.Request)

// WithLimitReachedHandler will configure the Middleware to use the given LimitReachedHandler.
func WithLimitReachedHandler(handler LimitReachedHandler) Option {
	return func(middleware *Middleware) {
		middleware.OnLimitReached = handler
	}
}

// DefaultLimitReachedHandler is the default LimitReachedHandler used by a new Middleware.
func DefaultLimitReachedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Limit exceeded", http.StatusTooManyRequests)
}

// KeyGetter will define the rate limiter key given the gin Context.
type KeyGetter func(r *http.Request) string

// WithKeyGetter will configure the Middleware to use the given KeyGetter.
func WithKeyGetter(handler KeyGetter) Option {
	return func(middleware *Middleware) {
		middleware.KeyGetter = handler
	}
}

// DefaultKeyGetter is the default KeyGetter used by a new Middleware.
// It returns the Client IP address.
func DefaultKeyGetter() func(r *http.Request) string {
	return func(r *http.Request) string {
		IPAddress := r.Header.Get("X-Real-Ip")
		if IPAddress == "" {
			IPAddress = r.Header.Get("X-Forwarded-For")
		}
		if IPAddress == "" {
			IPAddress = r.RemoteAddr
		}
		return IPAddress
	}
}

// WithExcludedKey will configure the Middleware to ignore key(s) using the given function.
func WithExcludedKey(handler func(string) bool) Option {
	return func(middleware *Middleware) {
		middleware.ExcludedKey = handler
	}
}
