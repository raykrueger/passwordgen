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
passwordgen                       # e.g. Prosy-Locket-Lube-Cup
passwordgen -words 5 -min 3 -max 8  # five words, 3-8 letters each
passwordgen -sep _ -lower         # prosy_locket_lube_cup
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

## Note

Strength depends on choosing words uniformly at random from the full pool.
This tool does that for you — picking words yourself defeats the purpose.
