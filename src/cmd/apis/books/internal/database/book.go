package database

import "time"

type Book struct {
	ID          string
	UserID      string
	Name        string
	Summary     string
	CreatedTime time.Time
}
