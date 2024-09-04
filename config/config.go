package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"reflect"
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

func checkConfigData(cfg Config) error {
	errMsg := "Please fill in these variables in the config file:"

	val := reflect.ValueOf(cfg)
	typ := reflect.TypeOf(cfg)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		if field.IsZero() {
			errMsg += " " + fieldName
		}
	}

	if errMsg != "Please fill in these variables in the config file:" {
		return errors.New(errMsg)
	}
	return nil
}
