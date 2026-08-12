package sink

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"example.com/eventbatch/internal/event"
)

// Memory records batches in process for tests and dry-runs.
type Memory struct {
	mu      sync.Mutex
	Batches []event.Batch
}

func (m *Memory) Write(b event.Batch) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := event.Batch{
		Seq:     b.Seq,
		Created: b.Created,
		Events:  append([]event.Event(nil), b.Events...),
	}
	m.Batches = append(m.Batches, cp)
	return nil
}

// FileWriter writes one JSON file per batch.
type FileWriter struct {
	Dir string
}

func (f FileWriter) Write(b event.Batch) error {
	if err := os.MkdirAll(f.Dir, 0o755); err != nil {
		return err
	}
	name := filepath.Join(f.Dir, fmt.Sprintf("batch-%04d.json", b.Seq))
	raw, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(name, raw, 0o644)
}
