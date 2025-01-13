package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/spencerrais/kafka_golang_testing/internal/logger"
)

type Consumer struct {
	Brokers    []string
	InputTopic string
	GroupID    string
	Process    func([]byte) error
}

func NewConsumer(brokers []string, inputTopic, groupID string, process func([]byte) error) *Consumer {
	return &Consumer{
		Brokers:    brokers,
		InputTopic: inputTopic,
		GroupID:    groupID,
		Process:    process,
	}
}

func (c *Consumer) Start() error {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Version = sarama.V2_8_0_0

	client, err := sarama.NewConsumerGroup(c.Brokers, c.GroupID, config)

	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	handler := &ConsumerHandler{Process: c.Process}
	return client.Consume(ctx, []string{c.InputTopic}, handler)
}

type ConsumerHandler struct {
	Process func([]byte) error
}

func (h *ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log := logger.GetGlobalLogger("logs")
	for message := range claim.Messages() {
		if err := h.Process(message.Value); err != nil {
			err_str := fmt.Sprintf("Error processing message: %v", err)
			log.Error(err_str)
		}
		session.MarkMessage(message, "")
	}
	return nil
}
