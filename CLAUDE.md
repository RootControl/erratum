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

`ERRATUM_NOALT=1` disables the alternate screen.

**`ERRATUM_PLAIN=1` is the harness for anything scripted.** It skips Bubble Tea
and plays over stdin/stdout, one block of text per command:

```sh
go build -o erratum .
printf 'x board\nx board\nfile A1\n' | ERRATUM_PLAIN=1 ./erratum
```

`smoke.py` drives the real TUI through a pty and is only useful for checking
that the TUI renders. Its output is every frame Bubble Tea painted,
concatenated — a non-alt-screen TUI repaints the whole viewport per keystroke,
so lines repeat and `tail` shows an earlier frame. That looks exactly like
commands being swallowed and is not; the older note in this file claiming the
first command is eaten was a misreading of the same artefact.

## Layout

- `internal/content/` — `arrival.json`, `antechamber.json`, `room.json`,
  `vaux_study.json`, `medicae.json`, `adjutant.json`, `arch.json` and
  `choir.json` (embedded, and loaded in that order) plus `strings.json` for the
  defaults every scene inherits, and the loader; `npc.go` holds the dialogue
  types. All authored prose is meant to be
  data, dialogue included.
- `internal/game/` — `state.go` (state, clock, flags, scenes), `engine.go` (the
  verb parser and every command), `dialogue.go` (ask / show / acts / beats),
  `archive.go` (the Vault and the case-notes).
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
   you've read it. `content.match` is loose on purpose (article-stripped, then
   aliases, then ≥3-char prefix or suffix) so a clue is never lost to
   vocabulary. Objects, people, topics, exits and clues all use it.
7. **Conversation is not a menu.** Topics are found by naming them, exactly
   like objects. `NPC.ListsTopics` is false for Vaux and `TOPICS` gets an
   in-voice refusal. The field exists only because 1.10 needs one listed
   conversation whose list can quietly lose a line at `RECONCILIATION` 10.
8. **A clue the player didn't ask for gets no bracket.** Beats (Vaux's
   mnemonic, A9) grant silently: the prose carries it and nothing labels it as
   evidence. `TestMnemonicIsNotAnnounced` asserts the string `A9` never reaches
   the screen when it arrives. A player who isn't listening loses it, which is
   the same thesis as the boots.
9. **A9 can never reach the cogitator**, for the same reason A3 cannot: the
   filing template has fields for a *statement* and a *transcript*, and a sound
   a man made without knowing he made it is neither.
10. **The first suspicion of the record is wrong, on purpose.** The Crown's
    verdict in 1.1 is struck two hours before the door was forced, and the
    explanation — a pre-generated watch form — is completely true. It is the
    only place in the game where the archive looks wrong and is not, and it is
    there so the player learns to distrust their own paranoia before they learn
    to trust it. Do not add a second tell to it.
11. **Four clues the record fails, four different ways.** A3 (the boots) — the
    form has no field for the truth. A9 (the mnemonic) — the utterance was
    never made by anybody who knew they made it. A14 (Bez's illiteracy and the
    sealed bequest) — the form has *exactly* the right field, pre-filled, nine
    seconds' work, and filing it would end her career, so the refusal is the
    player's and the text has to admit the machine would have taken it. A16
    (Fenn at the threshold) — the record accepts it, keeps it, and quietly
    holds *savant fatigue, recommend rest cycle* instead. Nothing announces
    that substitution; a player who only filed it simply has nothing.
12. **Crase is not lying in 1.2.** *"He'd stopped taking cases"* is false and
    she believes it, and the scene renders it as a plain answer with nothing
    marking it — no bracket, no tell, no hedge in the prose.
    `TestCraseBelievesTheLie` asserts the output carries no marker at all. It is
    the same rule as invariant 1, applied to a person instead of an archive.

## Conventions

- Flags are `SCREAMING_SNAKE` strings in `State.Flags`, using the same names as
  the act documents (`ROOM_READ_FULL`, `A1_FILED`, `CIPHER_HALF_ROOM`). Filing
  sets `<ID>_FILED`; noting sets `<ID>_PHYSICAL`. There is no enum and no
  translation layer — keep it that way.
