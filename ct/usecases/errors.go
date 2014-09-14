package usecases

import (
	"errors"
)

var ErrUserAlreadyExists = errors.New("Username unavailable")
var ErrAccessDenied = errors.New("Access denied")
