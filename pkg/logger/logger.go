package logger

import (
	"log"

	"go.uber.org/zap"
)

// Logger zap logger instance
var Logger *zap.Logger

func init() {
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatalln(err)
	}
	Logger = l
}
