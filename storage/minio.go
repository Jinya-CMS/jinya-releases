package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"jinya-releases/config"
	"strings"
)

func getMinioClient() (*minio.Client, error) {
	endpoint := strings.TrimPrefix(strings.TrimPrefix(config.LoadedConfiguration.StorageUrl, "http://"), "https://")
	accessKeyID := config.LoadedConfiguration.StorageAccessKey
	secretAccessKey := config.LoadedConfiguration.StorageSecretKey
	useSSL := strings.HasPrefix(config.LoadedConfiguration.StorageUrl, "https://")

	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
}

func SaveFile(path string, reader io.Reader, size int64, contentType string) error {
	client, err := getMinioClient()
	if err != nil {
		return err
	}

	_, err = client.PutObject(context.Background(), config.LoadedConfiguration.StorageBucket, path, reader, size, minio.PutObjectOptions{ContentType: contentType})

	return err
}

func GetFile(path string) (io.ReadCloser, string, int64, error) {
	client, err := getMinioClient()
	if err != nil {
		return nil, "", 0, err
	}

	object, err := client.GetObject(context.Background(), config.LoadedConfiguration.StorageBucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", 0, err
	}

	objectStat, err := object.Stat()
	if err != nil {
		return nil, "", 0, err
	}

	return object, objectStat.ContentType, objectStat.Size, nil
}

func DeleteFile(path string) error {
	client, err := getMinioClient()
	if err != nil {
		return err
	}

	return client.RemoveObject(context.Background(), config.LoadedConfiguration.StorageBucket, path, minio.RemoveObjectOptions{})
}
