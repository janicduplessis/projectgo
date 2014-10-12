package usecases

import (
	"errors"
)

var (
	// Registration
	ErrUserAlreadyExists   = errors.New("Username unavailable")
	ErrInvalidRegisterInfo = errors.New("Invalid registration info")

	// Invalid ids
	ErrInvalidChannelId = errors.New("Invalid channel id")
	ErrInvalidClientId  = errors.New("Invalid client id")

	// Others
	ErrAccessDenied = errors.New("Access denied")
	ErrNoChannel    = errors.New("User is not in any channel")
)
