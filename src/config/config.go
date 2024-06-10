package config

import (
	"github.com/spf13/viper"
)

var (
	Cfg *Config
)

type Config struct {
	MqURL        string `mapstructure:"MQ_URL"`
	MqRoutingKey string `mapstructure:"MQ_ROUTING_KEY"`
}

func ReadConfig() error {
	viper.AutomaticEnv()

	Cfg = &Config{
		MqURL:        viper.GetString("MQ_URL"),
		MqRoutingKey: viper.GetString("MQ_ROUTING_KEY"),
	}

	return nil
}
