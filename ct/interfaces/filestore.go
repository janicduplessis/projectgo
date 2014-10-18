package interfaces

import "errors"

type FileStore interface {
	Create(filePath string, data []byte) error
	Open(filePath string) ([]byte, error)
}

var ErrNoFile = errors.New("The file does not exists")
