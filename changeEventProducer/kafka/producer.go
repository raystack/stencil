package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/cactus/go-statsd-client/v5/statsd"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/proto"
)

type KafkaProducer struct {
	hostName     string
	timeout      int
	producer     *kafka.Producer
	statsdClient statsd.Statter
}

const (
	SuccessCountMetric   = "kafka.SuccessCount"
	FailureCountMetric   = "kafka.FailureCount"
	RetryExhaustedMetric = "kafka.RetryExhaustedMetric"
)

func NewKafkaProducer(hostName string, timeout int, statsdClient statsd.Statter) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": hostName,
		"message.timeout.ms": timeout})
	if err != nil {
		log.Printf("Failed to initialise Kafka Producer - %s", err.Error())
		return nil, err
	}

	return &KafkaProducer{producer: producer, hostName: hostName, timeout: timeout, statsdClient: statsdClient}, nil
}

func (kp *KafkaProducer) PushMessagesWithRetries(topic string, protoMessage proto.Message, retries int, retryInterval time.Duration) error {
	messageBytes, err := proto.Marshal(protoMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal Protobuf message: %s", err.Error())
	}

	for i := 0; i < retries; i++ {
		err := kp.PushMessages(messageBytes, topic)
		if err != nil {
			metricError := kp.statsdClient.Inc(FailureCountMetric, 1, 1)
			if metricError != nil {
				log.Printf("Failed to increase Failure metric - %s", metricError.Error())
			}
			time.Sleep(retryInterval)
			continue
		}
		err = kp.statsdClient.Inc(SuccessCountMetric, 1, 1)
		if err != nil {
			log.Printf("Failed to increase Success metric - %s", err.Error())
		}
		return nil
	}
	err = kp.statsdClient.Inc(RetryExhaustedMetric, 1, 1)
	if err != nil {
		log.Printf("Failed to increase retryExhausted metric - %s", err.Error())
	}
	return fmt.Errorf("failed to produce message after %d retries", retries)
}

func (kp *KafkaProducer) PushMessages(messageBytes []byte, topic string) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          messageBytes,
	}
	deliveryChan := make(chan kafka.Event)
	err := kp.producer.Produce(message, deliveryChan)
	if err != nil {
		log.Printf("Error in producing messages- %s", err.Error())
		return err
	}
	deliveryReport := <-deliveryChan
	if m, ok := deliveryReport.(*kafka.Message); ok && m.TopicPartition.Error != nil {
		log.Printf("Error in topic partitioning- %s", err.Error())
		return err
	}
	return nil
}
