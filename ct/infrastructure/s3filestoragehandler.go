package infrastructure

import (
	"fmt"
	"path"
	"strings"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"

	"github.com/janicduplessis/projectgo/ct/config"
	"github.com/janicduplessis/projectgo/ct/interfaces"
)

// S3FileStorageHandler handles reading and writing files to amazon s3
type S3FileStorageHandler struct {
	s3 *s3.S3
}

// Init initializes the connection to amazon
func (handler *S3FileStorageHandler) Init() {

	auth := aws.Auth{
		AccessKey: config.S3AccessKey,
		SecretKey: config.S3SecretKey,
	}
	region := aws.USEast
	handler.s3 = s3.New(auth, region)
}

// Create  creates a file
func (handler *S3FileStorageHandler) Create(filePath string, data []byte) error {
	bucket := handler.s3.Bucket(config.S3Bucket)
	return bucket.Put(filePath, data, fmt.Sprintf("image/%s", path.Ext(filePath)), s3.BucketOwnerFull)
}

// Open opens a file
func (handler *S3FileStorageHandler) Open(filePath string) ([]byte, error) {
	bucket := handler.s3.Bucket(config.S3Bucket)
	data, err := bucket.Get(filePath)
	if err != nil {
		if strings.Contains(err.Error(), "The specified key does not exist") {
			return nil, interfaces.ErrNoFile
		}
		return nil, err
	}

	return data, err
}
