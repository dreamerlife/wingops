package config

import "github.com/spf13/viper"

type Config struct {
	HTTPAddr string
}

func Load() Config {
	viper.SetDefault("http.addr", ":8080")
	viper.SetEnvPrefix("WINGOPS")
	viper.AutomaticEnv()

	return Config{
		HTTPAddr: viper.GetString("http.addr"),
	}
}
