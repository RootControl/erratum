# ACT II — THE UNRECORDED

**Dramatic question:** *Who can rewrite the record?* Four open threads let the player exhaust every human answer — a corrupt cardinal, a prefect's clerks, a xenos cult, a rival with a grievance — and find each one too small to have done it. The act must let the player be satisfyingly, expensively right about the wrong man at least once. It closes with what Vosk left them: a routing address the archive denies, and a sentence that kills whoever reads it. **Answer:** *Not who. What.*

---

## SCENE INDEX

| # | Scene | Location | Type | Gated by | Est. min |
|---|---|---|---|---|---|
| 2.0 | THE STANDING BOARD | Vosk's chambers, Adamant Crown | Hub | Act I complete | 10 |
| 2.0b | THE HOLDING CAGE *(optional, hidden)* | Lexicanum Vault, Stack Theta | Ending E recovery | ROOM_READ_PARTIAL | 20 |
| 2.1 | THE INVENTORY OF ALETH FENN | Safehouse, Crown sub-level 4 | Systems / dread | — | 20 |
| 2.2 | THE SHIELDED ROOM | Cael's workshop | Negotiation | B1 | 15 |
| 2.3 | THE DELTA RUN | Lexicanum Vault, Stack Delta | Dilemma | B1 + Act I shelf-code | 30 |
| 2.4 | THE TOLL OF NAMES | Rust Choir: gate, cistern, annexe | Faction / evidence | — | 40 |
| 2.5 | THE SPILL | Rust Choir, the Rows | Moral (unrewarded) | B4 | 15 |
| 2.6 | THE UNBLINKING EYE | Cathedral: crypt, confessional | Interrogation | — | 45 |
| 2.7 | THE GALLERY BRIEF | Crase's office, then ambient | Fork | Two threads touched | 20 |
| 2.8 | THE FERRUM DESCENT | Ferrum Gallery, coolant sub-galleries | Loss / interrogation | DECOY_PURSUED (or Mott) | 60 |
| 2.9 | THREE DOORS | Choir Spire, cylinder vault | Confrontation / branch | Two threads complete | 40 |
| 2.10 | THE LONG WAY ROUND | Choir Spire *or* Stack Delta | Hard route | !CIPHER_HALF_VAUX | 30 |
| 2.11 | THE ASH AND THE SENTENCE | Cael's workshop; Vosk's chambers | Act turn | ROOM_READ_* + CIPHER_COMPLETE | 35 |

---

### 2.0b THE HOLDING CAGE *(optional — Ending E's late window)*
**Location** — Lexicanum Vault, Stack Theta, effects holding.
**Present** — Under-Clerk Habrach; a wire cage; nine crates.
**Player goal** — Read the signs you did not read while the body was still warm.
**Entry gate** — `ROOM_READ_PARTIAL` only. Invisible to a player who already has `ROOM_READ_FULL`; never signposted to anyone.
**Beats** —
1. The Crown reclaimed Vosk's chambers at Act I close. His effects were crated, manifested and shelved. The player can requisition the crates — an ordinary, boring, entirely legitimate act that nobody obstructs.
2. **The manifest has been amended.** It lists seven crates of nine. The regicide board is on it as *"gaming set, incomplete."* The boots are not on it at all.
3. The objects are still in the cage. The archive amended the *list*; it has never once been able to touch a boot. This is the clearest statement the game ever makes about what the enemy can and cannot do, and it is made by a discrepancy in a stock manifest.
4. Reading the signs off the originals restores whatever was missed — but the crates were packed by someone who did not know the arrangement mattered, so the *positions* are gone. The player must reconstruct them from the pict-captures in their own Act I case-file, which means the reconstruction is only as good as what they filed, and what they filed is editable.
5. Cost: a Kesh countersign to open the cage (`RECONCILIATION` +8), four days, and Crase learns the player went back for a dead man's boots.
**Clues gained** — restores any unclaimed A-clue from 1.4 at one `VERITY` less than its Act I value.
**Branches** —
| Player choice | Flag set | Near | Far |
|---|---|---|---|
| Reconstruct all five signs | `ROOM_READ_FULL` | Four days gone; `CRASE_STANDING` −1 | **Ending E reopens** |
| Take the boots and leave | `VOSK_BOOTS` | Nothing | Ending D epilogue: Bez keeps them |
| Never come here | — | — | Four endings, silently |
**Sample text** —
> Seven crates, says the manifest. You count nine.
>
> You count them again, because the number is not interesting enough to be wrong. Nine crates, seven listed, and the two that don't exist are sitting in the cage in front of you with his name stencilled on them in somebody's careful hand.
>
> Whatever has been editing this hive can reach the paper. It cannot reach the boots. You stand in the cold of Stack Theta for a while, holding that, turning it over, finding it a great deal more frightening than the alternative.
**Failure handling** — Purely optional. No hint, no marker, no dialogue prompt. A player who never returns to the chambers' effects loses Ending E permanently and is never told.
**Flags set** — `ROOM_READ_FULL` (conditional), `VOSK_BOOTS`, `RECONCILIATION` +8.

---

### DESIGNER NOTE — CRASE SAVES THE PLAYER'S LIFE
The bible requires one literal life-save from Crase, and Act I spends its Crase capital on the requisition-embarrassment. It lands here, in whichever branch the player took:

- **`DECOY_PURSUED`** — in 2.8, in the coolant sub-galleries, Crase is there before the player is, having read the same brief and moved faster. She puts a round through a purestrain that the player never saw and does not mention it afterwards. If the player raises it she is embarrassed. `CRASE_STANDING` +2.
- **`DECOY_REFUSED`** — the sump, returning from 2.5. Two men with a length of pipe and no records, which means no motive the Arbites will ever reconstruct. Crase's people take them off the player forty seconds in. She was having the player followed. She says so, without apology, because she believes protection is what she owes them. `CRASE_STANDING` +2, `SURVEILLED` confirmed.

Both readings must survive Act III: she saved them because she is fond of them, **and** because an unedited Inquisitor with a rosette is the most valuable object in the sector and she was not going to lose it to a pipe.

---

## CROSS-THREAD INTERFERENCE

**I-1 — Vhail broken early sterilises the Vault.** `VHAIL_BROKEN` before any Vault action briefs Kesh inside forty hours. Habrach and the shelf clerks go from indiscreet to courteous and useless; every requisition is honoured, on time, clean. B23 becomes unobtainable and must instead be handed over by Kesh in Act III, pre-edited, costing 1 `VERITY`.

**I-2 — The Decoy eats the sump.** `DECOY_PURSUED` advances the clock fourteen days; the flush order is dated day 40. Return after it and 2.5 becomes its aftermath: the Rows sluiced, `MOTT_FAVOURS` halved, Wick not among the survivors, Mott still cheerfully doing business, which is worse.

**I-3 — Fenn is the Cathedral's lie-detector.** Vhail's alibi turns on dates. With `FENN_INTACT` he recites the alms ledger from un-networked memory, catches a two-year discrepancy, and one lever opens Vhail. At `FENN_DRIFT ≥ 3` he agrees with Vhail, warmly and in detail, and the player needs two.

