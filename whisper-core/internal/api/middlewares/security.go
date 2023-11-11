package middlewares

import (
	"net/http"
	"strings"

	"github.com/kyomawolf/EventWhisper/whisper-core/internal/configuration"
	log "github.com/sirupsen/logrus"
)

type SecurityMiddleware struct {
	Config *configuration.Config
}

func NewSecurityMiddleware(config *configuration.Config) *SecurityMiddleware {
	return &SecurityMiddleware{
		Config: config,
	}
}

func (m *SecurityMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Debug("The logger middleware is executing!")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Debug("No Authorization header found")
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Unauthorized"))
		if err != nil {
			log.Errorf("Error writing response: %v", err)
		}
		return
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 {
		log.Debug("Invalid Authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Unauthorized"))
		if err != nil {
			log.Errorf("Error writing response: %v", err)
		}
		return
	}

	if authHeaderParts[0] != "Bearer" {
		log.Debug("Invalid Authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Unauthorized"))
		if err != nil {
			log.Errorf("Error writing response: %v", err)
		}
		return
	}

	if authHeaderParts[1] != m.Config.ApiKey {
		log.Debug("Invalid Authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Unauthorized"))
		if err != nil {
			log.Errorf("Error writing response: %v", err)
		}
		return
	}

	next.ServeHTTP(w, r)
}
