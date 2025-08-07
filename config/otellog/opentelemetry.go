package otellog

import (
	"context"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func NewOTELLoggerProvider(ctx context.Context) (*sdklog.LoggerProvider, error) {
	exporter, err := otlploghttp.New(ctx, otlploghttp.WithEndpointURL(viper.GetString("otel.endpoint")))
	if err != nil {
		return nil, err
	}
	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)
	return provider, nil
}
