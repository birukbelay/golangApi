package config

import (
	"log"
	"os"
	"strconv"
)

type PostgresConfig struct {
	Dbname string
	Username string
	Pass string
}
type MongoConfig struct{
	DbName string
	Port string
}
type Config struct{
	Postgres PostgresConfig
	Port  int
	SignKey string
}

func New() *Config{
	return &Config{
		Postgres:PostgresConfig{
			Dbname: getEnv("PG_DB_NAME", ""),
			Username: getEnv("PG_USER_NAME",""),
			Pass: getEnv("PG_PASS",""),

		},
		Port: getEnvAsInt("PORT",8080),
		SignKey : getEnv("SIGN_KEY","mySecrato"),
	}
}

func getEnv(key string, defaultVal string) string  {
	if value, exists :=os.LookupEnv(key); exists{
		return value
	}
	log.Panicf("could Not get value of %s", key)
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}