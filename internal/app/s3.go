package app

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	ErrAWSAccessKeyIdIsRequired     = errors.New("AWS_ACCESS_KEY_ID is required")
	ErrAWSSecretAccessKeyIsRequired = errors.New("AWS_SECRET_ACCESS_KEY is required")
	ErrAWSRegionIsRequired          = errors.New("AWS_REGION is required")
	ErrAWSEndpointUrlIsRequired     = errors.New("AWS_ENDPOINT_URL is required")
)

func createS3Client() (*s3.Client, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKey == "" {
		return nil, ErrAWSAccessKeyIdIsRequired
	}
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secretKey == "" {
		return nil, ErrAWSSecretAccessKeyIsRequired
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, ErrAWSRegionIsRequired
	}
	endpointURL := os.Getenv("AWS_ENDPOINT_URL")
	if endpointURL == "" {
		return nil, ErrAWSEndpointUrlIsRequired
	}

	client := s3.New(s3.Options{
		Credentials:  credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:       region,
		BaseEndpoint: &endpointURL,
		UsePathStyle: true,
	})

	return client, nil
}

func createS3PresignClient(s3Client *s3.Client) *s3.PresignClient {
	return s3.NewPresignClient(s3Client)
}
