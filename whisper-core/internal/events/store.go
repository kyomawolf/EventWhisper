package events

import (
	"os"
	"sort"

	"encoding/json"
	"io/ioutil"

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

func (s *EventStore) CreateMatches(event Event, interests []string) (int, error) {
	matches := 0

	for _, i := range interests {
		for _, e := range event.Interest {
			if e == i {
				matches++
			}
		}
	}

	return matches, nil
}

func (s *EventStore) FindBestMatches(interests []string) ([]Event, error) {

	var bestMatches []EventMatches

	for _, e := range s.Events {
		matchCount, err := s.CreateMatches(e, interests)
		if err != nil {
			return nil, err
		}

		match := EventMatches{
			Event:      e,
			MatchCount: matchCount,
		}

		bestMatches = append(bestMatches, match)
	}

	sort.Slice(bestMatches, func(i, j int) bool {
		return bestMatches[i].MatchCount > bestMatches[j].MatchCount
	})

	return []Event{bestMatches[0].Event, bestMatches[1].Event, bestMatches[2].Event}, nil
}