**I-4 — The unrecorded courier changes Vaux.** Any approach routed through the archive reaches 2.9 at `VAUX_SUSPICION` 2. A Rust Choir runner (1 favour) reaches it at 0, and unlocks alliance at `CRASE_STANDING` ≤ +2 instead of ≤ 0.

**I-5 — Crase's approval is Vaux's evidence.** Each point of `CRASE_STANDING` above +2 deletes one persuasion option in 2.9 and adds a line to his indictment. The player is punished here for playing the politics well.

**I-6 — Surveillance closes the sump.** `CRASE_SURVEILLED` makes the Sluice Gate refuse the player on sight; re-entry costs 1 favour per visit to lose the tail, and Mott says so aloud, in front of the retinue, every time.

---

## SCENES

### 2.0 THE STANDING BOARD
**Location** — Vosk's chambers, unsealed.
**Present** — Player; retinue; the cogitator.
**Player goal** — Turn Act I's contradiction into a plan of work.
**Entry gate** — Act I complete.
**Beats** —
1. The cogitator boots with four of Vosk's case-files open: FENN/CUSTODY, SUMP-VARIANCE, ECCL/ALMS, GALLERY-9. He never left a file open. These have been open eleven weeks.
2. The two stores are made explicit: **case-file** (searchable, editable by the enemy) and **ledger** (handwritten, unsearchable, never wrong).
3. First real filing choice: the Act I dispatch clue goes to one or the other.
4. Threads unlock as destinations. No recommended order, no markers, no arrows.

**Clues gained** —
| ID | Evidence | Store | Unlocks |
|---|---|---|---|
| B0 | Vosk's four dormant case-files | Case-file | All four threads; proof he chose them |

**Interrogation / dialogue** — Ilaria will not cross the threshold. The air is wrong in here, she says, and is not being poetic: Vosk stripped the room's vox, so it has no shadow.

**Branches** —
| Player choice | Flag set | Near | Far |
|---|---|---|---|
| File to cogitator | `STORE_HABIT_FILED` | +1 query tier | Act III: three clues arrive pre-edited |
| Keep a ledger | `STORE_HABIT_LEDGER` | Slower, no search | `VERITY` accrues; Act III resistance |

**Sample text** —
> Four files, still open. He never left a file open — he closed them the way other men lock doors, at the end of a shift, finished or not, because an open file is a thing someone else can read.
>
> These have been open eleven weeks. Long enough that the cogitator has stopped asking whether you want to save them.
>
> He was pointing, for eleven weeks, and nobody in this hive was looking at his desk.

**Failure handling** — Unfailable; closed files reopen on next boot, without comment.
**Flags set** — `ACT2_OPEN`, store-habit flag.

---

### 2.1 THE INVENTORY OF ALETH FENN
**Location** — Safehouse, Crown sub-level 4. Two metres of rockcrete, no vox-plate. Vosk chose it.
**Present** — Fenn; player; optionally Bez.
**Player goal** — Learn what Fenn holds and what using him costs.
**Entry gate** — None.
**Beats** —
1. Fenn recites on demand, tirelessly, badly organised, delighted to be asked. He does not know what matters. That is the point.
2. Bez's vox-bead is live by default. Unless the player kills it, Fenn sits in weak relay coverage, and nothing appears to happen.
3. Forty minutes on, second query: **he corrects himself mid-sentence.** He is not distressed; he believes he made an error and cannot be shown otherwise.
4. `FENN_DRIFT` (0–6) appears as an equipment readout labelled `CONCORDANCE OF RECORD: 6/6`. Never relabel it.

**Clues gained** —
| ID | Evidence | Store | Unlocks |
|---|---|---|---|
| B1 | The Kadmon recitation and its retraction (record terminates 735 — ninety-one years ago) | Ledger if hand-written; case-file if dictated | Kadmon as a name; 2.2; drift |
| B2 | Drift-vs-coverage correlation | Ledger | Localises the mechanism to the relay |

*B1 dictated to the cogitator loses the population figure within six hours of play. Do not warn the player.*

**Interrogation / dialogue** — ship both strings from one node, the second overwriting the first in scrollback:

> **FENN (first pass):** "Kadmon Veil. Vardan subsector, coreward chain. Tithe grade Exactis Prima, assessed nine hundred million four hundred and twelve thousand souls, tithe paid in full 731, paid in full 732, paid in full 733, paid in full 73—"
>
> **FENN (correcting, same breath, cheerful):** "—there is no Kadmon Veil. The Vardan subsector holds four worlds and I have named all four. I'm sorry. Nine hundred and eleven million is a population figure but it isn't attached to anything. Sometimes a number gets loose. Shall I go again from the top?"

**Branches** —
| Player choice | Flag set | Near | Far |
|---|---|---|---|
| Kill all vox here | `FENN_ISOLATED` | Drift frozen; consultable only here, slowly | `FENN_INTACT` preserved |
| Keep him mobile | — | +1 drift per field scene | Opens 2.2, 2.3 |
| Forbid the word Kadmon | `FENN_SUPPRESSED` | He complies, then forgets he was told | +1 drift at once |

**Sample text** —
> He says nine hundred million, four hundred and twelve thousand the way a man reads a shipping manifest — a number too exact to be an estimate and too calm to be a lie. Then his jaw does something — nothing, really, a hitch you'd miss if you weren't watching his mouth — and he says there is no such world, and asks whether you'd like him to start again.
>
> You would not.
>
> You write the number down. Paper, not the cogitator. Your hand isn't steady and it doesn't matter, because handwriting is the only script in this hive that nothing else can read.

**Failure handling** — If no second query is run, Fenn contradicts himself unprompted at the first field scene under coverage. Only the clue's cleanliness can be lost.
**Flags set** — `FENN_INTACT`, `FENN_DRIFT` = 0, `KADMON_HEARD`.

---

### 2.2 THE SHIELDED ROOM
**Location** — Magos Yttrix Cael's workshop. Warm, tidy, smells of oil and tea.
**Present** — Cael; Enginseer-Adjunct Orven Tass; Fenn if brought.
**Player goal** — Make Fenn fieldable.
**Entry gate** — B1.
**Beats** —
1. Cael has had the idea already: a half-built mnemonic shroud on the bench, dated three weeks ago. No explanation offered; the player has no reason yet to want one.
2. It works — drift at one-third rate. A real gift, freely given.
3. The addendum: calibration telemetry, explained fondly and at length by Tass. It is a channel and it goes to Cael.
4. He also upgrades the player's query tier, unasked. Genuinely useful. Also a window.

**Clues gained** —
| ID | Evidence | Store | Unlocks |
|---|---|---|---|
| B3 | Shroud telemetry addendum | Ledger | `CAEL_SUSPECTED`; barter item in 2.11 |

**Interrogation / dialogue** — Cael explains mechanism honestly and exhaustively, and cannot be drawn on why he built a shroud against relay proximity before being asked. Pressed twice: *"Because a savant is expensive and Cantor is loud."* True. Insufficient. A reconciled man running on the residue of a reason.

