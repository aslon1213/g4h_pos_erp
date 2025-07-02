package configs

import (
	"crypto/rand"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

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

type ServerConfig struct {
	Host               string          `mapstructure:"host"`
	Port               string          `mapstructure:"port"`
	SecretSymmetricKey string          `mapstructure:"secret_symmetric_key"`
	TokenExpiryHours   int             `mapstructure:"token_expiry_hours"`
	AdminDocsUsers     []AdminDocsUser `mapstructure:"admin_docs_users"`
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
	filename := "config"
	if strings.ToLower(os.Getenv("ENVIRONMENT")) != "production" {
		filename = "config.local"
	}
	log.Info().Str("ENVIRONMENT", os.Getenv("ENVIRONMENT")).Str("filename", filename).Msg("ENVIRONMENT")
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
	log.Info().Interface("config", config).Msg("Config loaded")

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
