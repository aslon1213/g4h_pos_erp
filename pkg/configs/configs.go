package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB     DBConfig     `mapstructure:"database"`
	Redis  RedisConfig  `mapstructure:"redis"`
	Server ServerConfig `mapstructure:"server"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