**Branches** —
| Player choice | Flag set | Near | Far |
|---|---|---|---|
| Accept | `FENN_SHIELDED` | Drift ÷3 | Cael sees everywhere Fenn goes |
| Refuse | — | Fenn stays home or burns fast | Ending E route stays clean |
| Accept, cut telemetry | `FENN_SHIELDED`, `CAEL_DECEIVED` | Same benefit, no window | He never mentions it, ever, which the player remembers in Act III |

**Sample text** —
> "It wants recalibrating every eleven hours or the field wanders." He turns the harness so you can see the seating. "Tass will handle it. He'll want a telemetry line — nothing exotic. Position and field-strength, twenty bytes a burst."
>
> Position.
>
> "It's very small," he says, and he is right, and he pours the tea before you can answer, and it is good tea, and you have not eaten in a day and a half.

**Failure handling** — Skippable; Thread A stays playable but slower, and 2.3 becomes a straight sacrifice rather than a managed one.
**Flags set** — `FENN_SHIELDED`, `CAEL_SUSPECTED`, `QUERY_TIER` +1.

---

### 2.3 THE DELTA RUN
**Location** — Lexicanum Vault, Stack Delta. Repeaters every eighty metres; no vox-shadow anywhere.
**Present** — Fenn (required); Clerk-Ordinate Selm Habrach.
**Player goal** — Cross-reference the Act I shelf-code against Fenn's copy before the Vault's copy is corrected.
**Entry gate** — B1 + the regicide-board shelf-code.
**Beats** —
1. The Vault's copy is correct and boring. Fenn's is not.
2. Three sub-queries are offered; each costs a drift tick (three under shroud); each feeds a different thread. **The player can afford two.** One thread goes unhelped and the player picks which.
3. At drift 3 Fenn begins correcting the *player's* statements; at 5 he is polite about it; at 6 he is a pleasant man with a good memory for opera.
4. Habrach hovers, alarmed and indiscreet, and can be worked for B23 unless I-1 has fired.

**Clues gained** —
| ID | Evidence | Store | Unlocks |
|---|---|---|---|
| B23 | Kesh's variance-reduction memo | Ledger | `KESH_IDENTIFIED`; Act III |
| B7 | True Cathedral alms dates | Ledger | Lever A in 2.6 |
| B6a | Sump population variance, 19 yrs | Ledger | Corroborates B4 |
| B18 | Vosk's outgoing dispatch headers | Ledger | Fenn route in 2.10 |

**Interrogation / dialogue** — Habrach: *"The Lord-Prefect doesn't withhold, sir. Nothing is withheld. Things are reconciled."* He says it as a comfort. He has been saying it to himself for thirty years.

**Branches** —
| Player choice | Flag set | Near | Far |
|---|---|---|---|
| Two sub-queries | — | Drift 2–4; one thread unhelped | `FENN_INTACT` survives |
| Three | `FENN_BURNED_PARTIAL` | All threads helped | Cannot be the fallback in 2.10 |
| Past four | `FENN_INTACT` = false | Everything, now | Forecloses Ending E; 2.10 loses a route |

**Sample text** —
> "Alms disbursement, Cathedral of the Unblinking Eye, year 811," Fenn says, and stops, and looks at the ceiling like a man checking for rain.
>
> The repeater above the aisle clicks twice. It has done that every eighty seconds for two hundred years.
>
> "Year 811," he says again, warmly. "A good year. Nothing unusual in it."
>
> There were nine thousand tonnes unaccounted in 811. You know because he told you forty seconds ago, in a room with no vox in it, and you wrote it down.

**Failure handling** — A blanked Fenn costs 2.6 a second lever and forces the Ilaria route at 2.10. No soft-lock.
**Flags set** — `FENN_DRIFT`, sub-query flags.

---

### 2.4 THE TOLL OF NAMES
**Location** — The Rust Choir: condemned outflow pipe, recitation cistern, dry annexe.
**Present** — Mott Nine-Fingers; gate-warden Sesh; the name-bank elder; Bez (recommended; she hates it).
**Player goal** — Enter a settlement with every reason to kill you; learn who signs an erasure.
**Entry gate** — None; the annexe needs `MOTT_FAVOURS` ≥ 1.
**Beats** —
1. **The society, walked through rather than explained.** No entitlement, so no rations: the Choir eats *spill* — condemned bulk from the Gallery commissaries and Cathedral alms surplus that must physically be written off somewhere. They are fed by the paperwork of their own non-existence.
2. **The name-bank.** The erased trade in the recently dead who are still on file; a corpse with a live ration card is the most valuable object down here, and bodies are kept breathing on paper for decades. One elder's whole office is remembering which dead man each living man currently is.
3. **The Choir.** At shift-end four hundred people recite names and the year each was struck. Stored in mouths, because mouths are the one medium the enemy cannot reach.
4. Mott examines the rosette, hands it back, explains it means nothing here — a courtesy, not a threat. He cannot be prosecuted, surveilled, conscripted or protected. `MOTT_FAVOURS` moves one way only: coin needs a ledger, favours don't. He never collects in Act II.
5. **The purge order.** Vellum, nineteen years old, stolen from Vault intake the day he was struck, because he was the clerk meant to file it. It authorises removal of *record* for Sump District Nine. Not of people. The people were left where they were.
6. Signed **Ordinate Casmir Lyle**, seat and serial given. The player queries it: no Casmir Lyle, no such seat, ever. **The query is logged**, and forty minutes later an unrequested courtesy note confirms the null result — from a clerk who does not exist either.
7. A second document: a forward-dated flush order for the Rows, day 40, same seat, same hand.

**Clues gained** —
| ID | Evidence | Store | Unlocks |
|---|---|---|---|
| B6 | The oral census — 3,114 names | **Neither** — memory only | Act III leverage; unusable as evidence |
| B4 | Purge order, sig. Casmir Lyle | Ledger (vellum) | 2.5; corroborates Fenn |
| B5 | Null result + unrequested confirmation | Case-file | The archive answering a question nobody asked |
| B4b | Flush order, day 40 | Ledger | The clock |

*Try to transcribe B6 and the Choir stops singing.*

**Interrogation / dialogue** — Mott trades for favours and nothing else: not coin, rank, or threat. Asked what he was, he answers once, flat — *"Adept third class, Vault intake, nineteen years"* — and elaborates only above 2 favours. On the order: *"It's the only paper in the sector with my name on it. Not the front — the back. Intake clerk's mark, third from the bottom. That's the whole of me now: three letters on the back of the thing that killed me."*

**Branches** —
| Player choice | Flag set | Near | Far |
|---|---|---|---|
| Surrender sidearm at the gate | `CHOIR_TRUSTED` | Full access | `ILARIA_TRUST` +1 |
| Force entry with the rosette | `CHOIR_HOSTILE` | Access, contempt | Favour costs doubled all act |
| File B4 to the cogitator | `B4_FILED` | Searchable | Signature reads *[UNSIGNED]* by 2.9; the vellum is unchanged. Never comment |

