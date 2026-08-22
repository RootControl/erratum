# CLAUDE.md

THE ERRATUM — a playable vertical slice of a terminal investigation game
(Warhammer 40,000, grimdark noir). Go + Bubble Tea. Module
`github.com/RootControl/erratum`.

The repo is half design document, half prototype. `STORY_BIBLE.md` and
`ACT_I.md` / `ACT_II.md` / `ACT_III.md` are the source of truth for content and
vocabulary; `main.go` + `internal/` implement only scenes 1.3, 1.4, 1.9 and
1.10 of Act I. When the code and the act documents disagree about a name, a
flag or a fact, the documents win — and `STORY_BIBLE.md` §12 CANON LOCKS wins
over the rest of the bible.

## Commands

```sh
go mod tidy          # no go.sum is committed; nothing outside internal/game builds without this
go run .             # play it
go test ./internal/game/
gofmt -l .           # should print nothing
```

`internal/game` and `internal/content` have no third-party imports, so the test
suite runs without `go mod tidy`. `internal/ui` and `main.go` do not.

`ERRATUM_NOALT=1` disables the alternate screen. `smoke.py` drives the TUI
through a pty for headless checks: build first (`go build -o erratum .`) and it
picks up the binary beside it, or set `ERRATUM_BIN` to point elsewhere. Note
that its terminal-capability handshake echoes the OSC reply into the input
line, so the first command in a run is usually swallowed — pad the front of the
command list rather than trusting it.

## Layout

- `internal/content/` — `room.json` (embedded) plus its loader. All authored
  prose is meant to be data.
- `internal/game/` — `state.go` (state, clock, flags), `engine.go` (the verb
  parser and every command), `archive.go` (the Vault and the case-notes).
- `internal/ui/` — the cogitator, a Bubble Tea model. Viewport + text input.

Some long prose still lives as Go constants (`readingFull`, `readingPartial`,
`craseFallback`, `contradiction()`, the Vault records). That is drift, not a
pattern — new authored text belongs in `room.json`.

## Design invariants

These are load-bearing. `internal/game/playthrough_test.go` asserts most of
them, one test per guarantee, and those assertions are meant to be permanent.

1. **The archive never announces itself.** `Archive.Correct()` substitutes a
   plausible novice mistake into the player's own case-notes — same timestamp,
   same session, same author — and prints nothing. Never emit words like
   `DELETED` or `ERROR`, never flash the screen, never signal a correction.
2. **`internal/ui` has no vocabulary for "something is wrong."** No glitch
   effects, no colour shifts, no spacing changes on state. Keep it plain.
3. **Two evidence stores, never explained.** `FILE` puts a clue in the
   cogitator: actionable, requisitionable, and readable by the enemy. `NOTE`
   puts it in the player's coat: inert and safe. The game must not tell the
   player which is which.
4. **The notebook is append-only.** `State.Notebook` is a `[]string` with no
   edit path anywhere in the engine. Do not add one.
5. **A3 (the boots) can never reach the cogitator.** Its `store` is
   `"notebook"`; the filing template offers only distress / disorder /
   self-neglect. This is the one clue the record cannot hold.
6. **Free-form examination.** Nothing highlights itself. First `x` gives the
   careless read, second `x` gives the careful read and the clue; a third says
   you've read it. `Room.Find` matches loosely (aliases, then ≥3-char prefix)
   so a clue is never lost to vocabulary.

## Conventions

- Flags are `SCREAMING_SNAKE` strings in `State.Flags`, using the same names as
  the act documents (`ROOM_READ_FULL`, `A1_FILED`, `CIPHER_HALF_ROOM`). Filing
  sets `<ID>_FILED`; noting sets `<ID>_PHYSICAL`. There is no enum and no
  translation layer — keep it that way.
- Counters live in `State.Ints`: `VERITY` (independent proof), `RECONCILIATION`
  (corruption / exposure to the archive), `CRASE_STANDING`.
- `State.Do` returns `(out []string, timed bool)`. `timed` advances the clock.
  Inventory, help, `status`, reading your own notebook and rejected commands
  are free — the body is on a timer and checking what you hold must not cost.
- Clock constants in `state.go`: 40s per timed action, the body leaves at eight
  minutes, the archive reconciles 5 actions after the requisition, Crase's
  fallback fires at 9. These numbers are tuned; change them deliberately.
- Prose style: second person, short sentences under pressure, Imperial
  vocabulary used casually and never explained. See `STORY_BIBLE.md` §9.
