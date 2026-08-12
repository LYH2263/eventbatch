package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"example.com/eventbatch/internal/batcher"
	"example.com/eventbatch/internal/ingest"
	"example.com/eventbatch/internal/pipeline"
	"example.com/eventbatch/internal/sink"
)

func main() {
	input := flag.String("input", "", "JSONL events file")
	batchSize := flag.Int("batch", 100, "max events per batch")
	outDir := flag.String("out", "", "optional directory to write batch JSON files")
	flag.Parse()

	if *input == "" {
		fmt.Fprintln(os.Stderr, "usage: batcher -input events.jsonl [-batch N] [-out dir]")
		os.Exit(2)
	}

	events, err := ingest.FromFile(*input)
	if err != nil {
		log.Fatal(err)
	}

	var writer pipeline.Writer = &sink.Memory{}
	if *outDir != "" {
		writer = sink.FileWriter{Dir: *outDir}
	}

	n, err := (pipeline.Runner{
		Cfg:    batcher.Config{MaxEvents: *batchSize},
		Writer: writer,
	}).Run(events)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %d batches from %d events\n", n, len(events))
}
