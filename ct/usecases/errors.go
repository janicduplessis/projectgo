package usecases

import (
	"errors"
)

var ErrUserAlreadyExists = errors.New("The user with the specified name already exists")
