package batcher

import (
	"time"

	"example.com/eventbatch/internal/event"
)

// Config controls batch splitting.
type Config struct {
	MaxEvents int
	Now       func() time.Time
}

func (c Config) withDefaults() Config {
	if c.MaxEvents <= 0 {
		c.MaxEvents = 100
	}
	if c.Now == nil {
		c.Now = time.Now
	}
	return c
}

// Split partitions events into batches of at most MaxEvents.
func Split(events []event.Event, cfg Config) []event.Batch {
	cfg = cfg.withDefaults()
	if len(events) == 0 {
		return nil
	}

	var (
		out []event.Batch
		buf []event.Event
		seq int
	)

	flush := func() {
		if len(buf) == 0 {
			return
		}
		seq++
		// Copy buf into a fresh slice: each batch must own its backing
		// array, independent of buf. Reusing buf (via buf[:0]) would
		// alias the arrays, so later appends overwrite earlier batches'
		// events — the batches "bleed" into each other.
		cp := make([]event.Event, len(buf))
		copy(cp, buf)
		out = append(out, event.Batch{
			Seq:     seq,
			Events:  cp,
			Created: cfg.Now().UTC(),
		})
		buf = buf[:0]
	}

	for _, ev := range events {
		buf = append(buf, ev)
		if len(buf) >= cfg.MaxEvents {
			flush()
		}
	}
	flush()
	return out
}
