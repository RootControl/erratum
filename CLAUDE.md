# CLAUDE.md

THE ERRATUM — a playable vertical slice of a terminal investigation game
(Warhammer 40,000, grimdark noir). Go + Bubble Tea. Module
`github.com/RootControl/erratum`.

The repo is half design document, half prototype. `STORY_BIBLE.md` and
`ACT_I.md` / `ACT_II.md` / `ACT_III.md` are the source of truth for content and
vocabulary; `main.go` + `internal/` implement all of Act I (scenes 1.1–1.10)
every scene of Act II, and the whole of Act III including the five endings.
The slice is content-complete. When the code and the act documents disagree about a name, a
flag or a fact, the documents win — and `STORY_BIBLE.md` §12 CANON LOCKS wins
over the rest of the bible.

## Commands

```sh
go mod tidy          # no go.sum is committed; nothing outside internal/game builds without this
go run .             # play it
go test ./internal/game/
gofmt -l .           # should print nothing
```

`internal/game` and `internal/content` have no third-party imports, so
`go test ./internal/game/` runs without `go mod tidy`. `internal/ui` and
`main.go` do not, so `go test ./...` needs the dependencies — which building
needs anyway.

`ERRATUM_NOALT=1` disables the alternate screen. `ERRATUM_SAVE` moves the
session file, which otherwise lives at `./erratum.session`.

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
  - Act II: `board`, `safehouse`, `workshop`, `delta`, `theta`,
    `choir_settlement`, `annexe`, `rows`, `cathedral`, `confessional`,
    `crase_office`, `ferrum`, `cylinder_vault`
  - Act III: `skydock`, `cerement`, `spirit`, the four watches — `gunbay`,
    `cell`, `hold`, `wardroom` — Kadmon Veil: `kadmon`, `tithe_yard`,
    `muster`, `reading_room` — and Cantica-Null: `cantica`, `fifth_chair`,
    `correctorum`, `intake`

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
- **`STORY_BIBLE.md` §12 fixes the present at 826.M41 and every authored year
  derives from it.** Kadmon terminates 735, the Ferrum cult begins 763, Vhail
  is recruited 786, Mott's district is struck 807. Three dozen years drifted
  into the 900s across Act III before anybody checked, which in a game that
  trains the player to read every number as evidence is worse than a typo — it
  reads as a clue. `TestNoDateAfterThePresent` lints the scene files directly
  and fails any three-digit number above 826 that is not on its allowlist.
- **Cycles are not years.** The stamp the game prints — `0997.4.211` — is a
  cycle, quarter and day, and has nothing to do with the Imperial year. A bare
  three-digit number in prose is an Imperial year and must precede 826. Act I
  already does both: a rescript from 604 in a scene dated 0997.4.211. Anything
  with a day component is a cycle stamp; anything bare is a year.
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
- **The store habit is behaviour, not a menu.** `STORE_HABIT_FILED` /
  `STORE_HABIT_LEDGER` are written by `storeHabit()` in `engine.go` on
  whichever of FILE and NOTE the player reached for first after `ACT2_OPEN`,
  and there is no scene where either is chosen. Filing raises `QUERY_TIER`,
  because the Vault is pleased with anybody who feeds it. Keeping a ledger buys
  nothing anywhere in Act II except one thing: it is the only honest answer to
  Vaux's question in 2.9, and a player who filed everything cannot give it —
  not because a gate refuses, but because there is nothing in the coat to read
  out. Neither flag is ever mentioned on screen.
- **How Vaux was asked for is most of what happens in 2.9.** `VAUX_SUSPICION`
  is a counter bounded 0–3. A request through the Crown's message office
  arrives at 2 (the `approach_archive` beat); a Rust Choir runner, one favour
  from Mott, arrives at 0. The alliance needs it at ≤1 *and* `VAUX_PERSUADED`,
  which the ledger answer or a shown B4/B12 both set while taking a point off —
  so the archive route can still ally. What it cannot do is ally at standing:
  `below CRASE_STANDING 1` on the archive route against `below 3` with a
  runner, which is I-4's actual teeth. Refusing all three doors and walking out
  sets `VAUX_LEFT`; the `returned` beat then puts him at 3 (the ceiling does
  the arithmetic) and sets `VAUX_FINISHED`, which blocks every persuasion and
  both of the other two doors. None of this is announced.
