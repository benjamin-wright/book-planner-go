package redis

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type ConnectConfig struct {
	Host     string
	Port     int
	Database int
}

func ConfigFromEnv() (ConnectConfig, error) {
	empty := ConnectConfig{}

	host, ok := os.LookupEnv("REDIS_HOST")
	if !ok {
		return empty, errors.New("failed to lookup REDIS_HOST env var")
	}

	portString, ok := os.LookupEnv("REDIS_PORT")
	if !ok {
		return empty, errors.New("failed to lookup REDIS_PORT env var")
	}

	port, err := strconv.Atoi(portString)
	if err != nil {
		return empty, fmt.Errorf("failed to convert REDIS_PORT: %+v", err)
	}

	databaseString, ok := os.LookupEnv("REDIS_DATABASE")
	if !ok {
		return empty, errors.New("failed to lookup REDIS_DATABASE env var")
	}

	database, err := strconv.Atoi(databaseString)
	if err != nil {
		return empty, fmt.Errorf("failed to convert REDIS_DATABASE: %+v", err)
	}

	return ConnectConfig{
		Host:     host,
		Port:     port,
		Database: database,
	}, nil
}
