package middleware

import (
	"fmt"
	"net/http"
	"queueit/pkg/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("[REQ] %s %s from %s", r.Method, r.RequestURI, r.RemoteAddr))
		next.ServeHTTP(w, r)
	})
}
