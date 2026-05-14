// Package fieldextract turns unstructured log lines into structured key/value
// maps using Go regular expressions with named capture groups.
//
// # Basic usage
//
//	e, err := fieldextract.New(`(?P<level>\w+)\s+(?P<msg>.*)`)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fields, ok := e.Extract("ERROR disk full")
//	if ok {
//		fmt.Println(fields["level"]) // ERROR
//		fmt.Println(fields["msg"])   // disk full
//	}
//
// # Integration with output formatter
//
// The map returned by Extract can be serialised directly by the JSON formatter
// in internal/output, enabling structured log pipelines without a schema
// definition step.
//
// # Error handling
//
// New returns an error if the pattern is syntactically invalid or contains no
// named capture groups. A pattern with only unnamed groups is rejected because
// the resulting map would have no meaningful keys.
package fieldextract
