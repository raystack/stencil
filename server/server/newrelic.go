package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/odpf/stencil/server/config"
	"github.com/spf13/viper"
)

func getNewRelicMiddleware() gin.HandlerFunc {
	var newRelicConfig config.NewRelicConfig
	err := viper.UnmarshalKey("newrelic", &newRelicConfig)
	if err != nil {
		log.Fatal(err)
	}
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
