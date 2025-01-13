package processing

import (
	"encoding/json"

	"github.com/spencerrais/kafka_golang_testing/internal/logger"
	"github.com/spencerrais/kafka_golang_testing/internal/models"
)

func ProcessUserLogin(data []byte) error {
	log := logger.GetGlobalLogger("logs")

	var msg models.UserLoginMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Error("Failed to unmarshal user login message")
		return err
	}

	// Update shared state
	ProcessUserLoginMessage(msg)
	return nil
}
