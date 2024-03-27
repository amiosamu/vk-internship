package service

import (
	"fmt"
)

var (
	ErrCannotSignToken  = fmt.Errorf("cannot sign token")
	ErrCannotParseToken = fmt.Errorf("cannot parse token")

	ErrUserAlreadyExists = fmt.Errorf("already exists")
	ErrCannotCreateUser  = fmt.Errorf("cannot create user")
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrCannotGetUser     = fmt.Errorf("cannot get user")

	ErrCannotGetAdvertisement    = fmt.Errorf("cannot get advertisement")
	ErrCannotCreateAdvertisement = fmt.Errorf("cannot create advertisement")
	ErrAdvertisementNotFound     = fmt.Errorf("advertisement not found")
)
