# CLAUDE.md

THE ERRATUM — a playable vertical slice of a terminal investigation game
(Warhammer 40,000, grimdark noir). Go + Bubble Tea. Module
`github.com/RootControl/erratum`.

The repo is half design document, half prototype. `STORY_BIBLE.md` and
`ACT_I.md` / `ACT_II.md` / `ACT_III.md` are the source of truth for content and
vocabulary; `main.go` + `internal/` implement all of Act I (scenes 1.1–1.10)
and Act II scenes 2.0–2.8. Act II 2.0b and 2.9–2.11 are not built, and Act III
is not started. When the code and the act documents disagree about a name, a
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

- `internal/content/` — one JSON file per scene, embedded and loaded in
  chronological order, which is what a scene with no `date` inherits from:
  - Act I: `arrival`, `antechamber`, `room`, `vaux_study`, `medicae`,
    `adjutant`, `arch`, `choir`
  - Act II: `board`, `safehouse`, `workshop`, `delta`, `choir_settlement`,
    `annexe`, `rows`, `cathedral`, `confessional`, `crase_office`, `ferrum`

  Plus `strings.json` for the defaults every scene inherits, `content.go` (the
  loader and its validator) and `npc.go` (the dialogue types). All authored
  prose is data, dialogue included.
- `internal/game/` — `state.go` (state, clock, flags, scenes), `engine.go` (the
  verb parser and every command), `dialogue.go` (ask / show / acts / beats),
  `archive.go` (the Vault and the case-notes).
- `internal/ui/` — the cogitator, a Bubble Tea model. Viewport + text input.

Some long prose still lives as Go constants (`readingFull`, `readingPartial`,
`craseFallback`, `contradiction()`, the Vault records). That is drift, not a
pattern — new authored text belongs in the scene's JSON, or in `strings.json`
if more than one scene needs it.

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
11. **Five clues the record fails, five different ways.** A3 (the boots) — the
    form has no field for the truth. A9 (the mnemonic) — the utterance was
    never made by anybody who knew they made it. A14 (Bez's illiteracy and the
    sealed bequest) — the form has *exactly* the right field, pre-filled, and
    filing it would end her career, so the refusal is the player's. A16 (Fenn
    at the threshold) — the record accepts it, keeps it, and quietly holds
    *savant fatigue, recommend rest cycle*. B6 (the oral census) — `store:
    "neither"`: it cannot be recorded in any medium at all, because the medium
    is four hundred mouths and writing it down stops them. FILE and NOTE give
    the same refusal, and that is the point.
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
- Counters live in `State.Ints` and are declared in `content.KnownCounters`:
  `VERITY` (independent proof), `RECONCILIATION` (exposure to the archive),
  `CRASE_STANDING`, `BEZ_TRUST`, `ILARIA_TRUST`, `FENN_DRIFT`, `MOTT_FAVOURS`,
  `QUERY_TIER`. `State.bounds` holds each to the range the bible gives it.
- **The loader refuses content that would fail quietly.** `Game.validate()`
  rejects a clue reference that resolves to nothing, an exit to a scene that
  is not there, a `store` that is not one of the four, a `clock_at` that is
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
- **The Case Précis is Act I's terminator and Act II's door.** It fires on its
  own (`precisDelay` actions after the requisition, or `precisFloor` actions
  regardless), intercepts every command except reading your own coat, and has
  three answers — one of which, ANNOTATE, is unavailable to a player whose
  notebook is empty. Resolving it sets `ACT1_COMPLETE` and moves the player to
  `standing_board`.
- **One `Gate` type, embedded everywhere.** Topics, acts, shows, beats and
  exits all gate on `requires` / `blocks` / `has` / `at_least` / `below` —
  counters included, because a savant three points into drift is a different
  man and the content has to be able to say so. Exits add `requires_any` /
  `has_any` for the one disjunction the act documents ask for. Add a condition
  once, in `content.Gate`.
