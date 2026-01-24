package config

import (
	"reflect"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/spf13/viper"
)

func New() *domain.Config {
	config := &domain.Config{}
	err := Load(config)
	if err != nil {
		panic(err)
	}

	return config
}

func Load(config *domain.Config) error {
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	variableNames := getTags("mapstructure", domain.Config{})

	for _, v := range variableNames {
		if err := viper.BindEnv(v); err != nil {
			return err
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Handle error if is not equal ConfigFileNotFoundError
			return err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	return nil
}

func getTags(tagName string, obj any) []string {
	var tags []string
	envVarType := reflect.TypeOf(obj)

	for i := 0; i < envVarType.NumField(); i++ {
		field := envVarType.Field(i)
		tag := field.Tag.Get(tagName)
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags
}
