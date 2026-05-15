package masker_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/masker"
)

// TestMasker_CreditCard verifies that credit-card-like numbers are masked
// while surrounding text is preserved.
func TestMasker_CreditCard(t *testing.T) {
	m, err := masker.New(`card=(?P<number>\d{4}-\d{4}-\d{4}-\d{4})`, "[CC-REDACTED]")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	cases := []struct {
		input, want string
	}{
		{"card=1234-5678-9012-3456 ok", "card=[CC-REDACTED] ok"},
		{"no card here", "no card here"},
		{"card=1111-2222-3333-4444 card=5555-6666-7777-8888", "card=[CC-REDACTED] card=[CC-REDACTED]"},
	}

	for _, tc := range cases {
		got := m.Apply(tc.input)
		if got != tc.want {
			t.Errorf("Apply(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

// TestMasker_BearerToken verifies Authorization header values are masked.
func TestMasker_BearerToken(t *testing.T) {
	m, err := masker.New(`Authorization:\s*Bearer\s+(?P<tok>\S+)`, "<TOKEN>")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	line := `GET /api Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.payload.sig 200`
	got := m.Apply(line)

	if strings.Contains(got, "eyJhbGciOiJIUzI1NiJ9") {
		t.Errorf("token not masked; got: %s", got)
	}
	if !strings.Contains(got, "<TOKEN>") {
		t.Errorf("expected <TOKEN> placeholder; got: %s", got)
	}
}
