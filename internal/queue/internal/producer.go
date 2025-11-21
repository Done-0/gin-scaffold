// Package internal provides Kafka producer implementation
// Author: Done-0
// Created: 2025-08-25
package internal

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/Done-0/gin-scaffold/configs"
)

// Producer handles Kafka message production
type Producer struct {
	producer sarama.SyncProducer
}

// NewProducer creates a new Kafka producer instance
func NewProducer(config *configs.Config) (*Producer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 3

	producer, err := sarama.NewSyncProducer(config.KafkaConfig.Brokers, kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{producer: producer}, nil
}

// Send sends a message to the specified topic
func (p *Producer) Send(ctx context.Context, topic string, key, value []byte) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	if key != nil {
		msg.Key = sarama.ByteEncoder(key)
	}

	return p.producer.SendMessage(msg)
}

// Close closes the producer connection
func (p *Producer) Close() error {
	return p.producer.Close()
}