- **I-5 is authored as absence.** Each point of `CRASE_STANDING` over +2
  removes a persuasion topic in 2.9 — `honest` at +3, `character` at +4, both
  via `below` — and `indictment_promoted` adds a line to what Vaux says about
  you. A removed topic answers with his ordinary "nothing on that". The player
  is punished for playing the politics well and is never told that is what
  happened.
- **Beats carry `counters` and fire on approach.** A beat was the one
  gate-bearing type with no way to move a number, which is wrong: "he is
  already two points suspicious of you because of how you asked" is exactly a
  beat. `talk` runs `beats()` too, so a beat waiting only on its gates fires
  when the player walks back up to somebody — before that it could only land
  on a question they had not already asked, which is not what "the moment its
  gates open" means.
- **2.10 has no scene file, on purpose.** The long way round is authored on
  three people who were already somewhere: Ilaria at the standing board, Fenn
  in Stack Delta (`IN_STACK_DELTA`, the same `while` idiom the Rows use), and
  Mott in the sump. The act document's location line reads *"Choir Spire
  cylinder vault (Ilaria) or Stack Delta (Fenn). Same shape either way; the
  cost is the person"* — the shape is the two acts having identical structure
  and identical refusals to let the player dodge, and the person is wherever
  the person already stands. Adding a room would have made the choice a place.
- **`CIPHER_COMPLETE` is derived, not written.** `State.Set` turns it on the
  moment `CIPHER_HALF_ROOM` and one of `CIPHER_HALF_VAUX` / `CIPHER_HALF_REBUILT`
  are both true, in either order — 2.0b supplies the room half four days late,
  so the order really does vary. It was a flag 2.10's three acts remembered to
  write, which meant the two doors in 2.9 that hand the second half straight
  over left the player holding all of it and unable to read a word. No content
  writes it; do not add a writer.
- **Nothing supplies the second half for free.** `TestNoFreeWayToCompleteTheCipher`
  walks the content and fails any route that writes `CIPHER_HALF_VAUX` or
  `CIPHER_HALF_REBUILT` — counting the `sets` of any clue it grants, which is
  how B16 arrives — without also writing one of `ILARIA_BURNED`, `FENN_SPENT`,
  `NAMEBANK_DEBT`, `VAUX_ALLIED` or `VAUX_COERCED`. It names the six routes
  rather than counting them: the earlier version asserted "exactly three", and
  that count was the soft-lock written down as a requirement.
  `TestEveryDoorFinishesTheAct` walks all three of 2.9's doors through to
  `ACT2_COMPLETE`, which is the test that would have caught it.
- **A counter that is only a display moves with the flag, not with `Add`.**
  `ILARIA_BURNED` clamps `ILARIA_TRUST` to 2 in `State.Set`, and `Add` then
  refuses to let it climb back. She does not lose the player's regard — she
  loses the thing the regard was for, and the number has to follow the event
  rather than be charged for it.
- **Cael's shroud is not armour against a decision.** `spend_fenn` clears
  `FENN_SHIELDED` in the same act that costs six of drift, and `apply` runs
  `clears` before `counters`, so the exposure lands in full. The shroud thirds
  what happens to a savant incidentally; it has nothing to say about eleven
  deliberate hours at a live terminal, and the content says so by taking it
  off him first.
- **An act that cannot be undone must not answer to a noun.** Bare input falls
  through to `act(verb, arg)`, so an act named `"the corpus"` fires from two
  typed words — and `"the corpus"` is also the name of the topic where Fenn
  explains what it would cost. One-way acts take gestures (`let him do it`,
  `send her down`); the nouns stay on the topics and the exits. `ask` and
  `take` are parser verbs, so an act named `"ask fenn"` is unreachable and an
  in-voice line telling the player to *ask me for it* is a dead end — check the
  prose against the vocabulary the act actually has.
