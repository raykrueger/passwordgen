// Command passwordgen generates memorable passphrases by joining random
// dictionary words. Randomness comes from crypto/rand so the output is
// suitable for use as a password.
package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
	"regexp"
	"strings"
)

func main() {
	var (
		count = flag.Int("words", 4, "number of words in the passphrase")
		min   = flag.Int("min", 4, "minimum number of letters in each word")
		max   = flag.Int("max", 6, "maximum number of letters in each word")
		sep    = flag.String("sep", "-", "separator between words")
		dict   = flag.String("dict", "", "path to word list (default: built-in list)")
		lower  = flag.Bool("lower", false, "do not capitalize words")
		digits = flag.Int("digits", 0, "append this many random digits to the end")
	)
	flag.Parse()

	if *min < 1 || *max < *min {
		fmt.Fprintln(os.Stderr, "passwordgen: require 1 <= min <= max")
		os.Exit(1)
	}
	if *digits < 0 {
		fmt.Fprintln(os.Stderr, "passwordgen: digits must not be negative")
		os.Exit(1)
	}

	src := strings.NewReader(embeddedWords)
	var r io.Reader = src
	if *dict != "" {
		f, err := os.Open(*dict)
		if err != nil {
			fmt.Fprintln(os.Stderr, "passwordgen:", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	}

	words, err := loadWords(r, *min, *max)
	if err != nil {
		fmt.Fprintln(os.Stderr, "passwordgen:", err)
		os.Exit(1)
	}
	if len(words) == 0 {
		fmt.Fprintf(os.Stderr, "passwordgen: no %d-%d letter words found\n", *min, *max)
		os.Exit(1)
	}

	phrase, err := passphrase(words, *count, *sep, !*lower)
	if err != nil {
		fmt.Fprintln(os.Stderr, "passwordgen:", err)
		os.Exit(1)
	}

	if *digits > 0 {
		tail, err := randomDigits(*digits)
		if err != nil {
			fmt.Fprintln(os.Stderr, "passwordgen:", err)
			os.Exit(1)
		}
		phrase += *sep + tail
	}

	fmt.Println(phrase)
}

// loadWords returns every word from r made of between min and max lowercase
// ASCII letters (inclusive).
func loadWords(r io.Reader, min, max int) ([]string, error) {
	re := regexp.MustCompile(fmt.Sprintf("^[a-z]{%d,%d}$", min, max))
	var words []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		if w := s.Text(); re.MatchString(w) {
			words = append(words, w)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return words, nil
}

// passphrase joins count words chosen uniformly at random from words using a
// cryptographically secure source.
func passphrase(words []string, count int, sep string, capitalize bool) (string, error) {
	parts := make([]string, count)
	for i := range parts {
		w, err := choose(words)
		if err != nil {
			return "", err
		}
		if capitalize {
			w = strings.ToUpper(w[:1]) + w[1:]
		}
		parts[i] = w
	}
	return strings.Join(parts, sep), nil
}

// choose returns a random element of words using crypto/rand.
func choose(words []string) (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
	if err != nil {
		return "", err
	}
	return words[n.Int64()], nil
}

// randomDigits returns a string of n decimal digits (0-9) chosen with
// crypto/rand. Leading zeros are preserved.
func randomDigits(n int) (string, error) {
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < n; i++ {
		d, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		b.WriteByte('0' + byte(d.Int64()))
	}
	return b.String(), nil
}
