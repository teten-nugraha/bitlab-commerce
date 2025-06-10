package eventbus

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaEventBus struct {
	writer *kafka.Writer
}

func NewKafkaEventBus(brokers []string) *KafkaEventBus {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		Async:        true,
	}

	return &KafkaEventBus{
		writer: writer,
	}
}

func (k *KafkaEventBus) Publish(topic string, event interface{}) error {
	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Topic: topic,
			Value: message,
		},
	)

	if err != nil {
		log.Printf("failed to write message to kafka: %v", err)
		return err
	}

	return nil
}

func (k *KafkaEventBus) Close() error {
	return k.writer.Close()
}
