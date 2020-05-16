package config

import "time"

type Config struct {
	AppConfig   App   `mapstructure:"app"`
	RedisConfig Redis `mapstructure:"redis"`
}

type App struct {
	BaseUrl          string        `mapstructure:"base-url"`
	Port             int16         `mapstructure:"port"`
	BodyLimit        string        `mapstructure:"body-limit"`
	RequestStoreTime time.Duration `mapstructure:"request-store-time"`
	ClientUrl        string        `mapstructure:"client_url"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int16  `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}
