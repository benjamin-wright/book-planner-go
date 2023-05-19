package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/types"
	"ponglehub.co.uk/book-planner-go/src/pkg/postgres"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/wasm/validation"
)

var ErrNoUser = errors.New("user not found")
var ErrUserExists = errors.New("user already exists")
var ErrComplexity = errors.New("password didn't meet complexity requirements")

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

func (c *Client) DeleteAllUsers() error {
	_, err := c.conn.Exec(context.Background(), `DELETE FROM users`)
	if err != nil {
		return fmt.Errorf("failed to clear existing users; %+v", err)
	}

	return nil
}

func (c *Client) AddUser(name string, password string) error {
	if !validation.CheckPasswordComplexity(password) {
		return ErrComplexity
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %+v", err)
	}
	hash := string(bytes)

	_, err = c.conn.Exec(context.Background(), `INSERT INTO users("name", "password") VALUES ($1, $2)`, name, hash)

	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" && pgerr.ConstraintName == "users_name_key" {
			return ErrUserExists
		}
		return fmt.Errorf("failed to add user to database: %+v", err)
	}

	return nil
}

var ErrPasswordMismatch = errors.New("password mismatch")

func (c *Client) CheckPassword(name string, password string) (*types.User, error) {
	rows, err := c.conn.Query(context.Background(), `SELECT "id", "password" FROM users WHERE "name" = $1`, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from database: %+v", err)
	}
	defer rows.Close()

	numUsers := 0
	user := types.User{
		Name: name,
	}
	var passwordHash string

	for rows.Next() {
		numUsers += 1
		if err = rows.Scan(&user.ID, &passwordHash); err != nil {
			return nil, fmt.Errorf("failed to parse new user ID: %+v", err)
		}
	}

	if numUsers == 0 {
		return nil, ErrNoUser
	}

	if numUsers > 1 {
		return nil, fmt.Errorf("expected 1 user, got %d", numUsers)
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return nil, ErrPasswordMismatch
	}

	return &user, nil
}

func (c *Client) GetUser(name string) (*types.User, error) {
	rows, err := c.conn.Query(context.Background(), `SELECT "id" FROM users WHERE "name" = $1`, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from database: %+v", err)
	}
	defer rows.Close()

	numUsers := 0
	user := types.User{
		Name: name,
	}

	for rows.Next() {
		numUsers += 1
		if err = rows.Scan(&user.ID); err != nil {
			return nil, fmt.Errorf("failed to parse new user ID: %+v", err)
		}
	}

	if numUsers == 0 {
		return nil, ErrNoUser
	}

	if numUsers > 1 {
		return nil, fmt.Errorf("expected 1 user, got %d", numUsers)
	}

	return &user, nil
}

func (c *Client) ListUsers() ([]types.User, error) {
	rows, err := c.conn.Query(context.Background(), `SELECT "id", "name", "password" FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from database: %+v", err)
	}
	defer rows.Close()

	users := []types.User{}

	for rows.Next() {
		user := types.User{}
		if err = rows.Scan(&user.ID, &user.Name, &user.Password); err != nil {
			return nil, fmt.Errorf("failed to parse new user ID: %+v", err)
		}

		users = append(users, user)
	}

	return users, nil
}
