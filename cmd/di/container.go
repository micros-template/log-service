package di

import (
	"github.com/micros-template/log-service/config/logger"
	mq "github.com/micros-template/log-service/config/message-queue"
	"github.com/micros-template/log-service/config/otellog"
	"github.com/micros-template/log-service/internal/domain/handler"
	"github.com/micros-template/log-service/internal/domain/service"
	_mq "github.com/micros-template/log-service/internal/infrastructure/message-queue"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	// logger
	if err := container.Provide(logger.New); err != nil {
		panic("Failed to provide logger: " + err.Error())
	}
	// nats connection
	if err := container.Provide(mq.New); err != nil {
		panic("Failed to provide message queue: " + err.Error())
	}
	// jetstream
	if err := container.Provide(jetstream.New); err != nil {
		panic("Failed to provide jetstream: " + err.Error())
	}
	// nats infra
	if err := container.Provide(_mq.NewNatsInfrastructure); err != nil {
		panic("Failed to provide message queue infra: " + err.Error())
	}
	// otellogger
	if err := container.Provide(otellog.NewOTELLoggerProvider); err != nil {
		panic("Failed to provide otellog provider: " + err.Error())
	}
	// logs subscriber service
	if err := container.Provide(service.NewLogSubscriberService); err != nil {
		panic("Failed to provide log subscriber service " + err.Error())
	}

	// Logs subscriber handler
	if err := container.Provide(handler.NewLogSubscriberHandler); err != nil {
		panic("Failed to provide log subscriber handler" + err.Error())
	}

	return container
}
