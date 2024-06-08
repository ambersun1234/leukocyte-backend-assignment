package config

import (
	"leukocyte/src/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Cfg *Config
)

type Kubernetes struct {
	ConfigUrl string `mapstructure:"config_url"`
	InCluster bool   `mapstructure:"in_cluster"`
}

type MessageQueue struct {
	Url        string `mapstructure:"url"`
	RoutingKey string `mapstructure:"routing_key"`
}

type Config struct {
	Kubernetes   Kubernetes   `mapstructure:"kubernetes"`
	MessageQueue MessageQueue `mapstructure:"message_queue"`
}

func ReadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logger.Entry.Error("Error reading config file.", zap.Error(err))
		return err
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		logger.Entry.Error("Unable to decode into struct.", zap.Error(err))
		return err
	}

	return nil
}