**Sample text** —
> Four hundred people in a cistern, and not one of them singing a hymn.
>
> "Ovric Hald, struck 807. Marra Hald, struck 807. Boy Hald, no given name, struck 807." A pause the width of a breath, and four hundred voices take the next one.
>
> It runs forty minutes. You do the arithmetic without meaning to: at four seconds a name, this is a settlement that spends most of an hour a day being its own archive.

**Failure handling** — Hostile entry is still entry. Earn no favours and Mott offers B4 unprompted on a third visit, for two. The clue is not optional; only its price is.
**Flags set** — `MOTT_MET`, `MOTT_FAVOURS`, `CHOIR_TRUSTED`/`CHOIR_HOSTILE`, `PURGE_ORDER_SEEN`, `LYLE_QUERIED`, `CLOCK_DAY40_ARMED`.

---

### 2.5 THE SPILL
**Location** — The Rows. Nine hundred people in a drainage terrace.
**Present** — Wick (a girl of about nine; never had a name in any file); Mott.
**Player goal** — None. No objective marker. This is where the game finds out what the player is.
**Entry gate** — B4.
**Beats** —
1. The Rows flush on day 40. Nobody disbelieves it; they watched it happen three sumps over.
2. Four ways to help, all real, all costly, none rewarded: **defer the flush** (2 `CRASE_STANDING`; Crase does it instantly without asking why, which feels wonderful now and like a receipt later; buys 60 days); **divert alms surplus** (needs B7 and a Cathedral lever — spent here, it is not available in 2.6); **empty Vosk's discretionary account** (Act III's writ becomes `TRANSPORT_DEBT`); **forge death-cards** for 200 children into the name-bank (3 favours, and Bez walks them past a checkpoint, costing her a week of lung she does not have).
3. Wick asks the player's name, repeats it until she has it, and says she'll put it in the Choir. Not as thanks. As a filing.
4. **Load-bearing developer note: attach no clue, unlock, stat or asset to any of this.** `CHOIR_FED` and `CHOIR_DEFERRED` are read once each, as one sentence each, in the Act III epilogue. Bolt a reward on in revision and the act loses its spine.

**Clues gained** — None. Deliberately none.

**Interrogation / dialogue** — Mott, if the player helps: *"That was stupid."* Beat. *"Do it again in forty days, will you."* He never thanks them; gratitude is a record of a transaction and they don't keep those.

**Branches** —
| Player choice | Flag set | Near | Far |
|---|---|---|---|
| Defer | `CHOIR_DEFERRED` | −2 Crase | One epilogue line |
| Feed them | `CHOIR_FED` | Resource loss per method | One epilogue line; Bez's cough worsens |
| Walk away | — | Nothing | Day-40 aftermath; Wick not among the survivors |

**Sample text** —
> "What's your name," she says. Not a question — the flat delivery of someone taking down a shipping mark.
>
> You tell her. She says it back, wrong the first time, right the second, then a third time with her eyes shut, and you understand that you have just been filed.
>
> "You're in the Choir now. That's forever. That's the good one." Then she goes back to sorting spoiled protein by smell, because that's her shift, and she's nine, and she's on time for it.

**Failure handling** — No failure state. The function is entirely moral.
**Flags set** — `CHOIR_FED`, `CHOIR_DEFERRED`, `WICK_KNOWN`.

---

### 2.6 THE UNBLINKING EYE
**Location** — Cathedral of the Unblinking Eye: nave, junction crypt, Vhail's private confessional.
**Present** — Sister Oath-Warden Bettine Ruse; Deacon Pol Ivrit; Cardinal Ostrem Vhail; optionally Fenn (I-3).
**Player goal** — Find where confession goes, then break the man who sells it.
**Entry gate** — None; the confessional needs B8 plus one lever.
**Beats** —
1. Fifty-one confessionals. Check wall thickness against the filed floorplan: 400mm where the 641 plan says 210.
2. The junction crypt — fifty-one lines into one trunk, and the trunk is Ferrum Gallery vox-manufactorum stock, a hundred and forty years newer than the Cathedral.
3. Ruse is not corrupt; she was told it is a fire-suppression main and has believed it thirty-one years. Ivrit is corrupt, cheaply, and flips for the shipping code on one threat or one favour.
4. Vhail denies expansively for ten minutes; he is good at it. The lever lands and he moves to justification without passing through admission.
5. **The trap.** The confession is recorded in his own live confessional, feeding the trunk. The evidence that destroys him reaches his handlers in real time and they do not care, because he is not the author.
6. He cannot name the Correctorum or Cantica-Null, does not know where the data goes, has never had authority to shut the trunk off. **Everything he gives is readable off a Vault requisition slip in twenty minutes — ship that query, unmarked, available all act, so a diligent player can find that out.**

**Clues gained** —
| ID | Evidence | Store | Unlocks |
|---|---|---|---|
| B8 | Junction schematic | Ledger | The confessional; `VHAIL_EXPOSED` |
| B11 | Tithe relief, 0.4% under assessed since 786 | Case-file | Motive; Lever B |
| B10 | Feed prefix `CN-∅ / PRIORITY-BLACK` | Case-file | Corroborates B20 — **redundant by design** |
| B9 | Vhail's confession | Case-file (live-fed) | `VHAIL_BROKEN`; standing; almost nothing else |

**Interrogation / dialogue** — Ruse: *"It's the fire main."* Shown the schematic: *"It's the fire main."* Third time, after a pause: *"What do I do."* The game offers no good answer. Vhail's confession, ship verbatim:

> **VHAIL:** "Sins are not secrets. That is the thing your kind never grasp. A secret is a possession. A sin is a *symptom*. You do not protect a symptom, you report it, and the physician does what physicians do.
>
> "They came to me in 786 — clerks, not lords, clerks with a warrant and a discount — and they said the Ordo needs the shape of this hive's guilt. Not names. The shape. Heresy is a disease, disease has a curve, and they wanted the curve.
>
> "And then it was names, because a curve is made of names, and I gave those too. Forty years of them.
>
> "Now ask me what it bought. A heresy-incidence rate of zero point zero zero one. Across a *sector*. In 786 we burned two hundred a month in the Ferrum yards. Last year, nine. Nine, in a hive of two hundred million.
>
> "I have not lost one soul I could have saved. I gave up the ones I could not, so the physician could find them faster, and if that is a crime then the Emperor keep me — but do not stand in my own confessional and tell me I did not know what I was doing. I knew exactly. I was the only man in this hive who was *counting*."

Pressed on where the data goes he becomes honestly bewildered — the first true thing he does — and says *"To the Ordo. To your people."* He believes it. He has never had cause not to.

**Branches** —
| Player choice | Flag set | Near | Far |
|---|---|---|---|
| Warn Ruse | `RUSE_WARNED` | She resigns her wardenship | Act III: she smuggles one message free |
| Break and denounce | `VHAIL_DENOUNCED` | Ecclesiarchy standing +2; I-1 fires | Replaced in a month by a man who supplies the same data quietly |
| Break and hold | `VHAIL_LEVERAGED` | Cathedral access; alms diversion for 2.5 | He tries to have the player killed in Act III, badly |
| Kill him | `VHAIL_DEAD` | Satisfying | +12 `RECONCILIATION`; the trunk keeps running |

