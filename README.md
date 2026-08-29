# THE ERRATUM — vertical slice

**All of Act I is playable**, 1.1 through 1.10 — six hours, the rosette, the
locked room, the synthesis, the Vaux interview, the door and the body, nothing
to requisition, forty-one seconds, the Vault requisition and the Act I turn —
and it closes on the Case Précis and opens into **Act II 2.0 through 2.8**.
Go + Bubble Tea.

It exists to prove eleven things are buildable before Act II is written:

1. **Free-form examination.** No highlights, no glints, no prompts. Objects are
   found by naming them, and looking twice is not the same as looking once.
2. **Two evidence stores.** `FILE` puts a clue in the cogitator: indexed,
   actionable, and reachable by the enemy. `NOTE` puts it in your coat: inert,
   useless, and yours. The game never explains the difference.
3. **A flag engine** speaking the same SCREAMING_SNAKE vocabulary as the act
   documents, so design and code don't need a translation layer.
4. **An archive that edits your case-notes** — a real mutable store, not a
   scripted cutscene. `Archive.Correct()` mutates the data and says nothing.
5. **Conversation without a menu.** Topics are named, not listed, exactly like
   objects. What opens Vaux is not authority but accuracy: cite the
   seal-integrity register you actually pulled and he answers a different
   question. Push on a door he has closed and he produces, without noticing,
   six syllables of a dead man's cipher — with no bracket, no highlight and no
   indication that it was the most important thing said in the room.
6. **Scenes, and a clock that runs between them.** The lift to tier four is
   three actions each way. The body does not wait while you are in it — but it
   does wait through the ceremony in 1.2, because the eight minutes belong to
   the body and start when you can first reach him.
7. **A record that looks wrong and isn't.** The Crown's verdict is struck two
   hours before the door was forced. It is a pre-generated watch form, the
   explanation is true, and it is the only false alarm in the game — there so
   the player's first suspicion of the archive is mistaken, and they learn to
   distrust their paranoia before they learn to trust it.
8. **A clue the record would happily take.** The boots have no field on the
   form and the mnemonic was never said by anybody who knew they said it — but
   Bez's illiteracy has *exactly* the right field, pre-filled, nine seconds'
   work, and filing it would end a twenty-six-year career. The machine does not
   refuse. You do.
9. **Evidence that is a person, and a place.** Ask Fenn about the eighty-four
   days inside the Crown's Mechanicus shielding and his recall is perfect. Walk
   him twenty metres to the gallery arch and ask again: he stops mid-sentence,
   says there were no such dates, and apologises. Walk him back in and it
   returns, whole, with no memory of ever having lost it. Nothing comments on
   it. The only instrument that demonstrates the mechanism is a man, and the
   demonstration only works where nothing can record it.
10. **An object the archive cannot reach.** The Choir cuts wax redundancies of
    every transcription: no index, no address, no countersignature, nothing a
    cogitator can amend. Cylinder six holds forty-one seconds of your
    astropath's own voice speaking a language she does not know. Taking it
    gives you the one thing in the sector that will still disagree with the
    record in eleven days, and costs you the trust of the only person who can
    carry a message the enemy cannot read. There is no version where you get
    both.
11. **A person who is wrong and is not lying.** Ask Crase what Vosk was working
   on and she says *"nothing, he'd stopped taking cases"*, easily and
   completely, because she has been reconciled and from inside her head it is
   true. Nothing marks the answer. It is the archive's other invariant, worn by
   somebody the player likes.

## Run

```sh
go mod tidy
go run .
```

For anything scripted, skip the TUI entirely:

```sh
go build -o erratum .
printf 'x board\nx board\nfile A1\n' | ERRATUM_PLAIN=1 ./erratum
```

The game opens in a doorway you are not allowed through: `x writ`, `x writ`,
`note timestamp`, `ask about the form`, `ask about grief`, `ask about sealed`,
`tired` or `nothing`, `go on`.

Then the witnesses assemble: `x crase`, `x vaux`, `ask about his last cycle`,
`ask why not vaux`, `take rosette`, `promise` or `defy`, `go in`.

Then, in the chambers: `x board`, `x board`, `file A1`, `x boots`, `x boots`, `note A3`,
`read room backward`, `query C9-F4-D2`, `note A20`, wait five turns,
`query C9-F4-D2` again.

