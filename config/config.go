package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Host     string `mapstructure:"host"`
	UserName string `mapstructure:"userName"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Port     int    `mapstructure:"port"`
	SSLMode  string `mapstructure:"sslmode"`
}

func LoadConfig() (Config, error) {
	var config Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("error reading configuration file: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("error unpacking configuration: %w", err)
	}

	if err := checkConfigData(config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func checkConfigData(cfg Config) (err error) {
	errMsg := "Empty values in config file:"
	if cfg.Port == 0 {
		errMsg += " port"
	}
	if cfg.UserName == "" {
		errMsg += " userName"
	}
	if cfg.Password == "" {
		errMsg += " password"
	}
	if cfg.DBName == "" {
		errMsg += " dbname"
	}
	if cfg.SSLMode == "" {
		errMsg += " sslmode"
	}
	if cfg.Host == "" {
		errMsg += " host"
	}
	if errMsg != "Empty values in config file:" {
		err = errors.Errorf(errMsg)
	}
	return err
}
