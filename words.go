package main

import _ "embed"

// embeddedWords is a list of 3-8 letter lowercase English words, derived from
// the system dictionary at build time. It is the default source so the binary
// works without /usr/share/dict/words.
//
//go:embed words/words.txt
var embeddedWords string
