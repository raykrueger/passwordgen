# passwordgen

**Strong passwords you can actually remember.** Random character strings like
`k7$Qm2!x` are secure but impossible to recall; the passwords people can
remember are weak. `passwordgen` splits the difference: it assembles
passphrases from random common words — `Sphinx-Cosmic-Scroll-County` — that are
easy to memorize yet carry as much entropy as a random 8-character password.

Words are drawn from the [EFF Diceware wordlist](https://www.eff.org/dice)
using Go's `crypto/rand`, so every passphrase is genuinely unpredictable. The
wordlist is compiled into the binary, so there are no external files or network
calls at runtime.

If you've seen [xkcd 936 ("correct horse battery staple")](https://xkcd.com/936/),
this is that idea as a command-line tool — with the crucial detail that the
words are chosen by a cryptographic RNG, not by a human.

## Installation

Requires **Go 1.24+**.

```sh
go install github.com/raykrueger/passwordgen@latest
```

Or build from a checkout:

```sh
git clone https://github.com/raykrueger/passwordgen.git
cd passwordgen
go build -o passwordgen .
```

## Quick start

```sh
$ passwordgen
Sphinx-Cosmic-Scroll-County
```

That's a passphrase of four random words, each 4–6 letters, capitalized and
dash-separated — roughly 46 bits of entropy. Run it again for a new one.

## Usage

```sh
passwordgen                          # Level-Smile-Chair-Slaw
passwordgen -words 5                 # five words instead of four
passwordgen -min 3 -max 9            # full wordlist range, ~13 bits/word
passwordgen -sep _ -lower            # level_smile_chair_slaw
passwordgen -digits 2                # Level-Smile-Chair-Slaw-17
```

### Flags

| Flag      | Default           | Description                             |
|-----------|-------------------|-----------------------------------------|
| `-words`  | `4`               | number of words in the passphrase       |
| `-min`    | `4`               | minimum letters per word                |
| `-max`    | `6`               | maximum letters per word                |
| `-sep`    | `-`               | separator between words                 |
| `-lower`  | `false`           | do not capitalize words                 |
| `-digits` | `0`               | append this many random digits          |
| `-dict`   | *(built-in list)* | path to an external word list           |

## Strength

Each word is chosen uniformly at random from the pool, so entropy is
`words × log₂(pool size)`. The defaults (`-min 4 -max 6`) use a ~2,766-word
subset ≈ 11.4 bits/word; widening to `-min 3 -max 9` uses the full 7,776-word
list = exactly 12.9 bits/word (canonical Diceware).

| Command                        | Entropy    |
|--------------------------------|------------|
| `passwordgen`                  | ~46 bits   |
| `passwordgen -words 5`         | ~57 bits   |
| `passwordgen -min 3 -max 9`    | ~52 bits   |
| `passwordgen -words 6 -min 3 -max 9` | ~77 bits |

The strength depends entirely on the random source — which is why you should
let the tool pick. Choosing "memorable" words yourself collapses the pool and
the entropy with it.

### Satisfying "must contain a number" policies

`-digits N` appends N cryptographically-random digits after the final
separator (`-digits 2` → `...-Slaw-17`). Each digit adds only ~3.3 bits, far
less per character than a word — so use the smallest N a policy requires, and
reach for `-words 5` when you actually want more strength.

## NIST alignment

Passphrases like these are exactly what modern federal guidance recommends.
[NIST SP 800-63B](https://pages.nist.gov/800-63-3/sp800-63b.html) (Digital
Identity Guidelines) reversed decades of bad password advice, and `passwordgen`
lines up with it:

- **Length over complexity** — NIST prioritizes long passwords and explicitly
  drops mandatory composition rules (forced mixes of upper/lower/digits/
  symbols). A four-word passphrase is long and memorable without them.
- **Encourages passphrases** — NIST calls out multi-word secrets as a good
  practice; that is the entire premise here.

The one thing NIST-style guidance can't do for you is composition mandates from
*older* or non-conforming policies that still demand a digit or symbol — see
[`-digits`](#satisfying-must-contain-a-number-policies) for that.

## Word list

The built-in list (`words/words.txt`) is the EFF "large" Diceware wordlist:
7,776 common, unambiguous English words of 3–9 letters, chosen specifically for
passphrases. Dice indices are stripped; one word per line. Because it spans
only 3–9 letters, `-min`/`-max` outside that range return nothing unless you
supply a larger list with `-dict`.

To regenerate from source:

```sh
curl -fsS https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt \
  | awk '{print $2}' | grep -E '^[a-z]+$' | sort -u > words/words.txt
```

## Development

```sh
go test ./...    # run the test suite
go vet ./...     # static checks
```

## License

The code is licensed under the [MIT License](LICENSE).

The bundled wordlist (`words/words.txt`) is a derivative of the
[EFF Diceware "large" wordlist](https://www.eff.org/dice) by the Electronic
Frontier Foundation, licensed **CC BY 3.0 US**. Dice indices were stripped;
the words are unchanged. See [`words/LICENSE`](words/LICENSE) for the required
attribution.
