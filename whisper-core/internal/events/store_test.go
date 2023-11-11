package events

import (
	"testing"

	"github.com/kyomawolf/EventWhisper/whisper-core/internal/configuration"
	"github.com/kyomawolf/EventWhisper/whisper-core/internal/identities"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	store, err := NewEventStore(&configuration.Config{})
	assert.Nil(t, err)

	events := []Event{
		{
			ID: "1",
			Interest: []string{
				"foo",
				"bar",
			},
			Url: "http://example.com/1",
		},
		{
			ID: "2",
			Interest: []string{
				"foo",
				"baz",
			},
			Url: "http://example.com/2",
		},
		{
			ID: "3",
			Interest: []string{
				"foo",
				"bar",
				"baz",
			},
			Url: "http://example.com/3",
		},
	}

	for _, event := range events {
		_, err := store.InsertEvent(event)
		assert.Nil(t, err)
	}

	identity := identities.Identity{
		Interest: []string{
			"foo",
			"bar",
		},
	}

	matches, _ := store.CreateMatches(events[0], identity.Interest)
	assert.Equal(t, matches, 2)

	matches2, _ := store.CreateMatches(events[1], identity.Interest)
	assert.Equal(t, matches2, 1)

	matches3, _ := store.CreateMatches(events[2], identity.Interest)
	assert.Equal(t, matches3, 2)
}

func TestFindOrdered(t *testing.T) {
	store, err := NewEventStore(&configuration.Config{})
	assert.Nil(t, err)

	events := []Event{
		{
			ID: "1",
			Interest: []string{
				"foo",
				"bar",
			},
			Url: "http://example.com/1",
		},
		{
			ID: "2",
			Interest: []string{
				"food",
				"baz",
			},
			Url: "http://example.com/2",
		},
		{
			ID: "3",
			Interest: []string{
				"foo",
				"bat",
				"baz",
			},
			Url: "http://example.com/3",
		},
	}

	for _, event := range events {
		_, err := store.InsertEvent(event)
		assert.Nil(t, err)
	}

	identity := identities.Identity{
		Interest: []string{
			"foo",
			"bar",
		},
	}

	ordered, _ := store.FindBestMatches(identity.Interest)
	assert.Equal(t, ordered[0].ID, "2")
}
