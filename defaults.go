package main

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/connector/datadogconnector"
	"github.com/open-telemetry/opentelemetry-collector-contrib/connector/spanmetricsconnector"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awscloudwatchlogsexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsemfexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsxrayexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/loadbalancingexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logzioexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/sapmexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/awsproxy"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/ecsobserver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/sigv4authextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatorateprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbytraceprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricsgenerationprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsecscontainermetricsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxrayreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/debugexporter"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/extension/ballastextension"
	"go.opentelemetry.io/collector/extension/zpagesextension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.uber.org/multierr"
)

// Components register OTel components
func Components() (otelcol.Factories, error) {
	var errs error
	extensionsList := []extension.Factory{
		awsproxy.NewFactory(),
		ecsobserver.NewFactory(),
		healthcheckextension.NewFactory(),
		pprofextension.NewFactory(),
		sigv4authextension.NewFactory(),
		zpagesextension.NewFactory(),
		ballastextension.NewFactory(),
		filestorage.NewFactory(),
	}

	extensions, err := extension.MakeFactoryMap(extensionsList...)
	if err != nil {
		errs = multierr.Append(errs, err)
	}

	receiverList := []receiver.Factory{
		awsecscontainermetricsreceiver.NewFactory(),
		awscontainerinsightreceiver.NewFactory(),
		awsxrayreceiver.NewFactory(),
		statsdreceiver.NewFactory(),
		prometheusreceiver.NewFactory(),
		kafkareceiver.NewFactory(),
		jaegerreceiver.NewFactory(),
		zipkinreceiver.NewFactory(),
		otlpreceiver.NewFactory(),
		filelogreceiver.NewFactory(),
	}

	receivers, err := receiver.MakeFactoryMap(receiverList...)
	if err != nil {
		errs = multierr.Append(errs, err)
	}

	processorList := []processor.Factory{
		groupbytraceprocessor.NewFactory(),
		tailsamplingprocessor.NewFactory(),
		attributesprocessor.NewFactory(),
		resourceprocessor.NewFactory(),
		probabilisticsamplerprocessor.NewFactory(),
		spanprocessor.NewFactory(),
		filterprocessor.NewFactory(),
		metricstransformprocessor.NewFactory(),
		resourcedetectionprocessor.NewFactory(),
		metricsgenerationprocessor.NewFactory(),
		cumulativetodeltaprocessor.NewFactory(),
		deltatorateprocessor.NewFactory(),
		batchprocessor.NewFactory(),
		memorylimiterprocessor.NewFactory(),
		k8sattributesprocessor.NewFactory(),
	}
	processors, err := processor.MakeFactoryMap(processorList...)
	if err != nil {
		errs = multierr.Append(errs, err)
	}

	// enable the selected exporters
	exporterList := []exporter.Factory{
		awsemfexporter.NewFactory(),
		prometheusremotewriteexporter.NewFactory(),
		prometheusexporter.NewFactory(),
		fileexporter.NewFactory(),
		kafkaexporter.NewFactory(),
		sapmexporter.NewFactory(),
		signalfxexporter.NewFactory(),
		datadogexporter.NewFactory(),
		logzioexporter.NewFactory(),
		loggingexporter.NewFactory(),
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
		awsxrayexporter.NewFactory(),
		loadbalancingexporter.NewFactory(),
		awscloudwatchlogsexporter.NewFactory(),
		debugexporter.NewFactory(),
	}

	exporters, err := exporter.MakeFactoryMap(exporterList...)
	if err != nil {
		errs = multierr.Append(errs, err)
	}
	connectorList := []connector.Factory{
		datadogconnector.NewFactory(),
		spanmetricsconnector.NewFactory(),
	}

	connectors, err := connector.MakeFactoryMap(connectorList...)
	if err != nil {
		errs = multierr.Append(errs, err)
	}
	factories := otelcol.Factories{
		Extensions: extensions,
		Receivers:  receivers,
		Processors: processors,
		Exporters:  exporters,
		Connectors: connectors,
	}

	return factories, errs
}