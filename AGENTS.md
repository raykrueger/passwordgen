# AGENTS.md

Single-binary Go CLI that generates memorable passphrases from random words.

## Layout

Flat `package main` at the repo root (no `cmd/`):
- `main.go` — flag parsing in `run() error`, plus `loadWords`, `passphrase`, `choose`, `randomDigits`.
- `words.go` — `//go:embed words/words.txt` into `embeddedWords` (the default wordlist).
- `words/words.txt` — the embedded EFF Diceware list; see below.
- `main_test.go` — table-driven, parallel tests.

## Commands

```sh
go build -o passwordgen .   # build
go test -race ./...         # test (all logic is unit-testable, no I/O to mock)
go vet ./...                # must stay clean
gofmt -l .                  # must print nothing
go run .                    # try it
```

`staticcheck` currently can't run here: its binary predates the Go 1.26 toolchain and errors out. `gofmt` + `go vet` are the gates.

**These checks are not automated** — there is no CI workflow, pre-commit hook, or Makefile. Run `gofmt -l .`, `go vet ./...`, and `go test -race ./...` manually before every commit; nothing else will catch a regression.

## Invariants (don't break these)

- **Randomness is always `crypto/rand`, never `math/rand`.** This is the whole security premise; every random draw goes through `choose`/`randomDigits`.
- **Keep the `run() error` pattern.** `main` only calls `run` and exits. Do not add `os.Exit` inside helpers — the deferred `f.Close()` in `run` would be skipped.
- **`choose` panics on an empty slice** (`crypto/rand.Int` bound must be > 0); callers guard `len(words) == 0` first.

## Wordlist

`words/words.txt` is the EFF "large" Diceware list (words only, dice indices stripped), licensed **CC BY 3.0 US** — do not edit by hand and keep `words/LICENSE` attribution intact. Regenerate with:

```sh
curl -fsS https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt \
  | awk '{print $2}' | grep -E '^[a-z]+$' | sort -u > words/words.txt
```

It spans only 3–9 letter words, so `-min`/`-max` outside that range return nothing unless `-dict` points at a larger list. `TestEmbeddedWords` asserts a floor on the 4–6 letter count; update it if the list changes.

## Conventions

- Go 1.26; use modern idioms (`for range n`, `slices.Equal`, builtins). Don't name variables `min`/`max` — shadows builtins (`minLen`/`maxLen` used instead).
- Conventional commits, signed (`git commit -s`). Module path is `github.com/raykrueger/passwordgen`; tag releases (`v0.1.0` exists) so `go install ...@latest` resolves.
