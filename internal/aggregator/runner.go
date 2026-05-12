package aggregator

import (
	"context"
	"fmt"
	"io"
)

// Source is the minimal interface consumed by Run — satisfied by reader.Source.
type Source interface {
	Lines() <-chan string
	Err() error
	Close() error
}

// RunConfig holds options for a streaming aggregation run.
type RunConfig struct {
	Source  Source
	Agg     *Aggregator
	Out     io.Writer
	Context context.Context
}

// Run reads lines from cfg.Source, feeds them into cfg.Agg, then writes the
// sorted results to cfg.Out in a simple "KEY\tCOUNT" text format.
// It respects cfg.Context cancellation.
func Run(cfg RunConfig) error {
	ctx := cfg.Context
	if ctx == nil {
		ctx = context.Background()
	}
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case line, ok := <-cfg.Source.Lines():
			if !ok {
				break loop
			}
			cfg.Agg.Add(line)
		}
	}
	if err := cfg.Source.Err(); err != nil {
		return fmt.Errorf("aggregator run: source error: %w", err)
	}
	for _, e := range cfg.Agg.Results() {
		if _, err := fmt.Fprintf(cfg.Out, "%s\t%d\n", e.Key, e.Count); err != nil {
			return fmt.Errorf("aggregator run: write error: %w", err)
		}
	}
	return nil
}
