// Package multiline provides a line folder that collapses multi-line log
// events (e.g. Java stack traces or wrapped JSON blobs) into a single
// logical string before downstream processing.
//
// Usage:
//
//	f, err := multiline.New(`^\d{4}-\d{2}-\d{2}`, "\n")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, raw := range lines {
//		if event, ok := f.Add(raw); ok {
//			process(event)
//		}
//	}
//	if event, ok := f.Flush(); ok {
//		process(event)
//	}
//
// The start pattern marks the first line of every new event. Any line
// that does not match is treated as a continuation and joined to the
// current event using the configured join string.
package multiline
