package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/micros-template/log-service/pkg/dto"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type (
	LogSubscriberService interface {
		SendLog(msg dto.LogMessage) error
	}
	logSubscriberService struct {
		logger     zerolog.Logger
		otelLogger *slog.Logger
	}
)

func NewLogSubscriberService(logger zerolog.Logger, provider *sdklog.LoggerProvider) LogSubscriberService {
	otelLogger := otelslog.NewLogger("main-logger", otelslog.WithLoggerProvider(provider))
	return &logSubscriberService{
		logger:     logger,
		otelLogger: otelLogger,
	}
}

func (l *logSubscriberService) SendLog(msg dto.LogMessage) error {
	localTime := time.Now().Format(time.RFC3339)
	switch msg.Type {
	case "INFO":
		l.otelLogger.InfoContext(context.Background(), msg.Msg,
			slog.String("service.name", msg.Service),
			slog.String("protocol", msg.Protocol),
			slog.String("local_time", localTime),
		)
	case "WARN":
		l.otelLogger.WarnContext(context.Background(), msg.Msg,
			slog.String("service.name", msg.Service),
			slog.String("protocol", msg.Protocol),
			slog.String("local_time", localTime),
		)
	case "ERR":
		l.otelLogger.ErrorContext(context.Background(), msg.Msg,
			slog.String("service.name", msg.Service),
			slog.String("protocol", msg.Protocol),
			slog.String("local_time", localTime),
		)
	default:
		l.otelLogger.DebugContext(context.Background(), msg.Msg,
			slog.String("service.name", msg.Service),
			slog.String("protocol", msg.Protocol),
			slog.String("local_time", localTime),
		)
	}
	return nil
}
