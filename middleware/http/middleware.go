package http

import (
	"net/http"
	"strconv"

	limiter "github.com/trinhdaiphuc/rate-limit"
)

// Middleware is the middleware for basic http.Handler.
type Middleware struct {
	Limiter        *limiter.Limiter
	OnError        ErrorHandler
	OnLimitReached LimitReachedHandler
	KeyGetter      KeyGetter
	ExcludedKey    func(string) bool
}

// NewMiddleware return a new instance of a basic HTTP middleware.
func NewMiddleware(options ...Option) *Middleware {
	middleware := &Middleware{
		Limiter:        limiter.New(),
		OnError:        DefaultErrorHandler,
		OnLimitReached: DefaultLimitReachedHandler,
		KeyGetter:      DefaultKeyGetter(),
		ExcludedKey:    nil,
	}

	for _, option := range options {
		option(middleware)
	}

	return middleware
}

func (middleware *Middleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := middleware.KeyGetter(r)
		if middleware.ExcludedKey != nil && middleware.ExcludedKey(key) {
			h.ServeHTTP(w, r)
			return
		}

		context, err := middleware.Limiter.Process(r.Context(), key)
		if err != nil {
			middleware.OnError(w, r, err)
			return
		}

		w.Header().Add(limiter.XRateLimitLimit, strconv.FormatInt(context.Limit, 10))
		w.Header().Add(limiter.XRateLimitRemaining, strconv.FormatInt(context.Remaining, 10))
		w.Header().Add(limiter.XRateLimitReset, strconv.FormatInt(context.Reset, 10))

		if context.Reached {
			middleware.OnLimitReached(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}
