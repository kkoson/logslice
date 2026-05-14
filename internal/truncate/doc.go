// Package truncate provides a Truncator that shortens log lines exceeding a
// configurable maximum byte length.
//
// Usage:
//
//	tr, err := truncate.New(200, "...")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, line := range lines {
//		fmt.Println(tr.Apply(line))
//	}
//
// The Truncator is safe for concurrent use; Apply has no side-effects on
// shared state.
package truncate
