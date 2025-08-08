package mock

import (
	"context"

	"10.1.20.130/dropping/log-management/pkg/dto"

	"github.com/stretchr/testify/mock"
)

type LogEmitterMock struct {
	mock.Mock
}

func (m *LogEmitterMock) EmitLog(ctx context.Context, msg dto.LogMessage) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}
