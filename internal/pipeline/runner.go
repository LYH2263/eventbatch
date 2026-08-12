package pipeline

import (
	"fmt"

	"example.com/eventbatch/internal/batcher"
	"example.com/eventbatch/internal/event"
	"example.com/eventbatch/internal/sink"
)

type Writer interface {
	Write(b event.Batch) error
}

// Runner splits events and writes batches to a sink.
type Runner struct {
	Cfg    batcher.Config
	Writer Writer
}

func (r Runner) Run(events []event.Event) (int, error) {
	if r.Writer == nil {
		r.Writer = &sink.Memory{}
	}
	batches := batcher.Split(events, r.Cfg)
	for _, b := range batches {
		if err := r.Writer.Write(b); err != nil {
			return 0, fmt.Errorf("write batch %d: %w", b.Seq, err)
		}
	}
	return len(batches), nil
}