The shutter names a cell across the Crown gap: `go spire`,
`ask dross about access`, `play cylinder`, `x log`, `x log`, `decline`,
`take cylinder` — then `go back` and `query 0997.4-1188`.

One floor down there is a savant who does not think it has happened yet:
`go fenn`, `ask about vosk`, `ask about the dates`, `go arch`,
`ask about the dates`, `go back`, `ask about the dates`, `note the correction`.

When they take him, follow him down: `go medicae`, `ask bez about the door`,
`read roster`, `x kit`, `x kit`, `show the boots to bez`.

Or go and see Vaux, which costs you the boots if you have not read them yet:
`x register`, `x register`, `go vaux`, `ask vaux about the seal`,
`ask about the caseload`, `ask about being passed over`, `note phrase`,
`show the boots to vaux`.

`ERRATUM_NOALT=1` disables the alternate screen, which is how `smoke.py`
drives it headlessly.

## Test

```sh
go test ./internal/game/
```

Ninety-four tests, one per design guarantee:

| Test | Asserts |
|---|---|
| `TestCarefulPlaythrough` | five signs + the boots → `ROOM_READ_FULL`, `CIPHER_HALF_ROOM`, and A21 provable |
| `TestTrustingPlaythrough` | file everything, write nothing → the edit happens and **cannot** be proven |
| `TestBootsCannotBeFiled` | the filing template offers only distress / disorder / self-neglect |
| `TestBodyTimerClosesTheBoots` | eight minutes, then A3 is gone for good |
| `TestVaultRequiresFiling` | you cannot requisition against a code you only told your coat about |
| `TestAccuracyOpensVaux` | citing the register you pulled gets a different man and a different answer |
| `TestMnemonicCannotBeFiled` | A9 is a sound; the cogitator has no field for one |
| `TestMnemonicIsNotAnnounced` | the string `A9` never reaches the screen when A9 arrives |
| `TestDispositionIsOneOfN` | exactly one of detained / shown-room / unchallenged is true on the way out, including for the player who did nothing |
| `TestTravelCostsTheBody` | the lift is three actions each way and the medicae do not wait |
| `TestClueTravelsBetweenScenes` | A3 is picked up on tier nine and put on a desk on tier four |
| `TestFirstSuspicionIsWrong` | the verdict precedes discovery by two hours and the explanation is true |
| `TestA0LandsWhereverYouPutIt` | the first silent notebook test: file it or note it, the game does not care |
| `TestTheQuestionIsRecordedEitherWay` | *did he seem tired* — and "declined to characterise" is written down too |
| `TestGriefBuysTheLeadASceneEarly` | kindness in 1.1 gets the Choir lead early, and it is never delivered twice |
| `TestBezSaysItTwice` | she puts the sentence in the room twice and is still not "met" |
| `TestActOneOpensInOrder` | 1.1 → 1.2 → 1.3 at 04h08, 04h42, 06h12, with the body untouched |
| `TestVauxHasReadYourStatement` | the easy true thing said in 1.1 is the strongest line in the file in 1.5 |
| `TestThePrecisCannotBePutDown` | four true paragraphs describing a day that did not happen, and no way to defer it |
| `TestThePrecisHasThreeAnswers` | sign, annotate, refuse — and annotate needs a notebook you kept |
| `TestActTwoOpensOnTheFarSide` | eight days later, same room, four files still open |
| `TestIlariaWillNotComeIn` | the room has no shadow; she is describing an instrument reading |
| `TestTheKadmonRecitation` | nine hundred million four hundred and twelve thousand, and then no such world |
| `TestItIsNotTheRoom` | he does it again inside two metres of rockcrete, so it is not the shielding |
| `TestDriftIsEquipmentCondition` | `CONCORDANCE OF RECORD: n/6`, never relabelled |
| `TestIsolationIsATrade` | preserve the instrument by making it almost useless |
| `TestTheAccusationVanishesAtTwentyFive` | past the threshold the option does not exist, and there is no gap |
| `TestTheSmuggledVisitorHasNoForce` | no cordon, no writ, and the game says where the decision was made |
| `TestThePurgeOrderOpensLyle` | the vellum is what makes a man who does not exist askable |
| `TestNoTwoCluesShareAName` | the finder cannot see a collision, so there must not be one |
| `TestEveryGatedFlagHasAWriter` | a gate on a flag nothing sets is a line that silently never happens |
| `TestSarnHasTwoSubjectsAndAThird` | she will not discuss the brood, and has wanted to say the third thing for twenty years |
| `TestThePurgesAreBookkeeping` | eleven years, following tithe reconciliation, not discovery |
| `TestTheDeathCascadeIsDeterministic` | Fenn if brought, else Bez, else Vaux, else Ilaria |
| `TestBezIsCalmAboutItIfSheKnows` | and the calm is worse, because you told her the number |
| `TestWithdrawalIsAlwaysAvailable` | nobody dies, the cult expands, and no clue was behind the assault |
| `TestTheRefuserCanBuyTheMeeting` | two favours, and no one on the wet side of a coolant door |
| `TestTheVictoryBecomesARecordsPurge` | forty thousand relocated, and every part of it correct |
| `TestTheWallsAreTheWrongThickness` | 400mm where the 641 plan says 210, and nobody measures a wall |
| `TestRuseSaysItThreeTimes` | *"It's the fire main"*, twice, and then *"What do I do"* |
| `TestIvritSellsCheaply` | he picks you because you are the first person he can safely be useful to |
| `TestVhailJustifiesWithoutAdmitting` | denial straight to justification, with no admission in between |
| `TestTheConfessionalIsTheTrap` | taken inside the instrument, and nothing happens, because he is not the author |
| `TestTheQuerySlipMakesHimRedundant` | the whole thread is readable off an unmarked requisition in twenty minutes |
| `TestTheTrunkKeepsRunning` | denounce, hold or kill — it runs that evening either way |
| `TestTheGalleryBriefIsRealAndTimed` | the dossier is true, and that is not why it is on the desk |
| `TestFedByTheirOwnNonExistence` | eleven thousand people eat what two accounting systems must write off somewhere |
| `TestTheGateIsATrade` | surrender the sidearm, or force it and pay double for the rest of the act |
| `TestTheCensusCannotBeRecorded` | `store: "neither"` — four hundred mouths, and writing it down stops them |
| `TestThePurgeOrderRemovesRecordOnly` | nobody was killed, nobody was moved, and that was sufficient |
| `TestTheAnnexeCostsAFavour` | he is perfectly friendly about it and does not move |
| `TestTheArchiveAnswersAQuestionNobodyRepeated` | a nil return, confirmed, unasked, by a clerk who does not exist |
| `TestHelpingBuysNothing` | four ways to help, all costly, and not one attaches a clue or a stat |
| `TestTheCalendarOnlyMovesForward` | a scene's date is the earliest it can happen on |
| `TestCaelBuiltItFirst` | the docket is three weeks old and the request field is empty |
| `TestTheShroudIsAGiftWithALine` | drift at a third the rate, and position every eleven hours to a bench |
| `TestThePlayerCanAffordTwo` | three sub-queries offered, two affordable, nothing counting them down |
| `TestPastFourThereIsNothingLeft` | `FENN_INTACT` goes, and the game never says what that cost |
| `TestHabrachHandsOverTheMethod` | the method is published, unrestricted, and has been read four times |
| `TestTheSequenceNeverHadIt` | 1186, 1187, 1189, 1190 — not a hole; somebody would have had to make a hole |
| `TestReconciliationRemovesRatherThanCloses` | past its threshold a subject is absent, not refused |
| `TestTheSpireOpensTwoWays` | A2 or Crase's lead, and with neither there is nothing to go and ask |
| `TestDrossGrantsWhatHeCanSurvive` | a release is a decision; a retrieval is a line in the day-book |
| `TestFortyOneSeconds` | her own voice, off the lathe, in a language she does not know |
| `TestTheDispatchNumberOpensTheVault` | the second door into 1.9, edited exactly as silently |
| `TestTheArrestSheOffers` | detained, cleared, or neither — and the Crown is pleased by one of them |
| `TestTheCylinderIsATrade` | the only uneditable object in Act I, or her trust; never both |
| `TestTheGapsAreNotForAskingAbout` | named directly she shuts, and raises it herself on her own terms |
| `TestTheThreshold` | perfect inside, corrected outside, whole again inside — same man, three answers |
| `TestTheCorrectionIsNotAnnounced` | no bracket, no comment, no dialogue option about what just happened |
| `TestFilingTheCorrectionLosesIt` | filed, A16 becomes *savant fatigue, recommend rest cycle* |
| `TestFennNeedsTheCalendar` | a recurrence needs three points; without A5 he asks which dates |
| `TestFennIsOnePerson` | one man in two rooms, one greeting, one conversation |
| `TestFennIsShieldedOrFielded` | shielded, fielded, or undecided — and undecided is not lost |
| `TestIntrosHappenOnce` | a scene's opening beats do not replay when you walk back in |
| `TestAutopsyFloorIsUnmissable` | A11 and A12 for standing in the room: the act's guaranteed floor |
| `TestTierOneOpensWhenTheyTakeHim` | you cannot autopsy a man who is still on the floor upstairs |
| `TestTheBequestIsNeverFiled` | the cogitator has the right field, pre-filled — the refusal is the player's |
| `TestTheRosterIsTwoAnswers` | read it aloud, or name what you just saw; one of the two, once |
| `TestSheLiesAboutTheCough` | no pressure opens her, and the player is shown the decision to lie |
| `TestOrderingTheExamCosts` | fourteen months, confirmed, and she thanks the medicae |
| `TestBezConfirmsTheBootsWithoutKnowing` | she saw it on a ship eleven years ago and thought he was tired |
| `TestAskingTheWrongPersonSaysSo` | a question put to Bez is never answered by Crase |
| `TestQueryNeedsATerminal` | you cannot requisition from a room with no cogitator in it |
| `TestReadRoomIsAboutThisRoom` | signs are facts about where you found them; carrying them elsewhere reads nothing |
| `TestBareGoNamesTheDoors` | a bare `GO` names the ways out instead of guessing and spending three actions |
| `TestPeopleReadLikeObjects` | the second look at somebody is not the first look again, and it is free |
| `TestThereIsNoRefusing` | the door is the refusal; the parser never says "you can't do that" |
| `TestChoirLeadCannotBeMissed` | the one useful sentence in Act I is part of the ceremony, not the conversation |
| `TestKeshWarningIsThreeAnswers` | promise, defy, or say nothing — and saying nothing sets nothing |
| `TestCraseBelievesTheLie` | she is wrong, she is not lying, and no marker says so |
| `TestWhyNotVauxCrossesScenes` | asked in 1.2, paid off in 1.5, one flag, no translation layer |
| `TestClockJumpsAndTheBodyWaits` | scenes are hours apart; the eight minutes are the body's, not the clock's |

