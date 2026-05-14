package linecount_test

import (
	"sync"
	"testing"

	"github.com/yourorg/logslice/internal/linecount"
)

func TestNew_ZeroValues(t *testing.T) {
	c := linecount.New()
	if c.Seen() != 0 || c.Emitted() != 0 || c.Dropped() != 0 {
		t.Fatalf("expected all zeros, got seen=%d emitted=%d dropped=%d",
			c.Seen(), c.Emitted(), c.Dropped())
	}
}

func TestIncSeen(t *testing.T) {
	c := linecount.New()
	c.IncSeen()
	c.IncSeen()
	if c.Seen() != 2 {
		t.Fatalf("expected 2, got %d", c.Seen())
	}
}

func TestIncEmitted(t *testing.T) {
	c := linecount.New()
	c.IncSeen()
	c.IncSeen()
	c.IncSeen()
	c.IncEmitted()
	c.IncEmitted()
	if c.Emitted() != 2 {
		t.Fatalf("expected 2, got %d", c.Emitted())
	}
	if c.Dropped() != 1 {
		t.Fatalf("expected 1 dropped, got %d", c.Dropped())
	}
}

func TestSnapshot(t *testing.T) {
	c := linecount.New()
	for i := 0; i < 10; i++ {
		c.IncSeen()
	}
	for i := 0; i < 7; i++ {
		c.IncEmitted()
	}
	s := c.Snapshot()
	if s.Seen != 10 {
		t.Errorf("Seen: want 10, got %d", s.Seen)
	}
	if s.Emitted != 7 {
		t.Errorf("Emitted: want 7, got %d", s.Emitted)
	}
	if s.Dropped != 3 {
		t.Errorf("Dropped: want 3, got %d", s.Dropped)
	}
}

func TestCounter_ConcurrentSafety(t *testing.T) {
	c := linecount.New()
	const goroutines = 50
	const perGoroutine = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				c.IncSeen()
				c.IncEmitted()
			}
		}()
	}
	wg.Wait()
	expected := int64(goroutines * perGoroutine)
	if c.Seen() != expected {
		t.Errorf("Seen: want %d, got %d", expected, c.Seen())
	}
	if c.Emitted() != expected {
		t.Errorf("Emitted: want %d, got %d", expected, c.Emitted())
	}
	if c.Dropped() != 0 {
		t.Errorf("Dropped: want 0, got %d", c.Dropped())
	}
}