**Sample text** —
> He weeps for eleven minutes. It is not performance — you have seen performance and this is not it, and that is the worst of it, because a performance you could file.
>
> When he's done he asks whether the nine is a real number. You tell him it is. He says *nine*, twice, like a man setting down something heavy, and you understand that you have broken a genuinely wicked man wide open and learned a shipping prefix you could have read off a requisition slip.

**Failure handling** — Miss the crypt and Ivrit sells B8 for a bribe on a second visit. Without a lever Vhail stonewalls permanently; nothing required is gated behind him.
**Flags set** — `VHAIL_EXPOSED`, `CATHEDRAL_MAPPED`, `VHAIL_BROKEN` + sub-state, `RECONCILIATION` +4.

---

### 2.7 THE GALLERY BRIEF
**Location** — Crase's office; then, on refusal, ambient across three subsequent scenes.
**Present** — Crase. Then nobody, which is the point.
**Player goal** — Decide about the Ferrum Gallery.
**Entry gate** — Two threads touched.
**Beats** —
1. Crase has been *useful* for hours — a writ fixed, a bad informant flagged, and if the Sluice Gate went hostile, the player pulled out of an Arbites hold at 03:00.
2. She lays out the cult: fourth-generation infestation in Shift Nine, into three shift-councils and the coolant guild. Her evidence is real, thorough and correct. She asks; she does not order.
3. **On refusal:** −2 `CRASE_STANDING`, no dialogue attached. Requisitions run 40% slower. Two of Vosk's contacts stop answering. The Crown lift log shows the player's rosette used at hours they did not use it. A mono-tasked shelving servitor follows them one whole scene at exactly forty metres, which it is not supposed to be able to do.
4. Crase stays warm throughout — she sends a note asking after Bez's cough. She is worried about Bez *and* having the player followed, and both are true.

**Clues gained** —
| ID | Evidence | Store | Unlocks |
|---|---|---|---|
| B13a | Crase's Gallery dossier | Case-file | 2.8 |

**Interrogation / dialogue** — Confronted about the tail, Crase admits it gently, immediately, and does not stop: *"You'd do the same. You will do the same."*

**Branches** —
| Player choice | Flag set | Near | Far |
|---|---|---|---|
| Accept | `DECOY_PURSUED` | Clock +14 days; I-2 fires | +3 Crase; forecloses Ending E |
| Refuse | `DECOY_REFUSED`, `CRASE_SURVEILLED` | Surveillance begins; I-6 fires | −2 Crase; Ending E open |
| Delegate to Vaux | `DECOY_DELEGATED` | Requires `VAUX_ALLIED` | He is unavailable in 2.10 and comes back changed |

**Sample text** —
> The lift log says your rosette went to sub-level nine at 04:12. You were asleep at 04:12, in a chair, in a room with the door bolted, and you know that because you woke with the bolt still across.
>
> There is no sub-level nine in the Adamant Crown. There are eight. You have counted them. You count them again on the way down and there are eight.

**Failure handling** — The fork expires on day 25; Crase reassigns it and logs a refusal. Surveillance cues repeat until acknowledged; none is required.
**Flags set** — `DECOY_*`, `CRASE_SURVEILLED`, `SUBLEVEL_NINE_SEEDED`.

---

### 2.8 THE FERRUM DESCENT
**Location** — Coolant sub-galleries beneath Shift Nine; the Shift Nine chapel. Wet, four degrees, no vox.
**Present** — Bez; Vaux if allied; Gallery enforcers who will not survive; Primus Odalie Sarn, shift-foreman, forty-one years on the line, not remotely a monster in her own account.
**Player goal** — Break a real cult that is really killing people, and take the one true thing it knows.
**Entry gate** — `DECOY_PURSUED`. If `DECOY_REFUSED`, the Sarn interview alone is purchasable from Mott for 2 favours as a smuggled meeting.
**Beats** —
1. Fourteen days of legwork compressed into three interrogations and a map. The cult is competent, embedded, fourth-generation, and killing people this month.
2. The assault runs as three tense decision-points, not a firefight. Violence is rare, fast and expensive.
3. **The death — deterministic, in this order, so it can be tested.** If Fenn was brought (players do, for the map): **Fenn** drowns in a flooding sluice, implant unrecoverable. Else if `BEZ_ALIVE`: **Bez** holds the sluice door — she has been coughing for eleven scenes, and if `BEZ_KNOWS_DIAGNOSIS` she is calm about it and says so, which is worse. Else if Vaux is allied: **Vaux**, and `CIPHER_HALF_VAUX` dies with him unless already transferred. Else Ilaria is blinded a second way and `ILARIA_TRUST` caps at 3.
4. **Sarn's information, which nobody else in the game has.** The cult keeps an *oral* genealogy because written membership rolls get corrected; they began in 763 and nobody remembers deciding to. Over forty years members vanish from the *record* before they vanish from the Gallery — and not the leadership. The literate ones. Both attempts at a written history, and the clerks who made them.
5. Purges arrive on an eleven-year cycle and Sarn has the dates. They do not correlate with discovery, infestation density, or Arbites activity. They correlate with sector tithe-reconciliation cycles.
6. The victory is real: brood destroyed, shift saved, +3 Crase, an Ordo commendation shown again in Act III at the worst possible moment. Then the Arbites purge the Gallery — forty thousand relocated, records closed — and the player watches the purge they earned become a records purge they did not.

**Clues gained** —
| ID | Evidence | Store | Unlocks |
|---|---|---|---|
| B12 | Erasure selects for record-keepers | Ledger | Reframes Fenn, Mott, Vosk — and the player |
| B13 | 11-year purge cycle = tithe-reconciliation cycle | Ledger | Corroborates B4b: the flush is scheduled, not reactive |
| B14 | "The unwritten inspector" | Ledger | Act III dialogue with the Amanuensis |

**Interrogation / dialogue** — Sarn will not discuss the brood, the Patriarch, or her own physiology. She discusses the erasures at length, because she has wanted to tell somebody for twenty years, and wants it noted that her shift has never missed quota.

**Branches** —
| Player choice | Flag set | Near | Far |
|---|---|---|---|
| Full assault | `GALLERY_CLEARED` | Death per beat 3 | Commendation; Ending E closed |
| Use the Erratum on Sarn | `ERRATUM_UNUSED` = false | Nobody dies; instant win | **Forecloses Ending E.** No warning prompt |
| Withdraw | `GALLERY_ABORTED` | Nobody dies; cult expands | Crase's regard curdles; Act III casualties rise |

**Sample text** —
> "We had a historian. Twice. The first in 763; I was here for the second. Good hand. Wrote everything down and hid it well.
>
> "Both times, the man stopped existing before he stopped breathing. Ration card first. Then the shift roster. Then the wife's marriage record, which is how we noticed, because the wife was still standing there.
>
> "After the second one we stopped writing. We say it now. That's all. We say it, and we make the children say it back."

