package usecases

import (
	"errors"
)

var (
	// Registration
	ErrUserAlreadyExists   = errors.New("Username unavailable")
	ErrInvalidRegisterInfo = errors.New("Invalid registration info")

	// Others
	ErrAccessDenied     = errors.New("Access denied")
	ErrInvalidChannelId = errors.New("Invalid channel id")
)
