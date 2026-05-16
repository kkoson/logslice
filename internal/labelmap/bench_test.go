package labelmap

import "testing"

var sink map[string]string

func BenchmarkMap_Match(b *testing.B) {
	m, err := New(`(?P<level>\w+)\s+(?P<host>\S+)\s+(?P<msg>.+)`, nil)
	if err != nil {
		b.Fatal(err)
	}
	line := "ERROR web-01 disk quota exceeded"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink, _ = m.Map(line)
	}
}

func BenchmarkMap_NoMatch(b *testing.B) {
	m, err := New(`(?P<level>ERROR)\s+(?P<msg>.+)`, nil)
	if err != nil {
		b.Fatal(err)
	}
	line := "DEBUG nothing to see here"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink, _ = m.Map(line)
	}
}

func BenchmarkMap_WithOverrides(b *testing.B) {
	overrides := map[string]string{"service": "logslice", "env": "prod", "region": "us-east-1"}
	m, err := New(`(?P<level>\w+)\s+(?P<msg>.+)`, overrides)
	if err != nil {
		b.Fatal(err)
	}
	line := "WARN high memory usage detected"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink, _ = m.Map(line)
	}
}