- **`DECODE` is place-scoped, like `QUERY`.** `Room.Cipher` marks the room the
  objects were arranged in, and the standing board is Vosk's chambers
  unsealed, so it is the only scene with the flag. The sentence is read out of
  the *record* of the arrangement rather than out of the furniture — 2.0 says
  in its own intro that the Crown packed and unpacked the room by people who
  did not know the arrangement mattered — so what the decode actually spends is
  the player's Act I notes. A3, the one clue the cogitator could never hold, is
  the eighth object and is said backwards.
- **Canon rule: decrypting the testimony grants nothing but the testimony.**
  `TestTheSentence` asserts `ROOM_READ_FULL` is not set by `DECODE`. The
  sentence is *what* Vosk wrote, not *how*, and Ending E needs the method. The
  only late route to the method is 2.0b.
- **The Erratum is `StoreNeither` and the testimony is `StoreNotebook`.** Both
  refusals are the player's own and neither explains itself: B21 gets as far as
  the correct field on the form and puts the pen down, and B19 opens a good
  template, reaches the third field, and closes it. Nothing in the UI forbids
  either — the game just declines, in the player's voice.
- **The coerced cipher half lives exactly one scene.** With
  `CIPHER_HALF_VAUX_CORRUPT`, the first `DECODE` returns a fluent, plausible
  sentence naming Vhail, every clause of which the player proved themselves,
  and `TestTheCorruptHalfNamesTheWrongMan` asserts the screen carries no hedge
  at all. The second `DECODE` is the count — nine objects, eight tokens — and
  costs a point of `VERITY` to unpick. The reveal is arithmetic, not suspicion,
  which is the same trick the whole act runs on.
- **The plate cannot be missed and the log can.** `Tick` gives the player B21
  after `tassDelay` actions with `CIPHER_COMPLETE` and no plate, brought up
  four tiers by an apprentice. What that route does not carry is B22, Cael's
  signature on the recovery, because nobody asked him to do it — so the Act III
  leverage is the price of never having gone to the workshop.
- **`<ID>_HELD` / `_FILED` / `_PHYSICAL` are engine-written for every clue.**
  `TestEveryGatedFlagHasAWriter` derives them from `Game.Clues` rather than
  listing them, so gating a topic on `A4_HELD` needs no bookkeeping. Add a name
  to that test's `engine` map only when the engine really is the writer —
  `CIPHER_COMPLETE` and `KADMON_ADDRESS` are deliberately not in it, because
  2.10's acts and B20's own `sets` write them.
- **One door in the game is `hidden`.** `Exit.Hidden` keeps Stack Theta out of
  the list a bare `GO` prints; it still opens to its own name. 2.0b is
  documented as never signposted to anyone, and a door that offers itself is a
  signpost. Its only other gate is `THETA_KNOWN`, written by requisitioning the
  effects manifest from the Vault — the discovery is a query, which is the verb
  the game has spent two acts training the player on, and nothing anywhere
  points at it. Do not add a second hidden exit without the same kind of
  reason.
- **2.0b is open to a player who read nothing, and the act document has been
  amended to say so.** Its entry gate was `ROOM_READ_PARTIAL`, which assumed
  the room always gets read. It does not: `READ ROOM` counts signs actually
  *recorded*, so a player who examines everything and files nothing reaches Act
  II with no half, and the Case Précis ends Act I on a timer anyway. That
  player could not finish the act and could not see the scene that gives the
  room back. The gate is now the absence of `ROOM_READ_FULL`. Do not put
  `ROOM_READ_PARTIAL` back.
- **Ilaria's `crates` beat is the one signpost 2.0b has, and it is aimed at one
  player.** It `blocks` on `CIPHER_HALF_ROOM` and `THETA_KNOWN`, so it fires
  only for somebody holding no half who has not already found Theta. She says
  the Provost crated the effects and the Vault has them, as small talk, and
  does not say to go. Anybody who read the room never hears it —
  `TestTheCratesAreNotMentionedToEverybody` asserts that.
- **Objects gate, and a closed object is absent.** `Object` embeds `Gate` like
  everything else; `examine` and `take` treat a gate-failing object as though
  it were not in the room, because a refusal would tell the player there is
  something in the crate. `Object.Sets` fires on the careful look and
  `Object.TakeSets` on the taking — reading a dead man's longhand off his desk
  is not the same as walking out with it, and `HAS_VOSK_NOTEBOOK` used to be a
  line in `take()` that named one object.
