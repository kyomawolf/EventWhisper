package middlewares

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/EventWhisper/EventWhisper/whisper-core/internal/configuration"
)

func Logger(config *configuration.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.DebugContext(r.Context(), "The logger middleware is executing!")

			if strings.ToLower(config.LogLevel) == "debug" {
				for name, values := range r.Header {
					for _, value := range values {
						slog.DebugContext(r.Context(), "HEADERS", "name", name, "value", value)
					}
				}
			}

			t := time.Now()
			next.ServeHTTP(w, r)
			slog.DebugContext(r.Context(), "Execution time", "time", time.Since(t).String())
		})
	}
}
