package minio_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-pkg/pkg/minio"
	"testing"
)

const (
	// credential MinIO
	endpoint  = "localhost:9000"
	accessKey = "root"     //"your-access-key"
	secretKey = "password" //"your-secret-key"
	useSSL    = false

	bucketName      = "chum-bucket"
	objectName      = "notes.txt"
	filePath        = "samples/upload_files/"
	destinationFile = "samples/retrieve_files/"
)

func TestNewMinIOAuthService(t *testing.T) {
	result, err := minio.NewMinIOAuthService(endpoint, accessKey, secretKey, useSSL)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestUploadFile(t *testing.T) {
	var (
		ctx = context.Background()
	)

	service, err := minio.NewMinIOAuthService(endpoint, accessKey, secretKey, useSSL)
	result, err := service.UploadFile(ctx, objectName, bucketName, filePath)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetFile(t *testing.T) {
	var (
		ctx = context.Background()
	)

	service, err := minio.NewMinIOAuthService(endpoint, accessKey, secretKey, useSSL)
	err = service.RetrieveFile(ctx, objectName, bucketName, destinationFile)

	assert.NoError(t, err)
}
