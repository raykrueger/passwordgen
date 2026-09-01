# passwordgen

Generate memorable passphrases by joining random dictionary words.

Randomness comes from `crypto/rand`, so with the default settings
(4 words from the ~4,000-word pool of four-letter words) each passphrase
carries roughly 48 bits of entropy — comparable to a random 8-character
password, but far easier to remember.

## Build

```sh
go build -o passwordgen .
```

## Usage

```sh
passwordgen                       # e.g. Gray-Bolt-Lake-Vine
passwordgen -words 5 -length 5    # five five-letter words
passwordgen -sep _ -lower         # gray_bolt_lake_vine
```

### Flags

| Flag      | Default                  | Description                        |
|-----------|--------------------------|------------------------------------|
| `-words`  | `4`                      | number of words in the passphrase  |
| `-length` | `4`                      | number of letters in each word     |
| `-sep`    | `-`                      | separator between words            |
| `-dict`   | `/usr/share/dict/words`  | path to the word list              |
| `-lower`  | `false`                  | do not capitalize words            |

## Note

Strength depends on choosing words uniformly at random from the full pool.
This tool does that for you — picking words yourself defeats the purpose.
