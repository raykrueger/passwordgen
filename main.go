// Command passwordgen generates memorable passphrases by joining random
// dictionary words. Randomness comes from crypto/rand so the output is
// suitable for use as a password.
package main

import (
	"bufio"
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
	"regexp"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "passwordgen:", err)
		os.Exit(1)
	}
}

// run parses flags, builds a passphrase, and prints it. It returns an error
// rather than calling os.Exit so that deferred cleanup runs.
func run() error {
	var (
		count  = flag.Int("words", 4, "number of words in the passphrase")
		minLen = flag.Int("min", 4, "minimum number of letters in each word")
		maxLen = flag.Int("max", 6, "maximum number of letters in each word")
		sep    = flag.String("sep", "-", "separator between words")
		dict   = flag.String("dict", "", "path to word list (default: built-in list)")
		lower  = flag.Bool("lower", false, "do not capitalize words")
		digits = flag.Int("digits", 0, "append this many random digits to the end")
	)
	flag.Parse()

	if *minLen < 1 || *maxLen < *minLen {
		return errors.New("require 1 <= min <= max")
	}
	if *digits < 0 {
		return errors.New("digits must not be negative")
	}

	var r io.Reader = strings.NewReader(embeddedWords)
	if *dict != "" {
		f, err := os.Open(*dict)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}

	words, err := loadWords(r, *minLen, *maxLen)
	if err != nil {
		return err
	}
	if len(words) == 0 {
		return fmt.Errorf("no %d-%d letter words found", *minLen, *maxLen)
	}

	phrase, err := passphrase(words, *count, *sep, !*lower)
	if err != nil {
		return err
	}

	if *digits > 0 {
		tail, err := randomDigits(*digits)
		if err != nil {
			return err
		}
		phrase += *sep + tail
	}

	fmt.Println(phrase)
	return nil
}

// loadWords returns every word from r made of between minLen and maxLen
// lowercase ASCII letters (inclusive).
func loadWords(r io.Reader, minLen, maxLen int) ([]string, error) {
	re := regexp.MustCompile(fmt.Sprintf("^[a-z]{%d,%d}$", minLen, maxLen))
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
// cryptographically secure source. words must not be empty.
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

// choose returns a random element of words using crypto/rand. words must not be
// empty; crypto/rand.Int panics on a non-positive bound.
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
	for range n {
		d, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		b.WriteByte('0' + byte(d.Int64()))
	}
	return b.String(), nil
}
