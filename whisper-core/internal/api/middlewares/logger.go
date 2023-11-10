package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/kyomawolf/EventWhisper/whisper-core/internal/configuration"
	log "github.com/sirupsen/logrus"
)

type LoggerMiddleware struct {
	Config *configuration.Config
}

func NewLoggerMiddleware(config *configuration.Config) *LoggerMiddleware {
	return &LoggerMiddleware{
		Config: config,
	}
}

func (m *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Debug("The logger middleware is executing!")

	log.Debugf("LogLevel: %v", m.Config.LogLevel)
	if strings.ToLower(m.Config.LogLevel) == "debug" {
		for name, values := range r.Header {
			for _, value := range values {
				log.Debugf("%v : %v", name, value)
			}
		}
	}

	t := time.Now()
	next.ServeHTTP(w, r)
	log.Debugf("Execution time: %s ", time.Since(t).String())
}
