package identities

import (
	"context"

	"github.com/google/uuid"
	"github.com/kyomawolf/EventWhisper/whisper-core/internal/configuration"
	log "github.com/sirupsen/logrus"
)

type IdentityStore struct {
	Config     *configuration.Config
	Identities []Identity
}

func NewIdentityStore(config *configuration.Config) (*IdentityStore, error) {
	log.Info("Creating Identity store")

	return &IdentityStore{
		Config:     config,
		Identities: []Identity{},
	}, nil
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
