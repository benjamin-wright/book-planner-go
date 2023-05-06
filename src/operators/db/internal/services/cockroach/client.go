package cockroach

import (
	"fmt"

	"ponglehub.co.uk/book-planner-go/src/pkg/postgres"
)

type Client struct {
	conn *postgres.AdminConn
}

func New(database string, namespace string) (*Client, error) {
	cfg := postgres.ConnectConfig{
		Host:     fmt.Sprintf("%s.%s.svc.cluster.local", database, namespace),
		Port:     26257,
		Username: "root",
	}

	conn, err := postgres.NewAdminConn(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cockroach db at %s: %+v", database, err)
	}

	return &Client{
		conn: conn,
	}, nil
}

func remove[T comparable](slice []T, element T) []T {
	for idx, elem := range slice {
		if elem == element {
			slice[idx] = slice[len(slice)-1]
			return slice[:len(slice)-1]
		}
	}

	return slice
}

func (c *Client) ListDBs() ([]string, error) {
	databases, err := c.conn.ListDatabases()
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %+v", err)
	}

	databases = remove(databases, "postgres")
	databases = remove(databases, "system")

	return databases, nil
}

func (c *Client) CreateDB(name string) error {
	err := c.conn.CreateDatabase(name)
	if err != nil {
		return fmt.Errorf("failed to create database %s: %+v", name, err)
	}

	return nil
}