`TestTrustingPlaythrough` also asserts the archive never prints the words
`DELETED` or `ERROR`. That assertion is the thesis of the whole game and it
should stay in the suite forever.

## Layout

```
internal/content/   arrival, antechamber, room, vaux_study, medicae,
                    adjutant, arch, choir (Act I); board, safehouse, workshop,
                    delta, choir_settlement, annexe, rows, cathedral,
                    confessional, crase_office, ferrum (Act II); plus
                    strings.json for shared defaults, and the loader — all
                    authored prose is data, dialogue included
internal/game/      state, flags, scenes, command engine, dialogue, archive
internal/ui/        the cogitator. Deliberately has no vocabulary for
                    "something is wrong": no glitch effects, no colour shifts.
```

## What is not built yet

**Act II scenes 2.0b and 2.9–2.11**: the holding cage, Three Doors, the long
way round, and the act turn. Most of the systems those scenes need are in — scene
dates, counter gates, drift as equipment condition with its own accrual rules,
Reconciliation removing options rather than closing them, per-scene
conversations for a person who appears in both acts. Still to build: the
Erratum as a one-use item, `MOTT_FAVOURS`, and the four remaining
Reconciliation symptoms.

Also: save/load, the player backgrounds (Lexmechanic / Provost / Penitent), the
Reconciliation symptom system, and any content past the Act I turn. The clock
is 40 seconds per timed action; the body leaves eight minutes after you first
reach it and does not negotiate.

Nothing is stubbed. Every flag Act I reads, Act I writes, and the handoff to
Act II is complete: `ROSETTE_TAKEN` and `ACT1_TURN_SEEN` always true, all four
retinue alive and met, A11 and A12 guaranteed, and `ROOM_READ_FULL` set or
unset inside ninety minutes and never revisitable.

## Note on `go.mod`

Built in a sandbox that could reach github.com but not `proxy.golang.org` or
`golang.org`, which needed `replace golang.org/x/... => github.com/golang/...`
directives. Those are stripped from the shipped `go.mod`; `go mod tidy` on a
normal machine will resolve everything directly.
