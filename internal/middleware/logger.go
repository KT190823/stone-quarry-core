package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger logs all incoming HTTP requests with Method, Path, Status Code and Duration
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Colorize status code for terminal readability
		statusColor := "\033[32m" // Green for 2xx
		if rw.statusCode >= 400 && rw.statusCode < 500 {
			statusColor = "\033[33m" // Yellow for 4xx
		} else if rw.statusCode >= 500 {
			statusColor = "\033[31m" // Red for 5xx
		}
		resetColor := "\033[0m"

		fmt.Printf("[HTTP] %s%3d%s | %-6s %-35s | %s\n",
			statusColor, rw.statusCode, resetColor,
			r.Method, r.URL.Path, duration.Round(time.Microsecond),
		)
	})
}
