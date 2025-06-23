package database

import (
	"aslon1213/magazin_pos/pkg/configs"
	"context"
	"fmt"
	"strconv"

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

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", config.DB.Username, config.DB.Password, config.DB.Host, config.DB.Port)

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
