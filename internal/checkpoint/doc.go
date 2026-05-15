// Package checkpoint provides durable read-position tracking for logslice.
//
// When processing large or continuously-growing log files it is useful to
// remember the byte offset of the last successfully processed line so that
// logslice can resume from that position after a restart instead of
// re-reading the entire file.
//
// # Basic usage
//
//	cp, err := checkpoint.New("/var/run/logslice/app.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Restore the previous offset before opening the source.
//	state := cp.Get()
//
//	// After processing each line update the checkpoint.
//	if err := cp.Save(logPath, newOffset); err != nil {
//		log.Print(err)
//	}
//
// Saves are atomic on POSIX systems: the state is written to a temporary
// file and then renamed into place, so a crash mid-write cannot produce a
// truncated or corrupt checkpoint file.
package checkpoint
