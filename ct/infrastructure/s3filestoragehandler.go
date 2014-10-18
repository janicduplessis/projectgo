package infrastructure

import (
	"io"
	"path"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"

	"github.com/janicduplessis/projectgo/ct/config"
)

// S3FileStorageHandler handles reading and writing files to amazon s3
type S3FileStorageHandler struct {
	s3 *s3.S3
}

// ReadWriteCloser for s3
type s3ReadWriteCloser struct {
	fileName string
	bucket   *s3.Bucket
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
func (handler *S3FileStorageHandler) Create(path string) (io.WriteCloser, error) {
	bucket := handler.s3.Bucket(config.S3Bucket)
	writer := &s3ReadWriteCloser{
		fileName: path,
		bucket:   bucket,
	}
	return writer, nil
}

// Open opens a file
func (handler *S3FileStorageHandler) Open(path string) (io.ReadCloser, error) {
	bucket := handler.s3.Bucket(config.S3Bucket)
	reader := &s3ReadWriteCloser{
		fileName: path,
		bucket:   bucket,
	}
	return reader, nil
}

func (wc *s3ReadWriteCloser) Read(p []byte) (int, error) {
	data, err := wc.bucket.Get(wc.fileName)
	if err != nil {
		return 0, err
	}

	copy(data, p)

	return len(p), nil
}

func (wc *s3ReadWriteCloser) Write(p []byte) (int, error) {
	err := wc.bucket.Put(wc.fileName, p, "image/"+path.Ext(wc.fileName), s3.BucketOwnerFull)
	return len(p), err
}

func (wc *s3ReadWriteCloser) Close() error {
	return nil
}
