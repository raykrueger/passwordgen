// Command passwordgen generates memorable passphrases by joining random
// dictionary words. Randomness comes from crypto/rand so the output is
// suitable for use as a password.
package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
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
		sep   = flag.String("sep", "-", "separator between words")
		dict  = flag.String("dict", "/usr/share/dict/words", "path to word list")
		lower = flag.Bool("lower", false, "do not capitalize words")
	)
	flag.Parse()

	if *min < 1 || *max < *min {
		fmt.Fprintln(os.Stderr, "passwordgen: require 1 <= min <= max")
		os.Exit(1)
	}

	words, err := loadWords(*dict, *min, *max)
	if err != nil {
		fmt.Fprintln(os.Stderr, "passwordgen:", err)
		os.Exit(1)
	}
	if len(words) == 0 {
		fmt.Fprintf(os.Stderr, "passwordgen: no %d-%d letter words found in %s\n", *min, *max, *dict)
		os.Exit(1)
	}

	phrase, err := passphrase(words, *count, *sep, !*lower)
	if err != nil {
		fmt.Fprintln(os.Stderr, "passwordgen:", err)
		os.Exit(1)
	}
	fmt.Println(phrase)
}

// loadWords returns every word in the file at path made of between min and max
// lowercase ASCII letters (inclusive).
func loadWords(path string, min, max int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	re := regexp.MustCompile(fmt.Sprintf("^[a-z]{%d,%d}$", min, max))
	var words []string
	s := bufio.NewScanner(f)
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
