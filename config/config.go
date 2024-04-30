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

// GRPCConfig grpc options
type GRPCConfig struct {
	MaxRecvMsgSizeInMB int `default:"10"`
	MaxSendMsgSizeInMB int `default:"10"`
}

// Kafka Producer Config
type KafkaProducerConfig struct {
	BootstrapServer string
	Retries         int `default:"5"`
	Timeout         int `default:"5000"`
}

// StatsDConfig
type StatsDConfig struct {
	Address string
	Prefix  string
}

// SchameChangeConfig
type SchemaChangeConfig struct {
	KafkaTopic string
	Depth      int32
	Enable     bool `default:"false"`
}

// Config Server config
type Config struct {
	Port string `default:"8080"`
	// Timeout represents graceful shutdown period. Defaults to 60 seconds.
	Timeout       time.Duration `default:"60s"`
	CacheSizeInMB int64         `default:"100"`
	GRPC          GRPCConfig
	NewRelic      NewRelicConfig
	DB            DBConfig
	KafkaProducer KafkaProducerConfig
	StatsD        StatsDConfig
	SchemaChange  SchemaChangeConfig
}
