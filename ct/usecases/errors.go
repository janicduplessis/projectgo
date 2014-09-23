package usecases

import (
	"errors"
)

var (
	ErrUserAlreadyExists = errors.New("Username unavailable")
	ErrAccessDenied      = errors.New("Access denied")
	ErrInvalidChannelId  = errors.New("Invalid channel id")
)
