package configs

import (
	"github.com/spf13/viper"
	"os"
	"path"
)

type Config struct {
	AppHost string `mapstructure:"APP_HOST"`
	AppPort uint16 `mapstructure:"APP_PORT"`
	Secret  string `mapstructure:"SECRET"`

	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     uint16 `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
}

func (c *Config) LoadConfig() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path.Join(pwd, "configs"))
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return err
	}

	return nil
}
