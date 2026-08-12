package event

import "time"

// Event is one telemetry / audit record.
type Event struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Payload   map[string]string `json:"payload"`
	Timestamp time.Time         `json:"ts"`
}

// Batch is a group of events flushed together.
type Batch struct {
	Seq     int
	Events  []Event
	Created time.Time
}
