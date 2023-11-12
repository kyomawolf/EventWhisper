package notifications

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/kyomawolf/EventWhisper/whisper-core/internal/configuration"
	"github.com/kyomawolf/EventWhisper/whisper-core/internal/events"
	"github.com/kyomawolf/EventWhisper/whisper-core/internal/identities"
	log "github.com/sirupsen/logrus"
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

type TouchpointMsgModel struct {
	Identity identities.Identity `json:"identity"`
	Msg      string              `json:"msg"`
}

func (h *NotificationHandler) SendMsg(identity identities.Identity, msg string) error {
	body := TouchpointMsgModel{
		Identity: identity,
		Msg:      msg,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := "https://touchpoints.eventwhisper.de/telegram/sendmsg"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ToDo")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (h *NotificationHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	log.Debug("Running GetNotification")

	identities, err := h.IdentityStore.ReadAllIdentities()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			return
		}
		return
	}

	events, err := h.EventStore.ReadAllEvents()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write([]byte("Internal server error"))
		if e != nil {
			return
		}
		return
	}

	log.Debugf("Found %v identities", len(identities))
	log.Debugf("Found %v events", len(events))

	for _, identity := range identities {

		log.Debugf("Sending notification to %v", identity.Name)
		log.Debugf("Interests: %v", identity.Interests)

		msg := "Hello " + identity.Name + ",\n"
		msg += "Wir haben ein paar spannende Events für dich. Eventuell ist ja etwas dabei, worauf du Lust hast."
		h.SendMsg(identity, msg)

		for _, event := range events {
			msgEvent := event.Title + "\n"
			msgEvent += "am " + event.StartTime + "\n\n"
			msgEvent += event.Description + "\n\n"
			msgEvent += event.Url + "\n"

			log.Debugf("Sending event %v", event.Title)
			h.SendMsg(identity, msgEvent)
		}

	}
}
