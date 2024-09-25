package cmd

import (
	"context"

	"github.com/elastic/beats/v7/x-pack/filebeat/fbreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/confmap/provider/httpprovider"
	"go.opentelemetry.io/collector/confmap/provider/httpsprovider"
	"go.opentelemetry.io/collector/confmap/provider/yamlprovider"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/debugexporter"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	"go.opentelemetry.io/collector/receiver"
)

func components() (otelcol.Factories, error) {
	receivers, err := receiver.MakeFactoryMap(
		fbreceiver.NewFactory(),
		filelogreceiver.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, nil
	}

	exporters, err := exporter.MakeFactoryMap(
		debugexporter.NewFactory(),
		elasticsearchexporter.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, nil
	}

	processors, err := processor.MakeFactoryMap(
		batchprocessor.NewFactory(),
		memorylimiterprocessor.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, nil
	}

	return otelcol.Factories{
		Receivers:  receivers,
		Exporters:  exporters,
		Processors: processors,
	}, nil

}

func OtelCmd() *cobra.Command {
	command := &cobra.Command{
		Short: "Run this to start filebeat as a otel",
		Use:   "otel",
		RunE: func(cmd *cobra.Command, args []string) error {
			info := component.BuildInfo{
				Command:     "otel",
				Description: "Beats OTeL",
				Version:     "9.0.0",
			}

			set := otelcol.CollectorSettings{
				BuildInfo: info,
				Factories: components,
				ConfigProviderSettings: otelcol.ConfigProviderSettings{
					ResolverSettings: confmap.ResolverSettings{
						URIs: []string{"file:/Users/vihasmakwana/Desktop/Vihas/elastic/elastic-agent/otel.yml"},
						ProviderFactories: []confmap.ProviderFactory{
							fileprovider.NewFactory(),
							httpprovider.NewFactory(),
							httpsprovider.NewFactory(),
							yamlprovider.NewFactory(),
						},
					},
				},
			}

			col, err := otelcol.NewCollector(set)
			if err != nil {
				return err
			}
			return col.Run(context.Background())
		},
	}
	return command
}
