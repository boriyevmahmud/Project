package config

import (
	"os"

	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	Environment      string // develop, staging, production
	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string
	LogLevel         string
	RPCPort          string
	UserServiceHost  string
	UserServicePort  int
	KafkaHost        string
	KafkaPort        int
}

// Load loads environment vars and inflates Config
func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "posts"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "1"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 9000))
	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "localhost"))
	c.RPCPort = cast.ToString(getOrReturnDefault("RPC_PORT", ":7000"))

	c.KafkaHost = cast.ToString(getOrReturnDefault("KAFKA_HOST", "127.0.0.1"))
	c.KafkaPort = cast.ToInt(getOrReturnDefault("KAFKA_PORT", 9092))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
