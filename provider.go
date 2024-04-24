package main

import (
	"context"
	"fmt"
	"os"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/converter/expandconverter"
	"go.opentelemetry.io/collector/confmap/provider/envprovider"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/confmap/provider/httpprovider"
	"go.opentelemetry.io/collector/confmap/provider/yamlprovider"
	"go.opentelemetry.io/collector/otelcol"
)

type configProvider struct {
	file string
}

var _ otelcol.ConfigProvider = (*configProvider)(nil)

func NewProvider(file string) otelcol.ConfigProvider {
	return &configProvider{file: file}
}

func (p *configProvider) Get(ctx context.Context, factories otelcol.Factories) (*otelcol.Config, error) {
	config, err := LoadConfig(p.file, factories)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return nil, err
	}
	fmt.Printf("Resolved configuration: %#v\n", config)

	tp, err := component.NewType("datadog")
	if err != nil {
		return nil, err
	}
	id := component.NewID(tp)

	// Add datadog exporter
	config.Exporters[id] = &datadogexporter.Config{
		Metrics: datadogexporter.MetricsConfig{
			TCPAddrConfig: confignet.TCPAddrConfig{
				Endpoint: "https://api.datadoghq.com",
			},
			DeltaTTL: 3600,
			ExporterConfig: datadogexporter.MetricsExporterConfig{
				ResourceAttributesAsTags:           false,
				InstrumentationScopeMetadataAsTags: false,
			},
			HistConfig: datadogexporter.HistogramConfig{
				Mode:             "distributions",
				SendAggregations: false,
			},
			SumConfig: datadogexporter.SumConfig{
				CumulativeMonotonicMode:        datadogexporter.CumulativeMonotonicSumModeToDelta,
				InitialCumulativeMonotonicMode: datadogexporter.InitialValueModeAuto,
			},
			SummaryConfig: datadogexporter.SummaryConfig{
				Mode: datadogexporter.SummaryModeGauges,
			},
		},

		Traces: datadogexporter.TracesConfig{
			TCPAddrConfig: confignet.TCPAddrConfig{
				Endpoint: "https://trace.agent.datadoghq.com",
			},
			IgnoreResources: []string{},
		},

		Logs: datadogexporter.LogsConfig{
			TCPAddrConfig: confignet.TCPAddrConfig{
				Endpoint: "https://http-intake.logs.datadoghq.com",
			},
		},

		HostMetadata: datadogexporter.HostMetadataConfig{
			Enabled:        true,
			HostnameSource: datadogexporter.HostnameSourceConfigOrSystem,
		},
		TagsConfig: datadogexporter.TagsConfig{
			Hostname: "my-awesome-hostname",
		},
		API: datadogexporter.APIConfig{
			Key:  configopaque.String(os.Getenv("DD_API_KEY")),
			Site: "datadoghq.com",
		},
	}

	err = config.Validate()
	if err != nil {
		fmt.Printf("Error validating configuration: %v\n", err)
		return nil, err
	}
	for n, c := range config.Service.Pipelines {
		fmt.Printf("Pipeline %s: %v\n", n, c)
		c.Exporters = append(c.Exporters, id)
	}

	return config, nil
}

// Watch blocks until any configuration change was detected or an unrecoverable error
// happened during monitoring the configuration changes.
//
// Error is nil if the configuration is changed and needs to be re-fetched. Any non-nil
// error indicates that there was a problem with watching the config changes.
//
// Should never be called concurrently with itself or Get.
func (p *configProvider) Watch() <-chan error {
	ch := make(chan error)
	return ch
}

// Shutdown signals that the provider is no longer in use and the that should close
// and release any resources that it may have created.
//
// This function must terminate the Watch channel.
//
// Should never be called concurrently with itself or Get.
func (p *configProvider) Shutdown(ctx context.Context) error {
	return nil
}

// LoadConfig loads a config.Config
func LoadConfig(fileName string, factories otelcol.Factories) (*otelcol.Config, error) {
	// Read yaml config from file
	set := confmap.ProviderSettings{}
	provider, err := otelcol.NewConfigProvider(otelcol.ConfigProviderSettings{
		ResolverSettings: confmap.ResolverSettings{
			URIs: []string{fileName},
			Providers: makeMapProvidersMap(
				fileprovider.NewWithSettings(set),
				envprovider.NewWithSettings(set),
				yamlprovider.NewWithSettings(set),
				httpprovider.NewWithSettings(set),
			),
			Converters: []confmap.Converter{expandconverter.New(confmap.ConverterSettings{})},
		},
	})
	if err != nil {
		return nil, err
	}
	cfg, err := provider.Get(context.Background(), factories)
	if err != nil {
		return nil, err
	}
	err = cfg.Validate()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func makeMapProvidersMap(providers ...confmap.Provider) map[string]confmap.Provider {
	ret := make(map[string]confmap.Provider, len(providers))
	for _, provider := range providers {
		ret[provider.Scheme()] = provider
	}
	return ret
}
