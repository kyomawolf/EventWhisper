package events

import (
	"os"

	"encoding/json"
	"io/ioutil"

	"github.com/EventWhisper/EventWhisper/whisper-core/internal/configuration"
	"github.com/google/uuid"
)

type EventStore struct {
	Config *configuration.Config
	Events []Event `json:"events"`
}

func NewEventStore(config *configuration.Config) (*EventStore, error) {
	store := &EventStore{
		Config: config,
		Events: []Event{},
	}

	path := config.DBFilePath + "/events.json"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, return empty store
			return store, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, &store.Events)
	if err != nil {
		return nil, err
	}

	return store, nil

}

type EventInsertError string

var (
	ErrEventAlreadyExists EventInsertError = "Event already exists"
	ErrGeneralInsertError EventInsertError = "General insert error"
)

func (s *EventStore) SaveDataToJsonFile() error {
	path := s.Config.DBFilePath + "/events.json"

	// Create directory if it does not exist
	err := os.MkdirAll(s.Config.DBFilePath, 0755)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(s.Events)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *EventStore) InsertEvent(event Event) (*Event, *EventInsertError) {

	for _, m := range s.Events {
		if m.Url == event.Url {
			return &m, &ErrEventAlreadyExists
		}
	}

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

type EventMatches struct {
	Event      Event
	MatchCount int
}