- Counters live in `State.Ints`: `VERITY` (independent proof), `RECONCILIATION`
  (corruption / exposure to the archive), `CRASE_STANDING`.
- **The loader refuses content that would fail quietly.** `Game.validate()`
  rejects a clue reference that resolves to nothing, an exit to a scene that
  is not there, a `store` that is not one of the three, a `clock_at` that is
  not a time, a disposition default outside its own group, and duplicate
  topic / act / beat ids. Every one of those has already shipped broken in this
  repo at least once, and none of them raise anything at runtime — they just
  make a line of the game silently not happen. Add a check here when you add a
  reference.
- **Verbs that need a place stay in that place.** `QUERY` requires
  `"terminal": true` on the scene — there is exactly one cogitator in the slice
  and Vaux's study, which has a capped socket painted over, is pointedly not
  it. `READ ROOM` synthesises the room you are standing in: `State.Signs()`
  only counts held signs that the *current* scene defines. Both used to run
  anywhere and reach for chambers-only prose.
- **Prose authored as a struct field is read as a struct field.** `Look`,
  `LookNoBody`, `Intro` and `BodyRemoved` are fields on `content.Room`;
  everything else is a key in its `strings` table. Reaching for a field through
  `Room.S()` returns `[missing string: …]` and renders it — which is how `LOOK`
  and the body's removal both spent months printing a marker instead of their
  text. If you add authored prose, put it in `strings` unless it needs a field.
  `playFrom` in the test package fails any playthrough whose output contains
  `[missing string:`, so every test in the suite guards this.
- Shared defaults live in `strings.json`; a scene overrides by defining the same
  key. Keep scene-specific voice in the scene — `read_room_locked` reads as a
  1.4 beat in the chambers and as a plain refusal everywhere else, and that is
  two entries, not one clever one.
- Compare `Clue.Store` against `content.StoreNotebook` / `StoreCogitator` /
  `StoreEither`, never the string literals.
- `State.Do` returns `(out []string, ticks int)` — how many timed actions the
  command cost. Observing the world costs 1 (examining an object *or* a
  person); checking what you already hold is free (0) — inventory, help,
  `status`, your own notebook, and rejected commands. The body is on a timer
  and looking in your own coat must not cost. Travel and custody cost more than 1,
  set per exit / per act in the content.
- Mutually exclusive flags (`VAUX_DETAINED` | `_SHOWN_ROOM` | `_UNCHALLENGED`)
  are declared as an NPC `disposition` in the JSON and written through
  `State.SetOne`, which clears the siblings. They are still flat
  `SCREAMING_SNAKE` bools — the struct is the only thing that knows they are an
  enum, and there is still no translation layer.
- Dialogue gates read `requires`/`blocks` (flags), `has` (a clue in hand) and
  `holds` (a clue actually recorded in a store). Citing a register you read is
  knowledge; producing one is a document. Two topics may share a name: the
  engine takes the first whose gates the player has met, which is how accuracy
  gives Vaux a different answer. A closed topic is never spent — accuracy can
  open it retroactively.
- `State.Grant` puts a clue in hand; `FILE`/`NOTE` set `<ID>_FILED` /
  `<ID>_PHYSICAL` and both set `<ID>_HELD`, so `_HELD` means recorded, not
  merely heard.
- Clock constants in `state.go`: 40s per timed action, the archive reconciles 5
  actions after the requisition, Crase's fallback fires at 9. The lift to tier
  four is 3 actions each way and detaining Vaux is 6. These numbers are tuned;
  change them deliberately.
- **Two clocks.** `State.Clock` is hive time in seconds since midnight; it runs
  forward through a scene and jumps on entering the next one, to that scene's
  `clock_at` (1.2 opens at `04h42`, 1.3 at `06h12`). The body has its own
  counter, which starts the first time the player reaches the scene marked
  `"body": true` and does not stop when they walk out of it. Eight minutes
  means eight minutes of being able to reach him, not eight minutes of Act I.
