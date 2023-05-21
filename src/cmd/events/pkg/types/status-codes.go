package types

type StatusCode int64

const (
	Created StatusCode = iota
	AlreadyExists
	ServerError
)
