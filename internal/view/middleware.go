package view

import (
	"context"
	"github.com/major1201/goutils"
	"github.com/major1201/k8s-mutator/pkg/httputils"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// TimeoutMiddleware limits the time for every request
func TimeoutMiddleware(duration time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequestIDMiddleware set the "X-Request-Id" for each request, let's you trace each request
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := goutils.UUID()
		w.Header().Set("X-Request-Id", requestID)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), httputils.CtxRequestID, requestID)))
	})
}

// DurationMiddleware prints the duration for each request
func DurationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		log().Debug("request duration", zap.String("request_id", httputils.RequestID(r)), zap.Duration("duration", time.Now().Sub(startTime)))
	})
}

// LogMiddleware prints each request message
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log().Info("http request",
			zap.String("path", r.RequestURI),
			zap.String("method", r.Method),
			zap.String("raddr", r.RemoteAddr),
			zap.String("request_id", httputils.RequestID(r)),
		)
		next.ServeHTTP(w, r)
	})
}
