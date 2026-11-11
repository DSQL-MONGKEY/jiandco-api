package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"jianco-api/config"
)

var ErrNotConfigured = errors.New("storage is not configured")

type Client struct {
	client        *s3.Client
	bucket        string
	publicBaseURL string
	enabled       bool
}

func New(cfg config.R2Config) (*Client, error) {
	if !cfg.Enabled() {
		return &Client{enabled: false}, nil
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("load r2 config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(cfg.Endpoint)
		options.UsePathStyle = true
	})

	return &Client{
		client:        client,
		bucket:        cfg.Bucket,
		publicBaseURL: cfg.PublicBaseURL,
		enabled:       true,
	}, nil
}

func (c *Client) Enabled() bool {
	return c != nil && c.enabled
}

func (c *Client) Upload(ctx context.Context, key string, body io.Reader, contentType string, size int64) (string, error) {
	if !c.Enabled() {
		return "", ErrNotConfigured
	}

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(c.bucket),
		Key:           aws.String(key),
		Body:          body,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	})
	if err != nil {
		return "", fmt.Errorf("upload to r2: %w", err)
	}

	return c.URL(key), nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	if !c.Enabled() || key == "" {
		return nil
	}

	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("delete from r2: %w", err)
	}

	return nil
}

func (c *Client) GetObject(ctx context.Context, key string) ([]byte, string, error) {
	if !c.Enabled() {
		return nil, "", ErrNotConfigured
	}

	output, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, "", fmt.Errorf("get from r2: %w", err)
	}
	defer output.Body.Close()

	body, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read r2 object: %w", err)
	}

	contentType := "application/octet-stream"
	if output.ContentType != nil && *output.ContentType != "" {
		contentType = *output.ContentType
	}
	return body, contentType, nil
}

func (c *Client) URL(key string) string {
	if c == nil || key == "" {
		return ""
	}
	if c.publicBaseURL == "" {
		return key
	}
	return strings.TrimRight(c.publicBaseURL, "/") + "/" + strings.TrimLeft(key, "/")
}

func BuildKey(folder, ownerID, filename string) string {
	extension := strings.ToLower(filepath.Ext(filename))
	if extension == "" {
		extension = ".bin"
	}
	return fmt.Sprintf("%s/%s/%s%s", strings.Trim(folder, "/"), ownerID, uuid.NewString(), extension)
}
