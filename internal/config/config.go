package config

import (
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
)

type Configuration struct {
	AppConfig   App   `mapstructure:"app"`
	RedisConfig Redis `mapstructure:"redis"`
	Rate        Rate  `mapstructure:"rate"`
}

type App struct {
	Port             int16         `mapstructure:"port"`
	BodyLimit        string        `mapstructure:"bodyLimit"`
	RequestStoreTime time.Duration `mapstructure:"requestStoreTime"`
	ClientUrl        string        `mapstructure:"clientUrl"`
	ForbiddenHeaders []string      `mapstructure:"forbiddenHeaders"`
	HistoryCount     int           `mapstructure:"historyCount"`
	MaxRequestSize   int           `mapstructure:"maxRequestSize"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int16  `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Rate struct {
	Limit  int64         `mapstructure:"limit"`
	Period time.Duration `mapstructure:"period"`
}

func InitConfig() *Configuration {
	viper.SetConfigFile("configs/config.yml")
	viper.SetConfigType("yml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Configuration file not found...")
	}

	var configuration *Configuration

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return configuration
}
