# passwordgen

Generate memorable passphrases by joining random dictionary words.

Randomness comes from `crypto/rand`. By default it draws 4 words of 4–6
letters each from the system dictionary. Widening the length range grows
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
| `-dict`  | `/usr/share/dict/words`  | path to the word list                 |
| `-lower` | `false`                  | do not capitalize words               |

## Note

Strength depends on choosing words uniformly at random from the full pool.
This tool does that for you — picking words yourself defeats the purpose.
