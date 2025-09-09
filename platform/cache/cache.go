package cache

import (
	"github.com/aslon1213/g4h_pos_erp/pkg/configs"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

type Cache struct {
	RedisClient *redis.Client
}

func New() *Cache {
	log.Debug().Msg("Initializing Redis connection")

	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	log.Info().Str("host", config.Redis.Host).Str("port", config.Redis.Port).Int("db", config.Redis.Database).Msg("Connecting to Redis")
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: config.Redis.Password,
		DB:       config.Redis.Database,
	})

	// ping the redis client
	_, err = client.Ping().Result()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ping Redis")
	}

	log.Info().Msg("Successfully connected to Redis")
	return &Cache{
		RedisClient: client,
	}
}