- **The Cathedral is a whole thread the player did not need.** Everything Vhail
  gives up under forty-five minutes of interrogation is readable off
  `QUERY CATHEDRAL ALMS` — unmarked, ungated, available from the first hour of
  Act II, and pointed at by nothing. The scene must stay worth playing anyway,
  and the slip must never be signposted. That is UI-lie #4 in the act notes and
  the whole answer to *who can rewrite the record* is that it was never Vhail.
- **The death in 2.8 is a deterministic cascade, in this order, so it can be
  tested**: Fenn if brought, else Bez, else Vaux if allied, else Ilaria. It is
  five acts sharing one name — `pickAct` takes the first whose gates pass, the
  same rule topics use — so the outcome is authored, ordered and readable in
  the content rather than hidden in Go. `clears` on a line of dialogue is how a
  death takes `FENN_INTACT` and `_ALIVE` with it.
- **No two clues may answer to the same name.** `TestNoTwoCluesShareAName`
  enforces it, because the finder cannot see the difference and will silently
  record the wrong thing. `MatchClues` returns every match and the caller picks
  one the player is actually holding, which is a second line of defence, not
  a licence to collide.
- **Every flag the content gates on must be written by something.**
  `TestEveryGatedFlagHasAWriter` enforces it. A gate on a flag nothing sets is
  a line of the game that silently never happens, which is exactly what
  `PURGE_ORDER_SEEN` was for two scenes. Flags written by the engine rather
  than by content are listed in that test; add to the list only when the
  engine really is the writer.
- **Crase asks at two threads.** `State.threadsTouched()` counts a held clue
  from each of Vosk's four open files; the second one summons the player by
  runner. Refusing the Gallery brief is the only route that keeps Ending E
  open, costs two points of standing, starts the surveillance, and seeds
  `SUBLEVEL_NINE_SEEDED` for the Reconciliation-40 symptom. No dialogue is
  attached to the refusal and she stays warm throughout.
- **2.5 attaches nothing, and that is load-bearing.** The Rows have no clue,
  no unlock, no stat and no asset. `CHOIR_FED` and `CHOIR_DEFERRED` are read
  once each, as one sentence each, in the Act III epilogue. Every way to help
  costs something real — standing, a Cathedral lever, Vosk's account, a week of
  Bez's lung — and buys the player nothing. Bolt a reward on and the act loses
  its spine.
- **`MOTT_FAVOURS` moves one way and doubles when hostile.** Coin needs a
  ledger; favours do not. Forcing the Rust Choir's gate with the rosette sets
  `CHOIR_HOSTILE` and everything costs twice for the rest of the act, and
  nothing ever tells the player that is why. Handled in `State.Add`, beside the
  drift rule.
- **Dates are a floor, not a fact about where you are.** A scene's `date` is
  the earliest it can happen on: later advances the calendar and resets the
  time of day, earlier or equal only runs the clock forward. `TimeString` reads
  `State.date`, never `Room.Date`, or walking back to the hub rewinds the act.
  Dates compare as strings, which holds only because days are written three
  digits wide — keep them that way.
- **`FENN_DRIFT` has two rules attached, both in `State.Add`.** Under Cael's
  shroud it accrues at a third the rate (`driftPartial` carries the
  remainder). At 3 it sets `FENN_BURNED_PARTIAL`; at 5 it clears
  `FENN_INTACT` and forecloses Ending E. Neither threshold announces itself,
  and the only thing the player ever sees is the readout going down.
- **Reconciliation removes; it does not close.** A topic with `hidden_above`
  set stops *matching* past that threshold, so the person answers as though
  they had never heard of the subject and there is no gap to notice. Never
  render a removed option as unavailable, greyed, or numbered around.
- `FENN_DRIFT` is displayed only as `CONCORDANCE OF RECORD: n/6` in `status`,
  which is equipment-condition framing applied to a man. **Never relabel it**,
  and never show the counter's own name.
- A person authored in two scenes is two conversations (`NPC.Key()`), so
  `npc_refs` must be scoped as `"scene/id"` once a name exists twice — Fenn is
  authored in `adjutant` and again in `safehouse`, and the loader rejects an
  ambiguous bare reference with the fix in the message.
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
