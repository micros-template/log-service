package handler

import (
	"encoding/json"

	"10.1.20.130/dropping/log-management/internal/domain/dto"
	"10.1.20.130/dropping/log-management/internal/domain/service"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
)

type (
	LogSubscriberHandler interface {
		LogHandler(msg jetstream.Msg) error
	}
	logSubscriberHandler struct {
		subsService service.LogSubscriberService
		logger      zerolog.Logger
	}
)

func NewLogSubscriberHandler(svc service.LogSubscriberService, logger zerolog.Logger) LogSubscriberHandler {
	return &logSubscriberHandler{
		subsService: svc,
		logger:      logger,
	}
}

func (s *logSubscriberHandler) LogHandler(msg jetstream.Msg) error {
	var msgData dto.LogMessage
	err := json.Unmarshal(msg.Data(), &msgData)
	if err != nil {
		s.logger.Error().Err(err).Msg("error unmarshal")
		return err
	}
	if err = s.subsService.SendLog(msgData); err != nil {
		s.logger.Error().Err(err).Msg("failed to send log")
		return err
	}
	return nil
}