**Failure handling** — Withdrawal always available; no later clue depends on the cult. Untaken, B12/B13 surface at Cantica-Null, too late to use. The Mott route exists so refusers are not punished with an information hole.
**Flags set** — `DECOY_RESOLVED`, `GALLERY_CLEARED`/`_ABORTED`, `SARN_INTERVIEWED`, `PATTERN_LITERACY`, death flags.

---

### 2.9 THREE DOORS
**Location** — Choir Spire, wax-cylinder redundancy vault. Forty thousand cylinders and no live vox. Vaux chose it deliberately and says so.
**Present** — Interrogator Dmetri Vaux; Bez if brought (she will not intervene, and remembers).
**Player goal** — Survive being told the truth by someone who dislikes you, then resolve him.
**Entry gate** — Two threads complete.
**Beats** —
1. Vaux has run an off-book file on Crase for two years. He shows B15; he hands it over only in alliance.
2. He delivers the indictment of the player. **It must be correct.**
3. He asks one question the player cannot answer well: *what have you learned this act that you did not file?*
4. He is wrong about exactly one thing — he believes Crase is the author. The same shape of error as Vhail, at a higher grade. The player should see the pattern.
5. The three doors open. Whatever the player picks, he gets the same closing speech; reuse the string in all branches.

**Clues gained** —
| ID | Evidence | Store | Unlocks |
|---|---|---|---|
| B15 | Two-year file on Crase | Ledger (only if allied) | Correct data, wrong conclusion |
| B16 | `CIPHER_HALF_VAUX` | Ledger | `CIPHER_COMPLETE` |

**Interrogation / dialogue** —

> **VAUX:** "He found you in a records office. A man who needed a key went and got one, and the key is very proud of how well it turns.
>
> "Then he died, and within six hours *she* put the rosette in your hand. In front of witnesses. Witnesses she chose. Did you notice she chose them? No. You noticed the rosette.
>
> "Everything you have found this month, somebody gave you. The Gallery — she gave you. The Cathedral is on every open filing in this hive and you needed six weeks to look at it. The savant is being calibrated by a Magos who built your terminal.
>
> "I am not saying you are corrupt. I am saying corruption would be *slower* than what is being done to you. She does not need you turned. She needs you pointed. You have been pointed for a month, you have gone off exactly where she aimed every single time, and you have filed the results where she can read them."

At `CRASE_STANDING` ≥ +3 he adds: *"She promoted you again last week. You didn't ask why."* With `STORE_HABIT_LEDGER` the player can answer beat 3 truthfully, and Vaux stops talking for six full seconds. That is the alliance's real hinge.

**Branches** —
| Player choice | Flag set | Near | Far |
|---|---|---|---|
| **Alliance** — `VAUX_SUSPICION` ≤ 1 plus the ledger answer or physical evidence (B4/B12) shown | `VAUX_ALLIED`, `CIPHER_HALF_VAUX` | He hands the cipher over unprompted — carried eleven weeks, waiting for someone worth giving it to. `ILARIA_TRUST` +1 | Ending E open; he survives Act III or dies well |
| **Execution** — every charge true and documented; Crase signs within the hour and thanks the player for their discretion | `VAUX_ALIVE` = false | Cipher half lost; +2 Crase; **+10 `RECONCILIATION`** | 2.10 mandatory; Ending E closed |
| **Blackmail** — via B9, his unauthorised file, or a forty-year-old censure | `VAUX_COERCED`, `CIPHER_HALF_VAUX_CORRUPT` | He complies; the half is complete but deliberately mis-ordered | Found out only in 2.11; repair costs a second session and 1 `VERITY`; he informs Crase in Act III |

**Sample text** —
> "I was his apprentice for forty years and I know precisely what I am. I am the one he did not choose. I have had eleven weeks to be bitter about that and I have used them, and I am finished, and here is what's left.
>
> "He chose you because you cannot leave a discrepancy alone. That is not a compliment. That is a *specification*. He picked you out of a drawer, and I have spent forty years wanting to be picked out of that drawer, and I am telling you as a professional courtesy that being the tool is not the honour you think it is."

**Failure handling** — Refuse all three and Vaux leaves Cantor for three weeks, returning at `VAUX_SUSPICION` 3 (blackmail only). No path deletes the cipher without offering 2.10.
**Flags set** — `VAUX_CONFRONTED`, `VAUX_RESOLVED`, `VAUX_SUSPICION`, branch flags.

---

### 2.10 THE LONG WAY ROUND
**Location** — Choir Spire cylinder vault (Ilaria) *or* Stack Delta (Fenn). Same shape either way; the cost is the person.
**Present** — Ilaria or Fenn.
**Player goal** — Rebuild the missing cipher half without Vaux.
**Entry gate** — `CIPHER_HALF_ROOM` held; `CIPHER_HALF_VAUX` unobtainable.
**Beats** —
1. The cipher was never written down. Vosk *used* it — dictated dispatches in it for decades. The rote-sound survives.
2. **Ilaria route.** Forty hours of raw cylinder, unshielded, inside the Choir, including her own forty-one seconds over and over. She suggests it. Cost: `ILARIA_BURNED` — trust caps at 2 and in Act III she cannot carry a message, which removes the one thing she was for. She knows, and volunteers anyway.
3. **Fenn route.** Requires B18. Sustained deliberate relay exposure at a live Vault terminal to brute-force the header corpus. Cost: `FENN_INTACT` = false, permanently. Ending E foreclosed.
4. The player chooses. No third option, no skill check, no clever way to avoid the price.

**Clues gained** —
| ID | Evidence | Store | Unlocks |
|---|---|---|---|
| B17 | Rote-sound phoneme set | Ledger | `CIPHER_COMPLETE` |
| B18b | Header-derived key | Ledger | `CIPHER_COMPLETE` |

**Interrogation / dialogue** — Fenn, offered the choice, says yes immediately and asks what he should try to remember afterwards, which is not a question anyone can answer.

**Branches** —
| Player choice | Flag set | Near | Far |
|---|---|---|---|
| Spend Ilaria | `ILARIA_BURNED` | Cipher complete | Act III loses its courier |
| Spend Fenn | `FENN_INTACT` = false | Cipher complete | Ending E closed |

**Sample text** —
> "Forty-one seconds," Ilaria says. "I've been afraid of forty-one seconds for eleven weeks. Forty hours is the same thing eight hundred times. That's arithmetic."
>
> It is not arithmetic and you both know it.
>
> "You'll want to be somewhere else," she says. "Not out of kindness. If you're in the room I'll perform, and performing takes longer."

**Failure handling** — This scene exists to prevent a soft-lock. If both are unavailable, Mott sells a third route for 4 favours: a name-bank elder, once a Choir Spire transcriptionist, supplies the phonemes in exchange for the player's own record entering the name-bank. `RECONCILIATION` +15, and the archive now has a hook in them.
**Flags set** — `CIPHER_COMPLETE`; `ILARIA_BURNED` / `FENN_INTACT`=false / `NAMEBANK_DEBT`.

---

