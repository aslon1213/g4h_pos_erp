package configs

import (
	"crypto/rand"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var ENVT_TYPE_LOGGED bool = false

type Config struct {
	DB     DBConfig     `mapstructure:"database"`
	Redis  RedisConfig  `mapstructure:"redis"`
	Server ServerConfig `mapstructure:"server"`
	S3     S3Config     `mapstructure:"s3"`
}

type DBConfig struct {
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	Username       string `mapstructure:"username"`
	Password       string `mapstructure:"password"`
	Database       string `mapstructure:"database"`
	MaxConnections uint64 `mapstructure:"max_connections"`
	MinPoolSize    uint64 `mapstructure:"min_pool_size"`
	Auth           bool   `mapstructure:"auth"`
	ReplicaSet     string `mapstructure:"replica_set"`
	URL            string `mapstructure:"url"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type ProxyConfig struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
	Addr string `mapstructure:"addr"`
	APIKey string `mapstructure:"api_key"`
}

type ServerConfig struct {
	Host               string          `mapstructure:"host"`
	Port               string          `mapstructure:"port"`
	SecretSymmetricKey string          `mapstructure:"secret_symmetric_key"`
	TokenExpiryHours   int             `mapstructure:"token_expiry_hours"`
	AdminDocsUsers     []AdminDocsUser `mapstructure:"admin_docs_users"`
	Proxy              []ProxyConfig   `mapstructure:"proxy"`
}

type S3Config struct {
	Region          string `mapstructure:"region"`
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	ImageBucket     string `mapstructure:"image_bucket"`
}

type AdminDocsUser struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func LoadConfig(path string) (*Config, error) {
	envFrom := strings.ToLower(os.Getenv("ENV_FROM"))

	if envFrom == ".env" {
		return loadConfigFromEnv(path)
	}
	return loadConfigFromYaml(path)
}

func loadConfigFromEnv(path string) (*Config, error) {
	if !ENVT_TYPE_LOGGED {
		log.Info().Str("ENV_FROM", ".env").Msg("Loading config from environment variables")
		ENVT_TYPE_LOGGED = true
	}





	

	maxConnections, _ := strconv.ParseUint(os.Getenv("DATABASE_MAX_CONNECTIONS"), 10, 64)
	minPoolSize, _ := strconv.ParseUint(os.Getenv("DATABASE_MIN_POOL_SIZE"), 10, 64)
	dbAuth, _ := strconv.ParseBool(os.Getenv("DATABASE_AUTH"))
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	tokenExpiryHours, _ := strconv.Atoi(os.Getenv("SERVER_TOKEN_EXPIRY_HOURS"))

	log.Info().Str("DATABASE_HOST", os.Getenv("DATABASE_HOST")).Str("DATABASE_PORT", os.Getenv("DATABASE_PORT")).Str("DATABASE_NAME", os.Getenv("DATABASE_NAME")).Str("DATABASE_MAX_CONNECTIONS", os.Getenv("DATABASE_MAX_CONNECTIONS")).Str("DATABASE_MIN_POOL_SIZE", os.Getenv("DATABASE_MIN_POOL_SIZE")).Str("DATABASE_AUTH", os.Getenv("DATABASE_AUTH")).Str("DATABASE_REPLICA_SET", os.Getenv("DATABASE_REPLICA_SET")).Str("DATABASE_URL", os.Getenv("DATABASE_URL")).Str("REDIS_HOST", os.Getenv("REDIS_HOST")).Str("REDIS_PORT", os.Getenv("REDIS_PORT")).Str("REDIS_DATABASE", os.Getenv("REDIS_DATABASE")).Str("SERVER_TOKEN_EXPIRY_HOURS", os.Getenv("SERVER_TOKEN_EXPIRY_HOURS")).Msg("Loading config from environment variables")

	config := &Config{
		DB: DBConfig{
			Host:           os.Getenv("DATABASE_HOST"),
			Port:           os.Getenv("DATABASE_PORT"),
			Username:       os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
			Password:       os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
			Database:       os.Getenv("DATABASE_NAME"),
			MaxConnections: maxConnections,
			MinPoolSize:    minPoolSize,
			Auth:           dbAuth,
			ReplicaSet:     os.Getenv("DATABASE_REPLICA_SET"),
			URL:            os.Getenv("DATABASE_URL"),
		},
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     strconv.Itoa(redisPort),
			Password: os.Getenv("REDIS_PASSWORD"),
			Database: redisDB,
		},
		Server: ServerConfig{
			Host:               os.Getenv("SERVER_HOST"),
			Port:              os.Getenv("SERVER_PORT"),
			SecretSymmetricKey: os.Getenv("SERVER_SECRET_SYMMETRIC_KEY"),
			TokenExpiryHours:   tokenExpiryHours,
		},
		S3: S3Config{
			Region:          os.Getenv("S3_REGION"),
			Endpoint:        os.Getenv("S3_ENDPOINT"),
			AccessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
			ImageBucket:     os.Getenv("S3_IMAGE_BUCKET"),
		},
	}

	if strings.ToLower(os.Getenv("LOG_ENV")) == "true" {
		log.Info().Interface("config", config).Msg("Config loaded from env")
	}

	return config, nil
}

func loadConfigFromYaml(path string) (*Config, error) {
	filename := "config"
	if strings.ToLower(os.Getenv("ENVIRONMENT")) != "production" {
		filename = "config.local"
	}
	if !ENVT_TYPE_LOGGED {
		log.Info().Str("ENVIRONMENT", os.Getenv("ENVIRONMENT")).Str("filename", filename).Msg("Loading config from YAML")
		ENVT_TYPE_LOGGED = true
	}
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// log environment
	if strings.ToLower(os.Getenv("LOG_ENV")) == "true" {
		log.Info().Interface("config", config).Msg("Config loaded")
	}

	return &config, nil
}

func GenerateSecretSymmetricKey() []byte {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		log.Fatal().Err(err).Msg("Failed to generate symmetric key")
	}
	log.Info().Str("secret_symmetric_key", string(key)).Int("Length", len(key)).Msg("secret_symmetric_key")
	return key
}
