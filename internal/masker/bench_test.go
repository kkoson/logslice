package masker

import "testing"

var sink string

func BenchmarkApply_NoMatch(b *testing.B) {
	m, _ := New(`password=(?P<pw>\S+)`, "***")
	line := "2024-01-15T10:00:00Z INFO user logged in successfully"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = m.Apply(line)
	}
}

func BenchmarkApply_SingleMatch(b *testing.B) {
	m, _ := New(`password=(?P<pw>\S+)`, "***")
	line := "2024-01-15T10:00:00Z INFO password=s3cr3t user=admin"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = m.Apply(line)
	}
}

func BenchmarkApply_MultipleMatches(b *testing.B) {
	m, _ := New(`(?P<tok>tok_[A-Za-z0-9]{16})`, "[REDACTED]")
	line := "tok_AAAAAAAAAAAAAAAA tok_BBBBBBBBBBBBBBBB tok_CCCCCCCCCCCCCCCC"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = m.Apply(line)
	}
}
