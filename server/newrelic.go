package server

import (
	"log"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func getNewRelic(config *Config) *newrelic.Application {
	newRelicConfig := config.NewRelic
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(newRelicConfig.AppName),
		newrelic.ConfigLicense(newRelicConfig.License),
		newrelic.ConfigEnabled(newRelicConfig.Enabled),
	)
	if err != nil {
		log.Fatal(err)
	}
	return app
}
