package mocks

import (
	"context"

	"github.com/micros-template/log-service/pkg/dto"

	"github.com/stretchr/testify/mock"
)

type LogEmitterMock struct {
	mock.Mock
}

func (m *LogEmitterMock) EmitLog(ctx context.Context, msg dto.LogMessage) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}
