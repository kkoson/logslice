package grok

import (
	"testing"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New("", nil)
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_UnknownPattern(t *testing.T) {
	_, err := New("%{NOPE}", nil)
	if err == nil {
		t.Fatal("expected error for unknown pattern token")
	}
}

func TestNew_InvalidRegexp(t *testing.T) {
	extra := map[string]string{"BAD": `(?P<BAD>[`}
	_, err := New("%{BAD}", extra)
	if err == nil {
		t.Fatal("expected error for invalid regexp")
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := New("%{IP} %{LOGLEVEL}", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_NoMatch(t *testing.T) {
	p, _ := New("%{IP}", nil)
	if got := p.Parse("not an ip"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestParse_IPMatch(t *testing.T) {
	p, err := New("%{IP}", nil)
	if err != nil {
		t.Fatal(err)
	}
	fields := p.Parse("request from 192.168.1.1 ok")
	if fields == nil {
		t.Fatal("expected match")
	}
	if got := fields["IP"]; got != "192.168.1.1" {
		t.Fatalf("IP: want 192.168.1.1, got %q", got)
	}
}

func TestParse_MultipleFields(t *testing.T) {
	p, err := New(`%{TIMESTAMP} %{LOGLEVEL} %{GREEDYDATA}`, nil)
	if err != nil {
		t.Fatal(err)
	}
	line := "2024-01-15T12:00:00Z ERROR disk full"
	fields := p.Parse(line)
	if fields == nil {
		t.Fatal("expected match")
	}
	if fields["LOGLEVEL"] != "ERROR" {
		t.Fatalf("LOGLEVEL: want ERROR, got %q", fields["LOGLEVEL"])
	}
	if fields["GREEDYDATA"] != "disk full" {
		t.Fatalf("GREEDYDATA: want 'disk full', got %q", fields["GREEDYDATA"])
	}
}

func TestParse_CustomPattern(t *testing.T) {
	extra := map[string]string{"REQID": `(?P<REQID>req-[0-9a-f]+)`}
	p, err := New("%{REQID}", extra)
	if err != nil {
		t.Fatal(err)
	}
	fields := p.Parse("trace req-deadbeef done")
	if fields == nil {
		t.Fatal("expected match")
	}
	if fields["REQID"] != "req-deadbeef" {
		t.Fatalf("REQID: want req-deadbeef, got %q", fields["REQID"])
	}
}
