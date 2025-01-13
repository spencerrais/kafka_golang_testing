package main

import (
	"fmt"
	"os"

	"github.com/spencerrais/kafka_golang_testing/internal/config"
	"github.com/spencerrais/kafka_golang_testing/internal/kafka"
	"github.com/spencerrais/kafka_golang_testing/internal/logger"
	"github.com/spencerrais/kafka_golang_testing/internal/processing"
)

func main() {
	// Create logs directory if it doesn't exist
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logger.GetGlobalLogger(logDir).Error("Failed to create logs directory")
		return
	}

	// Set up global logger
	log := logger.GetGlobalLogger(logDir)
	log.Info("Starting consumer")

	// Load configuration
	cfg := config.LoadConfig()
	for k, v := range cfg.KafkaBrokers {
		log.Info(fmt.Sprintf("Kafka broker %d: %s", k, v))
	}
	log.Info(cfg.GroupID)
	log.Info(cfg.InputTopic)
	log.Info(cfg.OutputTopic)

	// Kafka producer
	producer, err := kafka.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		err_str := fmt.Sprintf("Failed to create producer: %v", err)
		log.Fatal(err_str)
	}
	defer producer.Close()
	log.Info("Producer created")

	// Start periodic insight publishing
	log.Info("Starting user login insight publisher")
	go processing.StartUserLoginInsightPublisher(producer, cfg.OutputTopic)

	// Start user-login consumer
	log.Info("Starting user login consumer")
	userLoginConsumer := kafka.NewConsumer(cfg.KafkaBrokers, "user-login", "user-login-group", processing.ProcessUserLogin)
	if err := userLoginConsumer.Start(); err != nil {
		err_str := fmt.Sprintf("User login consumer error: %v", err)
		log.Fatal(err_str)
	}
}
