package nexusv

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

// TelemetryEvent represents an anonymous usage event.
type TelemetryEvent struct {
	Event     string `json:"event"`
	Template  string `json:"template"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

const (
	telemetryURL = "https://telemetry.stillsystems.io/v1/event"
)

// SendTelemetry sends an anonymous usage event if telemetry is enabled.
func SendTelemetry(cfg Config, event, template string, version string) {
	if !cfg.Telemetry.Enabled {
		return
	}

	payload := TelemetryEvent{
		Event:     event,
		Template:  template,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Version:   version,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	// Send in background to not block the user
	go func() {
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Post(telemetryURL, "application/json", bytes.NewBuffer(data))
		if err == nil {
			defer resp.Body.Close()
		}
	}()
}

