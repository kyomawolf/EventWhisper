package identities

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kyomawolf/EventWhisper/whisper-core/internal/configuration"

	log "github.com/sirupsen/logrus"
)

type IdentityHandler struct {
	Config *configuration.Config
	Store  *IdentityStore
}

func NewIdentityHandler(config *configuration.Config, store *IdentityStore) *IdentityHandler {
	return &IdentityHandler{
		Config: config,
		Store:  store,
	}
}

func (h *IdentityHandler) GetIdentity(w http.ResponseWriter, r *http.Request) {
	log.Debug("The GetIdentity handler is executing!")

	vars := mux.Vars(r)
	sub := vars["sub"]

	identity, err := h.Store.GetIdentity(r.Context(), sub)
	if err != nil {
		log.Errorf("Error getting Identity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			log.Errorf("Error writing response: %v", e)
		}
		return
	}

	json, err := json.Marshal(identity)
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

func (h *IdentityHandler) GetAllIdentities(w http.ResponseWriter, r *http.Request) {
	log.Debug("The GetAllIdentities handler is executing!")

	identities, err := h.Store.ReadAllIdentities()
	if err != nil {
		log.Errorf("Error getting Identities: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			log.Errorf("Error writing response: %v", e)
		}
		return
	}

	json, err := json.Marshal(identities)
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

func (h *IdentityHandler) PostIdentity(w http.ResponseWriter, r *http.Request) {
	log.Info("The PostIdentity handler is executing!")

	var identity Identity

	err := json.NewDecoder(r.Body).Decode(&identity)
	if err != nil {
		log.Errorf("Error decoding json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_, e := w.Write([]byte("Bad request"))
		if e != nil {
			log.Errorf("Error writing response: %v", e)
		}
		return
	}

	insertedIdentity, err := h.Store.InsertIdentity(r.Context(), identity)
	if err != nil {
		log.Errorf("Error inserting Identity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			log.Errorf("Error writing response: %v", e)
		}
		return
	}

	log.Debug("Saving data to json file")
	h.Store.SaveDataToJsonFile()

	json, err := json.Marshal(insertedIdentity)
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
