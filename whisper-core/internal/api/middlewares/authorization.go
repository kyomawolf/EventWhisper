package middlewares

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/EventWhisper/EventWhisper/whisper-core/internal/configuration"
)

func Authorization(config *configuration.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			slog.DebugContext(r.Context(), "The logger middleware is executing!")

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				slog.DebugContext(r.Context(), "No Authorization header found")
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte("Unauthorized"))
				if err != nil {
					slog.ErrorContext(r.Context(), "Error writing response: %v", err)
				}
				return
			}

			authHeaderParts := strings.Split(authHeader, " ")
			if len(authHeaderParts) != 2 {
				slog.DebugContext(r.Context(), "Invalid Authorization header")
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte("Unauthorized"))
				if err != nil {
					slog.ErrorContext(r.Context(), "Error writing response: %v", err)
				}
				return
			}

			if authHeaderParts[0] != "Bearer" {
				slog.DebugContext(r.Context(), "Invalid Authorization header")
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte("Unauthorized"))
				if err != nil {
					slog.ErrorContext(r.Context(), "Error writing response: %v", err)
				}
				return
			}

			if authHeaderParts[1] != config.ApiKey {
				slog.DebugContext(r.Context(), "Invalid Authorization header")
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte("Unauthorized"))
				if err != nil {
					slog.ErrorContext(r.Context(), "Error writing response: %v", err)
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
