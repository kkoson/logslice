// Package window provides a sliding time-window counter for rate-aware log
// processing.
//
// A Window divides a fixed duration into a configurable number of buckets.
// Each call to Add records events in the current (newest) bucket. As time
// advances, stale buckets are rotated out so that Total always reflects only
// the events that occurred within the last window.Size duration.
//
// Typical usage:
//
//	w, err := window.New(time.Minute, 12) // 12 × 5-second buckets
//	if err != nil {
//		log.Fatal(err)
//	}
//	w.Add(1)          // record one event now
//	fmt.Println(w.Total()) // events in the last minute
//
// All methods are safe for concurrent use.
package window
