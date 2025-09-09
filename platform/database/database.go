package database

import (
	"context"
	"strconv"

	"github.com/aslon1213/g4h_pos_erp/pkg/configs"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewDB() *mongo.Client {
	log.Debug().Msg("Initializing database connection")

	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	uri := ""
	if config.DB.Auth {
		uri = config.DB.URL
	} else {
		uri = config.DB.URL
	}

	// uri := "mongodb://localhost:27017/?replicaSet=myReplicaSet"

	log.Debug().Str("uri", uri).Str("max_connections", strconv.FormatUint(config.DB.MaxConnections, 10)).Str("min_pool_size", strconv.FormatUint(config.DB.MinPoolSize, 10)).Msg("Connecting to MongoDB")

	client, err := mongo.Connect(options.Client().ApplyURI(uri).SetMaxConnecting(config.DB.MaxConnections).SetMinPoolSize(config.DB.MinPoolSize))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	// Check if the database is connected
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ping MongoDB")
	}

	log.Info().Msg("Successfully connected to MongoDB")
	return client
}

func NewSession(client *mongo.Client) (*mongo.Session, error) {
	session, err := client.StartSession()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start session")
		return nil, err
	}
	return session, nil
}

func NewTransaction(ctx context.Context, session *mongo.Session) error {
	if err := session.StartTransaction(); err != nil {
		log.Error().Err(err).Msg("Failed to start transaction")
		return err
	}
	return nil
}

func StartTransaction(client *mongo.Client) (*mongo.Session, context.Context, error) {
	session, err := NewSession(client)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start transaction")
		return nil, nil, err
	}
	ctx := mongo.NewSessionContext(context.Background(), session)
	if err := session.StartTransaction(); err != nil {
		log.Error().Err(err).Msg("Failed to start transaction")
		return nil, nil, err
	}
	return session, ctx, nil
}
