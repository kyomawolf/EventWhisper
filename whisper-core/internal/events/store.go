package events

import (
	"github.com/google/uuid"
	"github.com/kyomawolf/EventWhisper/whisper-core/internal/configuration"
)

type EventStore struct {
	Config *configuration.Config
	Events []Event `json:"events"`
}

func NewEventStore(config *configuration.Config) (*EventStore, error) {
	return &EventStore{
		Config: config,
		Events: []Event{},
	}, nil
}

func (s *EventStore) InsertEvent(event Event) (*Event, error) {
	event.ID = uuid.New().String()
	s.Events = append(s.Events, event)

	return &event, nil
}

func (s *EventStore) ReadAllEvents() ([]Event, error) {
	return s.Events, nil
}

func (s *EventStore) ReadEvent(id string) (*Event, error) {
	var model Event

	for _, m := range s.Events {
		if m.ID == id {
			return &m, nil
		}
	}

	return &model, nil
}
