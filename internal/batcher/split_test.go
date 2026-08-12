package batcher_test

import (
	"fmt"
	"testing"
	"time"

	"example.com/eventbatch/internal/batcher"
	"example.com/eventbatch/internal/event"
)

func ev(id string) event.Event {
	return event.Event{
		ID:        id,
		Type:      "t",
		Payload:   map[string]string{"k": id},
		Timestamp: time.Unix(1, 0).UTC(),
	}
}

func ids(b event.Batch) []string {
	out := make([]string, len(b.Events))
	for i, e := range b.Events {
		out[i] = e.ID
	}
	return out
}

func TestSplitEmpty(t *testing.T) {
	if got := batcher.Split(nil, batcher.Config{MaxEvents: 2}); got != nil {
		t.Fatalf("got %#v", got)
	}
}

func TestSplitCounts(t *testing.T) {
	in := []event.Event{ev("a"), ev("b"), ev("c"), ev("d"), ev("e")}
	got := batcher.Split(in, batcher.Config{MaxEvents: 2})
	if len(got) != 3 {
		t.Fatalf("batches=%d want 3", len(got))
	}
	if len(got[2].Events) != 1 {
		t.Fatalf("last batch len=%d want 1", len(got[2].Events))
	}
}

func TestSplitPreservesEachBatchContents(t *testing.T) {
	in := make([]event.Event, 0, 6)
	for i := 1; i <= 6; i++ {
		in = append(in, ev(fmt.Sprintf("e%d", i)))
	}

	got := batcher.Split(in, batcher.Config{MaxEvents: 2, Now: func() time.Time { return time.Unix(0, 0).UTC() }})
	if len(got) != 3 {
		t.Fatalf("batches=%d want 3", len(got))
	}

	want := [][]string{
		{"e1", "e2"},
		{"e3", "e4"},
		{"e5", "e6"},
	}
	for i := range want {
		gotIDs := ids(got[i])
		if len(gotIDs) != len(want[i]) {
			t.Fatalf("batch[%d]=%v want %v", i, gotIDs, want[i])
		}
		for j := range want[i] {
			if gotIDs[j] != want[i][j] {
				t.Fatalf("batch[%d]=%v want %v (batch0=%v batch1=%v batch2=%v)",
					i, gotIDs, want[i], ids(got[0]), ids(got[1]), ids(got[2]))
			}
		}
	}
}
