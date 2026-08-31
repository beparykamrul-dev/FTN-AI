package httpx

import (
	"context"
	"net/http"
)

type correlationKey struct{}

func CorrelationID(r *http.Request) string {
	return r.Header.Get(RequestIDHeader)
}

func WithCorrelationID(r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), correlationKey{}, CorrelationID(r)))
}

func Correlation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, WithCorrelationID(r))
	})
}
