package kafka

import (
	"log"

	"github.com/cactus/go-statsd-client/v5/statsd"
)

type MetricsCollector struct {
	StatsdClient statsd.Statter
}

func NewMetricsCollector(statsdAddr string) (*MetricsCollector, error) {
	statsdClient, err := statsd.NewClient(statsdAddr, "")
	if err != nil {
		log.Printf("Failed to initialise statsd client- %s", err.Error())
		return nil, err
	}

	return &MetricsCollector{
		StatsdClient: statsdClient,
	}, nil
}
