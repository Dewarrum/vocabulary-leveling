package subtitles

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	gosubs "github.com/martinlindhe/subtitles"
)

const (
	FailedToUpload   = "failed to upload"
	FailedToDownload = "failed to download"
	FailedToParse    = "failed to parse"
)

type FileStorage struct {
	s3Client        *s3.Client
	s3PresignClient *s3.PresignClient
}

func NewFileStorage(s3Client *s3.Client, s3PresignClient *s3.PresignClient) *FileStorage {
	return &FileStorage{
		s3Client:        s3Client,
		s3PresignClient: s3PresignClient,
	}
}

func (f *FileStorage) Upload(videoId uuid.UUID, body io.Reader, contentType string, context context.Context) error {
	_, err := f.s3Client.PutObject(context, &s3.PutObjectInput{
		Bucket:      aws.String("default"),
		Key:         aws.String(fmt.Sprintf("%s/subtitles", videoId)),
		Body:        body,
		ContentType: &contentType,
	})
	if err != nil {
		return errors.Join(err, errors.New(FailedToUpload))
	}

	return nil
}

func (f *FileStorage) Download(videoId string, context context.Context) (*gosubs.Subtitle, error) {
	response, err := f.s3Client.GetObject(context, &s3.GetObjectInput{
		Bucket: aws.String("default"),
		Key:    aws.String(fmt.Sprintf("%s/subtitles", videoId)),
	})
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToDownload))
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToDownload))
	}

	subtitle, err := gosubs.Parse(responseBody)
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToParse))
	}

	return &subtitle, nil
}