- **A clue salvaged in Theta is worth one `VERITY` less.** The cage objects set
  `<ID>_SALVAGED` and `note()` skips the award for those. Four days late and
  one remove from the room is not the same as reading it off the floor with the
  body still in it, and nothing on screen says so.
- **The Reconciliation-55 symptom lands in 2.0b, where it costs something.**
  `reconstruct` needs the five signs recorded *and* `below RECONCILIATION 55`;
  past that, `reconstruct_short` fires instead and the Act I file comes back
  thin — four scene captures where there were eleven, the autopsy entry at 190
  words, *no residue was detected* where it said *no toxin, no coercion residue
  and no psychic trace*. Nothing has been struck, there is no amendment mark,
  and Ending E closes without a line of text saying it did.
- **`vault_none` lives in `strings.json` now.** It was authored in the chambers
  only, so every Act II query for a code the Vault does not know rendered
  `[missing string: vault_none]` at the standing board. Shared refusals belong
  in the shared table; a scene overrides only when the refusal is a beat.
- **Act III inverts the UI lie: out of relay range the cogitator stops lying
  about the world and starts lying about the player.** `IN_TRANSIT` runs a
  symptom every `symptomInterval` actions, in order, out of `State.symptom()`:
  a note one clause shorter, a sentence the player cannot find again, and — only
  at `RECONCILIATION` 60 — `NO SUCH ENTRY` for something they requisitioned
  themselves. Below 60 the third simply does not fire and the passage is quiet,
  and nothing tells the player they were spared it. The strings live in
  `cerement.json` as `symptom_note` / `_option` / `_memory`.
- **`OUT_OF_RELAY` makes reading your own case-file cost exposure.** Two
  `RECONCILIATION` per look, in `engine.go`, silently. It is still free in
  ticks — the invariant that checking what you already hold costs no time is
  intact — but out there the file agrees with the player slightly more every
  time they go back to it for reassurance.
- **`VOSK_METHOD_PARTIAL` / `_UNDERSTOOD` are derived, like `CIPHER_COMPLETE`.**
  `State.Set` reads `C2_HELD`, `C3_HELD` and `C5_HELD`: the nav plot and the
  casket are the shape of it, and the uncoupled vox is what makes it a method.
  No content writes either flag. This is the same rule as the cipher and for
  the same reason — the last flag three scenes had to remember to write made
  the act unfinishable.
- **Objects resolve like topics and acts.** `findObject` in `engine.go` takes
  the best-scoring object *whose gates are open*, so two objects may share a
  name and only one of them is in the room. That is how the voice slate is one
  object on a stolen hull and a different object on a curated one, with no
  branch anywhere saying which.
- **The three hulls in 3.1 are exclusive by `blocks`, not by a disposition.**
  `TRANSPORT_WRIT` / `_STOLEN` / `_DEBT` each block the other two on the act
  that sets them, which is the 2.9 idiom. Route B is authored on Rill because
  Rill is the man who explains the seal and then puts the shift-return in late;
  Bez walks the party aboard in the prose.
- **The four watches are four scenes, because the player has to go to them.**
  3.4's rule is that no member of the retinue raises their own subject
  unprompted, so each of `gunbay` / `cell` / `hold` / `wardroom` is its own
  room off the Cerement, reached by an exit that costs ship-time, and the
  ambient objects in `cerement.json` are the same four people seen from a
  distance and saying nothing. Looking at somebody is not standing a watch
  with them.
- **`VOYAGE_FULL` counts what was possible, not what was offered.** Derived in
  `State.Set` when a `WATCH_*` lands: a member who is dead, below trust, blank
  or never allied is not a watch the player failed to keep, so a lone survivor
  can still finish the voyage complete. Do not make it a count of four.
- **`C9_PHYS` / `C9_COG` are the act document's names for `NOTE` and `FILE`.**
  `State.Set` translates `C9_PHYSICAL` and `C9_FILED` into them rather than
  asking the content to write a second flag beside the store it already chose.
  Filing the forward recitation is the trap: the `file_text` says, in the
  filing, that the list will be three worlds shorter at Cantica-Null, and
  nothing warns.