### 2.11 THE ASH AND THE SENTENCE
**Location** — Cael's workshop; then Vosk's chambers at the twelfth hour, shutter at four degrees.
**Present** — Cael; Tass; then the player and whoever is left.
**Player goal** — Recover what Vosk burned, and read the corpse.
**Entry gate** — `ROOM_READ_FULL` or `_PARTIAL`; `CIPHER_COMPLETE` for the second half. If the ash was missed in Act I, Bez collected it — Vosk taught her to.
**Beats** —
1. Cael's spectral imager reads carbon differential off burned parchment. He is thrilled to be asked; a favour, freely given, on a Tuesday. If `CAEL_DECEIVED` he asks for B3 back as the price — the player admitting they cut the telemetry, which he already knew.
2. The plate resolves in four passes. He watches three, comments on the ink chemistry, steps out for tea during the fourth. **The fourth pass is the Erratum:** twenty-one glyphs in a pre-Imperial hand. Neither he nor Tass reads pre-Imperial; neither looks at it for more than a second and neither ever will.
3. **One of the murderers hands the player the murder weapon in a felt case and tells them to mind the edges.** He does not know. He has been reconciled and the hole in him is exactly this object's shape. The plate is physical, has no cogitator representation, and the archive does not know it exists.
4. In the chambers, the room's objects become cipher input — board position, shutter angle, calendar, boots, ink. Nine objects, and the boots said *backwards*.
5. If `CIPHER_HALF_VAUX_CORRUPT`, the first decryption yields a fluent, plausible, entirely false sentence naming Vhail. The player will believe it. It holds one scene, until the count will not reconcile: nine objects, eight tokens.
6. The true decryption, glyph-group by glyph-group. Then the query: *no such routing address; query logged.* Forty seconds later, unrequested: *confirmed: no such routing address.* Fenn, if intact, says "Kadmon Veil" aloud from another room and does not know why.

**The decoded sentence, verbatim (read in reverse object-order):**

> THE OBJECTS IN THIS ROOM ARE MY LAST FILING AND I HAVE MADE THEM CAREFULLY. I HELD THE INQUISITION'S SEAT ON THE CONCORDANCE EIGHTY-ONE YEARS AND SOMEBODY IS SITTING IN IT NOW. I SIGNED KADMON VEIL — NINE HUNDRED MILLION FOUR HUNDRED AND TWELVE THOUSAND, STRUCK OVER A TITHE DISPUTE, AND THEY ARE STILL DOWN THERE PAYING. WHAT KILLED ME IS NOT A MAN AND CANNOT BE ARRESTED. IT IS AN ENGINE AND IT KEEPS ITS BOOKS AT ROUTING CANTICA-NULL, PRIORITY BLACK, IN ORBIT OF A WORLD THE INDEX SAYS WAS NEVER SETTLED. THE INDEX IS THE CRIME. I COULD NOT WRITE THIS, SPEAK THIS, CARRY THIS, OR TELL ANY LIVING SOUL, SO I DID THE ONLY THING LEFT AND MADE MYSELF THE DOCUMENT. THE NINTH OBJECT IS THE WEAPON. IT IS IN THE ASH. DO NOT READ IT. DO NOT LET ANYONE READ IT. I AM SORRY ABOUT THE ROOM, AND ABOUT THE ELEVEN YEARS, AND ABOUT WHAT I AM MAKING YOU. — H.V.

**Clues gained** —
| ID | Evidence | Store | Unlocks |
|---|---|---|---|
| B21 | **THE ERRATUM** — imaging plate, 21 glyphs | Physical only; unfilable | The item |
| B22 | Cael's imaging log, his signature on the recovery | Ledger | Proves he handled it and did not know; Act III leverage |
| B19 | The decoded testimony | Ledger — **the game refuses to file it** | Act III |
| B20 | `CANTICA-NULL / PRIORITY-BLACK`; Kadmon named | Ledger | `KADMON_ADDRESS`; departure |

**Interrogation / dialogue** — Cael, handing it over: *"Whatever it is, it's older than the ink. Older than the idea of the ink. Vosk collected things like this — a mind for dead scripts, which I always thought was a waste in a man who could have done real work."* He is being fond about the thing that stopped his friend's heart.

**Branches** —
| Player choice | Flag set | Near | Far |
|---|---|---|---|
| Take the plate | `ERRATUM_OBTAINED`, `ERRATUM_UNUSED` | The item | All endings |
| Tell Cael what it is | `CAEL_WARNED` | He goes grey and asks the player to leave | Act III: he defects once, briefly, and it costs him |
| Show B19 to the retinue | `TESTIMONY_SHARED` | +1 `ILARIA_TRUST`; Bez asks what a tithe dispute is | Act III: they come willingly |

**Sample text** —
> *I am sorry about the room, and about the eleven years, and about what I am making you.*
>
> You read it four times. The fourth time you notice the tense — *making*, not *made* — and that he did it on purpose, because Hadrian Vosk never once in his life left a tense loose.
>
> He is not apologising for what he did. He is apologising for what you are going to do, which he worked out in advance and was patient enough to be sorry about ahead of time.

**Failure handling** — Unmissable; if the player never visits Cael, Tass brings the developed plate to the chambers unbidden. If the corrupt-cipher path stalls more than one scene, Ilaria flags the object-count discrepancy aloud.
**Flags set** — `ERRATUM_OBTAINED`, `ERRATUM_UNUSED` = true, `KADMON_ADDRESS`, `ACT2_COMPLETE`, `CIPHER_COMPLETE`.

> **Canon rule — no retroactive `ROOM_READ_FULL`.** Decrypting the testimony sets `CIPHER_COMPLETE` and nothing else. A player who read three signs in Act I gets the sentence, gets the address, gets the weapon — and still does not have `ROOM_READ_FULL`, because the sentence is *what Vosk wrote*, not *how he wrote it*, and Ending E needs the method, not the message. The only way to acquire it late is 2.0b.

---

## ACT-LEVEL SYSTEMS NOTES

### The Erratum as an item
One object, one use. Delivery: leave the plate where a specific person will read it, or hand it over. Reading it aloud in the pre-Imperial phonology is disallowed at UI level, with no explanation given. Using it on a human being sets `ERRATUM_UNUSED` = false and **permanently forecloses Ending E**. The game never says so. The requirement list is never visible in-game.

**Three places the player will want to use it, all easy, all satisfying:**
1. **Vhail**, after 2.6 — loathsome, proven, and about to be replaced by an identical supplier. Leave the plate in his private confessional; he reads everything left there.
2. **Sarn**, in 2.8 — using it here *saves a named retinue member's life*. The humane choice by every measure the player can see, and the cruellest offer in the game. No warning prompt.
3. **Kesh**, at her desk, over tea, in the Act II/III seam. She will have shown the player the numbers by then. She never raises her voice. She has deleted four hundred million people and she has a teapot.
*(A fourth, for the Act III writer: Crase at Cantica-Null.)*

### Reconciliation rising
Sources: +4 filing a false conclusion (2.6), +10 executing Vaux, +12 killing Vhail, +15 name-bank debt, +2 per case-file query of a struck name, +3 per Cael upgrade. Careless finish 55–70; careful 15–25. **Nothing below is ever announced.**

