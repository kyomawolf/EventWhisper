package events

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/EventWhisper/EventWhisper/whisper-core/internal/configuration"
	"github.com/gorilla/mux"
)

type EventHandler struct {
	Config *configuration.Config
	Store  *EventStore
}

func NewEventHandler(config *configuration.Config, store *EventStore) *EventHandler {
	return &EventHandler{
		Config: config,
		Store:  store,
	}
}

func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.Store.ReadAllEvents()
	if err != nil {
		slog.ErrorContext(r.Context(), "Error getting events", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			slog.ErrorContext(r.Context(), "Error writing response", "error", e)
		}
		return
	}

	json, err := json.Marshal(events)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error marshalling json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			slog.ErrorContext(r.Context(), "Error writing response", "error", e)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error writing response", "error", err)
	}
}

func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	slog.DebugContext(r.Context(), "The GetEvent handler is executing!")

	vars := mux.Vars(r)
	eventId := vars["eventid"]

	event, err := h.Store.ReadEvent(eventId)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error getting Event", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			slog.ErrorContext(r.Context(), "Error writing response", "error", e)
		}
		return
	}

	json, err := json.Marshal(event)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error marshalling json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			slog.ErrorContext(r.Context(), "Error writing response", "error", e)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error writing response", "error", err)
	}
}

func (h *EventHandler) PostEvent(w http.ResponseWriter, r *http.Request) {
	slog.DebugContext(r.Context(), "The PostEvent handler is executing!")

	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error decoding json", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		_, e := w.Write([]byte("Bad request"))
		if e != nil {
			slog.ErrorContext(r.Context(), "Error writing response", "error", e)
		}
		return
	}

	insertedEvent, errInsert := h.Store.InsertEvent(event)
	if errInsert == &ErrEventAlreadyExists {
		slog.ErrorContext(r.Context(), "Error inserting Event", "error", err)
		w.WriteHeader(http.StatusConflict)
		_, e := w.Write([]byte("Event already exists"))
		if e != nil {
			slog.ErrorContext(r.Context(), "Error writing response", "error", e)
		}
		return
	}
	if errInsert != nil {
		slog.ErrorContext(r.Context(), "Error inserting Event", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			slog.ErrorContext(r.Context(), "Error writing response", "error", e)
		}
		return
	}

	h.Store.SaveDataToJsonFile()

	json, err := json.Marshal(insertedEvent)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error marshalling json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			slog.ErrorContext(r.Context(), "Error writing response", "error", e)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error writing response", "error", err)
	}
}
