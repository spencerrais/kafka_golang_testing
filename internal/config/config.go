package config

import (
	"os"
	"strings"
)

type Config struct {
	KafkaBrokers []string
	InputTopic   string
	OutputTopic  string
	GroupID      string
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// parseBrokers splits the KAFKA_BROKERS environment variable into a slice of strings
func parseBrokers(brokers string) []string {
	return strings.Split(brokers, ",")
}

func LoadConfig() Config {
	return Config{
		KafkaBrokers: parseBrokers(getEnv("KAFKA_BROKERS", "kafka:9092")),
		InputTopic:   getEnv("INPUT_TOPIC", "user-login"),
		OutputTopic:  getEnv("OUTPUT_TOPIC", "processed-insights"),
		GroupID:      getEnv("GROUP_ID", "user-login-group"),
	}
}
