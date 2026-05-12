// Package aggregator implements log line aggregation by regex capture groups.
//
// # Overview
//
// The aggregator package lets callers count how often each distinct value of a
// named regex capture group appears across a stream of log lines.  Results are
// returned sorted by frequency (descending) with alphabetical tie-breaking.
//
// # Basic usage
//
//	agg, err := aggregator.New(`level=(?P<level>\w+)`, "level")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, line := range lines {
//		agg.Add(line)
//	}
//	for _, e := range agg.Results() {
//		fmt.Printf("%s: %d\n", e.Key, e.Count)
//	}
//
// # Streaming
//
// For streaming use cases, Run() accepts any Source (e.g. a reader.FileSource)
// and writes formatted results to an io.Writer once the source is exhausted or
// the context is cancelled.
package aggregator
