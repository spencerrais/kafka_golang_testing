package processing

import "sync"

import "github.com/spencerrais/kafka_golang_testing/internal/models"

var (
	deviceCounts     = make(map[string]int)
	appVersionCounts = make(map[string]int)
	uniqueUsers      = make(map[string]struct{})
	totalLogins      int
	lock             sync.Mutex
)

// ProcessUserLoginMessage updates the shared state for a login message
func ProcessUserLoginMessage(msg models.UserLoginMessage) {
	lock.Lock()
	defer lock.Unlock()

	// Increment totals and counts
	totalLogins++
	deviceCounts[msg.DeviceType]++
	appVersionCounts[msg.AppVersion]++

	// Track unique users
	uniqueUsers[msg.UserID] = struct{}{}
}
