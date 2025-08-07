package server

import (
	"context"
	"log"

	"10.1.20.130/dropping/log-management/internal/domain/handler"
	mq "10.1.20.130/dropping/log-management/internal/infrastructure/message-queue"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.uber.org/dig"
)

type LogSubscriber struct {
	Container       *dig.Container
	ConnectionReady chan bool
}

func (s *LogSubscriber) Run(ctx context.Context) {
	err := s.Container.Invoke(func(
		logger zerolog.Logger,
		sl handler.LogSubscriberHandler,
		js jetstream.JetStream,
		mq mq.Nats,
		_mq *nats.Conn,
	) {
		defer _mq.Drain()
		err := mq.CreateOrUpdateNewStream(ctx, &jetstream.StreamConfig{
			Name:        viper.GetString("jetstream.stream.name"),
			Description: viper.GetString("jetstream.stream.description"),
			Subjects:    []string{viper.GetString("jetstream.subject.global")},
			MaxBytes:    10 * 1024 * 1024,
			Storage:     jetstream.FileStorage,
		})
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create or update log stream")
		}

		// consumer for email
		logCons, err := mq.CreateOrUpdateNewConsumer(ctx, viper.GetString("jetstream.stream.name"), &jetstream.ConsumerConfig{
			Name:          viper.GetString("jetstream.consumer.log"),
			Durable:       viper.GetString("jetstream.consumer.log"),
			FilterSubject: viper.GetString("jetstream.subject.log"),
			AckPolicy:     jetstream.AckExplicitPolicy,
			DeliverPolicy: jetstream.DeliverNewPolicy,
		})

		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create or update log consumer")
		}

		_, err = logCons.Consume(func(msg jetstream.Msg) {
			go func() {
				sl.LogHandler(msg)
				msg.Ack()
			}()
		})

		if err != nil {
			logger.Error().Err(err).Msg("Failed to consume log consumer")
			return
		}

		if s.ConnectionReady != nil {
			s.ConnectionReady <- true
		}

		logger.Info().Msg("subscriber for log is running")
		<-ctx.Done()
		logger.Info().Msg("Shutting down subscriber for log")
	})
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
}
