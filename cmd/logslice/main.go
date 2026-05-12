// Command logslice is a fast log filtering and aggregation tool
// with regex pipelines and structured output formats.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/logslice/internal/aggregator"
	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/reader"
)

type multiFlag []string

func (m *multiFlag) String() string  { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error { *m = append(*m, v); return nil }

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	var patterns multiFlag
	var invertPatterns multiFlag
	aggPattern := fs.String("agg", "", "aggregation regex with a named group (e.g. (?P<key>\\w+))")
	format := fs.String("format", "plain", "output format: plain, json, csv")
	inputFile := fs.String("file", "", "input file (defaults to stdin)")

	fs.Var(&patterns, "match", "filter pattern (repeatable); lines must match all patterns")
	fs.Var(&invertPatterns, "exclude", "exclude pattern (repeatable); lines matching are dropped")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Build source
	var src reader.Source
	var err error
	if *inputFile != "" {
		src, err = reader.NewFileSource(*inputFile)
		if err != nil {
			return err
		}
	} else {
		src = reader.NewStdinSource()
	}
	defer src.Close()

	// Build filter pipeline
	var filters []*filter.Filter
	for _, p := range patterns {
		f, err := filter.NewFilter(p, false)
		if err != nil {
			return fmt.Errorf("invalid match pattern %q: %w", p, err)
		}
		filters = append(filters, f)
	}
	for _, p := range invertPatterns {
		f, err := filter.NewFilter(p, true)
		if err != nil {
			return fmt.Errorf("invalid exclude pattern %q: %w", p, err)
		}
		filters = append(filters, f)
	}
	pipeline := filter.NewPipeline(filters...)

	// Build formatter
	fmt_, err := output.NewFormatter(*format, os.Stdout)
	if err != nil {
		return err
	}

	// Aggregation mode
	if *aggPattern != "" {
		agg, err := aggregator.New(*aggPattern)
		if err != nil {
			return err
		}
		return aggregator.Run(src, pipeline, agg, fmt_)
	}

	// Plain filter mode
	for src.Scan() {
		line := src.Text()
		if pipeline.Match(line) {
			if err := fmt_.Write(line, nil); err != nil {
				return err
			}
		}
	}
	return src.Err()
}