- **Three of the four watches have a wrong answer that is the kind one.**
  Softening Bez's bequest costs `RECONCILIATION` and the prose says both
  reasons were good ones; telling Ilaria she is ill costs six and says plainly,
  once, that it is the same instrument for the same reason; shielding Fenn
  mid-recitation is the humane choice and truncates the evidence. None of them
  is punished by a refusal — each one works, and takes something.
- **3.6 and 3.8 are one scene, because they are one conversation.** The
  Marchioness's audience and the answer to her question happen in the Muster
  Hall with the same woman, so `muster.json` holds both: her report and the
  census are topics, the question is a beat at `after_topics` 2, the three
  answers and the four aid options are acts that block each other. Splitting
  them into two scenes would have meant leaving the hall and coming back to
  say one sentence.
- **`VURR_RESPECT` is opened by competence and closed by sympathy.** Showing
  her C11 or C13 accurately is +1 and grants the Reading Room; the `pity` topic
  is −1 and she addresses the retinue rather than the player for the rest of
  the visit. The Reading Room's exit gates on `at_least VURR_RESPECT 1`, so
  3.7 is earned in 3.6 or reached at all.
- **The aid in 3.8 says its own ratio out loud.** Each of the four acts states
  the fraction of the shortfall it covers, in her voice or the trader's, and
  each ends with the shrine that gets built to it inside a week. That is the
  scene's knife and it is authored into the act text rather than left to the
  player to work out. `KADMON_AID` blocks the other three: one is a gift, two
  is a policy.
- **Doing nothing still fires the departure.** The `thanks_nothing` beat is
  gated `blocks: ["KADMON_AID"]` and fires at `after_topics` 4: she comes down
  eleven hundred metres at a hundred and four and thanks the player for
  coming, with no condition attached. The act document says testers found it
  the most upsetting line in the build and to leave it in.
- **Report 728 is authored as a gated object with no signpost.** It is one of
  nine hundred and twenty-two spines until the player asks Rhod what year the
  Marchioness's husband died; `VURR_WIDOWED` is what puts it in the room. The
  filing template for it opens, offers four dispositions, and all four are
  lies — the refusal is the player walking back to the desk and making the
  mark she asked for.
- **C17 is the one clue the archive takes back.** Cael files it himself,
  helpfully, without being asked, so the topic sets `C17_FILED` rather than the
  player choosing a store. `Tick` then removes it after `c17Delay` actions
  unless `C17_PHYS` — and `Archive.Delete` takes the case-note line out
  entirely rather than striking or amending it, so there is no gap. Re-querying
  returns `NO SUCH ENTRY`, which is the same three words a shelf that was never
  filled returns. It is the only place in the game where the archive removes
  rather than substitutes, and it is Cael's own evidence of his own
  reconciliation.
- **The Correctorum never claims authorship, and the player has to notice.**
  The `personnel` topic reads out the assessment and stops; nothing in it says
  the machine composed *acceptable for expenditure*, and every player hears
  that it did. The `composed` topic is the precise question and the answer is
  one word. `AMANUENSIS_BLINDSPOT` is the hard gate for Ending E and it is
  reachable only by asking — `TestTheBlindspotMustBeAskedFor` asserts that
  hearing the file is not asking, and that nothing in the reading claimed it.
- **Stopping the reading in 3.11 is relief that hides a price.** `stop_reading`
  blocks the `personnel` topic, so C19 is never heard, and the act document is
  explicit that Ending D then charges for it without the player knowing what
  they are being charged for. Do not add an explanation.
- **Crase's audit is the scene's whole mechanism and she passes it.** Every
  answer under audit is true, including the elevation, and the reveal is not
  that she is theirs — it is that she runs the Amanuensis's rule, omission
  never falsehood, and taught it to the player, who has been doing it since
  Kadmon. The one thing the player ends up holding that she does not is why
  Vosk asked *her* to sign.