- **≥ 25 — options vanish.** In 2.6 the option `> ACCUSE — you wired fifty-one confessionals into a vox trunk` is silently absent. The list renumbers so there is no gap.
- **≥ 40 — the retinue remembers what you don't.** Bez, unprompted: *"After sub-level nine. After you carried me out — and don't tell me you didn't, I was awake for most of it."* The player has no record; the query returns `NO SUCH OPERATION`. The Crown has eight sub-levels; the string was seeded in 2.7, before the symptom.
- **≥ 55 — the case-file shortens.** The Act I autopsy entry drops from 340 words to 190; *"no toxin, no coercion residue, and no psychic trace were detected"* becomes *"no residue was detected."* A ledger-keeping player still has the long version. Never point at it.
- **≥ 55 — retinue files thin.** Ilaria's entry loses the paragraph about the forty-one seconds; Fenn's loses the word *un-networked*.
- **≥ 70 — you go.** The rosette record shows an initial where the given name was. Two dialogue options are removed from every subsequent scene. The act summary screen does not contain the player's name anywhere.

### Verity and Standing
`VERITY` +1 per ledger-held clue the case-file version has since altered; Act II's honest maximum is 6 (B1, B4, B6, B12, B19, and one of B9/B23). Below 3, the Act III offer scene gives the player no rebuttals at all. `CRASE_STANDING` starts +1: decoy accepted +3, refused −2; flush deferred −2; Vaux executed +2; Vhail denounced +1. Above +3, I-5 fires — play the politics well and arrive at 2.9 with the worst possible hand.

### Where the UI lies in this act
1. B4's signature line reads `[UNSIGNED]` in the case-file after 2.9; the vellum is unchanged.
2. `CONCORDANCE OF RECORD: 6/6` is equipment-condition framing applied to a man.
3. Unrequested confirmations (2.4, 2.11) — the archive answering a question the player did not repeat.
4. The unmarked Vault query returning `CN-∅` all act, proving Vhail was redundant.
5. The autopsy word-count shrink at Reconciliation ≥ 55.
6. Crase's surveillance appears in no system view; inferable only from lift logs and latency.

---

## FLAG LEDGER

| Flag | R/W | Where | Notes |
|---|---|---|---|
| `ROOM_READ_FULL` / `_PARTIAL` | R/W | 2.0b only | **Never** retro-set by 2.11 |
| `CIPHER_HALF_ROOM` | R | 2.11 | From Act I |
| `CIPHER_HALF_VAUX` / `_CORRUPT` | W | 2.9 | Alliance / blackmail |
| `CIPHER_COMPLETE` | W | 2.10, 2.11 | Never soft-lockable |
| `VAUX_ALIVE`, `VAUX_ALLIED`, `VAUX_COERCED`, `VAUX_SUSPICION` (0–3), `VAUX_CONFRONTED`, `VAUX_RESOLVED` | R/W | 2.9 | `VAUX_ALLIED` = Ending E req. |
| `FENN_INTACT` | R/W | 2.1, 2.3, 2.8, 2.10 | Ending E req. |
| `FENN_SHIELDED`, `FENN_ISOLATED`, `FENN_SUPPRESSED`, `FENN_BURNED_PARTIAL` | W | 2.1–2.3 | |
| `FENN_DRIFT` (0–6) | R/W | 2.1, 2.3 | Shown as equipment condition |
| `ILARIA_TRUST` (0–5), `ILARIA_BURNED` | R/W | 2.4, 2.9–2.11 | Burned caps trust at 2 |
| `BEZ_ALIVE`, `BEZ_KNOWS_DIAGNOSIS` | R/W | 2.5, 2.8 | Changes her death line |
| `MOTT_FAVOURS` (int), `MOTT_MET` | R/W | 2.4, 2.5, 2.8, 2.10 | Owed; never collected in Act II |
| `CHOIR_TRUSTED` / `CHOIR_HOSTILE` | W | 2.4 | |
| `CHOIR_FED`, `CHOIR_DEFERRED`, `WICK_KNOWN` | W | 2.5 | Epilogue-only reads |
| `DECOY_REFUSED` / `_PURSUED` / `_DELEGATED` | W | 2.7 | Ending E requires REFUSED |
| `GALLERY_CLEARED` / `_ABORTED`, `DECOY_RESOLVED`, `SARN_INTERVIEWED`, `PATTERN_LITERACY` | W | 2.8 | |
| `VHAIL_BROKEN`, `_DENOUNCED`, `_LEVERAGED`, `_DEAD`, `VHAIL_EXPOSED`, `CATHEDRAL_MAPPED`, `RUSE_WARNED` | W | 2.6 | `VHAIL_BROKEN` fires I-1 |
| `KESH_IDENTIFIED` | W | 2.3 | Via B23 |
| `CAEL_SUSPECTED`, `CAEL_DECEIVED`, `CAEL_WARNED` | W | 2.2, 2.11 | |
| `ERRATUM_OBTAINED`, `ERRATUM_UNUSED` | W | 2.8, 2.11 | Ending E requires UNUSED |
| `KADMON_HEARD`, `KADMON_ADDRESS` | W | 2.1, 2.11 | |
| `CRASE_STANDING` (−5..+5), `CRASE_SURVEILLED`, `SUBLEVEL_NINE_SEEDED` | R/W | 2.5, 2.7, 2.9 | |
| `RECONCILIATION` (0–100), `VERITY` (int) | R/W | all | |
| `STORE_HABIT_FILED` / `_LEDGER`, `B4_FILED` | W | 2.0, 2.4 | |
| `NAMEBANK_DEBT` | W | 2.10 | Last-resort route |
| `PURGE_ORDER_SEEN`, `LYLE_QUERIED`, `CLOCK_DAY40_ARMED` | W | 2.4 | |
| `TESTIMONY_SHARED`, `ACT2_OPEN`, `ACT2_COMPLETE` | W | 2.0, 2.11 | |
| `TRANSPORT_WRIT` / `_STOLEN` / `_DEBT` | R | 2.5 | Constrained here, set in Act III |

---

## HANDOFF TO ACT III

1. `CIPHER_COMPLETE` is set. There is no exit from Act II without it.
2. `ERRATUM_OBTAINED` is set. `ERRATUM_UNUSED` may be either value; if false, Ending E is closed and Act III must never say why.
3. `KADMON_ADDRESS` is set. The player knows the words *Cantica-Null* and *Kadmon Veil*, and knows the archive denies both.
4. Exactly one of: Vaux allied, coerced, dead. If dead or coerced, one of `ILARIA_BURNED`, `FENN_INTACT`=false, or `NAMEBANK_DEBT` is set — Act III may assume the cipher was paid for with a person.
5. `RECONCILIATION` 15–75; `VERITY` 0–6; `CRASE_STANDING` −4..+5, and Crase is still, in her own account and most of her actions, the player's friend.
6. `MOTT_FAVOURS` is a positive integer never yet collected. Act III collects it once, at Cantica-Null.
7. Bez may be dead, Fenn blank, Ilaria burnt. At least one of the four retinue is intact. The game guarantees no more than that.
