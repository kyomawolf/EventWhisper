package events

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kyomawolf/EventWhisper/whisper-core/internal/configuration"
	log "github.com/sirupsen/logrus"
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
		log.Errorf("Error getting events: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			log.Errorf("Error writing response: %v", e)
		}
		return
	}

	json, err := json.Marshal(events)
	if err != nil {
		log.Errorf("Error marshalling json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			log.Errorf("Error writing response: %v", e)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		log.Errorf("Error writing response: %v", err)
	}
}

func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	log.Debug("The GetEvent handler is executing!")

	vars := mux.Vars(r)
	eventId := vars["eventid"]

	event, err := h.Store.ReadEvent(eventId)
	if err != nil {
		log.Errorf("Error getting Event: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			log.Errorf("Error writing response: %v", e)
		}
		return
	}

	json, err := json.Marshal(event)
	if err != nil {
		log.Errorf("Error marshalling json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			log.Errorf("Error writing response: %v", e)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		log.Errorf("Error writing response: %v", err)
	}
}

func (h *EventHandler) PostEvent(w http.ResponseWriter, r *http.Request) {
	log.Debug("The PostEvent handler is executing!")

	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Errorf("Error decoding json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_, e := w.Write([]byte("Bad request"))
		if e != nil {
			log.Errorf("Error writing response: %v", e)
		}
		return
	}

	insertedEvent, err := h.Store.InsertEvent(event)
	if err != nil {
		log.Errorf("Error inserting Event: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			log.Errorf("Error writing response: %v", e)
		}
		return
	}

	json, err := json.Marshal(insertedEvent)
	if err != nil {
		log.Errorf("Error marshalling json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			log.Errorf("Error writing response: %v", e)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		log.Errorf("Error writing response: %v", err)
	}
}
