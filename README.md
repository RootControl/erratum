# THE ERRATUM — vertical slice

A playable slice of Act I: the locked room (1.3), the synthesis (1.4), the Vault
requisition (1.9) and the Act I turn (1.10). Go + Bubble Tea.

It exists to prove four things are buildable before any more content is written:

1. **Free-form examination.** No highlights, no glints, no prompts. Objects are
   found by naming them, and looking twice is not the same as looking once.
2. **Two evidence stores.** `FILE` puts a clue in the cogitator: indexed,
   actionable, and reachable by the enemy. `NOTE` puts it in your coat: inert,
   useless, and yours. The game never explains the difference.
3. **A flag engine** speaking the same SCREAMING_SNAKE vocabulary as the act
   documents, so design and code don't need a translation layer.
4. **An archive that edits your case-notes** — a real mutable store, not a
   scripted cutscene. `Archive.Correct()` mutates the data and says nothing.

## Run

```sh
go mod tidy
go run .
```

Then: `x board`, `x board`, `file A1`, `x boots`, `x boots`, `note A3`,
`read room backward`, `query C9-F4-D2`, `note A20`, wait five turns,
`query C9-F4-D2` again.

`ERRATUM_NOALT=1` disables the alternate screen, which is how `smoke.py`
drives it headlessly.

## Test

```sh
go test ./internal/game/
```

Five tests, one per design guarantee:

| Test | Asserts |
|---|---|
| `TestCarefulPlaythrough` | five signs + the boots → `ROOM_READ_FULL`, `CIPHER_HALF_ROOM`, and A21 provable |
| `TestTrustingPlaythrough` | file everything, write nothing → the edit happens and **cannot** be proven |
| `TestBootsCannotBeFiled` | the filing template offers only distress / disorder / self-neglect |
| `TestBodyTimerClosesTheBoots` | eight minutes, then A3 is gone for good |
| `TestVaultRequiresFiling` | you cannot requisition against a code you only told your coat about |

`TestTrustingPlaythrough` also asserts the archive never prints the words
`DELETED` or `ERROR`. That assertion is the thesis of the whole game and it
should stay in the suite forever.

## Layout

```
internal/content/   room.json + loader — all authored prose is data
internal/game/      state, flags, command engine, and the archive
internal/ui/        the cogitator. Deliberately has no vocabulary for
                    "something is wrong": no glitch effects, no colour shifts.
```

## What the slice deliberately does not have

Save/load, Crase and Bez as present characters, scenes 1.1/1.2/1.5–1.8, the
Reconciliation symptom system, and any content past the Act I turn. The clock
is 40 seconds per timed action; the body leaves at eight minutes and does not
negotiate.

## Note on `go.mod`

Built in a sandbox that could reach github.com but not `proxy.golang.org` or
`golang.org`, which needed `replace golang.org/x/... => github.com/golang/...`
directives. Those are stripped from the shipped `go.mod`; `go mod tidy` on a
normal machine will resolve everything directly.
