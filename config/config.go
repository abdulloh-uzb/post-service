package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Environment         string
	PostgresHost        string
	PostgresPort        int
	PostgresDatabase    string
	PostgresUser        string
	PostgresPassword    string
	LogLevel            string
	RPCPort             string
	CustomerServiceHost string
	CustomerServicePort int
	RankingServiceHost  string
	RankingServicePort  int
}

func Load() Config {
	c := Config{}
	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "post_service"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "1234"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.RPCPort = cast.ToString(getOrReturnDefault("RPC_PORT", ":8000"))

	c.CustomerServiceHost = cast.ToString(getOrReturnDefault("CUSTOMER_HOST", "localhost"))
	c.CustomerServicePort = cast.ToInt(getOrReturnDefault("CUSTOMER_PORT", 9000))

	c.RankingServiceHost = cast.ToString(getOrReturnDefault("CUSTOMER_HOST", "localhost"))
	c.RankingServicePort = cast.ToInt(getOrReturnDefault("CUSTOMER_PORT", 1111))
	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaultValue
}
