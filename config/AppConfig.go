package config

import "time"

type Config struct {
	AppConfig   App   `mapstructure:"app"`
	RedisConfig Redis `mapstructure:"redis"`
}

type App struct {
	Env              string        `mapstructure:"env"`
	Port             int16         `mapstructure:"port"`
	BodyLimit        string        `mapstructure:"bodyLimit"`
	RequestStoreTime time.Duration `mapstructure:"requestStoreTime"`
	ClientUrl        string        `mapstructure:"clientUrl"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int16  `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}
