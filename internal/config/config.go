package config

import (
	"os"
)

type PGConfig struct {
	Conn     string
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type Config struct {
	Server ServerConfig
	PG     PGConfig
}

type ServerConfig struct {
	Address string
}


func New() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Address: os.Getenv("SERVER_ADDRESS"),
		},
		PG: PGConfig{
			Conn:     os.Getenv("POSTGRES_CONN"),
			Username: os.Getenv("POSTGRES_USERNAME"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Database: os.Getenv("POSTGRES_DATABASE"),
		},
	}, nil
}
