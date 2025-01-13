package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

type GlobalLogger struct {
	logger   *log.Logger
	file     *os.File
	basePath string
	mu       sync.Mutex
}

var (
	globalLogger *GlobalLogger
	once         sync.Once
)

// GetGlobalLogger initializes or retrieves the singleton logger
func GetGlobalLogger(basePath string) *GlobalLogger {
	once.Do(func() {
		globalLogger = &GlobalLogger{basePath: basePath}
		if err := globalLogger.updateLogFile(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
	})
	return globalLogger
}

func (g *GlobalLogger) updateLogFile() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Close the existing log file, if any
	if g.file != nil {
		g.file.Close()
	}

	// Generate a new log file name based on the current time
	now := time.Now()
	fileName := filepath.Join(g.basePath, fmt.Sprintf("%s.log", now.Format("2006-01-02_15")))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Create a multi-writer for console and file
	multiWriter := io.MultiWriter(file, os.Stdout)

	// Update the logger to write to the new file
	g.file = file
	g.logger = log.New(multiWriter, "", log.LstdFlags)
	return nil
}

// Log logs a message with the specified level
func (g *GlobalLogger) Log(level, message string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Check if the hour has changed, and update the log file if needed
	currentHour := time.Now().Format("2006-01-02_15")
	logFileHour := filepath.Base(g.file.Name())
	if logFileHour != fmt.Sprintf("%s.log", currentHour) {
		g.updateLogFile()
	}

	// Create a JSON log entry
	entry := LogEntry{
		Level:   level,
		Message: message,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling log entry: %v\n", err)
		return
	}

	g.logger.Println(string(data))
}

// Convenience methods for common log levels
func (g *GlobalLogger) Info(message string) {
	g.Log("INFO", message)
}

func (g *GlobalLogger) Debug(message string) {
	g.Log("DEBUG", message)
}

func (g *GlobalLogger) Error(message string) {
	g.Log("ERROR", message)
}

func (g *GlobalLogger) Fatal(message string) {
	g.Log("FATAL", message)
	os.Exit(1)
}
