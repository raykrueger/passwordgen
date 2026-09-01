package main

import _ "embed"

// embeddedWords is the EFF "large" Diceware wordlist (7,776 common,
// recognizable English words of 3-9 letters), with the dice indices stripped.
// It is the default source so the binary works without /usr/share/dict/words
// and only produces familiar words.
//
// Source: https://www.eff.org/dice (CC BY 3.0 US).
//
//go:embed words/words.txt
var embeddedWords string
