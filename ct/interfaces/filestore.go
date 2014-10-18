package interfaces

import (
	"io"
)

type FileStore interface {
	Create(path string) (io.WriteCloser, error)
	Open(path string) (io.ReadCloser, error)
}
