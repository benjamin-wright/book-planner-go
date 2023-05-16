package client

import (
	"fmt"

	"github.com/jackc/pgx/v4"
	"ponglehub.co.uk/book-planner-go/src/pkg/postgres"
)

type Client struct {
	conn *pgx.Conn
}

func New() (*Client, error) {
	config, err := postgres.ConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed getting config from environment: %+v", err)
	}

	conn, err := postgres.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to postgres: %+v", err)
	}

	return &Client{conn: conn}, nil
}
