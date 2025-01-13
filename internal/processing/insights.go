package processing

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/spencerrais/kafka_golang_testing/internal/logger"
)

func generateUserLoginInsights() map[string]interface{} {
	lock.Lock()
	defer lock.Unlock()

	insights := map[string]interface{}{
		"device_percentages":      map[string]float64{},
		"app_version_percentages": map[string]float64{},
		"unique_users":            len(uniqueUsers),
		"total_logins":            totalLogins,
		"timestamp":               time.Now().Unix(),
	}

	// Avoid division by zero
	if totalLogins == 0 {
		return insights
	}

	devicePercentages := make(map[string]float64)
	for device, count := range deviceCounts {
		devicePercentages[device] = (float64(count) / float64(totalLogins)) * 100
	}

	appVersionPercentages := make(map[string]float64)
	for version, count := range appVersionCounts {
		appVersionPercentages[version] = (float64(count) / float64(totalLogins)) * 100
	}

	insights["device_percentages"] = devicePercentages
	insights["app_version_percentages"] = appVersionPercentages

	// Reset state
	deviceCounts = make(map[string]int)
	appVersionCounts = make(map[string]int)
	uniqueUsers = make(map[string]struct{})
	totalLogins = 0

	return insights
}

func StartUserLoginInsightPublisher(producer sarama.SyncProducer, topic string) {
	go func() {
		log := logger.GetGlobalLogger("logs")
		for range time.Tick(1 * time.Minute) {
			data, err := json.Marshal(generateUserLoginInsights())
			if err != nil {
				err_str := fmt.Sprintf("Error marshalling insights: %v", err)
				log.Error(err_str)
				continue
			}
			// Make the json message a string and log it
			data_str := string(data)
			log.Info(fmt.Sprintf("Insights: %v", data_str))

			message := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.ByteEncoder(data),
			}

			if _, _, err := producer.SendMessage(message); err != nil {
				err_str := fmt.Sprintf("Error publishing insights: %v", err)
				log.Error(err_str)
			}
		}
	}()
}