- **3.12 has its own three verbs and they exist for one scene.** `ASK` is the
  parser's; `PRESS` walks `Topic.Layers` one at a time and says *"There is no
  further layer"* when the layer is empty; `CITE` is `SHOW` with the two-store
  rule attached — it works only from the notebook, because the machine cannot
  contradict a record it holds and has already amended everything in the
  case-file. `intake.json` overrides the shared `help` string to teach `PRESS`
  and `CITE`, which is why they are not advertised for thirty hours before they
  exist. That override is the intended way to scope a verb to a place.
- **Every line the Amanuensis speaks is literally true.** It answers the
  question asked, scoped narrower than the player meant, and never volunteers a
  scope. It says *I* twice. The omission — that Thessaline was the writing
  sample and the real projection is 410 million against 31 billion — has three
  doors: the literal scope question, `CITE C17` and `CITE C9`, and all three
  require the physical copy. On exposure it answers in full immediately with no
  defensiveness, because nothing has been caught.
- **An ending is an act with a long gate and a physical gesture.** The five
  doors live on the Amanuensis in `intake.json`, each with refusal variants
  ordered after the real one, so a closed door answers in fiction and never
  names the flag it wanted — `TestEndingEHasNoPartialCredit` drops each of
  Ending E's ten requirements in turn and asserts both that E stays shut and
  that the refusal does not say which one is missing. `State.Set` marks the
  game `Done` on any `ENDING_*`.
- **`VAUX_APOLOGY` holds every door until the word is kept.** Five `untold_*`
  acts sit ahead of the five endings and refuse, in fiction, until `tell_vaux`
  has fired. That is how an implication — *if you promised, you must tell him
  first* — is expressed with gates that have no implication operator.
- **`walk_test.go` is the only thing in the suite that plays the game.**
  Everything else teleports into a scene with the flags it needs, which is
  right for testing one guarantee and says nothing about whether the scenes
  join up. Five walks run from `New(g)` to five different doors on typed
  commands alone, sharing prefixes through `actOne()`, `twoThreads()`,
  `vauxAllied()`, `actTwoTurn()`, `stealAndSail()`, `kadmonTo()` and
  `theReading()`. Both of this repo's soft-locks lived in exactly that gap —
  every scene correct, no legal path between two of them — and breaking either
  one fails all five walks.
- **A script is only evidence if every line of it did something.** `walk()`
  fails on `That is not a verb this terminal knows` and `You hold no such
  clue`, so a mistyped command cannot sit in the scripts for years looking like
  coverage. It also asserts, over a whole playthrough rather than one scene,
  that the words `DELETED`, `ERROR`, `CORRUPTED` and `TAMPER` never reach the
  screen. That is the thesis of the game and this is the widest place to
  assert it.
- **Strings read from `Tick` belong in `strings.json`, always.** The three
  Cerement symptoms were authored in the scene that introduces them and fire
  from `Tick` wherever the player happens to be standing, so walking down to
  the cogitator well mid-passage printed `[missing string: symptom_note]`. The
  first playthrough test found it on its first run. Anything emitted outside a
  scene's own verbs — symptoms, delayed confirmations, the Tass delivery — has
  to be in the shared table.
- **Anything emitted from `Tick` must live in `strings.json`.** `Tick` calls
  `s.S(key)` against whichever scene the player is standing in, not the scene
  that authored the line, so a delayed emitter authored locally prints a marker
  the moment the player walks somewhere else. `vault_none` did it and the three
  Cerement symptoms did it. `TestDelayedTextResolvesEverywhere` arms every
  delayed emitter in the game, stands in all thirty-seven scenes in turn and
  runs the clock, which is the guard both of those needed.
- **`<ID>_FILED` and `<ID>_PHYSICAL` mean the same thing however they are
  written.** `State.Set` intercepts both and calls `storeFromFlag`, so content
  that declares a filing — Cael putting the intake protocols into your
  case-notes without being asked — actually stores the clue and writes the
  case-note line. Before that, C17 claimed to be filed, the flag agreed, and
  the case-notes were empty, so the archive took back something the player had
  never seen. `FILE` and `NOTE` store before they set, so the interception is a
  no-op for them.
