package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kyomawolf/EventWhisper/whisper-core/internal/api/middlewares"
	"github.com/kyomawolf/EventWhisper/whisper-core/internal/configuration"
	"github.com/kyomawolf/EventWhisper/whisper-core/internal/identities"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

type Server struct {
	Config          *configuration.Config
	Router          *negroni.Negroni
	IdentityHandler *identities.IdentityHandler
}

func NewServer(config *configuration.Config) (*Server, error) {

	identityStore, err := identities.NewIdentityStore(config)
	if err != nil {
		return nil, err
	}

	return &Server{
		Config:          config,
		Router:          nil,
		IdentityHandler: identities.NewIdentityHandler(config, identityStore),
	}, nil
}

func (s *Server) ConfigureRouter() error {
	router := mux.NewRouter()

	router.Path(fmt.Sprintf("%v/identity", s.Config.BasePath)).HandlerFunc(options).Methods("OPTIONS")
	router.Path(fmt.Sprintf("%v/identity", s.Config.BasePath)).HandlerFunc(s.IdentityHandler.PostIdentity).Methods("POST")

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	routerMiddleware := negroni.New()

	routerMiddleware.Use(&middlewares.CorsMiddleware{})
	routerMiddleware.Use(&middlewares.LoggerMiddleware{Config: s.Config})
	routerMiddleware.UseHandler(router)

	s.Router = routerMiddleware

	return nil
}

func (s *Server) Start() error {

	server := &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%v", s.Config.Port),
		ReadHeaderTimeout: 3 * time.Second, //nolint:gomnd // 3 seconds is a reasonable timeout
		Handler:           s.Router,
	}

	log.Infof("starting server on port %v", s.Config.Port)

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	_, err := w.Write([]byte("Not found."))
	if err != nil {
		log.Errorf("Error writing response: %v", err)
	}
}

func options(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "text")
	_, err := w.Write(nil)
	if err != nil {
		log.Errorf("Error writing response: %v", err)
	}
}
