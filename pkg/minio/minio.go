package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"os"
	"path/filepath"
)

type MinIOConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	UseSSL          bool   `json:"use_ssl"`
}

// UploadResult struct untuk mengelola hasil upload file
type UploadResult struct {
	Etag      string
	PublicURL string
}

// FileResult struct untuk menampung file
type FileResult struct {
	Source          string
	Bucket          string
	Object          string
	DestinationFile string
	Success         bool
	ErrorMessage    string
}

// MinIOAuthService mengimplementasikan file storage service untuk MinIO
type MinIOAuthService struct {
	minioClient *minio.Client
}

// NewMinIOAuthService menginisialisasi MiniIO Service
func NewMinIOAuthService(endpoint, accessKey, secretKey string, ssl bool) (*MinIOAuthService, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure:       ssl, // Set to true if using HTTPS
		BucketLookup: minio.BucketLookupAuto,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOAuthService{
		minioClient: minioClient,
	}, nil
}

// UploadFile mengupload file ke MinIO storage
func (m *MinIOAuthService) UploadFile(ctx context.Context, objectName string, bucketName string, filePath string) (*UploadResult, error) {
	// Membuka file
	file, err := os.Open(filepath.Join(filePath, objectName))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Mengambil statistik file
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Set optional metadata
	options := minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	}

	etag, err := m.minioClient.PutObject(ctx, bucketName, objectName, file, stat.Size(), options)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		Etag:      etag.ETag,
		PublicURL: etag.Location,
	}, nil
}

// RetrieveFile mengambil file dari MinIO storage
func (m *MinIOAuthService) RetrieveFile(ctx context.Context, objectName, bucketName, destinationPath string) error {
	filePath := filepath.Join(destinationPath, objectName)

	err := m.minioClient.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
