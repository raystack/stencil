package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/odpf/stencil/server/config"
)

func getNewRelicMiddleware(config *config.Config) gin.HandlerFunc {
	newRelicConfig := config.NewRelic
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(newRelicConfig.AppName),
		newrelic.ConfigLicense(newRelicConfig.License),
		newrelic.ConfigEnabled(newRelicConfig.Enabled),
	)
	if err != nil {
		log.Fatal(err)
	}
	return nrgin.Middleware(app)
}
