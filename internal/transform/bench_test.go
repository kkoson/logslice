package transform

import "testing"

var sink string

func BenchmarkApply_NoMatch(b *testing.B) {
	tr := MustNew(`\d{4}-\d{2}-\d{2}`, "<date>")
	line := "INFO no date in this log line at all"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = tr.Apply(line)
	}
}

func BenchmarkApply_Match(b *testing.B) {
	tr := MustNew(`\d{4}-\d{2}-\d{2}`, "<date>")
	line := "2024-06-15 INFO server started"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = tr.Apply(line)
	}
}

func BenchmarkChain_ThreeSteps(b *testing.B) {
	trs := []*Transformer{
		MustNew(`\d{4}-\d{2}-\d{2}`, "<date>"),
		MustNew(`\b(?:\d{1,3}\.){3}\d{1,3}\b`, "<ip>"),
		MustNew(`password=\S+`, "password=<redacted>"),
	}
	line := "2024-06-15 INFO login from 192.168.0.1 password=hunter2"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = Chain(trs, line)
	}
}
