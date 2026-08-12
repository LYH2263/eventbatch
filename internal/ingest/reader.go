package ingest

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"example.com/eventbatch/internal/event"
)

// Reader loads JSONL events from an io.Reader.
type Reader struct {
	r io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

func FromFile(path string) ([]event.Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewReader(f).ReadAll()
}

func (r *Reader) ReadAll() ([]event.Event, error) {
	sc := bufio.NewScanner(r.r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var out []event.Event
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev event.Event
		if err := json.Unmarshal(line, &ev); err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNo, err)
		}
		if ev.Timestamp.IsZero() {
			ev.Timestamp = time.Unix(0, 0).UTC()
		}
		out = append(out, ev)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
