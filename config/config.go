package config

import "time"

// NewRelicConfig contains the New Relic go-agent configuration
type NewRelicConfig struct {
	Enabled bool   `default:"false"`
	AppName string `default:"stencil"`
	License string
}

// DBConfig contains DB connection details
type DBConfig struct {
	ConnectionString string
}

//GRPCConfig grpc options
type GRPCConfig struct {
	MaxRecvMsgSizeInMB int `default:"10"`
	MaxSendMsgSizeInMB int `default:"10"`
}

//Config Server config
type Config struct {
	Port string `default:"8080"`
	// Timeout represents graceful shutdown period. Defaults to 60 seconds.
	Timeout  time.Duration `default:"60s"`
	GRPC     GRPCConfig
	NewRelic NewRelicConfig
	DB       DBConfig
}