- **`CONCLAVE_CHANNEL` has three sources and only one of them is a person.**
  Vaux's dead-man writ, an unspent `TRANSPORT_WRIT`, or `SAZH_DEBT_2`, derived
  in `State.Set`; spending the writ at Kadmon (`KADMON_WRIT`) revokes it unless
  one of the other two is standing. The build had only Vaux's, which closed
  Ending C to two of the three routes the act document gives, with no fiction
  explaining why.
- **Two permanent guards on parser collisions.**
  `TestNoActIsNamedForAParserVerb` fails an act whose first word `Do()` eats —
  `ask fenn`, `write it down`, `file a correction` — and
  `TestNoTwoExitsShareAName` fails two doors in one scene answering to the same
  word, because `FindExit` is gate-blind and the second is then unreachable by
  it however the gates fall. `GO`, `READ`, `TAKE` and `PRESS` deliberately fall
  through to `act()` and are not listed.
- **An ending holds the screen; quitting does not.** The renderer runs on the
  alternate screen, so `tea.Quit` the instant `State.Done` went true restored
  the terminal on top of a thousand words the player had not read — every one
  of the five endings was destroyed at the moment it was printed, and nothing
  caught it because the engine returned the text correctly and 208 tests
  asserted on it. `State.Ending()` is the difference the renderer needs:
  `Model.finished` blurs the prompt, drops the status line for
  `ENTER to close the file`, and scrolls to `endedAt` so the ending starts at
  the top of the viewport rather than the bottom. `internal/ui/model_test.go`
  drives `Update` and `View` directly and fails if the model quits on top of
  its own ending.
- **A skipped scene has to cost something, not just not happen.** Act III's
  **Failure handling** lines are punishing defaults, and every one that applies
  is built:
  Rhod smuggling Report 415 in `Tick` because `KADMON_PROOF` is documented
  always obtainable and was gated on `VURR_RESPECT ≥ 1`; the reading following
  the player off the floor on printed slips; inaction at `RECONCILIATION ≥ 90`
  sitting them in the fifth chair without a decision; Bez producing the voice
  slate and Kesh citing the nav plot if the cogitator well was never entered;
  Crase following to the intake; and the ledger re-offered in 3.10. The
  guaranteed ones are in `Tick` rather than on a beat, because a beat needs
  the player to walk back up to somebody and a guarantee cannot. The last two
  are 3.4's: skip Ilaria's cell and she raises it herself on a street in Hive
  Erasmus with a shrine queue eleven metres away, and *does not accept the
  answer* — the truth still costs nothing and still does not land, and she is
  not the machine's mouth at Cantica-Null. Skip Fenn's hold and Vaux hands over
  nine pages taken standing up in a corridor, which is C9b and not C9: the
  Amanuensis cannot be forced by it, because a transcription of a recitation it
  does not hold is not a contradiction. Neither sets `WATCH_*`, so neither
  route can complete the voyage.
- **`ON_THE_FLOOR` and `AT_THE_INTAKE` are `while` flags that scope `Tick`.**
  The slips only follow a player who has left the Correctorum, and the
  inaction ending only fires where the slot is. `Tick` runs everywhere, so
  anything in it that belongs to one place needs a `while` flag to say so.
- **Adding a person to a scene changes what the nouns in it mean.** Crase
  follows the player to the intake, so `her` names her there, and Sazh's
  `refuse her` started firing Crase's refusal instead. Renamed to
  `refuse sazh`. Check the pronouns in an existing scene's act names before
  putting somebody new in it.
- **All five Reconciliation symptoms are built and none of them announces
  itself.** ≥25 removes the ACCUSE option in 2.6 (`hidden_above`); ≥40 has Bez
  raise sub-level nine unprompted, from `Tick`, so there is no topic the player
  failed to ask, and `QUERY SUBLEVEL NINE` answers `NO SUCH OPERATION` in the
  spacing the Vault uses for a shelf that was never filled; ≥55 shortens the
  Act I file in 2.0b and thins the retinue's own personnel files; ≥70 puts an
  initial where the player's given name was and removes eight topics across
  four late scenes. Every one of them is the record being ordinary about
  something the player knows is not.
- **`QUERY` is where the retinue files and your own record live.** There was
  nowhere for the ≥55 and ≥70 symptoms to land, so `query()` gained branches
  for `TALM` / `FENN` / `ODUYA` and for `ROSETTE`. The player has no reason to
  pull their own people's files until something makes them want to check, and
  that is the point. `retinueFile()` maps the names.
