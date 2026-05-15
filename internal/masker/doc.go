// Package masker provides regex-based redaction of sensitive fields in log lines.
//
// A Masker is constructed with a regular expression that must contain one or
// more named capture groups.  When [Masker.Apply] is called, every occurrence
// of those named groups is replaced with a configurable mask string (default
// "***"), while the rest of the line — including any surrounding literal text
// captured by the overall match — is preserved unchanged.
//
// # Example
//
//	m, err := masker.New(`password=(?P<pw>\S+)`, "[REDACTED]")
//	if err != nil {
//		log.Fatal(err)
//	}
//	out := m.Apply("login password=hunter2 ok")
//	// out == "login password=[REDACTED] ok"
//
// Masker is safe for concurrent use after construction.
package masker
