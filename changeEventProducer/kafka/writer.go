package kafka

import (
	"context"
	"fmt"
	"github.com/cactus/go-statsd-client/v5/statsd"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

const (
	SuccessCountMetric = "kafka.SuccessCount"
	FailureCountMetric = "kafka.FailureCount"
)

type Writer struct {
	kafkaWriter  *kafka.Writer
	statsdClient statsd.Statter
}

func NewWriter(kafkaBrokerUrl string, timeout int, retries int, statsdClient statsd.Statter) *Writer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kafkaBrokerUrl),
		WriteTimeout: time.Second * time.Duration(timeout),
		MaxAttempts:  retries,
	}

	return &Writer{kafkaWriter: writer, statsdClient: statsdClient}
}

func (w *Writer) Close() error {
	return w.kafkaWriter.Close()
}

func (w *Writer) Write(topic string, protoMessage proto.Message) error {
	messageBytes, err := proto.Marshal(protoMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal Protobuf message: %s", err.Error())
	}

	kafkaMessage := kafka.Message{
		Topic: topic,
		Value: messageBytes,
	}
	return w.send(kafkaMessage)
}

func (w *Writer) send(message kafka.Message) error {
	err := w.kafkaWriter.WriteMessages(context.Background(), message)
	if err != nil {
		metricsErr := w.statsdClient.Inc(FailureCountMetric, 1, 1)
		if metricsErr != nil {
			log.Printf("Failed to increase Failure metric - %s", err.Error())
		}
		return err
	}
	err = w.statsdClient.Inc(SuccessCountMetric, 1, 1)
	if err != nil {
		log.Printf("Failed to increase Success metric - %s", err.Error())
	}
	return nil
}
