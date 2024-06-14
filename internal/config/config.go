package config

import (
	"github.com/charmbracelet/log"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	HTTPServer `yaml:"http_server"`
	Postgres
	NatsStreamingConfig
	Redis
}

type HTTPServer struct {
	Address string `yaml:"address" default:"localhost:8080"`
}

type NatsStreamingConfig struct {
	ClientAddress string `yaml:"client_address" default:"localhost:4222"`
	HttpAddress   string `yaml:"http_address" default:"localhost:4223"`
}

type Postgres struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     int
}

type Redis struct {
	Host     string
	Port     int
	Password string
}

func retrieveRedisConfig(config *Config) {
	var redis Redis
	redis.Host = os.Getenv("REDIS_HOST")
	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Fatal("failed to parse REDIS_PORT: " + err.Error())
	}
	redis.Port = port
	redis.Password = os.Getenv("REDIS_PASSWORD")
	config.Redis = redis
}

func retrievePostgresConfig(config *Config) {
	var postgres Postgres
	postgres.User = os.Getenv("POSTGRES_USER")
	postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	postgres.DBName = os.Getenv("POSTGRES_DB")
	postgres.Host = os.Getenv("POSTGRES_HOST")

	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Fatal("failed to parse POSTGRES_PORT: ", "err", err)
	}
	postgres.Port = port
	config.Postgres = postgres
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file: ", "err", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file not found: "+configPath, "err", err)
	}

	var config Config

	err = cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		log.Fatal("failed to read config: ", "err", err)
	}

	retrievePostgresConfig(&config)
	retrieveRedisConfig(&config)

	return &config
}
