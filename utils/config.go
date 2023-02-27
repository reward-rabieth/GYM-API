package utils

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	Logger = logger
}

type Configuration struct {
	Dbcfg     DbConfig `envconfig:"DB"`
	Servercfg ServerConfig
}

type DbConfig struct {
	User     string `required:"true" split_words:"true"`
	Name     string `required:"true" split_words:"true"`
	Password string `required:"true" split_words:"true"`
	SSLMode  string ` required:"true" default:"disable"`
}

type ServerConfig struct {
	Port string `envconfig:"PORT" default:"8000"`
}

func Loadconfig() Configuration {
	Logger.Info("loading configurations from enviroment variables")
	var cfg Configuration
	err := envconfig.Process("", &cfg)
	if err != nil {
		Logger.Fatal("failed to load enviroment variables %s", zap.Error(err))

	}
	return cfg

}
