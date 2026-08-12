package ingest_test

import (
	"strings"
	"testing"
	"time"

	"example.com/eventbatch/internal/ingest"
)

func TestReadAllJSONL(t *testing.T) {
	raw := strings.NewReader(`{"id":"1","type":"login","payload":{"u":"a"},"ts":"2024-01-01T00:00:00Z"}
{"id":"2","type":"logout","payload":{"u":"a"},"ts":"2024-01-01T01:00:00Z"}
`)
	got, err := ingest.NewReader(raw).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("len=%d", len(got))
	}
	if got[0].ID != "1" || got[1].Type != "logout" {
		t.Fatalf("got=%+v", got)
	}
	if !got[0].Timestamp.Equal(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("ts=%v", got[0].Timestamp)
	}
}
