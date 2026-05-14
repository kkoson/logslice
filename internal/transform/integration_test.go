package transform_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/transform"
)

// TestTransformer_RedactIP verifies that IP addresses are masked end-to-end.
func TestTransformer_RedactIP(t *testing.T) {
	tr, err := transform.New(
		`\b(?:\d{1,3}\.){3}\d{1,3}\b`,
		"<ip>",
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	cases := []struct {
		input, want string
	}{
		{"connected from 192.168.1.1", "connected from <ip>"},
		{"no address here", "no address here"},
		{"src=10.0.0.1 dst=10.0.0.2", "src=<ip> dst=<ip>"},
	}

	for _, tc := range cases {
		got := tr.Apply(tc.input)
		if got != tc.want {
			t.Errorf("Apply(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

// TestChain_MultiStepPipeline exercises a realistic multi-step redaction chain.
func TestChain_MultiStepPipeline(t *testing.T) {
	transformers := []*transform.Transformer{
		transform.MustNew(`\b(?:\d{1,3}\.){3}\d{1,3}\b`, "<ip>"),
		transform.MustNew(`password=\S+`, "password=<redacted>"),
		transform.MustNew(`(?i)error`, "ERR"),
	}

	input := "ERROR: login from 10.0.0.1 with password=s3cr3t"
	want := "ERR: login from <ip> with password=<redacted>"

	got := transform.Chain(transformers, input)
	if got != want {
		t.Fatalf("Chain result = %q; want %q", got, want)
	}
}
