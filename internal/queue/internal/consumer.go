// Package internal provides Kafka consumer implementation
// Author: Done-0
// Created: 2025-08-25
package internal

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/Done-0/gin-scaffold/configs"
)

// Handler defines the interface for processing consumed messages
type Handler interface {
	Handle(ctx context.Context, msg *sarama.ConsumerMessage) error
}

// Consumer handles Kafka message consumption
type Consumer struct {
	consumer sarama.ConsumerGroup
	handler  Handler
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewConsumer creates a new Kafka consumer instance
func NewConsumer(config *configs.Config, handler Handler) (*Consumer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(config.KafkaConfig.Brokers, config.KafkaConfig.ConsumerGroup, kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Consumer{
		consumer: consumerGroup,
		handler:  handler,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// Subscribe starts consuming messages from the specified topics
func (c *Consumer) Subscribe(topics []string) error {
	go func() {
		for {
			if err := c.consumer.Consume(c.ctx, topics, c); err != nil {
				return
			}
			if c.ctx.Err() != nil {
				return
			}
		}
	}()
	return nil
}

// Close gracefully shuts down the consumer
func (c *Consumer) Close() error {
	c.cancel()
	return c.consumer.Close()
}

// Setup implements sarama.ConsumerGroupHandler
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error { return nil }

// Cleanup implements sarama.ConsumerGroupHandler
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim implements sarama.ConsumerGroupHandler
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		c.handler.Handle(session.Context(), message)
		session.MarkMessage(message, "")
	}
	return nil
}
