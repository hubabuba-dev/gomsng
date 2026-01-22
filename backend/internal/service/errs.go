package service

import "errors"

var (
	ErrBlocked       = errors.New("User is blocked")
	ErrWrongCreds    = errors.New("User not exists or password is wrong")
	ErrTokenNotValid = errors.New("Not valid hash token")
)

var (
	ErrWrongChatType = errors.New("Wront chat type")
)
