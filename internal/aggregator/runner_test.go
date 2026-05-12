package aggregator_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/aggregator"
)

// fakeSource implements aggregator.Source for testing.
type fakeSource struct {
	lines []string
	ch    chan string
}

func newFakeSource(lines []string) *fakeSource {
	ch := make(chan string, len(lines))
	for _, l := range lines {
		ch <- l
	}
	close(ch)
	return &fakeSource{lines: lines, ch: ch}
}

func (f *fakeSource) Lines() <-chan string { return f.ch }
func (f *fakeSource) Err() error          { return nil }
func (f *fakeSource) Close() error        { return nil }

func TestRun_WritesResults(t *testing.T) {
	src := newFakeSource([]string{
		"level=ERROR details",
		"level=INFO details",
		"level=INFO details",
	})
	agg, _ := aggregator.New(`level=(?P<level>\w+)`, "level")
	var buf bytes.Buffer
	err := aggregator.Run(aggregator.RunConfig{
		Source:  src,
		Agg:     agg,
		Out:     &buf,
		Context: context.Background(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "INFO\t2") {
		t.Errorf("expected INFO\t2 in output, got:\n%s", output)
	}
	if !strings.Contains(output, "ERROR\t1") {
		t.Errorf("expected ERROR\t1 in output, got:\n%s", output)
	}
}

func TestRun_ContextCancelled(t *testing.T) {
	ch := make(chan string) // never sends
	type blockSource struct{ fakeSource }
	src := &fakeSource{ch: ch}
	agg, _ := aggregator.New(`(?P<k>\w+)`, "k")
	var buf bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	err := aggregator.Run(aggregator.RunConfig{
		Source:  src,
		Agg:     agg,
		Out:     &buf,
		Context: ctx,
	})
	if err != nil {
		t.Fatalf("unexpected error on cancelled context: %v", err)
	}
}
