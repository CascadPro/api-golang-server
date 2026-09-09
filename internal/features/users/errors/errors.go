package users_errors

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotActivated   = errors.New("user not activated")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUsernameExists     = errors.New("username already exists")
	ErrAvatarNotFound     = errors.New("avatar not found")
)
