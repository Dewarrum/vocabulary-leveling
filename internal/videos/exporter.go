package videos

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/app"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type Exporter struct {
	FileStorage  *FileStorage
	MessageQueue *MessageQueue
}

func NewExporter(dependencies *app.Dependencies) (*Exporter, error) {
	messageQueue, err := NewMessageQueue(dependencies.RabbitMqChannel)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to create message queue"))
	}

	return &Exporter{
		FileStorage:  NewFileStorage(dependencies.S3Client, dependencies.S3PresignClient),
		MessageQueue: messageQueue,
	}, nil
}

func (e *Exporter) Run(context context.Context) {
	log.Print("Starting video exporter")

	messages, err := e.MessageQueue.Consume()
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
		return
	}

	for message := range messages {
		log.Printf("Exporting video %s", message.VideoId)

		err = e.handleMessage(message, context)
		if err != nil {
			log.Printf("Failed to export video %s: %s", message.VideoId, err)
			continue
		}

		log.Printf("Video exported successfully %s", message.VideoId)
	}
}

func (e *Exporter) handleMessage(message ExportVideoMessage, context context.Context) error {
	directory := fmt.Sprintf("tmp/%s", message.VideoId)
	err := os.Mkdir(directory, 0755)
	if err != nil {
		return errors.Join(err, errors.New("failed to create directory"))
	}

	err = e.downloadVideo(message.VideoId, directory, context)
	if err != nil {
		return errors.Join(err, errors.New("failed to download video"))
	}

	err = exportVideo(directory)
	if err != nil {
		return errors.Join(err, errors.New("failed to run ffmpeg"))
	}

	err = e.uploadVideo(message.VideoId, directory, context)
	if err != nil {
		return errors.Join(err, errors.New("failed to upload video"))
	}

	err = clearDirectory(directory)
	if err != nil {
		return errors.Join(err, errors.New("failed to clear directory"))
	}

	return nil
}

func (e *Exporter) downloadVideo(videoId, directory string, context context.Context) error {
	log.Printf("Downloading video %s", videoId)

	response, err := e.FileStorage.Download(videoId, context)
	if err != nil {
		return errors.Join(err, errors.New("failed to download video"))
	}

	log.Printf("Downloaded video %s", videoId)

	path := fmt.Sprintf("%s/original", directory)
	fi, err := os.Create(path)
	if err != nil {
		return errors.Join(err, errors.New("failed to create file"))
	}
	defer fi.Close()

	_, err = io.Copy(fi, response.Body)
	if err != nil {
		return errors.Join(err, errors.New("failed to download video"))
	}

	log.Printf("Saved downloaded video %s", videoId)

	return nil
}

func exportVideo(directory string) error {
	log.Printf("Running ffmpeg for video %s", directory)

	cmd := exec.Command(
		"ffmpeg",
		"-i", fmt.Sprintf("%s/original", directory),
		"-g", "30",
		"-keyint_min", "30",
		"-sc_threshold", "0",
		"-seg_duration", "5",
		"-use_template", "0",
		"-init_seg_name", "output-init-stream$RepresentationID$.$ext$",
		"-media_seg_name", "output-chunk-stream$RepresentationID$-$Number%05d$.$ext$",
		"-f", "dash",
		fmt.Sprintf("%s/manifest.mpd", directory))

	err := cmd.Run()
	if err != nil {
		return errors.Join(err, errors.New("failed to run ffmpeg"))
	}

	log.Printf("Sucessfully exported video %s to dash format", directory)

	return nil
}

func (e *Exporter) uploadVideo(videoId, directory string, context context.Context) error {
	log.Printf("Uploading video %s", videoId)

	entries, err := os.ReadDir(directory)
	if err != nil {
		return errors.Join(err, errors.New("failed to read directory"))
	}

	for _, entry := range entries {
		if entry.IsDir() {
			log.Printf("Skipping file %s", entry.Name())
			continue
		}

		if isChunkStream(entry.Name()) {
			file, err := os.Open(fmt.Sprintf("%s/%s", directory, entry.Name()))
			if err != nil {
				return errors.Join(err, errors.New("failed to open file"))
			}

			err = e.FileStorage.UploadChunkStream(videoId, entry.Name(), file, context)
			if err != nil {
				return errors.Join(err, errors.New("failed to upload chunk stream"))
			}
			continue
		}

		if isInitStream(entry.Name()) {
			file, err := os.Open(fmt.Sprintf("%s/%s", directory, entry.Name()))
			if err != nil {
				return errors.Join(err, errors.New("failed to open file"))
			}

			err = e.FileStorage.UploadInitStream(videoId, entry.Name(), file, context)
			if err != nil {
				return errors.Join(err, errors.New("failed to upload init stream"))
			}
		}

		if entry.Name() == "manifest.mpd" {
			file, err := os.Open(fmt.Sprintf("%s/%s", directory, entry.Name()))
			if err != nil {
				return errors.Join(err, errors.New("failed to open file"))
			}

			err = e.FileStorage.UploadManifest(videoId, file, context)
			if err != nil {
				return errors.Join(err, errors.New("failed to upload manifest"))
			}
		}

		log.Printf("Skipping unknown file %s", entry.Name())
	}

	log.Printf("Uploaded video %s", videoId)
	return nil
}

func clearDirectory(directory string) error {
	err := os.RemoveAll(directory)
	if err != nil {
		return errors.Join(err, errors.New("failed to clear directory"))
	}

	return nil
}
