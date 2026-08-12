package pipeline_test

import (
	"testing"

	"example.com/eventbatch/internal/batcher"
	"example.com/eventbatch/internal/event"
	"example.com/eventbatch/internal/pipeline"
	"example.com/eventbatch/internal/sink"
)

func TestRunnerWritesIndependentBatches(t *testing.T) {
	events := []event.Event{
		{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"},
	}
	mem := &sink.Memory{}
	n, err := (pipeline.Runner{
		Cfg:    batcher.Config{MaxEvents: 2},
		Writer: mem,
	}).Run(events)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 || len(mem.Batches) != 2 {
		t.Fatalf("n=%d stored=%d", n, len(mem.Batches))
	}
	if mem.Batches[0].Events[0].ID != "1" || mem.Batches[0].Events[1].ID != "2" {
		t.Fatalf("batch0=%+v", mem.Batches[0].Events)
	}
	if mem.Batches[1].Events[0].ID != "3" || mem.Batches[1].Events[1].ID != "4" {
		t.Fatalf("batch1=%+v", mem.Batches[1].Events)
	}
}
