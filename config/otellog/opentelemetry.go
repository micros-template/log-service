package otellog

import (
	"context"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

func NewOTELLoggerProvider() (*sdklog.LoggerProvider, error) {
	ctx := context.Background()
	exporter, err := otlploghttp.New(ctx, otlploghttp.WithEndpointURL(viper.GetString("otel.endpoint")))
	if err != nil {
		return nil, err
	}
	res, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", "main-logger-service"),
		),
	)
	if err != nil {
		return nil, err
	}
	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(res),
	)
	return provider, nil
}
