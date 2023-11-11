package notifications

import (
	"github.com/kyomawolf/EventWhisper/whisper-core/internal/configuration"
	"github.com/kyomawolf/EventWhisper/whisper-core/internal/events"
	"github.com/kyomawolf/EventWhisper/whisper-core/internal/identities"
)

type NotificationHandler struct {
	Config        *configuration.Config
	EventStore    *events.EventStore
	IdentityStore *identities.IdentityStore
}

func NewNotificationHandler(config *configuration.Config, eventStore *events.EventStore, identityStore *identities.IdentityStore) *NotificationHandler {
	return &NotificationHandler{
		Config:        config,
		EventStore:    eventStore,
		IdentityStore: identityStore,
	}
}

// func (h *NotificationHandler) PostNotification(w http.ResponseWriter, r *http.Request) {

// 	ids, err := h.IdentityStore.ReadAllIdentities()
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		_, e := w.Write([]byte("Internal server error"))
// 		if e != nil {
// 			return
// 		}
// 		return
// 	}

// }
