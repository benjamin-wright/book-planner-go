package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"ponglehub.co.uk/book-planner-go/src/pkg/postgres"
)

type Client struct {
	conn *pgx.Conn
}

var ErrAlreadyExists = errors.New("already exists")

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

func (c *Client) DeleteAllBooks() error {
	_, err := c.conn.Exec(context.TODO(), `DELETE FROM books`)
	if err != nil {
		return fmt.Errorf("failed to delete all books: %+v", err)
	}

	return nil
}

func (c *Client) AddBook(book Book) error {
	_, err := c.conn.Exec(
		context.TODO(),
		`INSERT INTO books ("user_id", "name", "summary", "created_time") VALUES ($1, $2, $3, $4)`,
		book.UserID, book.Name, book.Summary, time.Now(),
	)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" && pgerr.ConstraintName == "books_user_id_name_key" {
			return ErrAlreadyExists
		}
		return fmt.Errorf("error creating book: %+v", err)
	}

	return nil
}

func (c *Client) GetBooks(user string) ([]Book, error) {
	rows, err := c.conn.Query(context.TODO(), `SELECT "id", "user_id", "name", "summary", "created_time" FROM books WHERE user_id = $1`, user)
	if err != nil {
		return nil, fmt.Errorf("error fetching books: %+v", err)
	}
	defer rows.Close()

	books := []Book{}
	for rows.Next() {
		book := Book{}

		err = rows.Scan(&book.ID, &book.UserID, &book.Name, &book.Summary, &book.CreatedTime)
		if err != nil {
			return nil, fmt.Errorf("failed to parse row: %+v", err)
		}
	}

	return books, nil
}
