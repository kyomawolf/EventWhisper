package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/EventWhisper/EventWhisper/whisper-core/internal/api/middlewares"
	"github.com/EventWhisper/EventWhisper/whisper-core/internal/configuration"
	"github.com/EventWhisper/EventWhisper/whisper-core/internal/events"
	"github.com/EventWhisper/EventWhisper/whisper-core/internal/identities"
	"github.com/EventWhisper/EventWhisper/whisper-core/internal/notifications"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Config          *configuration.Config
	IdentityHandler *identities.IdentityHandler
	EventHandler    *events.EventHandler
	NotifyHandler   *notifications.NotificationHandler
}

func NewServer(config *configuration.Config) (*Server, error) {

	identityStore, err := identities.NewIdentityStore(config)
	if err != nil {
		return nil, err
	}

	eventsStore, err := events.NewEventStore(config)
	if err != nil {
		return nil, err
	}

	return &Server{
		Config:          config,
		IdentityHandler: identities.NewIdentityHandler(config, identityStore),
		EventHandler:    events.NewEventHandler(config, eventsStore),
		NotifyHandler:   notifications.NewNotificationHandler(config, eventsStore, identityStore),
	}, nil
}

func (s *Server) Start() error {

	r := chi.NewRouter()

	r.Use(middlewares.Logger(s.Config))
	r.Use(middlewares.Cors(s.Config))
	r.Use(middlewares.Authorization(s.Config))

	r.Options(fmt.Sprintf("%v/identity", s.Config.BasePath), options)
	r.Post(fmt.Sprintf("%v/identity", s.Config.BasePath), s.IdentityHandler.PostIdentity)
	r.Get(fmt.Sprintf("%v/identity", s.Config.BasePath), s.IdentityHandler.GetAllIdentities)
	r.Options(fmt.Sprintf("%v/identity/{sub}", s.Config.BasePath), options)
	r.Get(fmt.Sprintf("%v/identity/{sub}", s.Config.BasePath), s.IdentityHandler.GetIdentity)

	r.Options(fmt.Sprintf("%v/events", s.Config.BasePath), options)
	r.Post(fmt.Sprintf("%v/events", s.Config.BasePath), s.EventHandler.PostEvent)
	r.Get(fmt.Sprintf("%v/events", s.Config.BasePath), s.EventHandler.GetAllEvents)
	r.Options(fmt.Sprintf("%v/events/{eventid}", s.Config.BasePath), options)
	r.Get(fmt.Sprintf("%v/events/{eventid}", s.Config.BasePath), s.EventHandler.GetEvent)

	r.Options(fmt.Sprintf("%v/notify", s.Config.BasePath), options)
	r.Get(fmt.Sprintf("%v/notify", s.Config.BasePath), s.NotifyHandler.GetNotification)
	r.Post(fmt.Sprintf("%v/notify", s.Config.BasePath), s.NotifyHandler.PostNotification)

	r.NotFound(http.HandlerFunc(notFoundHandler))

	server := &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%v", s.Config.Port),
		ReadHeaderTimeout: 3 * time.Second, //nolint:gomnd // 3 seconds is a reasonable timeout
		Handler:           r,
	}

	slog.Info("starting server", "port", s.Config.Port)

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	_, err := w.Write([]byte("Not found."))
	if err != nil {
		slog.ErrorContext(r.Context(), "Error writing response", "error", err)
	}
}

func options(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "text")
	_, err := w.Write(nil)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error writing response", "error", err)
	}
}
