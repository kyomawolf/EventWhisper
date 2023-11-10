package identities

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kyomawolf/EventWhisper/whisper-core/internal/configuration"
	log "github.com/sirupsen/logrus"
)

type IdentityStore struct {
	Config   *configuration.Config
	DBClient *mongo.Client
}

func NewIdentityStore(config *configuration.Config) (*IdentityStore, error) {
	log.Info("Creating Identity store")

	log.Debug("Connecting to database")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.DBConnection))
	if err != nil {
		return nil, err
	}

	log.Debug("Reading all collections")
	collections, err := client.Database(config.DatabaseName).ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		log.Warn("Could not read collections")
		return nil, err
	}

	log.Debug("Checking if Identitys collection exists")
	collectionExists := false
	for _, element := range collections {
		if element == "identities" {
			collectionExists = true
		}
	}

	if !collectionExists {
		log.Debug("Creating Identitys collection")
		e := client.Database(config.DatabaseName).CreateCollection(context.Background(), "Identitys")
		if e != nil {
			log.Warn("Could not create Identitys collection")
			return nil, err
		}
	}

	err = client.Disconnect(context.Background())
	if err != nil {
		log.Warn("Could not disconnect from database")
		return nil, err
	}

	return &IdentityStore{
		Config: config,
	}, nil
}

func (s *IdentityStore) InsertIdentity(ctx context.Context, Identity Identity) (*Identity, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(s.Config.DBConnection))
	if err != nil {
		log.Warn("Could not connect to database")
		return nil, err
	}

	collection := client.Database("stopmotion").Collection("Identitys")
	_, err = collection.InsertOne(ctx, Identity)
	if err != nil {
		log.Warn("Could not insert Identity")
		return nil, err
	}

	err = client.Disconnect(ctx)
	if err != nil {
		log.Warn("Could not disconnect from database")
		return nil, err
	}

	return &Identity, nil
}

// func (s *IdentityStore) GetAllIdentitys(ctx context.Context) ([]Identity, error) {
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(s.Config.DBConnection))
// 	if err != nil {
// 		return []Identity{}, err
// 	}

// 	collection := client.Database("stopmotion").Collection("Identitys")
// 	opts := options.Find().SetProjection(bson.D{bson.E{Key: "images", Value: 0}})
// 	cursor, err := collection.Find(ctx, bson.M{}, opts)
// 	if err != nil {
// 		return []Identity{}, err
// 	}

// 	var m []Identity
// 	for cursor.Next(ctx) {
// 		var model Identity
// 		e := cursor.Decode(&model)
// 		if e != nil {
// 			return []Identity{}, err
// 		}

// 		m = append(m, model)
// 	}

// 	err = client.Disconnect(ctx)
// 	if err != nil {
// 		log.Warn("Could not disconnect from database")
// 		return []Identity{}, err
// 	}

// 	return m, nil
// }

func (s *IdentityStore) GetIdentity(ctx context.Context, id string) (*Identity, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(s.Config.DBConnection))
	if err != nil {
		return nil, err
	}

	collection := client.Database("stopmotion").Collection("Identitys")
	var model Identity
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&model)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	err = client.Disconnect(ctx)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
