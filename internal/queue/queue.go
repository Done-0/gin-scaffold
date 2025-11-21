// Package queue provides Kafka message queue functionality
// Author: Done-0
// Created: 2025-08-25
package queue

import (
	"context"

	"github.com/IBM/sarama"

	"github.com/Done-0/gin-scaffold/configs"
	"github.com/Done-0/gin-scaffold/internal/queue/internal"
)

// Producer defines the interface for Kafka message production
type Producer interface {
	Send(ctx context.Context, topic string, key, value []byte) (partition int32, offset int64, err error)
	Close() error
}

// Consumer defines the interface for Kafka message consumption
type Consumer interface {
	Subscribe(topics []string) error
	Close() error
}

// Handler defines the interface for message processing
type Handler interface {
	Handle(ctx context.Context, msg *sarama.ConsumerMessage) error
}

// NewProducer creates a new Kafka producer
func NewProducer(config *configs.Config) (Producer, error) {
	return internal.NewProducer(config)
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(config *configs.Config, handler Handler) (Consumer, error) {
	return internal.NewConsumer(config, handler)
}
