package storages

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/Ankumeah/pkgeep/package-managers"
	"github.com/Ankumeah/pkgeep/utils"
)

type S3Storage struct{}

func init() { Stoarges[S3Storage{}.Name()] = S3Storage{} }

func (s S3Storage) Name() string { return "s3" }

func (s S3Storage) Save(
	packages []packagemanagers.Package,
	config map[string]any,
) error {
	ctx := context.Background()

	timeout := utils.GetFromConfig[float64]("timeout_seconds", config)
	if timeout < 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeout))
		defer cancel()
	}

	objectID := utils.GetFromConfig[string]("object_id", config)
	bucket := utils.GetFromConfig[string]("bucket", config)

	client, err := s3GetClient(config)
	if err != nil {
		return err
	}

	data := bytes.NewBuffer(nil)
	if err := json.NewEncoder(data).Encode(packages); err != nil {
		return err
	}

	if _, err := client.PutObject(
		ctx, bucket, objectID, data, int64(data.Len()),
		minio.PutObjectOptions{
			ContentType: "application/json",
		},
	); err != nil {
		return err
	}

	return nil
}

func (s S3Storage) Get(
	config map[string]any,
) ([]packagemanagers.Package, error) {
	ctx := context.Background()

	timeout := utils.GetFromConfig[float64]("timeout_seconds", config)
	if timeout < 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeout))
		defer cancel()
	}

	objectID := utils.GetFromConfig[string]("object_id", config)
	bucket := utils.GetFromConfig[string]("bucket", config)

	client, err := s3GetClient(config)
	if err != nil {
		return nil, err
	}

	obj, err := client.GetObject(ctx, bucket, objectID, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	packages := []packagemanagers.Package{}
	if err := json.NewDecoder(obj).Decode(&packages); err != nil {
		return nil, err
	}

	return packages, nil
}

func s3GetClient(config map[string]any) (*minio.Client, error) {
	location := utils.GetFromConfig[string]("location", config)
	url := utils.GetFromConfig[string]("url", config)
	accessKeyID := utils.GetFromConfig[string]("access_key_id", config)
	secretAccessKey := utils.GetFromConfig[string]("secret_access_key", config)
	ssl := utils.GetFromConfig[bool]("ssl", config)

	minioClient, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Region: location,
		Secure: ssl,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
