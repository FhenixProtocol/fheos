package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

type TelemetryCollector struct {
	isOn          bool
	telemetryPath string
	logger        log.Logger
	mutex         sync.Mutex
}

type FheOperationRequestTelemetry struct {
	TelemetryType string `json:"telemetryType"`
	OperationType string `json:"operationType"`
	ID            string `json:"id"`
	Handle        string `json:"handle"`
	Inputs        string `json:"inputs"`
}

type FheOperationUpdateTelemetry struct {
	TelemetryType  string `json:"telemetryType"`
	ID             string `json:"id"`
	InternalHandle string `json:"internalHandle"`
	Status         string `json:"status"`
}

func (t FheOperationUpdateTelemetry) SetStatus(status string) FheOperationUpdateTelemetry {
	t.Status = status
	return t
}

// NewTelemetryCollector creates a new telemetry collector
func NewTelemetryCollector(telemetryPath string) *TelemetryCollector {
	tc := &TelemetryCollector{}
	if telemetryPath != "" {
		tc.TurnOn(telemetryPath)
	}
	return tc
}

func (tc *TelemetryCollector) SetLogger(logger log.Logger) {
	tc.logger = logger
}

// IsCollectorOn returns whether the collector is active
func (tc *TelemetryCollector) IsCollectorOn() bool {
	return tc.isOn
}

// TurnOn activates the telemetry collector
func (tc *TelemetryCollector) TurnOn(telemetryPath string) error {
	tc.isOn = true
	tc.telemetryPath = telemetryPath

	// Create file if it doesn't exist
	if _, err := os.Stat(telemetryPath); os.IsNotExist(err) {
		if _, err := os.Create(telemetryPath); err != nil {
			return fmt.Errorf("failed to create telemetry file: %w", err)
		}
	}
	return nil
}

// TurnOff deactivates the telemetry collector
func (tc *TelemetryCollector) TurnOff() {
	tc.isOn = false
}

// AddTelemetry adds a telemetry entry to the file
func (tc *TelemetryCollector) AddTelemetry(telemetry interface{}) error {
	if !tc.isOn {
		return nil
	}

	data, err := json.Marshal(telemetry)
	if err != nil {
		return fmt.Errorf("failed to marshal telemetry: %w", err)
	}

	tc.logger.Info("Adding telemetry", "telemetry", string(data))
	entry := fmt.Sprintf("%s %s\n", time.Now().UTC().Format(time.RFC3339), string(data))

	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	f, err := os.OpenFile(tc.telemetryPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open telemetry file: %w", err)
	}
	
	defer f.Close()

	if _, err := f.WriteString(entry); err != nil {
		return fmt.Errorf("failed to append telemetry: %w", err)
	}

	return nil
}
