package main

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/otelcol"
)

func NewCollector(ctx context.Context, cfg otelcol.ConfigProvider) (*otelcol.Collector, error) {
	// Create a new Collector
	c, err := otelcol.NewCollector(
		otelcol.CollectorSettings{
			Factories: Components,

			BuildInfo: component.BuildInfo{
				Version: "0.0.1",
				Command: "otelcol",
			},
			// ConfigProvider: provider,
			ConfigProvider: cfg,
		},
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}