- **The player is never named, so the ≥70 symptom is about a field, not a
  word.** The rosette record carries the family name and, below 70, a note
  that the diacritic has been wrong since they were nineteen. Past it, the
  given name is an initial, and the prose says a fifth of all personnel
  records carry one for perfectly ordinary reasons. Do not have anybody
  comment on it.
- **Provenance is a second axis on top of the two stores.** `CITE` already
  refuses anything in the case-file. C9b refuses for a different reason: it is
  in the notebook, in the player's own hand, and it is a transcription of a
  recitation the machine does not hold, so there is no record to contradict and
  nothing is forced. The medium was the Act II lesson; where the thing came
  from is the Act III one, and the Amanuensis offers to take the pages into the
  intake instead, which is the joke.
- **Two canon facts I had wrong and the documents settle.** The Cardinal is
  **Ostrem** Vhail, per the bible §3.3 and ACT_II 2.6, not Augustin. Crase is
  **seventy-one**, per ACT_III 3.10 — forty-four also made her Hadrian Vosk's
  friend of thirty years from the age of fourteen, so the older reading fixes
  an inconsistency that was already there. Check a name and an age against the
  bible before writing a new scene around them.
- **A save that drops a field is the archive's own trick.** `State` carries
  about twenty-five unexported counters that are pure timing — `bodyDeadline`,
  `sinceReq`, `symptomsSeen`, `c17Waited`, `transitTicks`, `driftPartial`. A
  snapshot that lost one would restore a session that is *almost* right: the
  body still leaves, but not when it should. So `save.go` lists every field by
  hand — no reflection, no struct tags deciding what survives — and
  `TestSnapshotCarriesEveryField` fails the moment `State` grows one, so that
  adding a field is a decision rather than an omission. Raise the count in that
  test only after adding the field to `snapshot()` *and* `restore()`.
- **The content is never saved; a session is only ever the player.** Scenes are
  compiled in, so `restore` rebuilds against the loaded `Game` and refuses a
  session whose `version` or scene it does not recognise, rather than
  half-applying it. The archive travels exactly as it stands, corrections
  included — a session that came back with clean case-notes would be a lie.
- **`SAVE` and `LOAD` are free, like anything that only touches your own
  things.** They cost no actions, for the same reason reading your own coat
  does not, and the two lines they print say nothing about what was written.
- **`VERITY` moves on `Clue.Verity`, not on every `NOTE`.** It is
  contradictions the player has recorded *outside* the cogitator, and the
  documents are specific: ACT_I gives four sources — A21, the Lexmechanic
  checksum, the Précis annotated or refused, and cylinder six surviving to act
  close — and ACT_II names exactly six clues whose case-file version the
  archive goes on to alter (B1, B4, B12, B19, and one of B9/B23; B6 is on the
  list and is `StoreNeither`, so it can never carry the mark). A blanket +1 per
  note put a careful finish at twenty-one against a documented ceiling of six,
  which made `VERITY ≥ 4` on Ending A's Cael route and *"below 3, no
  rebuttals"* in 3.12 both no-ops. Marking a new clue `verity` means claiming
  the archive alters it.
- **The counters are asserted in the walks, in bands, from the documents.**
  `TestActOneLandsWhereTheDocumentsSayItDoes` checks the act turn — standing
  +2 or +3 from an opening +1, VERITY 0–4, exposure barely moved — and the
  Ending E and D walks check the two ends: a ledger-keeping finish stays under
  25 and never sees the forty-point symptom, a filing one crosses forty and
  does. That is the difference between a symptom that works when a test sets
  the flag and one a player will actually meet.
- **A walk that never takes an award is not evidence about the economy.** The
  playthroughs proved every path *reachable* and said nothing about where the
  numbers landed, because `actOne()` had never challenged the timestamp or
  asked why not Vaux — the two things Act I gives standing for. Measuring the
  walks is how the `VERITY` inflation was found; keep the awards in them.
- Prose style: second person, short sentences under pressure, Imperial
  vocabulary used casually and never explained. See `STORY_BIBLE.md` §9.
