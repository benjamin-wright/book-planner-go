package cockroach

import (
	"fmt"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/pkg/postgres"
)

type Client struct {
	conn     *postgres.AdminConn
	database string
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
		conn:     conn,
		database: database,
	}, nil
}

func (c *Client) Stop() {
	zap.S().Infof("Closing connection to DB %s", c.database)
	c.conn.Stop()
}

func isReservedDB(name string) bool {
	return name == "system" || name == "postgres"
}

func (c *Client) ListDBs() ([]Database, error) {
	names, err := c.conn.ListDatabases()
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %+v", err)
	}

	databases := []Database{}
	for _, name := range names {
		if isReservedDB(name) {
			continue
		}

		databases = append(databases, Database{
			DB:   c.database,
			Name: name,
		})
	}

	return databases, nil
}

func (c *Client) CreateDB(db Database) error {
	err := c.conn.CreateDatabase(db.Name)
	if err != nil {
		return fmt.Errorf("failed to create database %s: %+v", db.Name, err)
	}

	return nil
}

func (c *Client) DeleteDB(db Database) error {
	err := c.conn.DropDatabase(db.Name)
	if err != nil {
		return fmt.Errorf("failed to create database %s: %+v", db.Name, err)
	}

	return nil
}

func isReservedUser(name string) bool {
	return name == "" || name == "admin" || name == "root"
}

func (c *Client) ListUsers() ([]User, error) {
	names, err := c.conn.ListUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %+v", err)
	}

	users := []User{}
	for _, name := range names {
		if isReservedUser(name) {
			continue
		}

		users = append(users, User{
			DB:   c.database,
			Name: name,
		})
	}

	return users, nil
}

func (c *Client) CreateUser(user User) error {
	err := c.conn.CreateUser(user.Name)
	if err != nil {
		return fmt.Errorf("failed to create user %s: %+v", user, err)
	}

	return nil
}

func (c *Client) DeleteUser(user User) error {
	err := c.conn.DropUser(user.Name)
	if err != nil {
		return fmt.Errorf("failed to delete user %s: %+v", user, err)
	}

	return nil
}

func (c *Client) ListPermitted(db Database) ([]Permission, error) {
	permitted, err := c.conn.ListPermitted(db.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %+v", err)
	}

	permissions := []Permission{}
	for _, user := range permitted {
		if isReservedUser(user) {
			continue
		}

		permissions = append(permissions, Permission{
			DB:       c.database,
			Database: db.Name,
			User:     user,
		})
	}

	return permissions, nil
}

func (c *Client) GrantPermission(permission Permission) error {
	err := c.conn.GrantPermissions(permission.User, permission.Database)
	if err != nil {
		return fmt.Errorf("failed to grant permission: %+v", err)
	}

	return nil
}

func (c *Client) RevokePermission(permission Permission) error {
	err := c.conn.RevokePermissions(permission.User, permission.Database)
	if err != nil {
		return fmt.Errorf("failed to revoke permission: %+v", err)
	}

	return nil
}
