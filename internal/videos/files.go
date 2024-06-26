package videos

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/mpd"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
)

const (
	FailedToUpload   = "failed to upload"
	FailedToDownload = "failed to download"
	FailedToList     = "failed to list"
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
		Key:         aws.String(fmt.Sprintf("%s/original", videoId)),
		Body:        body,
		ContentType: &contentType,
	})
	if err != nil {
		return errors.Join(err, errors.New(FailedToUpload))
	}

	return nil
}

func (f *FileStorage) Download(videoId uuid.UUID, context context.Context) (*s3.GetObjectOutput, error) {
	result, err := f.s3Client.GetObject(context, &s3.GetObjectInput{
		Bucket: aws.String("default"),
		Key:    aws.String(fmt.Sprintf("%s/original", videoId)),
	})
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToDownload))
	}

	return result, nil
}

func (f *FileStorage) UploadChunkStream(videoId uuid.UUID, chunkStreamName string, body io.Reader, context context.Context) error {
	contentType := "video/iso.segment"
	_, err := f.s3Client.PutObject(context, &s3.PutObjectInput{
		Bucket:      aws.String("default"),
		Key:         aws.String(fmt.Sprintf("%s/chunks/%s", videoId, chunkStreamName)),
		Body:        body,
		ContentType: &contentType,
	})
	if err != nil {
		return errors.Join(err, errors.New(FailedToUpload))
	}

	return nil
}

func (f *FileStorage) ListChunkStreams(videoId uuid.UUID, context context.Context) (*s3.ListObjectsV2Output, error) {
	response, err := f.s3Client.ListObjectsV2(context, &s3.ListObjectsV2Input{
		Bucket: aws.String("default"),
		Prefix: aws.String(fmt.Sprintf("%s/chunks", videoId)),
	})
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToList))
	}

	return response, nil
}

func (f *FileStorage) PresignChunkStreams(videoId uuid.UUID, context context.Context) ([]string, error) {
	response, err := f.ListChunkStreams(videoId, context)
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToList))
	}

	var presignedUrls []string
	for _, chunkStream := range response.Contents {
		presignedUrl, err := f.s3PresignClient.PresignGetObject(context, &s3.GetObjectInput{
			Bucket: aws.String("default"),
			Key:    aws.String(*chunkStream.Key),
		})
		if err != nil {
			return nil, errors.Join(err, errors.New("failed to presign object"))
		}
		presignedUrls = append(presignedUrls, presignedUrl.URL)
	}

	return presignedUrls, nil
}

func (f *FileStorage) DownloadChunkStream(videoId uuid.UUID, chunkStreamName string, context context.Context) (*s3.GetObjectOutput, error) {
	result, err := f.s3Client.GetObject(context, &s3.GetObjectInput{
		Bucket: aws.String("default"),
		Key:    aws.String(fmt.Sprintf("%s/chunks/%s", videoId, chunkStreamName)),
	})
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToDownload))
	}

	return result, nil
}

func (f *FileStorage) UploadInitStream(videoId uuid.UUID, initStreamName string, body io.Reader, context context.Context) error {
	contentType := "video/iso.segment"
	_, err := f.s3Client.PutObject(context, &s3.PutObjectInput{
		Bucket:      aws.String("default"),
		Key:         aws.String(fmt.Sprintf("%s/init/%s", videoId, initStreamName)),
		Body:        body,
		ContentType: &contentType,
	})
	if err != nil {
		return errors.Join(err, errors.New(FailedToUpload))
	}

	return nil
}

func (f *FileStorage) ListInitStreams(videoId uuid.UUID, context context.Context) (*s3.ListObjectsV2Output, error) {
	response, err := f.s3Client.ListObjectsV2(context, &s3.ListObjectsV2Input{
		Bucket: aws.String("default"),
		Prefix: aws.String(fmt.Sprintf("%s/init", videoId)),
	})
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToList))
	}

	return response, nil
}

func (f *FileStorage) PresignInitStreams(videoId uuid.UUID, context context.Context) ([]string, error) {
	response, err := f.ListInitStreams(videoId, context)
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToList))
	}

	var presignedUrls []string
	for _, initStream := range response.Contents {
		presignedUrl, err := f.s3PresignClient.PresignGetObject(context, &s3.GetObjectInput{
			Bucket: aws.String("default"),
			Key:    aws.String(*initStream.Key),
		})
		if err != nil {
			return nil, errors.Join(err, errors.New("failed to presign object"))
		}
		presignedUrls = append(presignedUrls, presignedUrl.URL)
	}

	return presignedUrls, nil
}

func (f *FileStorage) DownloadInitStream(videoId uuid.UUID, initStreamName string, context context.Context) (*s3.GetObjectOutput, error) {
	result, err := f.s3Client.GetObject(context, &s3.GetObjectInput{
		Bucket: aws.String("default"),
		Key:    aws.String(fmt.Sprintf("%s/init/%s", videoId, initStreamName)),
	})
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToDownload))
	}

	return result, nil
}

func (f *FileStorage) UploadManifest(videoId uuid.UUID, body io.Reader, context context.Context) error {
	contentType := "application/dash+xml"
	_, err := f.s3Client.PutObject(context, &s3.PutObjectInput{
		Bucket:      aws.String("default"),
		Key:         aws.String(fmt.Sprintf("%s/manifest.mpd", videoId)),
		Body:        body,
		ContentType: &contentType,
	})
	if err != nil {
		return errors.Join(err, errors.New(FailedToUpload))
	}

	return nil
}

func (f *FileStorage) DownloadManifest(videoId uuid.UUID, context context.Context) (*mpd.MPD, error) {
	response, err := f.s3Client.GetObject(context, &s3.GetObjectInput{
		Bucket: aws.String("default"),
		Key:    aws.String(fmt.Sprintf("%s/manifest", videoId)),
	})
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToDownload))
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return mpd.Parse(responseBody)
}
