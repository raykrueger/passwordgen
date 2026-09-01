package main

import (
	"strings"
	"testing"
)

func TestLoadWords(t *testing.T) {
	const list = `cat
fish
gold
zebra
hippopotamus
Upper
d1git
he-llo
mind
`
	tests := []struct {
		name     string
		min, max int
		want     []string
	}{
		{"exactly four", 4, 4, []string{"fish", "gold", "mind"}},
		{"range three to five", 3, 5, []string{"cat", "fish", "gold", "zebra", "mind"}},
		{"only three", 3, 3, []string{"cat"}},
		{"none in range", 20, 30, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadWords(strings.NewReader(list), tt.min, tt.max)
			if err != nil {
				t.Fatalf("loadWords returned error: %v", err)
			}
			if !equal(got, tt.want) {
				t.Errorf("loadWords(%d, %d) = %v, want %v", tt.min, tt.max, got, tt.want)
			}
		})
	}
}

func TestLoadWordsRejectsNonLowercaseAlpha(t *testing.T) {
	const list = "Upper\nd1git\nhe-llo\nspace bar\ngood\n"
	got, err := loadWords(strings.NewReader(list), 1, 20)
	if err != nil {
		t.Fatalf("loadWords returned error: %v", err)
	}
	if !equal(got, []string{"good"}) {
		t.Errorf("loadWords = %v, want [good]", got)
	}
}

func TestPassphraseWordCount(t *testing.T) {
	words := []string{"alpha", "bravo", "charlie", "delta"}
	for _, count := range []int{1, 4, 10} {
		p, err := passphrase(words, count, "-", false)
		if err != nil {
			t.Fatalf("passphrase(count=%d) error: %v", count, err)
		}
		parts := strings.Split(p, "-")
		if len(parts) != count {
			t.Errorf("count=%d: got %d parts (%q), want %d", count, len(parts), p, count)
		}
	}
}

func TestPassphraseSeparator(t *testing.T) {
	words := []string{"same"}
	p, err := passphrase(words, 3, "_", false)
	if err != nil {
		t.Fatalf("passphrase error: %v", err)
	}
	if p != "same_same_same" {
		t.Errorf("got %q, want same_same_same", p)
	}
}

func TestPassphraseCapitalize(t *testing.T) {
	words := []string{"gold"}

	up, err := passphrase(words, 2, "-", true)
	if err != nil {
		t.Fatalf("passphrase error: %v", err)
	}
	if up != "Gold-Gold" {
		t.Errorf("capitalize=true: got %q, want Gold-Gold", up)
	}

	low, err := passphrase(words, 2, "-", false)
	if err != nil {
		t.Fatalf("passphrase error: %v", err)
	}
	if low != "gold-gold" {
		t.Errorf("capitalize=false: got %q, want gold-gold", low)
	}
}

func TestChooseWithinPool(t *testing.T) {
	words := []string{"one", "two", "six"}
	pool := map[string]bool{"one": true, "two": true, "six": true}
	for i := 0; i < 100; i++ {
		got, err := choose(words)
		if err != nil {
			t.Fatalf("choose error: %v", err)
		}
		if !pool[got] {
			t.Fatalf("choose returned %q, not in pool", got)
		}
	}
}

// TestChooseDistribution is a light sanity check that choose eventually hits
// every element of a small pool, i.e. it isn't stuck on one index.
func TestChooseDistribution(t *testing.T) {
	words := []string{"a", "b", "c", "d"}
	seen := map[string]bool{}
	for i := 0; i < 500; i++ {
		got, err := choose(words)
		if err != nil {
			t.Fatalf("choose error: %v", err)
		}
		seen[got] = true
	}
	if len(seen) != len(words) {
		t.Errorf("choose only produced %d of %d words: %v", len(seen), len(words), seen)
	}
}

func TestRandomDigits(t *testing.T) {
	for _, n := range []int{1, 2, 4, 8} {
		got, err := randomDigits(n)
		if err != nil {
			t.Fatalf("randomDigits(%d) error: %v", n, err)
		}
		if len(got) != n {
			t.Errorf("randomDigits(%d) = %q, want length %d", n, got, len(got))
		}
		for _, c := range got {
			if c < '0' || c > '9' {
				t.Errorf("randomDigits(%d) = %q, contains non-digit %q", n, got, c)
			}
		}
	}
}

func TestRandomDigitsZero(t *testing.T) {
	got, err := randomDigits(0)
	if err != nil {
		t.Fatalf("randomDigits(0) error: %v", err)
	}
	if got != "" {
		t.Errorf("randomDigits(0) = %q, want empty string", got)
	}
}

// TestRandomDigitsSpread checks that every digit 0-9 shows up over many draws,
// i.e. the generator isn't stuck on a subset.
func TestRandomDigitsSpread(t *testing.T) {
	seen := map[rune]bool{}
	for i := 0; i < 500; i++ {
		s, err := randomDigits(4)
		if err != nil {
			t.Fatalf("randomDigits error: %v", err)
		}
		for _, c := range s {
			seen[c] = true
		}
	}
	if len(seen) != 10 {
		t.Errorf("only produced %d distinct digits: %v", len(seen), seen)
	}
}

// TestEmbeddedWords guards against an empty or malformed built-in list.
func TestEmbeddedWords(t *testing.T) {
	words, err := loadWords(strings.NewReader(embeddedWords), 4, 6)
	if err != nil {
		t.Fatalf("loadWords on embedded list: %v", err)
	}
	if len(words) < 2000 {
		t.Errorf("embedded list has only %d words in 4-6 range, expected many more", len(words))
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
