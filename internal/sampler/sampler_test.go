package sampler

import (
	"math/rand"
	"testing"
)

func TestNew_InvalidRate(t *testing.T) {
	cases := []float64{0, -0.5, 1.1, 2.0}
	for _, r := range cases {
		_, err := New(r, nil)
		if err == nil {
			t.Errorf("expected error for rate %v, got nil", r)
		}
	}
}

func TestNew_ValidRate(t *testing.T) {
	s, err := New(0.5, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Rate() != 0.5 {
		t.Errorf("expected rate 0.5, got %v", s.Rate())
	}
}

func TestKeep_RateOne_KeepsAll(t *testing.T) {
	s, _ := New(1.0, rand.NewSource(0))
	for i := 0; i < 100; i++ {
		if !s.Keep("any line") {
			t.Fatal("rate=1.0 should keep every line")
		}
	}
}

func TestKeep_RateNearZero_KeepsFew(t *testing.T) {
	// Use a fixed seed for determinism.
	s, _ := New(0.01, rand.NewSource(1))
	kept := 0
	const total = 10000
	for i := 0; i < total; i++ {
		if s.Keep("line") {
			kept++
		}
	}
	// With rate 0.01 we expect ~100 kept; allow wide tolerance.
	if kept > 300 {
		t.Errorf("rate=0.01 kept too many lines: %d/%d", kept, total)
	}
}

func TestKeep_ApproximateRate(t *testing.T) {
	s, _ := New(0.5, rand.NewSource(99))
	kept := 0
	const total = 100000
	for i := 0; i < total; i++ {
		if s.Keep("line") {
			kept++
		}
	}
	ratio := float64(kept) / float64(total)
	if ratio < 0.48 || ratio > 0.52 {
		t.Errorf("expected ratio ~0.5, got %.4f", ratio)
	}
}
