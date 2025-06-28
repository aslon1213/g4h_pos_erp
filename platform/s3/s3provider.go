package s3provider

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aslon1213/go-pos-erp/pkg/configs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

type S3Client struct {
	Client *s3.Client
	Bucket string
}

type CredentialsProvider struct {
	AccessKeyID     string
	SecretAccessKey string
}

func (c *CredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     c.AccessKeyID,
		SecretAccessKey: c.SecretAccessKey,
	}, nil
}

func New() *S3Client {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	crendentials := &CredentialsProvider{
		AccessKeyID:     config.S3.AccessKeyID,
		SecretAccessKey: config.S3.SecretAccessKey,
	}
	return &S3Client{
		Client: s3.NewFromConfig(aws.Config{
			Region:      config.S3.Region,
			Credentials: crendentials,
		}),
		Bucket: config.S3.ImageBucket,
	}
}

func (s *S3Client) UploadFile(file *multipart.FileHeader, key string) error {
	f, err := file.Open()
	if err != nil {
		log.Error().Err(err).Msg("Failed to open file")
		return err
	}
	defer f.Close()

	_, err = s.Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to upload file to S3")
		return err
	}

	return nil
}

func (s *S3Client) DeleteFile(key string) error {
	_, err := s.Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete file from S3")
		return err
	}

	return nil
}

func (s *S3Client) GetFile(key string) (*s3.GetObjectOutput, error) {
	result, err := s.Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get file from S3")
		return nil, err
	}

	return result, nil
}

func (s *S3Client) GetFileURL(key string) (string, error) {
	// Construct the S3 URL
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.Bucket, key)
	return url, nil
}

func (s *S3Client) ListFiles(prefix string) ([]string, error) {
	result, err := s.Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.Bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to list files from S3")
		return nil, err
	}

	files := make([]string, 0)
	for _, obj := range result.Contents {
		url, err := s.GetFileURL(*obj.Key)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get file URL")
			return nil, err
		}
		files = append(files, url)
	}

	return files, nil
}
