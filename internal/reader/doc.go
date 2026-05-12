// Package reader provides log input sources for logslice.
//
// It abstracts reading log lines from different origins — files on disk or
// standard input — behind a common Source interface. Consumers receive lines
// over a string channel and errors over a separate error channel, enabling
// non-blocking, streaming processing of potentially large log files.
//
// # Usage
//
// File-based source:
//
//	src, err := reader.NewFileSource("/var/log/app.log")
//	if err != nil { ... }
//	defer src.Close()
//
//	lines, errs := src.Lines()
//	for line := range lines {
//	    // process line
//	}
//	if err := <-errs; err != nil { ... }
//
// Stdin source:
//
//	src := reader.NewStdinSource()
//	lines, errs := src.Lines()
//	for line := range lines {
//	    // process line
//	}
package reader
