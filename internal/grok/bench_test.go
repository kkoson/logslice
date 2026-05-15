package grok

import "testing"

var sink map[string]string

func BenchmarkParse_Match(b *testing.B) {
	p, err := New(`%{TIMESTAMP} %{LOGLEVEL} %{GREEDYDATA}`, nil)
	if err != nil {
		b.Fatal(err)
	}
	line := "2024-06-01T10:00:00Z INFO server started successfully"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = p.Parse(line)
	}
}

func BenchmarkParse_NoMatch(b *testing.B) {
	p, err := New(`%{TIMESTAMP} %{LOGLEVEL} %{GREEDYDATA}`, nil)
	if err != nil {
		b.Fatal(err)
	}
	line := "this line will never match the pattern"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = p.Parse(line)
	}
}

func BenchmarkNew_Compile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := New(`%{IP} %{LOGLEVEL} %{GREEDYDATA}`, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}
