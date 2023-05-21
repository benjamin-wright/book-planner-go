package types

type CreateBookEvent struct {
	UserID  string
	Name    string
	Summary string
}

type CreateBookResponseEvent struct {
	UserID string
	Name   string
	Status StatusCode
}
