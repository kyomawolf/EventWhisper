package identities

import (
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/EventWhisper/EventWhisper/whisper-core/internal/configuration"
	"github.com/gorilla/mux"
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
	ctx := r.Context()
	slog.ErrorContext(ctx, "The GetIdentity handler is executing!")

	vars := mux.Vars(r)
	sub := vars["sub"]

	identity, err := h.Store.GetIdentity(ctx, sub)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting Identity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			slog.ErrorContext(ctx, "Error writing response: %v", e)
		}
		return
	}

	json, err := json.Marshal(identity)
	if err != nil {
		slog.ErrorContext(ctx, "Error marshalling json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			slog.ErrorContext(ctx, "Error writing response: %v", e)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		slog.ErrorContext(ctx, "Error writing response: %v", err)
	}
}

func (h *IdentityHandler) GetAllIdentities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.ErrorContext(ctx, "The GetAllIdentities handler is executing!")

	identities, err := h.Store.ReadAllIdentities()
	if err != nil {
		slog.ErrorContext(ctx, "Error getting Identities: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			slog.ErrorContext(ctx, "Error writing response: %v", e)
		}
		return
	}

	json, err := json.Marshal(identities)
	if err != nil {
		slog.ErrorContext(ctx, "Error marshalling json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			slog.ErrorContext(ctx, "Error writing response: %v", e)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		slog.ErrorContext(ctx, "Error writing response: %v", err)
	}
}

func (h *IdentityHandler) PostIdentity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "The PostIdentity handler is executing!")

	var identity Identity

	err := json.NewDecoder(r.Body).Decode(&identity)
	if err != nil {
		slog.ErrorContext(ctx, "Error decoding json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_, e := w.Write([]byte("Bad request"))
		if e != nil {
			slog.ErrorContext(ctx, "Error writing response: %v", e)
		}
		return
	}

	insertedIdentity, err := h.Store.InsertIdentity(ctx, identity)
	if err != nil {
		slog.ErrorContext(ctx, "Error inserting Identity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			slog.ErrorContext(ctx, "Error writing response: %v", e)
		}
		return
	}

	slog.ErrorContext(ctx, "Saving data to json file")
	h.Store.SaveDataToJsonFile()

	json, err := json.Marshal(insertedIdentity)
	if err != nil {
		slog.ErrorContext(ctx, "Error marshalling json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			slog.ErrorContext(ctx, "Error writing response: %v", e)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		slog.ErrorContext(ctx, "Error writing response: %v", err)
	}
}
