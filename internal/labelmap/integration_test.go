package labelmap_test

import (
	"testing"

	"github.com/user/logslice/internal/labelmap"
)

func TestLabelMap_ApacheLog(t *testing.T) {
	pattern := `(?P<host>\S+)\s+\S+\s+\S+\s+\[(?P<time>[^\]]+)\]\s+"(?P<method>\w+)\s+(?P<path>\S+)[^"]*"\s+(?P<status>\d+)`
	m, err := labelmap.New(pattern, map[string]string{"source": "apache"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200`
	labels, ok := m.Map(line)
	if !ok {
		t.Fatal("expected match")
	}

	cases := map[string]string{
		"host":   "127.0.0.1",
		"method": "GET",
		"path":   "/apache_pb.gif",
		"status": "200",
		"source": "apache",
	}
	for k, want := range cases {
		if got := labels[k]; got != want {
			t.Errorf("labels[%q]: got %q, want %q", k, got, want)
		}
	}
}

func TestLabelMap_NoMatchReturnsNil(t *testing.T) {
	m, err := labelmap.New(`(?P<level>ERROR)`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := m.Map("INFO this is fine")
	if ok {
		t.Fatal("expected no match")
	}
}
