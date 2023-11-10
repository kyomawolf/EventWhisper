package middlewares

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type CorsMiddleware struct{}

func (*CorsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Info("The Cors middleware is executing!")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	next.ServeHTTP(w, r)
}
