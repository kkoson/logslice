package rotate

import (
	"testing"
)

var sink int

func BenchmarkWrite_NoRotation(b *testing.B) {
	w, err := New(Config{
		Dir:      b.TempDir(),
		Prefix:   "bench-",
		MaxBytes: 1 << 30, // 1 GiB — effectively never rotate
	})
	if err != nil {
		b.Fatalf("New: %v", err)
	}
	defer w.Close()

	line := []byte("2024-01-01T00:00:00Z level=info msg=\"benchmark log line\" key=value\n")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ := w.Write(line)
		sink += n
	}
}

func BenchmarkWrite_FrequentRotation(b *testing.B) {
	w, err := New(Config{
		Dir:      b.TempDir(),
		Prefix:   "bench-rot-",
		MaxBytes: 128,
	})
	if err != nil {
		b.Fatalf("New: %v", err)
	}
	defer w.Close()

	line := []byte("2024-01-01T00:00:00Z level=warn msg=\"rotation benchmark\"\n")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ := w.Write(line)
		sink += n
	}
}
