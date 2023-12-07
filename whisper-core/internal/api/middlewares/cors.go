package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/EventWhisper/EventWhisper/whisper-core/internal/configuration"
)

func Cors(config *configuration.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.DebugContext(r.Context(), "The Cors middleware is executing!")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			next.ServeHTTP(w, r)
		})
	}
}
