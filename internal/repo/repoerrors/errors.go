package repoerrors

import "errors"

var (
	ErrNotFound  = errors.New("user not found")
	CannotCreate = errors.New("cannot create user")
	ErrAlreadyExists = errors.New("already exists")
)
