package identities

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log/slog"

	"os"

	"github.com/EventWhisper/EventWhisper/whisper-core/internal/configuration"
	"github.com/google/uuid"
)

type IdentityStore struct {
	Config     *configuration.Config
	Identities []Identity
}

func NewIdentityStore(config *configuration.Config) (*IdentityStore, error) {
	slog.Info("Creating Identity store")

	store := &IdentityStore{
		Config:     config,
		Identities: []Identity{},
	}

	path := config.DBFilePath + "/identities.json"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, return empty store
			return store, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, &store.Identities)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *IdentityStore) SaveDataToJsonFile() error {
	path := s.Config.DBFilePath + "/identities.json"

	// Create directory if it does not exist
	err := os.MkdirAll(s.Config.DBFilePath, 0755)
	if err != nil {
		return err
	}

	slog.Debug(path)

	jsonData, err := json.Marshal(s.Identities)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *IdentityStore) InsertIdentity(ctx context.Context, identity Identity) (*Identity, error) {

	identity.Sub = uuid.New().String()
	s.Identities = append(s.Identities, identity)

	return &identity, nil
}

func (s *IdentityStore) GetIdentity(ctx context.Context, sub string) (*Identity, error) {

	var model Identity

	for _, m := range s.Identities {
		if m.Sub == sub {
			return &m, nil
		}
	}

	return &model, nil
}

func (s *IdentityStore) ReadAllIdentities() ([]Identity, error) {
	return s.Identities, nil
}

func (s *IdentityStore) ReadIdentityByChatId(chatId string) (*Identity, error) {
	for _, m := range s.Identities {
		if m.Channels[0].Specifics.ChatID == chatId {
			return &m, nil
		}
	}

	return nil, nil
}
