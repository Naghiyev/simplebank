package util

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DbDriver      string `mapstructure:"DB_DRIVER"`
	DataSource    string `mapstructure:"DATA_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return
	}
	err = viper.Unmarshal(&config)
	return

}
