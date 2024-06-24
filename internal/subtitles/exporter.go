package subtitles

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/app"
	"log"
	"strings"
	"time"
)

type Exporter struct {
	MessageQueue          *MessageQueue
	SubtitleCueRepository *SubtitleCueRepository
	FileStorage           *FileStorage
}

func NewExporter(dependencies *app.Dependencies) (*Exporter, error) {
	fileStorage := NewFileStorage(dependencies.S3Client, dependencies.S3PresignClient)
	messageQueue, err := NewMessageQueue(dependencies.RabbitMqChannel)
	if err != nil {
		return nil, err
	}

	return &Exporter{
		MessageQueue:          messageQueue,
		SubtitleCueRepository: newSubtitleCueRepository(dependencies.Postgres),
		FileStorage:           fileStorage,
	}, nil
}

func (e *Exporter) Run(context context.Context) {
	log.Printf("Starting subtitle exporter")

	messages, err := e.MessageQueue.Consume()
	if err != nil {
		log.Fatal(err)
		return
	}

	for message := range messages {
		err = e.handleMessage(message, context)
		if err != nil {
			log.Printf("Failed to handle message: %v", err)
		}
	}
}

func (e *Exporter) handleMessage(message ExportSubtitlesMessage, context context.Context) error {
	subtitle, err := e.FileStorage.Download(message.VideoId.String(), context)
	if err != nil {
		return err
	}

	for _, caption := range subtitle.Captions {
		empty := (time.Time{}).AddDate(-1, 0, 0)
		subtitleCue := newDbSubtitleCue(message.VideoId, strings.Join(caption.Text, "\n"), caption.Seq, caption.Start.Sub(empty).Milliseconds(), caption.End.Sub(empty).Milliseconds())
		err := e.SubtitleCueRepository.Insert(subtitleCue)
		if err != nil {
			return err
		}
	}

	return nil
}
