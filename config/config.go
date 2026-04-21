package config

import "time"

// DBConfig contains DB connection details
type DBConfig struct {
	ConnectionString string
}

// CORSConfig contains CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `default:"[\"*\"]"`
}

// Config Server config
type Config struct {
	Port string `default:"8080"`
	// Timeout represents graceful shutdown period. Defaults to 60 seconds.
	Timeout        time.Duration `default:"60s"`
	CacheSizeInMB  int64         `default:"100"`
	MaxRecvMsgSize int           `default:"10485760"` // 10 MB
	MaxSendMsgSize int           `default:"10485760"` // 10 MB
	CORS           CORSConfig
	DB             DBConfig
}
