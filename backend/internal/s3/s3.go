package s3

import (
	"context"
	"mime/multipart"
	"pixelbattle/internal/config"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Uploader interface {
	UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, objectName string) (string, error)
	GetPresignedURL(ctx context.Context, objectName string, expires time.Duration) (string, error)
}

type Client struct {
	cli    *minio.Client
	bucket string
}

func New(cfg config.Config) (*Client, error) {
	cli, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: cfg.Minio.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	err = cli.MakeBucket(ctx, cfg.Minio.Bucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := cli.BucketExists(ctx, cfg.Minio.Bucket)
		if errBucketExists != nil || !exists {
			return nil, err
		}
	}

	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::avatars/*"]
			}
		]
	}`

	err = cli.SetBucketPolicy(ctx, "avatars", policy)
	if err != nil {
		return nil, err
	}

	return &Client{cli: cli, bucket: cfg.Minio.Bucket}, nil
}

func (c *Client) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, objectName string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = c.cli.PutObject(ctx, c.bucket, objectName, file, fileHeader.Size, minio.PutObjectOptions{ContentType: fileHeader.Header.Get("Content-Type")})
	if err != nil {
		return "", err
	}
	return objectName, nil
}

func (c *Client) GetPresignedURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	u, err := c.cli.PresignedGetObject(ctx, c.bucket, objectName, expires, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
