package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	NATS     NATSConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Addr string
}

type PostgresConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}

type NATSConfig struct {
	URL string
}

type JWTConfig struct {
	Secret                string
	AccessTokenTTLMinutes int
}

func Load() Config {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	viper.SetDefault("server.addr", ":8080")
	viper.SetDefault("postgres.dsn", "host=localhost user=wingops password=wingops dbname=wingops port=5432 sslmode=disable")
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("nats.url", "nats://localhost:4222")
	viper.SetDefault("jwt.secret", "dev-secret-change-before-production")
	viper.SetDefault("jwt.access_token_ttl_minutes", 60)
	viper.SetEnvPrefix("WINGOPS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			panic(err)
		}
	}

	return Config{
		Server: ServerConfig{
			Addr: viper.GetString("server.addr"),
		},
		Postgres: PostgresConfig{
			DSN: viper.GetString("postgres.dsn"),
		},
		Redis: RedisConfig{
			Addr: viper.GetString("redis.addr"),
		},
		NATS: NATSConfig{
			URL: viper.GetString("nats.url"),
		},
		JWT: JWTConfig{
			Secret:                viper.GetString("jwt.secret"),
			AccessTokenTTLMinutes: viper.GetInt("jwt.access_token_ttl_minutes"),
		},
	}
}
