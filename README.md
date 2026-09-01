# passwordgen

Generate memorable passphrases by joining random dictionary words.

Randomness comes from `crypto/rand`. A word list is compiled into the
binary (via `//go:embed`), so it runs anywhere with no external files.
By default it draws 4 words of 4–6 letters each. Widening the length range grows
the word pool, which raises the entropy per word — a passphrase like this
is comparable to (or stronger than) a random 8-character password, but far
easier to remember.

## Build

```sh
go build -o passwordgen .
```

## Usage

```sh
passwordgen                       # e.g. Level-Smile-Chair-Slaw
passwordgen -words 5 -min 3 -max 8  # five words, 3-8 letters each
passwordgen -sep _ -lower         # level_smile_chair_slaw
passwordgen -digits 2             # Level-Smile-Chair-Slaw-17
```

### Flags

| Flag     | Default                  | Description                           |
|----------|--------------------------|---------------------------------------|
| `-words` | `4`                      | number of words in the passphrase     |
| `-min`   | `4`                      | minimum number of letters in each word|
| `-max`   | `6`                      | maximum number of letters in each word|
| `-sep`   | `-`                      | separator between words               |
| `-dict`  | *(built-in list)*        | path to an external word list         |
| `-lower` | `false`                  | do not capitalize words               |
| `-digits`| `0`                      | append this many random digits        |

## Word list

The built-in list (`words/words.txt`) is the [EFF "large" Diceware
wordlist](https://www.eff.org/dice) — 7,776 common, recognizable English
words of 3–9 letters, chosen specifically for passphrases (no obscure or
ambiguous entries). The dice indices are stripped; one word per line.

Because it only spans 3–9 letters, `-min`/`-max` outside that range yield
nothing unless you supply a larger list with `-dict`. To regenerate:

```sh
curl -fsS https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt \
  | awk '{print $2}' | grep -E '^[a-z]+$' | sort -u > words/words.txt
```

## Digits

Some password policies insist on a number. `-digits N` appends N
cryptographically-random digits (0–9) after the final separator, e.g.
`-digits 2` → `Level-Smile-Chair-Slaw-17`. Each digit adds ~3.3 bits of
entropy — far less per character than a word, so prefer adding a word
(`-words 5`) for real strength and use digits only to satisfy the policy.

## Note

Strength depends on choosing words uniformly at random from the full pool.
This tool does that for you — picking words yourself defeats the purpose.
