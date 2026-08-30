# THE ERRATUM — vertical slice

**All of Act I is playable**, 1.1 through 1.10 — six hours, the rosette, the
locked room, the synthesis, the Vaux interview, the door and the body, nothing
to requisition, forty-one seconds, the Vault requisition and the Act I turn —
and it closes on the Case Précis, runs **Act II 2.0 through 2.11** to the act
turn, and runs the whole of **Act III** to all five endings.
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

`SAVE` and `LOAD` suspend a session and take it back up. They cost no actions,
for the same reason reading your own coat costs none, and the file lives at
`./erratum.session` unless `ERRATUM_SAVE` says otherwise. The content is never
written to it — a session is only ever the player, and it carries the archive
exactly as it stands, corrections included.

`ERRATUM_NOALT=1` disables the alternate screen, which is how `smoke.py`
drives it headlessly.

## Test

```sh
go test ./...
```

`internal/game` and `internal/content` have no third-party imports, so
`go test ./internal/game/` runs without `go mod tidy`. `internal/ui` needs the
dependencies, which building needs anyway.

Two hundred and thirty-three tests — 230 in `internal/game`, three in `internal/ui`.
Five of them play the whole game from a cold start to five different endings on
typed commands alone; the rest are one per design guarantee:

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
| `TestApproachSetsSuspicionSilently` | how you asked to see Vaux is worth two points and the counter never names itself |
| `TestLedgerAnswerIsTheHinge` | the only answer to his question is one a thorough player cannot give |
| `TestAllianceNeedsMoreThanAsking` | he refuses the hand without saying what the missing thing was |
| `TestArchiveApproachRaisesTheStandingBar` | playing the politics well is what closes the alliance, not suspicion |
| `TestHighStandingDeletesPersuasion` | two ways to reach him are absent, not greyed, and a line is added to the indictment |
| `TestBlackmailNeverAnnouncesTheCorruption` | the coerced cipher half reads exactly like the given one |
| `TestRefusalReturnsHimAtThree` | three weeks in a transit hostel, one door left, and no gap where the others were |
| `TestStoreHabitIsNeverChosen` | the habit is the first reach for FILE or NOTE, and nothing says it was recorded |
| `TestLongWayOnlyOpensWhenTheHalfIsLost` | Ilaria does not offer the wax while there is still a short way round |
| `TestIlariaPricesItHerself` | she states what it ends, and that there is no clever third option |
| `TestSpendingIlariaCapsHerTrust` | what she spends is not the player's regard; the number follows the event |
| `TestFennAsksWhatToRemember` | he says yes before he can be talked into it, and asks the unanswerable question |
| `TestSpendingFennEndsHim` | eleven hours, six of drift, and a readout that goes to 0/6 without naming itself |
| `TestTheShroudDoesNotSoftenIt` | Cael's gift thirds what happens by accident and nothing that happens on purpose |
| `TestTheThirdRouteWaitsForBothToBeGone` | Mott's elder exists to prevent a soft-lock, not to be a third choice |
| `TestCaelNeverAsksWhatItIs` | one of the murderers hands over the murder weapon and says to mind the edges |
| `TestTheAshIsNeverActuallyMissed` | Bez collected it on the first morning, because Vosk taught her to |
| `TestTassBringsItUnbidden` | the plate is unmissable; the signature on the recovery is not |
| `TestDecodeNeedsTheRoomAndBothHalves` | the sentence is read out of furniture, and half a key is a different object |
| `TestTheSentence` | the tense is *making*, not *made*, and the method is not granted with the message |
| `TestTheCorruptHalfNamesTheWrongMan` | the false reading carries no hedge; what breaks it is counting to nine |
| `TestTheTestimonyCannotBeFiled` | the template is correct, the field is right, and the player closes it |
| `TestTheVaultConfirmsWhatNobodyAsked` | nobody confirms a nil return, and something in the hive does it twice |
| `TestTheErratumIsOneUse` | spending it on a person is silent, and nothing says what it just cost |
| `TestThetaIsNeverOffered` | the one hidden door in the game; a bare `GO` does not list it, for anybody |
| `TestTheArchiveCanReachThePaperAndNotTheBoots` | seven crates on the return, nine in the cage, and no line for footwear |
| `TestClosedCratesReadAsAbsent` | a sealed crate is not in the room; a refusal would say there was something in it |
| `TestSalvagedCluesAreWorthLess` | four days late is not the same as reading it off the floor, and nothing says so |
| `TestFourSignsIsNotFive` | he did not leave a message that could be partly read |
| `TestTheCaseFileComesBackShort` | the file is thin, nothing has been struck, and Ending E closes in silence |
| `TestVaultRefusesEverywhere` | the Vault's null answer was authored in one scene and rendered a marker in the rest |
| `TestEveryDoorFinishesTheAct` | all three of 2.9's doors reach the act turn; two of them used to dead-end |
| `TestOneHalfIsNeverComplete` | a cipher is both halves, in either order, and nothing else |
| `TestNoFreeWayToCompleteTheCipher` | every route to the second half spends somebody, named route by named route |
| `TestThetaRecoversAPlayerWithNothing` | a player who read nothing in Act I can still reach the act turn |
| `TestTheCratesAreNotMentionedToEverybody` | the one signpost 2.0b has is aimed at the one player who needs it |
| `TestThePassageIsPendedNotRefused` | six is more than five, and there is no hand on it to push against |
| `TestTheThreeHullsAreExclusive` | a writ, a theft or a debt; taking one closes the other two |
| `TestCraseHasNeverLiedToYou` | she says it without a hedge and does not defend it |
| `TestReadingTheClauseAloudCosts` | she agrees, and a layer of warmth goes and does not come back |
| `TestSazhsTallyCannotBeFiled` | forty-one favours, no names, and no template that will take it |
| `TestTheCerementLiesAboutThePlayer` | out of relay range the record edits the player instead of the world |
| `TestReadingYourOwnFileCostsInTransit` | reassurance costs exposure out there, and the cost is never shown |
| `TestTheMethodIsAssembledNotAnnounced` | plot, casket and uncoupled vox are the method; nothing says so |
| `TestTheWritRouteGetsACuratedHull` | his effects were forwarded, and the slate came back nineteen clean seconds |
| `TestTheWatchesMustBeVisited` | nobody raises their own subject; looking at somebody is not standing a watch |
| `TestBezAsksTheOnlyWaySheCan` | she cannot read, the game never says the word, and the ask costs her more than the lung |
| `TestSofteningTheBequestIsWhatTheyDo` | two small kindnesses, both good reasons, and it is the thing you have been chasing |
| `TestIlariaAsksToBeToldSheIsMad` | the honest answer is worse than insanity and she does not break |
| `TestHandlingIlariaIsTheSameInstrument` | the kind answer, applied by hand, to one person, for the same reason |
| `TestFennIsFinishingNotGuessing` | he is not prophesying; he is running the rule to the end of the sum |
| `TestFilingTheRecitationIsTheTrap` | the filing says what it will do, in the filing, and nothing warns |
| `TestVauxsApologyBuysTheChannel` | the price of one promise is the only sealed channel out of the sector |
| `TestVoyageFullCountsWhatWasPossible` | a lost retinue is not a skipped scene |
| `TestTheyHailYouCorrectly` | ninety-one years of drilling a challenge into a system where nothing arrives |
| `TestTheOwedIsNeverOpened` | one year and four months of food, sealed, on a world hungry since year forty |
| `TestVurrWillNotCorrectTerra` | she reports a false figure for sixty years rather than falsify, and the cogitator resolves it in nine seconds |
| `TestPityClosesHerAndCitationOpensHer` | sympathy is cut off cold; accurate citation is what earns the Reading Room |
| `TestTheThreeAnswers` | truth, lie or silence, and each one costs a different thing |
| `TestTheAidIsInsufficientAndExclusive` | every option says its own ratio out loud, and there is a shrine to it in a week |
| `TestSheThanksYouForComingAnyway` | do nothing at all and she still comes down to see you off |
| `TestTheUnsentReportsAreTheProof` | nine hundred and twenty-two reports, an immaculate index, and nobody has read one |
| `TestReport728MustBeLookedFor` | the one that is not like the others is a spine until you ask what year he died |
| `TestNobodyLooksUp` | no challenge, no search, and a berth placard with the polish worn round the letters |
| `TestKeshIndexesRatherThanDefends` | she keeps her count on paper because a screen is a rumour |
| `TestC17DeletesItself` | the one clue the archive takes back, with no gap where it was |
| `TestHandCopyingC17IsTheOnlyWayToKeepIt` | four hours of copying shapes while a Magos watches a person write |
| `TestCraseHasBeenRunningTheRuleAllAlong` | she has never lied, and that is worse, and she taught you how |
| `TestTheFiledWords` | he argues for it, well, in the structure he taught you on a Tuesday |
| `TestTheBlindspotMustBeAskedFor` | it never claimed authorship; the exploit is one precise question |
| `TestIlariaStopsAtExpenditure` | she asked to be an instrument rather than a suspect, and here is the invoice |
| `TestStoppingTheReadingHidesThePrice` | relief that lasts nine seconds and is charged for later, unexplained |
| `TestItAnswersNarrowerThanYouMeant` | every line is true, scoped narrower than the question, and PRESS gets one layer |
| `TestPressingTheChairToTheBottom` | three were considered, none retained, and the names are nowhere else |
| `TestTheOmissionHasThreeDoors` | CITE works from the notebook only; the case-file copy has already been amended |
| `TestItsWorstAnswersAreTrue` | "that has not been established" and "I have no information on that classification" |
| `TestTheWordToVauxHoldsTheDoor` | a promise made nine days out holds every ending until it is kept |
| `TestTheFiveDoors` | five endings, five keys, and the game ends on a physical action |
| `TestEndingEHasNoPartialCredit` | ten flags, no partial credit, and the refusal never names what is missing |
| `TestKadmonEndingNeedsAWeightToSpend` | exposure is the currency; a player who stayed clean cannot pay |
| `TestKeshsSentenceIsTheSolution` | a warning, an instruction, and the answer to a door she does not know exists |
| `TestPlaythroughEndingA` … `E` | five walks from `New(g)` to five doors, on typed commands and nothing else |
| `TestDelayedTextResolvesEverywhere` | anything `Tick` prints has to resolve in every scene, not just its own |
| `TestNoActIsNamedForAParserVerb` | an act named `write it down` is eaten by the parser and never reaches a person |
| `TestNoTwoExitsShareAName` | `FindExit` is gate-blind, so a shared name makes the second door unreachable |
| `TestTheChannelHasThreeSources` | Ending C's channel is a writ, a favour or a man, and spending the writ revokes it |
| `TestCaelsFilingReachesTheCaseNotes` | he says he filed it, so it is in the case-notes, so the archive can take it back |
| `TestAnEndingHoldsTheScreen` | the alternate screen used to restore itself on top of every ending |
| `TestQuittingStillQuits` | leaving is a player leaving, and still closes at once |
| `TestRhodSmugglesTheProof` | a nineteen-year-old decides it himself, and is beaten for it off-screen |
| `TestTheReadingFollowsYouOut` | walk out and it arrives under the cabin door, in order, numbered |
| `TestInactionSitsYouDown` | you do not choose the chair; you find you have already sat down |
| `TestTheHullsEvidenceSurfacesAnyway` | skip the cogitator well and Bez has the slate, the station has the plot |
| `TestTheLedgerIsReOffered` | she will not argue about a number you have not seen |
| `TestBezRemembersSubLevelNine` | she raises it unprompted, so there is no question you failed to ask |
| `TestTheRetinueFilesThin` | your own people's files come back shorter, with nothing struck |
| `TestAtSeventyYouGo` | an initial where the given name was, and a fifth of all records carry one |
| `TestAtSeventyOptionsAreAbsent` | removed, not closed: she answers as though she had never heard of it |
| `TestIlariaFiresOneSidedAtKadmon` | skip her cell and she asks on a street, and does not accept the answer |
| `TestVauxReconstructsItUncitably` | nine pages taken standing up in a corridor cannot force anything |
| `TestNoDateAfterThePresent` | the present is 826.M41 and an impossible date reads as a clue, not a typo |
| `TestSnapshotCarriesEveryField` | a save that quietly loses a timer is the archive's trick, done by the engine |
| `TestASuspendedSessionFinishesTheSameWay` | put it down mid-Act-II, take it up, and still reach the hidden door |
| `TestActOneLandsWhereTheDocumentsSayItDoes` | +2 or +3 with Crase, nought to four of proof, and the archive barely moved |
| `TestOnlyTheAlteredCluesAreProof` | VERITY is the six clues the archive alters, not a tally of everything noted |
| `TestTheArchiveTravelsAsItStands` | a session that came back with clean case-notes would be a lie |

