package identities

import (
	"encoding/json"
	"net/http"
	"path"

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

func (h *IdentityHandler) GetImageAsFile(w http.ResponseWriter, r *http.Request) {
	log.Debug("The GetImages handler is executing!")

	projectID := mux.Vars(r)["project_id"]
	if projectID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, e := w.Write([]byte("project_id is required"))
		if e != nil {
			log.Errorf("Error writing response: %v", e)
		}
		return
	}

	imgID := mux.Vars(r)["img_id"]
	if imgID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, e := w.Write([]byte("id is required"))
		if e != nil {
			log.Errorf("Error writing response: %v", e)
		}
	}

	imgDir := path.Join("_images/", projectID)
	imgPath := path.Join(imgDir, imgID+".png")

	http.ServeFile(w, r, imgPath)
}

func (h *IdentityHandler) PostIdentity(w http.ResponseWriter, r *http.Request) {
	log.Debug("The PostIdentity handler is executing!")

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
