package infrastructure

import (
	"io/ioutil"
	"os"

	"github.com/janicduplessis/projectgo/ct/interfaces"
)

type LocalFileStoreHandler struct {
}

func (handler *LocalFileStoreHandler) Create(filePath string, data []byte) error {
	return ioutil.WriteFile(filePath, data, 0777)
}

func (handler *LocalFileStoreHandler) Open(filePath string) ([]byte, error) {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil, interfaces.ErrNoFile
		}
		return nil, err
	}
	return ioutil.ReadFile(filePath)
}
