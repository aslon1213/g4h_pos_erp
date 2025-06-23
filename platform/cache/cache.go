package cache

import (
	"aslon1213/magazin_pos/pkg/configs"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

func NewCache() *redis.Client {
	log.Debug().Msg("Initializing Redis connection")

	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: config.Redis.Password,
		DB:       config.Redis.Database,
	})

	log.Info().Msg("Successfully connected to Redis")
	return client
}
