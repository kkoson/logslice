// Package linecount provides a thread-safe counter for tracking pipeline
// throughput statistics.
//
// A Counter records how many log lines were read from a source (Seen),
// how many survived all filters and were written to the output (Emitted),
// and derives the number that were dropped (Seen - Emitted).
//
// Typical usage:
//
//	c := linecount.New()
//	for src.Next() {
//		c.IncSeen()
//		if pipeline.Match(src.Line()) {
//			c.IncEmitted()
//			formatter.Write(src.Line())
//		}
//	}
//	fmt.Fprintf(os.Stderr, "seen=%d emitted=%d dropped=%d\n",
//		c.Seen(), c.Emitted(), c.Dropped())
//
// Snapshot returns an immutable Summary suitable for structured reporting
// without the risk of the values changing mid-read.
package linecount
