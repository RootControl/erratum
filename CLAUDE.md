# CLAUDE.md

THE ERRATUM — a playable vertical slice of a terminal investigation game
(Warhammer 40,000, grimdark noir). Go + Bubble Tea. Module
`github.com/RootControl/erratum`.

The repo is half design document, half prototype. `STORY_BIBLE.md` and
`ACT_I.md` / `ACT_II.md` / `ACT_III.md` are the source of truth for content and
vocabulary; `main.go` + `internal/` implement all of Act I (scenes 1.1–1.10)
and every scene of Act II. Act III is not started. When the code and the act documents disagree about a name, a
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
  - Act II: `board`, `safehouse`, `workshop`, `delta`, `theta`,
    `choir_settlement`, `annexe`, `rows`, `cathedral`, `confessional`,
    `crase_office`, `ferrum`, `cylinder_vault`

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
- Prose style: second person, short sentences under pressure, Imperial
  vocabulary used casually and never explained. See `STORY_BIBLE.md` §9.