`TestTrustingPlaythrough` also asserts the archive never prints the words
`DELETED` or `ERROR`. That assertion is the thesis of the whole game and it
should stay in the suite forever.

## Layout

```
internal/content/   arrival, antechamber, room, vaux_study, medicae,
                    adjutant, arch, choir (Act I); board, safehouse, workshop,
                    delta, choir_settlement, annexe, rows, cathedral,
                    theta, confessional, crase_office, ferrum,
                    cylinder_vault (Act II); skydock, cerement, spirit,
                    gunbay, cell, hold, wardroom, kadmon, tithe_yard,
                    muster, reading_room, cantica, fifth_chair,
                    correctorum, intake (Act III); plus
                    strings.json for shared defaults, and the loader — all
                    authored prose is data, dialogue included
internal/game/      state, flags, scenes, command engine, dialogue, archive
internal/ui/        the cogitator. Deliberately has no vocabulary for
                    "something is wrong": no glitch effects, no colour shifts.
```

## What is not built yet

Nothing in the three act documents, and nothing in `STORY_BIBLE.md` §12. All
thirty-one scenes, all five endings, every Reconciliation symptom and every
applicable punishing default are in, and the suite plays the game from a cold
start to five different doors.

Absent, and never specified anywhere: the **player backgrounds** (Lexmechanic /
Provost / Penitent), which the bible offers as "mechanical flavour, same hook".
The bible also says the player's name is player-set with a default of Ammon
Reyl; the game is written in second person and never names them, which is a
deliberate reading and not the bible's.

Nothing is stubbed. Every flag the content gates on is written by something —
`TestEveryGatedFlagHasAWriter` enforces it — and the two engine-derived flags
that used to be forgotten by content, `CIPHER_COMPLETE` and
`VOSK_METHOD_UNDERSTOOD`, are computed rather than declared.

The clock is 40 seconds per timed action; the body leaves eight minutes after
you first reach it and does not negotiate.

## Note on `go.mod`

Built in a sandbox that could reach github.com but not `proxy.golang.org` or
`golang.org`, which needed `replace golang.org/x/... => github.com/golang/...`
directives. Those are stripped from the shipped `go.mod`; `go mod tidy` on a
normal machine will resolve everything directly.