- `CRASE_STANDING` opens at +1, per the bible's systems notes. Act I is built to
  push an attentive player to +2 or +3 from there.
- **The act runs on one clock and every authored hour has to agree with it.**
  The door came off its seal at the twenty-third hour and fifty-five; the
  Crown's finding is struck two hours earlier, at the twenty-first and
  fifty-five. Scenes open at `04h08` (1.1), `04h42` (1.2) and `06h12` (1.3),
  which is what makes the chambers' *"six hours ago"* and Crase's *"four
  hours"* both true. An accidental discrepancy in this game does not read as a
  typo — it reads as a clue. Check new prose against these.
- A scene's `sets` writes its flags once on arrival (`ACT1_OPEN` in 1.1); its
  `grants` puts clues in the player's hands once, unasked, with the usual
  FILE/NOTE bracket. That is how the autopsy floor (A11, A12) is guaranteed:
  standing on the medicae deck is the entire requirement.
- Counters move through `counters` on a topic / act / show — `{"BEZ_TRUST": 2}` —
  and the loader rejects any name not in `content.KnownCounters`, because a typo
  there is an award that silently never lands. `State.Add` holds each counter
  inside the range the bible gives it (`bounds` in `state.go`).
- **Gates, everywhere, are the same four fields.** `requires` (flags all set),
  `blocks` (flags none set), `has` (clues in hand) and, on topics, `holds`
  (clues actually recorded in a store) — on topics, acts, shows and beats
  alike. Exits add `requires_any` / `has_any`, which are one disjunction
  between them, because the Choir Spire opens on A2 *or* on Crase's lead.
- A topic's `counters` reward being answered; `closed_counters` charge for
  being pushed on. They are separate because they are not the same number —
  raising the timestamp with Crase is worth a point, and naming Ilaria's gaps
  costs two.
- **A person can be in two scenes and stay one person.** A scene lists people
  inline in `npcs`, or names somebody authored elsewhere in `npc_refs` — Fenn
  follows you to the gallery arch and is one definition seen from two rooms.
  Conversation state is keyed by `NPC.Key()` (the scene that authored them),
  not by name, so Crase greets you once in 1.1 and again in 1.2 with a rosette
  in a handkerchief, while Fenn keeps one conversation across the threshold.
  Iterate `Room.People()`, never `Room.NPCs`.
- `while` on a scene holds flags that are true only while the player is
  standing in it (`OUTSIDE_SHIELDING` at the arch). Relay shielding is a
  property of a place, not an event.
- `revisit` replaces `intro` on every entry after the first, falling back to
  `look`. An intro carries the beats of arriving somewhere — Fenn standing up,
  Ostro reading out the autopsy — and those happen once.
- Several topics may share a name; the engine takes the first whose gates the
  player has met, and when none are open it answers with whichever variant
  actually has a `closed` line authored. Fenn's `dates` is three topics on one
  name: inside, outside, and inside again.
- The medicae deck is gated on `BODY_REMOVED`. You cannot autopsy a man who is
  still on the floor upstairs, and that is what makes the timer feel like a
  door rather than a punishment.
- **The body has a deadline, not a stopwatch.** `bodyDeadline` is set to
  `Clock + bodyLimit` the first time the player reaches the chambers, and any
  clock jump past it takes him — so crossing to the Choir Spire, which opens at
  `07h40`, means the medicae have already been. Checked on arrival as well as
  on tick, so travel narrates it.
- 1.9 has two doors, per its entry gate: `A1_FILED`, or `DISPATCH_ID_KNOWN` and
  `QUERY 0997.4-1188`. `Archive.Correct()` deletes the record under both
  addresses, because the record is gone and not merely unindexed by one of
  them.
- A beat with no `after_topics` and no `after_pressure` but a non-empty
  `requires` fires the moment its gates open, rather than waiting on the
  conversation. That is how the Choir lead is delivered by the ceremony in 1.2
  and skipped entirely for a player who was kind to Crase in 1.1.
- Prose style: second person, short sentences under pressure, Imperial
  vocabulary used casually and never explained. See `STORY_BIBLE.md` §9.
