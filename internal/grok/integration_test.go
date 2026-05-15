package grok_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/grok"
)

// TestParser_ApacheCommonLog verifies a realistic Apache Common Log line.
func TestParser_ApacheCommonLog(t *testing.T) {
	extra := map[string]string{
		"HTTPREQ": `(?P<HTTPREQ>[A-Z]+)`,
		"URI":     `(?P<URI>/\S*)`,
		"PROTO":   `(?P<PROTO>HTTP/\d\.\d)`,
		"STATUS":  `(?P<STATUS>\d{3})`,
		"BYTES":   `(?P<BYTES>\d+|-)`,
	}

	pattern := `%{IP} - - \[%{TIMESTAMP}\] "%{HTTPREQ} %{URI} %{PROTO}" %{STATUS} %{BYTES}`
	p, err := grok.New(pattern, extra)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	line := `10.0.0.1 - - [2024-03-01T08:00:00Z] "GET /index.html HTTP/1.1" 200 1234`
	fields := p.Parse(line)
	if fields == nil {
		t.Fatal("expected match")
	}

	cases := map[string]string{
		"IP":      "10.0.0.1",
		"HTTPREQ": "GET",
		"URI":     "/index.html",
		"STATUS":  "200",
		"BYTES":   "1234",
	}
	for field, want := range cases {
		if got := fields[field]; got != want {
			t.Errorf("%s: want %q, got %q", field, want, got)
		}
	}
}

// TestParser_NoMatchReturnsNil ensures non-matching lines produce no output.
func TestParser_NoMatchReturnsNil(t *testing.T) {
	p, err := grok.New(`%{IP} %{LOGLEVEL}`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := p.Parse(""); got != nil {
		t.Fatalf("expected nil for empty line, got %v", got)
	}
	if got := p.Parse("just some text"); got != nil {
		t.Fatalf("expected nil for non-matching text, got %v", got)
	}
}
