package pkg

import (
	"context"
	"encoding/json"
	"fmt"

	"10.1.20.130/dropping/log-management/pkg/dto"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
)

type (
	LogEmitter interface {
		EmitLog(ctx context.Context, msg dto.LogMessage) error
	}
	logEmitter struct {
		subject string
		logger  zerolog.Logger
		js      jetstream.JetStream
	}
)

func NewLogEmitter(js jetstream.JetStream, logger zerolog.Logger, streamName, streamDescription, globalSubject, subjectPrefix string) LogEmitter {
	cfg := &jetstream.StreamConfig{
		Name:        streamName,
		Description: streamDescription,
		Subjects:    []string{globalSubject},
		MaxBytes:    6 * 1024 * 1024,
		Storage:     jetstream.FileStorage,
	}
	_, err := js.CreateOrUpdateStream(context.Background(), *cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create or update JetStream Log Stream")
	}
	subject := subjectPrefix
	return &logEmitter{js: js, subject: subject, logger: logger}
}

func (l *logEmitter) EmitLog(ctx context.Context, msg dto.LogMessage) error {
	marshalledMsg, err := json.Marshal(msg)
	if err != nil {
		l.logger.Error().Err(err)
		return err
	}
	sub := fmt.Sprintf("%s.%s", l.subject, msg.Service)
	if _, err := l.js.Publish(ctx, sub, marshalledMsg); err != nil {
		l.logger.Error().Err(err).Msg("failed to publish message")
		return err
	}
	return nil
}
