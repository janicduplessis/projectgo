package infrastructure

import (
	"io"
	"os"
)

type LocalFileStoreHandler struct {
}

func (handler *LocalFileStoreHandler) Create(path string) (io.WriteCloser, error) {
	return os.Create(path)
}

func (handler *LocalFileStoreHandler) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}
