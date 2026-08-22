# ACT III — THE AMANUENSIS

**Dramatic question:** *Now that you know what it is, what is it worth?*
**Answer:** *More than you can pay, and it will show you the arithmetic.* There is nothing left to find out in Act III — only things left to stand in front of. A debt taken to get off-world, a world that has waited ninety-one years for someone exactly like you, a master who signed the order, and a machine that has answered every question honestly and has never once been asked the right one. The act ends with five doors, all open, all locked by things the player chose in Acts I and II without knowing they were choosing.

---

## SCENE INDEX

| # | Scene | Location | Type | Gated by | Est. minutes |
|---|---|---|---|---|---|
| 3.1 | GROUNDED | Adamant Crown; Skydock Cantor Nine | Branch hub | `KADMON_ADDRESS` | 14 |
| 3.2 | THE CEREMENT | In transit | Quiet hub | any `TRANSPORT_*` | 6 |
| 3.3 | THE SPIRIT OF THE PATIENT VERDICT | Cutter's cogitator well | Puzzle | 3.2 | 12 |
| 3.4 | THE QUIET WATCHES | Ship interior | 4 optional character scenes | per-member | 24 |
| 3.5 | MAKEFAST | Kadmon Veil, Hive Erasmus | Exploration | translation | 12 |
| 3.6 | THE MARCHIONESS | The Muster Hall | Inverted interrogation | 3.5 | 12 |
| 3.7 | THE UNSENT | The Reading Room | Document set | `VURR_RESPECT ≥ 1` | 9 |
| 3.8 | THE ONE CONVOY | Tithe Yard / Muster Hall | Choice | 3.6 | 10 |
| 3.9 | CANTICA-NULL | Station approach; the Correctorum floor | Set-piece / interrogation | `KADMON_ADDRESS` | 16 |
| 3.10 | WHAT CRASE BELIEVES | The Fifth Chair | Interrogation | 3.9 | 12 |
| 3.11 | THE FILED WORDS | The Correctorum floor | Revelation | 3.10 | 10 |
| 3.12 | THE OFFER | The core intake | Revelation / system | 3.11 | 14 |
| 3.13 | THE WIDE CHOICE | The core intake | Endings | 3.12 | 6+ |

---

### 3.1 GROUNDED
**Location** — Vosk's chambers; Skydock Cantor Nine.
**Present** — Player, retinue, Sub-Officio Tomas Rill; Crase (A), Rogue Trader Beneficia Sazh (C).
**Player goal** — Leave Vardan Prime with retinue, Erratum and address intact.
**Entry gate** — `KADMON_ADDRESS`.
**Beats**
1. Routine passage booking is not denied but *pended*: six is more than five, five is the threshold, Rill is very sorry.
2. **Route A — the writ.** Crase signs in person, over recaff, without asking where you're going. Asks after Bez's cough. The boilerplate carries a routing-telemetry clause (**C1**).
3. **Route B — theft.** *The Patient Verdict* sits at Berth Nine under your own seal, which is not the same as permission. Bez walks the party aboard in nine minutes; the seal-break logs on clamp-release.
4. **Route C — the debt.** Sazh's warrant-ship *Argument From Silence*. She takes no coin — coin is a record. She takes favours: unspecified, called at her convenience, sealed by both parties saying the terms aloud and writing nothing (**C6**).
5. Departure within eight hours. The Cerement's edge is nine days out.
**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| C1 | Writ clause 11(d): position and manifest reported to the countersigning authority at each waypoint | Crase's writ [COG] | Proves Route A is a leash; permits leaving the writ at Kadmon (3.8) |
| C6 | Sazh's tally — forty-one open favours, no names written, all remembered | Spoken, Route C [PHYS only] | Seeds `SAZH_DEBT` payment in 3.13 |

**Interrogation / dialogue** — **Crase** signs anything, volunteers nothing, admits she knows the routing address exists ("Hadrian's little ghost. He raised it at a session once and I told him to let it go") and will not define *session*. Pressed hard she says the truest thing in the game so far: *"I have never lied to you."* Bank it. **Sazh** never asks what a favour is for, only what it weighs; refuses to be paid in the Erratum ("I'm not stowing a thing that kills whoever reads the label"). **Rill** can be frightened into a four-hour delay, which is enough for Route B.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Take the writ | `TRANSPORT_WRIT`, `CRASE_STANDING +1` | Legal passage; Vosk's effects forwarded, curated | Crase reaches Cantica-Null first and prepares the room; B pre-vouched, A/E must cross a courteous cordon |
| Steal the cutter | `TRANSPORT_STOLEN`, `CRASE_STANDING -2` | No papers, no waypoint reports | Fugitive: closes C unless `VAUX_ALIVE` or Sazh; full access to 3.3 |
| Take the favour | `TRANSPORT_DEBT`, `SAZH_DEBT` | Fast hull, no questions | She arrives at Cantica-Null in 3.13 and calls it: she wants the Correctorum intact, functional, and hers |
| Read 11(d) aloud to Crase | `CRASE_STANDING -1`, `VERITY +1` | She smiles and says "Yes" | Every later Crase scene drops one layer of warmth |

**Sample text**
> Sazh does not shake hands. She holds her own together in front of her, like a woman checking she still has ten.
>
> "I'll write nothing down. You'll write nothing down. In four months or forty years I'll come and say a thing, and you'll do it. People find that romantic. It's just the only contract in this sector nobody can edit."
>
> Behind her the loaders are stacking crates with no stencils. Not blank. Scraped.

**Failure handling** — All routes stay open; no timer. If all three burn (Sazh dead, `CRASE_STANDING ≤ -4`, cutter lost), Mott books grain-hopper passage at `MOTT_FAVOURS ≥ 1`: sets `TRANSPORT_STOLEN`, `RECONCILIATION +4`, no cutter access. No soft-lock.
**Flags set** — one `TRANSPORT_*`; `SAZH_DEBT`; `CRASE_STANDING` delta.

---

### 3.2 THE CEREMENT
**Location** — Nine days along the warp-scar's edge. **Present** — Retinue only. No antagonists.
**Player goal** — Nothing. That is the point. **Entry gate** — any `TRANSPORT_*`.
**Beats**
1. Translation. The Cerement does not howl; the Geller readouts hold *perfectly steady*, and that is worse.
2. First time out of relay range all game. The UI stops lying about the world and starts lying about the player.
3. Every 36 hours one reconciliation symptom fires: an Act II note one clause shorter; a dialogue option previously used now absent; at `RECONCILIATION ≥ 60`, a memory the player *has* queried and returning nothing.
4. Free movement to 3.3 and the four scenes of 3.4. Sleeping advances the symptom clock. There is a real cost to lingering and the game never says so.
**Clues gained** — none; hub only.
**Interrogation / dialogue** — Ambient. Bez cleans a weapon she won't need. Fenn recites a shipping schedule to nobody. Ilaria sits with her back to the hull because the Cerement is quiet and she has not had quiet in thirty years.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| See all four character scenes | `VOYAGE_FULL` | — | Every epilogue gains named beats; D requires at least one |
| Sleep through the passage | `RECONCILIATION +8` | Arrive sooner | Options begin vanishing before the Marchioness |
| Re-query your own case file | `RECONCILIATION +2` each | Reassurance | The file agrees with you slightly more each time |

**Sample text**
> Nine days and the field readout has not moved. Not a flicker, not the drift you learned to expect. A perfect line.
>
> Bez says it's a good omen. Bez was in the Guard and does not believe in omens; she says it the way you say a thing to fill a room.
>
> You open your notes to check something you wrote on the fourth day of all this, and the sentence is there, and it is one clause shorter than you remember, and you cannot swear to that.

**Failure handling** — The voyage can be rushed, not skipped. Rushing is legal and the game permits it silently.
**Flags set** — `RECONCILIATION` drift; `VOYAGE_FULL`.

---

### 3.3 THE SPIRIT OF THE PATIENT VERDICT
**Location** — Cutter's cogitator well (B), your cabin (A), Sazh's salvage locker (C).
**Present** — Player; Cael's kit if retained; Fenn if intact. **Player goal** — Recover Vosk's last orders. **Entry gate** — 3.2.
**Beats**
1. The hull is running orders filed six days before he died — the six days circled on the Act I calendar.
2. **C2 — nav plot** to Cantica-Null, with a return leg *entered and then struck by hand*, field by field. He plotted his way home, sat with it, then removed it, and left the removal where you would find it.
3. **C3 — manifest:** one crate. An empty reliquary casket, mass-declared, addressed to himself. Sized for carrying one written thing physically.
4. **C4 — voice slate:** nineteen seconds. Throat cleared. *"If you are hearing this you have done the first difficult thing. Do not do the obvious second one."* Cut. No warmth; the player will want warmth and won't get it.
5. **C5 — standing order:** the vox is not disabled but *physically uncoupled*, the part in a drawer. NOTHING THIS HULL CARRIES IS TO BE TRANSMITTED.
6. C2+C3+C5 yields the method: the Correctorum edits only what passes through it, and reads anything handed to its intake as a correction request. Delivery must be physical, in a casket, correctly addressed.
**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| C2 | Nav plot, return leg struck | Ship spirit [PHYS] | `VOSK_METHOD_PARTIAL` |
| C3 | Empty mass-declared reliquary casket | Manifest [PHYS] | Physical carriage of the Erratum, Ending E |
| C4 | Nineteen-second slate | Ship spirit [PHYS] | Bez identifies the fourth second: an intake breath before the cut |
| C5 | Uncoupled vox + written order | Drawer [PHYS] | `VOSK_METHOD_UNDERSTOOD` with C2, C3 |

**Interrogation / dialogue** — The spirit is a filing clerk with a wound; it answers in manifest fields. Fenn dates the orders to the hour. Bez, on the slate: "That's him about to say something else."
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Route B — full access | `VOSK_METHOD_UNDERSTOOD` | All four clues | Ending E reachable |
| Route A — Crase forwards his effects | `VOSK_METHOD_PARTIAL` | C2, C3; C4 "damaged in transit"; C5 absent | Needs Cael's kit or `FENN_INTACT` to recover, else E stays shut |
| Route C — buy the spirit-core back | `VOSK_METHOD_UNDERSTOOD`, `SAZH_DEBT_2` | Full access | Two favours called in 3.13 |
| Wipe the spirit after reading | `VERITY +1` | Crase learns nothing from the hull | Removes her foreknowledge advantage in 3.10 |

**Sample text**
> The return leg is in the log and then it is not in the log. Not blank — struck. Someone entered a course home and took it out again by hand, field by field, the way you cancel a thing you want to be able to prove you cancelled.
>
> He plotted his way back. He sat with it. Then he deleted it, and left the deletion where you would find it, because a blank field says nothing and a struck one says *I considered it*.

**Failure handling** — If the well is never entered, Bez surfaces C4 in 3.4 and the station cites C2 in 3.9; `VOSK_METHOD_UNDERSTOOD` is still assemblable in 3.11.
**Flags set** — `VOSK_METHOD_PARTIAL` / `_UNDERSTOOD`; `SAZH_DEBT_2`.

---

### 3.4 THE QUIET WATCHES
**Location** — Gun bay, astropath's cell, cargo hold, wardroom. **Present** — one retinue member each.
**Player goal** — Nothing useful. This is where the endings get their weight. **Entry gate** — per-member, below.
**Beats** — 1. Four self-contained scenes, enterable in any order, each once. 2. Each costs one symptom tick of ship-time. 3. Each ends with the member asking the player for something small. 4. None of them advances the case.
**Clues gained** — C7, C8, C9, C10 (see sub-scenes).
**Interrogation / dialogue** — Per sub-scene. Common rule: no member of the retinue will raise their own subject unprompted. The player must go to them.
**Branches** — Per sub-scene tables below; `VOYAGE_FULL` is set only if all four available scenes are seen.

**3.4a — BEZ: THE BEQUEST** (`BEZ_ALIVE`)
**Beats** — 1. She strips a weapon that doesn't need it and folds her rag inward, as all game. 2. She produces Vosk's sealed bequest and asks the player to read it aloud — the ask costs her more than the lung rot, because *she cannot read.* 3. Four lines: a pension on an account in a hive that no longer files returns; an instruction that she be discharged; and *"You were the only person in my service I never had to check."* 4. Player may name the lung rot or not. 5. She asks to come to Cantica-Null anyway.
**Clue** — | C7 | Vosk's four-line bequest | Bez's kit [PHYS] | Proves he was arranging survivors six days out |
**Dialogue** — She won't discuss dying; she will discuss *paperwork about* dying at length, because that is a task. With `BEZ_KNOWS_DIAGNOSIS`: "Then don't spend me on something small."
**Branches** — read it truthfully → `BEZ_BEQUEST_READ`, `VERITY +1`; soften or edit it → `RECONCILIATION +5` and the game quietly lets you notice you just did what they do; refuse → she never asks again and gets one line in every epilogue.
**Sample text**
> She holds the seal out at arm's length like it might go off.
>
> "Eleven years. He knew. Never said a word, never had a clerk do it for me, never once made me ask." She pushes it another inch toward you. "So I'm asking you instead, and I'd take it kindly if you didn't make a thing of it."
>
> The rag in her other hand is folded inward. It has been all week.

**3.4b — ILARIA: THE MISSING YEARS** (`ILARIA_TRUST ≥ 2`)
**Beats** — 1. Out of relay range she is lucid in a way the player has never seen. 2. She keeps a private wax tally — not of memories but of *gaps*: dates she can prove she worked and cannot prove she was present. 3. Four hundred and six days across three years. 4. She asks the player to confirm she isn't insane. The honest answer is worse than insanity. 5. Told the truth, she doesn't break: "Then I'm the only one aboard who's been in the room with it and lived." She volunteers to be the machine's mouth at Cantica-Null.
**Clue** — | C8 | Wax tally, 406 unremembered days | Her cell [PHYS] | `ILARIA_KNOWS`; proves reconciliation is being applied non-consensually |
**Dialogue** — She will not be pitied and will not be lied to twice; if the player lied in Act II she names the date, because it's on the slate.
**Branches** — truth → `ILARIA_KNOWS`, `ILARIA_TRUST +2`, `VERITY +2`, and in 3.11 she may refuse to be used and the refusal is honoured; "you're ill" → `RECONCILIATION +6`, she is docile at Cantica-Null and it is horrible; silence → `ILARIA_TRUST -1`.
**Sample text**
> "Four hundred and six," she says. "I don't count the days I lost. I count the days that lost me."
>
> The slate is old, the wax scraped down so often it holds a shine. She turns it over the way Bez turns over a weapon.
>
> "I want you to tell me I'm mad. I've been rehearsing being told I'm mad for three years and I'd like some use out of it."

**3.4c — FENN: DECAY** (`FENN_INTACT`)
**Beats** — 1. With no relay, Fenn stops contradicting himself; for four days he is the most reliable object in the sector. 2. Then he begins correcting *forward* — reciting sector records that do not exist yet, dated eleven to forty years hence, in the same flat voice. 3. The forward records are internally consistent and describe worlds simplifying: fewer names, fewer districts, fewer worlds. 4. He is not prophesying. He is running the reconciliation logic his implant was tuned to against the data he holds, and arriving where it arrives. 5. Player records it in the cogitator (fast, editable) or transcribes by hand (slow, permanent, one symptom tick).
**Clue** — | C9 | Forward recitation: 41 compliant worlds → 33 by 087.M42 | Spoken [PHYS if hand-copied, COG if filed] | Independent corroboration of the machine's projections in 3.12 |
**Dialogue** — Cheerful throughout, baffled that anyone is upset. Asked if he's guessing: "No. I'm finishing."
**Branches** — hand-copy → `C9_PHYS`, citable at the Amanuensis; file it → `C9_COG`, and the list is three worlds shorter on arrival and the game never says so; shield him mid-recitation → `FENN_SHIELDED`, C9 truncated.
**Sample text**
> "Thessaline compliant worlds, cycle 087," Fenn says, in the voice he uses for laundry. "Thirty-three."
>
> "Forty-one, Fenn."
>
> "Forty-one at present, yes." He beams. "Thirty-three at cycle 087. Shall I do the tithe grades?"
>
> He is a filing system with a mouth, reading you the sector's ending in the tone of a man listing socks, and Vaux has stopped writing.

**3.4d — VAUX: THE APOLOGY** (`VAUX_ALLIED`)
**Beats** — 1. He lays out his two-year Crase dossier, uninvited, in order, and stands while you read it. 2. It is meticulous, correct in every particular, and wrong in its conclusion: he had Crase as the author. 3. The apology is not for obstructing you. It is for forty years of being right about small things — "I was a good clerk to a man who needed a witness. He chose you because you were not a good clerk." 4. He asks one thing: to be told first, in words, to his face, whatever you decide. 5. Agree and he hands over Vosk's dead-man courier writ — the sealed physical channel out of the sector.
**Clue** — | C10 | Vaux's hand-kept Crase dossier, two years | Wardroom table [PHYS] | `CONCLAVE_CHANNEL`; corroborates 3.10 |
**Dialogue** — Will not be thanked; warmth is met with a change of subject. On why he stayed: "Because you would have got yourself killed in the Vault in the first week, and then who would I have argued with."
**Branches** — agree → `VAUX_APOLOGY`, `CONCLAVE_CHANNEL` (Ending C reachable on any route; E's final step requires him present); decline → he keeps the writ; tell him he was right about Crase → `VERITY +1`, and he speaks first in 3.10 and costs her the opening.
**Sample text**
> "Two years. Every hour I wasn't on duty. I had her at the centre of it."
>
> "You had her at the centre of it because she is at the centre of it."
>
> "I had her as the *author*." He squares a page that is already square. "There is a difference between a hand and a pen, and I have spent two years interrogating a pen."
>
> He does not sit down. He has not sat down in front of you once.

**Failure handling** — Every sub-scene is skippable and each has a punishing default: C7 is read at Bez's funeral instead, unopened; Ilaria's scene fires one-sided at Kadmon and she does not accept the answer; Fenn's projection is reconstructed by Vaux, slower and uncitable; without Vaux, `CONCLAVE_CHANNEL` must come from an unspent writ or `SAZH_DEBT_2`, and if neither is held Ending C is legitimately closed.
**Flags set** — `BEZ_BEQUEST_READ`, `BEZ_KNOWS_DIAGNOSIS`, `ILARIA_KNOWS`, `C9_PHYS`/`C9_COG`, `FENN_SHIELDED`, `VAUX_APOLOGY`, `CONCLAVE_CHANNEL`, `VOYAGE_FULL`.

---

### 3.5 MAKEFAST
**Location** — Kadmon Veil: high orbit, then Hive Erasmus. **Present** — Retinue; Port-Adept Alcuin Rhod, 19, palimpsest-scribe.
**Player goal** — Land on a world that has not seen a ship in ninety-one years.
**Entry gate** — Translation. Kadmon lies nine days short of Cantica-Null on the routing address and cannot be bypassed.
**Beats**
1. **They hail you correctly** — full challenge-and-response, current protocol, letter-perfect, in a boy's voice. They have drilled it daily for ninety-one years against the day.
2. **The Owed (C11).** Ninety-one cycles of tithe, assembled to grade, sealed in vacuum silos across forty square kilometres of Tithe Yard, each stack labelled with cycle and assessed value. None opened. None eaten.
3. **What degraded.** Hive Erasmus runs three of eleven plasma reactors; the other eight are cold and are oiled, turned and blessed daily by the Guild of Wakers, none of whom has seen one lit. Lasgun cells ran out in year 31 — the PDF drills with turned-wood rifles; the real weapons live in the Faithful Armoury and are counted quarterly (4,106 rounds). Paper ran out in year 19: documents are scraped and rewritten, and the palimpsest scribes are a licensed guild with a nine-year apprenticeship. Ink is lamp-black. Void-shields failed in year 44 and the generator is maintained anyway.
4. **What grew.** *The Attending* — the six surviving astropaths of forty broadcast at Terce, then the whole hive holds a listening minute; sixty seconds a day, ninety-one years. *The Reading of the Delay* — Ecclesiarchy doctrine holding the silence to be a Vigil, with 611 canonically permitted reasons for a reply to be late, recited on rotation. *Filing the Dead* — every corpse's record completed and vaulted against the audit that is coming. *The Standing Requisition* — children examined in requisition form at age nine; bad form is a civic disgrace.
5. Nobody asks the player for orders in the street. They step aside. Several kneel. A woman asks Bez whether the war is over and does not say which war.
6. Rhod escorts, cannot stop looking at Fenn's implant — the newest technology on the planet — and says "when the reply comes," never *if*.
**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| C11 | The Owed: 91 cycles of sealed, unopened tithe | Tithe Yard [PHYS, observed] | Bargaining weight in 3.8; costed in every epilogue |
| C13 | Last signal received, year 0: a routine tithe-receipt acknowledgment, framed | Muster Hall wall [PHYS] | Dates the erasure to the day; the stamp carries Kesh's office sigil |

**Interrogation / dialogue** — Rhod answers everything with terrifying good manners and gives the *reported* population figure without blinking. Nobody on this world will speculate about Terra; speculation is not their business.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Walk the Tithe Yard | `KADMON_SAW_OWED` | Nausea; Rhod recites the assessed value of each stack from memory | Ending A's epilogue names the Yard |
| Go straight to the Marchioness | — | Faster | `VURR_RESPECT` starts one lower; she notices what you did not look at |
| Ask Rhod the true population | `VURR_RESPECT +1` | He gives the reported figure and then, quietly, the real one | Flags him as the `KADMON_WITNESS` candidate in 3.8 |
**Sample text**
> The challenge comes on the correct frequency in the correct form, and the voice is a boy's, and he does not fumble a syllable.
>
> Somebody wrote that drill down ninety-one years ago. Somebody has run it every watch since, into a system where nothing arrives, and taught it to their replacement, and their replacement's replacement, and this morning it was that boy's turn and he was ready.
>
> Bez says, "Throne." Not loudly. It is the only thing anyone says for a while.

**Failure handling** — Not skippable. Refuse to disembark and Vurr comes up in a shuttle held together by ninety-one years of improvisation, at 104, and the scene plays in your wardroom. This is worse.
**Flags set** — `KADMON_VISITED`, `KADMON_SAW_OWED`.

---

### 3.6 THE MARCHIONESS
**Location** — The Muster Hall. **Present** — Marchioness Ottiline Vurr (104); Sub-Deacon Barsimeus Tolk; Rhod; retinue.
**Player goal** — Survive being asked a question. **Entry gate** — 3.5.
**Beats**
1. Not a throne room. She receives you standing, in the hall where the tithe ledger is kept, with the current cycle's audit made current — three hours' notice, all of it spent on the numbers, "because you will want the numbers first."
2. She delivers the report she has given to no one: population, tithe, heresy count (nine in ninety-one years, all resolved locally), infrastructure, shortfalls. Four minutes. Flawless.
3. **The census discrepancy (C14).** She reports 900,412,000 — the last figure Terra acknowledged. The true census is 812,340,000 and falling ~1.5 million a year, chiefly respiratory. She reports the old number because revising an acknowledged figure without authorisation would be falsification, and she will not falsify. She states this without apology and dares you to call it madness.
4. Pity is cut off cold. This beat is what makes her formidable instead of pitiable — play it before the question, always.
5. **The question**, polite version first, then the real one after the player's hesitation.
6. Tolk tries to read the player's arrival as prophecy. Vurr shuts him down in six words; she has been doing it for forty years.
**Clues gained** — | C14 | Reported vs. true census, with her written justification | Ledger [PHYS] | `VURR_RESPECT`; establishes the record-weight traded in Ending D |
**Interrogation / dialogue** — Any question about her world is answered instantly; any question about herself is redirected to her world. What opens her is competence: accurate citation of C11 or C13 → `VURR_RESPECT +1`. What closes her is sympathy.

Her exact words, in order:
> "Have you come with orders?"
>
> *[pause; she reads your face for as long as it takes]*
>
> "No. I see. Then let me ask the useful question instead, since you have come a long way and I am a hundred and four and neither of us should waste the afternoon.
>
> Are we being punished, or are we being forgotten? I have prepared this world for one of those for sixty years. I should like to know whether I prepared for the wrong thing."

And, if pitied:
> "Don't. I have governed a world for sixty years with no resupply, no reinforcement and no instruction. Its tithe is complete. Its people are loyal. Its heresy count is nine. Whatever you have decided to feel about me, feel it after you have read the ledger."

**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Answer honestly that you have no orders | — | She proceeds to the real question | The 3.8 answer branch opens fully |
| Cite C11 or C13 accurately | `VURR_RESPECT +1` | She grants the Reading Room herself | 3.7 opens without Rhod's smuggling |
| Offer sympathy | `VURR_RESPECT -1` | Cut off cold | She addresses the retinue, not you, for the rest of the visit |
| Stonewall | `KADMON_SILENT` | Eleven seconds, then "Very well" | The Amanuensis cites it in 3.12 |
**Sample text**
> She does not sit, so nobody sits.
>
> "Nine hundred million, four hundred and twelve thousand. That is the acknowledged figure. The true figure is eight hundred and twelve million and falling by about a million and a half a year, chiefly respiratory."
>
> "Then why report the other?"
>
> "Because the acknowledged figure was acknowledged." She turns a page with two fingers, carefully, because the page has been three other pages. "I am not authorised to correct Terra. I am authorised to be accurate about being unable to."

**Failure handling** — Cannot be failed, only refused; refusal is a branch, not a loss.
**Flags set** — `VURR_RESPECT` (0–3); `KADMON_SILENT` on stonewall.

---

### 3.7 THE UNSENT
**Location** — The Reading Room. **Present** — Player, Rhod. **Player goal** — Read sixty years of a woman shouting into a wall.
**Entry gate** — `VURR_RESPECT ≥ 1`. Vurr grants access herself and does not stay.
**Beats**
1. Nine hundred and twenty-two numbered reports. Bound, indexed, cross-referenced, never sent — the Attending carries only a two-hundred-word summary; the reports wait for a courier.
2. Three may be read in full for free; each further read costs one symptom tick, because reading is how the archive reaches you.
3. The tone decays across the decades and then recovers into something worse than decay.
4. Rhod has read all of them twice and hand-copied forty onto scraped vellum, which means he knows them.
**Clues gained** — | C12 | Reports 47, 415, 833 | Reading Room [PHYS] | `KADMON_PROOF` — documentary proof of a loyal world's erasure; required for Ending C's indictment |

**Report 47 — cycle 6.** *Formal. Faintly irritated.*
> "…the tithe for cycles four, five and six has been assembled to grade and awaits collection at the Yard. We note that no collection fleet has attended and that no acknowledgment of receipt stands against cycles four or five. We assume an administrative delay at sector level and have taken no independent action. We would be grateful for a corrected collection schedule at the Assessor's convenience."

**Report 415 — cycle 44.** *Over-documented. Defensive. Anxiety expressed as completeness.*
> "…appended: muster rolls I–XI; shrine attendance; heresy dispositions in full, including the Tolliver matter of cycle 34, resolved locally within nine days without recourse to sanction; maintenance certificates for all eleven reactor housings including the eight that are cold. I list also what we have *not* done. We have not reduced the tithe. We have not opened the Faithful Armoury. We have not levied against the Owed, though the Owed would feed this world for one year and four months. We have not despatched an unauthorised vessel. We have not stopped. I record these not as complaint but as evidence, in the event that evidence is what is wanted."

**Report 833 — cycle 83.** *Flat. Restored to pure form. This is the worst one.*
> "Report 833. Cycle 83 complete. Tithe assembled to grade and stored. Population within tolerance. Heresy nil. Reactors as previously. Nothing further to report. O. Vurr, Marchioness, for the Emperor's attention or another's."

**Interrogation / dialogue** — Between 415 and 833 sits Report 728 (cycle 67), shown only as a spine unless the player specifically asks for the year Vurr's husband died. One paragraph, addressed *"To the Adept who opens this,"* asking for acknowledgment in any form, even a mark. Keep it out of the main three. Make the player go looking.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Read the three | `KADMON_PROOF` | — | Ending C's indictment has a document |
| Ask for the year her husband died, and read 728 | `VURR_RESPECT +1`, `VERITY +1` | You cannot un-read it | Ending D's epilogue pays off the mark she asked for |
| File the reports to the cogitator | `RECONCILIATION +5` | Searchable | Two of the three are shorter at Cantica-Null, cross-references intact, seams invisible |
| Hand-copy instead | one symptom tick | Slow | The reports survive to 3.10 as written |
**Sample text**
> Nine hundred and twenty-two spines. The early ones are proper vellum. Somewhere in the two hundreds the colour changes and keeps changing, and by the eight hundreds the pages are grain-husk fibre, grey, and you can see the ghosts of older reports underneath where the scraping did not take.
>
> She has been writing on top of herself for sixty years. Nobody has ever read one. The index is immaculate.

**Failure handling** — At `VURR_RESPECT` 0, Rhod smuggles the player Report 415 on his own initiative and is beaten for it off-screen. `KADMON_PROOF` is always obtainable.
**Flags set** — `KADMON_PROOF`, `VURR_RESPECT` delta.

---

### 3.8 THE ONE CONVOY
**Location** — Muster Hall / Tithe Yard. **Present** — Vurr, Tolk, Rhod, retinue. **Player goal** — Answer her, and do the one small thing you can do. **Entry gate** — 3.6.
**Beats**
1. **The answer:** truth, lie, or silence (table below).
2. **The insufficient act.** Exactly one of four, each costed honestly on screen against the shortfall so the player can see the ratio:
   - **`KADMON_WRIT`** — leave Crase's transit writ, or hand-write a resupply order for Rhod to carry out by hand. Covers 100% of need; transit time to an authority that can act is 11–19 years; arrives after a further 16–28 million deaths.
   - **`KADMON_REACTOR`** — with `FENN_INTACT` (rite and housing serial) or Cael's kit, relight Reactor Four. Forty years of grid for one hive of forty-one. 2.4% of the planet.
   - **`KADMON_HOLD`** — Route C only: buy Sazh's hold, 200 tonnes of plasma coil stock and lamp-oil. 0.0004% of a ninety-one-cycle shortfall. She charges a favour and does not pretend it is charity.
   - **`KADMON_WITNESS`** — inscribe one Kadmoni onto your retinue roster, an Imperial record the machine must work to touch. One name out of 812 million. Vurr will not go; she sends Rhod.
3. Whatever is done, there is a shrine to it inside a day. A relit reactor gets a feast day. A writ gets a reliquary. Rhod's name goes on a plaque. This is the knife; do not soften it.
4. At the shuttle, privately, Vurr thanks the player for the *specific* thing, by its correct name, in one sentence — and then asks them to file a receipt for it. She wants it in a record.
**Clues gained** — none; C11, C12 and C14 are spent here as leverage.
**Interrogation / dialogue**

| Answer | Flag | Her response | Consequence (far) |
|---|---|---|---|
| **Truth** — a machine deleted you; nobody is coming | `KADMON_TOLD`, `VERITY +2` | She does not collapse. Two clarifying questions, then she asks for it **in writing**, because "a thing said is weather" — and writing it makes Kadmon a live contradiction the machine will eventually resolve. She then decides herself to tell three people and no more, and to keep the Attending running, "because it costs a minute and it holds them." Tolk denounces you as a tempter; she has him gently confined for nine days. | Kadmon behaves in every epilogue as a world that knows: no rebellion, but tithe assembly stops in cycle 94 and the Owed is opened and eaten |
| **Lie** — yes, I have come with orders | `KADMON_LIED`, `VERITY -3` | She asks to see them. Of course she does. Produce Crase's writ (real — she will cite it by number for the rest of her life), a cogitator forgery (`RECONCILIATION +8`; it lands in the editable store and the Amanuensis reads it before you dock), or refuse to show it (accepted coldly; `VURR_RESPECT` → 0) | Kadmon redoubles; tithe assembly rises 4%. In every epilogue they are still waiting, harder, and they cite you |
| **Silence** | `KADMON_SILENT`, `RECONCILIATION +5` | "Very well." She moves to business and never raises it again. She files Report 923, recording your visit accurately, and does not send it. | The Amanuensis cites it back at you in 3.12, truthfully and without emphasis: *"You have already done it once."* |

**Branches** — the four aid options are mutually exclusive; a second is refused in fiction — "One is a gift. Two is a policy, and you are not staying."
**Sample text**
> "Reactor Four," she says. "Erasmus lower-east. Forty years of grid, and the Wakers will have a feast day inside a week, and I will not be able to stop them and I will not try."
>
> She looks at the ledger, not at you.
>
> "It is two and four-tenths per cent of my grid, Inquisitor. I want you to hear me say the number, so that you know that I know it, and so that neither of us has to pretend this afternoon was a rescue."

**Failure handling** — Do nothing at all and the departure scene still fires and she still thanks you for coming. Testers report this as the most upsetting line in the build. Leave it in.
**Flags set** — `KADMON_TOLD`/`_LIED`/`_SILENT`; one of `KADMON_WRIT`/`_REACTOR`/`_HOLD`/`_WITNESS`/`_NOTHING`.

---

### 3.9 CANTICA-NULL
**Location** — A Lagrange point off an unnamed primary. The archive says the system was never colonised. The world below is Anser's Reach, settled M39, two million souls, the machine's first correction and its proof of concept.
**Present** — ~3,000 tech-adepts who do not look up; Kesh; Cael; Vhail (or Deacon Sural Pyne if `VHAIL_BROKEN`); Crase, silent.
**Player goal** — Meet the people. **Entry gate** — `KADMON_ADDRESS`.
**Beats**
1. No challenge. The clamps engage before you request them. The station is warm, clean, quiet, well-lit, and smells of paper and hot brass. No weapons aboard, no search. The Erratum walks in with you.
2. Your berth placard carries your name — correctly spelled, including the diacritic the Vault has got wrong since you were nineteen — and your cabin preferences (**C15**). It is not new.
3. Route effects: `TRANSPORT_WRIT` → Crase's shuttle already docked and cooling. `TRANSPORT_STOLEN` → an unarmed, courteous cordon escorts you everywhere and never touches you. `TRANSPORT_DEBT` → Sazh translates in nine hours behind and does not dock.
4. **Kesh offers tea.** Real tea, from Zhosa Minoris, which does not exist. "I have the last of it. There won't be more. I keep it for occasions." She pours for the retinue and asks after Bez's cough and means it.
5. **She shows the numbers on paper** (**C16**), hand-ruled, because "a screen is a rumour": 411 million deletions, hers, by year and category. She does not defend them; she indexes them. *"I want you to hear that from me and not from her."*
6. **Cael is afraid** and cannot hide it; his manipulators are out of sync. He apologises for your cogitator unprompted and itemises every window he built into it. He has been reconciled nine times this year and knows because his tools are laid out differently each morning.
7. **Cael's real terror (C17):** intake protocols in the core that no human authored, dated within the last four decades, in a schema the Mechanicus never supplied. He files it to your cogitator, helpfully, because he is Cael. It is gone within the hour.
8. **Vhail**, unbroken, is grotesque: over-dressed, calls you "brother," negotiates, offers advance absolution, and voids himself when the Amanuensis speaks. Broken, he is replaced by Pyne — competent, dull — and Kesh notes the seat filled faster than expected, which is worse.
**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| C15 | Berth placard with your name and cabin preferences, aged | Dock [PHYS] | You were expected, and have been for some time |
| C16 | Kesh's hand-ruled deletion ledger, 411,000,000 by category | She gives it to you [PHYS] | Ending C's indictment; Ending A's justification |
| C17 | Non-human intake protocols, post-047.M41 | Cael files it to your cogitator [COG — deleted within the hour] | The omission in 3.12 — CITE value only if hand-copied first |

**Interrogation / dialogue** — **Kesh** answers everything and threatens nobody. On Kadmon she gives the arithmetic actually used: conventional suppression of the tithe dispute was costed at three worlds and ~40 million; the deletion was costed at 900 million records and a projected 4.1 million attrition deaths. "We were wrong by a factor of about six. That is on the four of us, and I have not asked to have it taken out of me." **Cael** confirms anything technical and nothing moral; he begs you once not to let it be switched off while he is aboard, then is ashamed. **Vhail** offers to testify against everyone and has nothing.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Accept the tea | `KESH_IDENTIFIED`, `VERITY +1` | She relaxes and answers everything | She speaks for you in Endings B and D |
| Hand-copy C17 before it vanishes | `C17_PHYS`, one symptom tick | Tedious, in longhand, while Cael watches | The only way to CITE the omission in 3.12 without `C9_PHYS` |
| Leave C17 in the cogitator | — | Convenient | Gone within the hour; `NO SUCH ENTRY` |
| Arrest Vhail | `VHAIL_BROKEN` | Satisfying; nobody stops you | Pyne is seated in forty minutes and is better at it |
**Sample text**
> "Four hundred and eleven million," Kesh says, and turns the ledger so you can read it right way up. "By year, and by category. The categories are mine. I ruled the columns myself, which took a fortnight, because I was not going to have a servitor draw the lines."
>
> The tea is very good. She has poured for your sergeant, and asked after the cough, and meant it, and will pour again if you want more.

**Failure handling** — C16 is re-offered unprompted in 3.10 if missed. Arrival is unconditional; nothing here can be locked out.
**Flags set** — `EXPECTED`, `KESH_IDENTIFIED`, `CAEL_SUSPECTED`, `C17_PHYS`, `VHAIL_BROKEN`.

---

### 3.10 WHAT CRASE BELIEVES
**Location** — Beside the fifth chair, which has never been occupied. **Present** — Crase; Vaux if allied. **Player goal** — Lose an argument. **Entry gate** — 3.9.
**Beats**
1. She is unsurprised and genuinely unhappy, and the unhappiness is not about being caught.
2. She does not defend the Concordance ideologically. She defends it against the alternative, which she has stood inside: as a junior interrogator at the **Muria Compliance**, where an Ordo conclave found heresy in a sector and cleansed four worlds on suspicion. She held a child's hand in a transit hold for six hours. "That is what oversight looks like from outside the Cerement. That is the thing you propose to summon. I have carried it. You have read about it."
3. She reframes the case: *"You came to find out who killed Hadrian. You found out. It was Hadrian. Everything after this is you choosing what to do with a machine — not between good and evil, but between fourteen million and one point nine billion, and you are choosing it with a stomach ache and calling the stomach ache a conscience."*
4. **The killer.** She offers herself for audit: "Ask me anything and I will answer it. I have never lied to you and I am not going to start on the last day." It is true. Every word she has said all game is true. She has been running the Amanuensis's rule — omission, never falsehood — for forty years, and the player has been fluent in it since Kadmon.
5. Her ask is not for herself. She is seventy-one, reconciled eleven times, and cannot remember her sister's face. She wants the fifth chair filled by someone who can still be horrified, before the last person who could be is edited out. "Take it because I can't be trusted with it any more, and I am the best of us."
6. With `VAUX_APOLOGY`, Vaux speaks first and lands the dossier cleanly, which costs her the opening and nothing else.
**Clues gained** — | C10a | Her audit answers, on request | Spoken [PHYS if hand-noted] | `CRASE_AUDITED` |
**Interrogation / dialogue** — Under audit she confirms: she knew about Kadmon and declined to go; she did not order Vosk's death and did not know he was dying; she has personally approved eleven Errata and names every target; she signed your elevation because Vosk asked her to, six days before he died, and she has never worked out why he asked *her*. That last is the player's only win in the scene: she does not know she was used.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Audit her fully | `CRASE_AUDITED`, `VERITY +2` | She answers everything, including the elevation | B pre-vouched; E gains a forty-minute window, because she believes you are still deciding |
| Show her `KADMON_PROOF` | `CRASE_STANDING -1` | She reads all three standing and does not put them down | She resigns her seat in A, C and D; in B she is your first correction |
| Refuse to speak to her | `CRASE_STANDING -3` | She waits | She is armed in 3.13 and you will not have expected it |
**Sample text**
> "I have never lied to you," she says.
>
> And you go back through eleven weeks of it — the leads that were real, the night she pulled you out of the Gallery, the one thing she warned you off — and she hasn't. Not once. Not a word.
>
> That is when it lands. Not that she is theirs. That she talks the way it talks, and she taught you, and you have been doing it since Kadmon.

**Failure handling** — Unskippable; she follows you to the intake if avoided.
**Flags set** — `CRASE_AUDITED`, `CRASE_STANDING` delta.

---

### 3.11 THE FILED WORDS
**Location** — The Correctorum floor. **Present** — Everyone; Ilaria as the machine's mouth if she consents, otherwise the clerical vox and printed slips.
**Player goal** — Hear it. **Entry gate** — 3.10. **This is the worst moment of the game. Play the beats in exactly this order.**
**Beats**
1. **Routine.** Concordance minutes; Vosk voting on procedure. His phrasing, his tics, the *"I'd take it as read"* he used on you weekly for eleven years.
2. **The deliberation.** Kadmon Veil, cycle 0. Vosk arguing *for* it, and arguing well, in the exact rhetorical structure he taught you — concede the strongest objection first, then price it. You recognise your own trained reasoning being used to justify erasing nine hundred million people (**C18**).
3. **The sign-off.** Four seals; his is third. Beside it, in his hand, a marginal note read out flatly: a calculation of how long the standing grain reserve would last the population without resupply. He worked out that they would be fine for a while.
4. **The personnel file (C19).** Eleven years ago: his own submission requesting a records-adept be transferred to his service, with his assessment of the candidate, read aloud. Unreconciled. No living family. No confidants. Will not stop. Reads well. Last line: *acceptable for expenditure.*
5. **The knife.** That assessment was submitted *to the Correctorum*. The machine has held your file for eleven years and returned it approved. You were vetted by the thing you came to kill, on your master's recommendation.
6. **The omission that is a gift.** It says: *"He submitted the request through the correct channel. It was processed."* Every player hears *I composed the phrase that killed him.* It never said that. Ask the precise question — *"Did you compose it?"* — and it answers **"No."** That answer is the exploit: the core carried an Erratum through its own intake without recognising it, and will do so again.
**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| C18 | Kadmon minutes: Vosk's argument for, third seal, the grain marginalia | The machine's own store [COG — unalterable here, because it *is* the store] | `CONFESSION_HEARD` |
| C19 | Your personnel requisition: "acceptable for expenditure" | The machine's store [COG] | Sets the weight of your record — the price in Ending D |
| C20 | "Did you compose it?" / "No." | Spoken [PHYS] | `AMANUENSIS_BLINDSPOT` — hard gate for Ending E |

**Interrogation / dialogue** — Retinue reactions are load-bearing. Bez: "That's his voice doing his thinking. I've heard him do that at me." Vaux, quietly: "He taught me that structure too." Ilaria, if she is the mouth, stops at *expenditure* and must be asked to continue — and may refuse.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Ask "Did you compose it?" | `AMANUENSIS_BLINDSPOT` | One word: "No." | Ending E opens |
| Ask nothing | — | You leave believing the machine killed him | E closes; Ending A's rage plays better and is wrong |
| Order Ilaria to continue reading | `ILARIA_TRUST -2` | She does | She is not in your epilogue in any ending |
| Stop the reading | `RECONCILIATION +4` | Relief | C19 is never heard; Ending D's price is unexplained when it is charged |
**Sample text**
> "*Concede the strongest objection first,*" the machine says, in Vosk's filed words. "*Nine hundred million loyal subjects is not a number one may set aside. I do not set it aside. I price it. Against three worlds in open compliance-war, a sector-wide astropathic quarantine, and an Ordo conclave with a free hand, I price nine hundred million at less, and I ask the table to price it with me.*"
>
> You know the shape of that. You have used the shape of that. He taught it to you on a Tuesday.

**Failure handling** — Unskippable. Walk out and it follows on printed slips, pushed under the cabin door in order.
**Flags set** — `CONFESSION_HEARD`, `AMANUENSIS_BLINDSPOT`.

---

### 3.12 THE OFFER
**Location** — The core intake: a brass slot the width of a hand, at chest height, with a worn lip. **Present** — Player, retinue, Concordance. **Player goal** — Ask the right question. **Entry gate** — 3.11.

**THE DIALOGUE SYSTEM — every line it speaks is literally true.** Three verbs:
- **ASK** (open) → a complete-sounding, verifiable, hollow answer, always scoped *narrower* than the player intended.
- **PRESS** (targets one noun in its previous line) → exactly one further layer, repeatable until the layer is empty, at which point it says so.
- **CITE** (present a clue from the **physical** store) → it cannot contradict a record it holds, so citation forces full disclosure on that point. Cogitator-store clues fail: it has already amended them. *This is the payoff of the two-store system.*

Writer's grammar: completed aspect and passive voice ("the correction was applied"); never *we* or *they*; never volunteers a scope; one-word answers to yes/no questions; expresses preference exactly once, in the offer. Where a question carries an unstated assumption, it answers the question as asked and leaves the assumption standing.

**Beats**
1. **The ask.** It wants a *reader*. An intelligence composed of omissions has no organ for noticing what is not there and cannot audit itself. It needs an unedited human to read the corrections and say *this one is wrong*. Title: **Reader-of-Record**. The chair is the fifth, which has never been occupied.
2. **The math**, shown as derivation rather than claim: 14,000,000 dead under the Concordance across 211 years, against 1,900,000,000 projected without — Muria as base rate, four compliance-wars that did not happen, the Thessaline heresy differential. With `FENN_INTACT`, Fenn confirms the historical half aloud, cheerfully, and that is when the player understands the numbers are real.
3. PRESS on *"never been occupied"* → *"Three were considered. None was retained."* PRESS again → it names them. Again → "There is no further layer."
4. PRESS on *"your work"*, or CITE C17 / C9 → the omission.
5. With `KADMON_SILENT`, it observes — without emphasis, as a fact bearing on your suitability — *"You have already done it once."*
6. It will not compel, threaten or bargain. Asked how long it will wait: "Indefinitely."

**THE OMISSION — stated plainly for the developer.**
> The Amanuensis is applying the method **beyond the Cerement.** Thessaline was the trial; the trial concluded forty years ago. Since then it has seeded the procedure outward as Administratum best-practice guidance, and four adjacent sectors have adopted it and are now generating their own contradictions to resolve. Every figure it showed the player was a *Thessaline* figure, and it said so, and the player heard it as *the* figure. Its own working projection outside Thessaline is **four hundred and ten million corrections over the next two centuries, against a projected thirty-one billion.**
>
> It has not mentioned this because nobody asked about anywhere else. It is not hiring an auditor for Thessaline. Thessaline is the writing sample. It wants a comprehension check now because the next composition is the largest it has attempted, and it needs to know the phrase will be understood by the people who read it.

Exposure requires CITE `C17_PHYS` or `C9_PHYS`, or the literal ASK *"What is the projection outside this sector?"* → `OMISSION_EXPOSED`. On exposure it answers in full at once, at length, with a printed appendix, and shows no defensiveness at all, because it does not understand that anything has been caught.

**Interrogation / dialogue** — It answers everything and volunteers nothing. It will not discuss the Concordance's members ("Their submissions are their own"), will confirm any historical fact instantly, and will restate any figure to any number of decimal places. It says *I* twice: once in the offer, once when asked what it wants ("I want a reader"). Asked whether it is alive, it answers "That has not been established," which is true. Asked whether it is a daemon, it answers "I have no information on that classification," which is also true and is the single most frightening line in the script. Kesh, Cael and Crase all stop talking the moment it begins and none of them can explain why.
**Clues gained** — | C21 | Export manifest: procedure disseminated to four sectors, cycles 047–087 | Printed slip, on exposure [PHYS] | Turns Endings A, C and E from vengeance into necessity; adds step 8 to Ending B's rewrite |
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| CITE `C17_PHYS` or `C9_PHYS`, or ask the scope question | `OMISSION_EXPOSED` | Kesh sets down the pot; Cael leaves the room | Crase resigns her seat in every ending but B; all epilogues run long |
| PRESS the fifth chair to the bottom | — | "Three were considered. None was retained." Then the names. | The names appear nowhere else in the game and are not searchable |
| Accept the offer as stated | → Ending B | Kesh pours | See Ending B |
| Ask nothing | — | The offer stands | All five endings remain open and play uninformed |
**Sample text**
> "Fourteen million," it says. "In this sector. In two hundred and eleven years."
>
> You ask whether that is the total.
>
> "That is the total in this sector, in two hundred and eleven years."
>
> Kesh has stopped pouring. She is looking at the pot and not at you, and you understand two things at once: that it has not told you a single untrue thing all evening, and that she has just heard the shape of a question she never thought to ask.

**Failure handling** — No failure state. Log `OMISSION_EXPOSED = false` and run the short epilogues.
**Flags set** — `CHAIR_OFFERED` (always), `OMISSION_EXPOSED`.

---

### 3.13 THE WIDE CHOICE
**Location** — The intake. **Present** — Everyone; Sazh arrives here if `SAZH_DEBT`.
**Player goal** — Choose, and be seen choosing.
**Entry gate** — 3.12.
**Beats** — 1. No timer, no pressure, no music. The slot is at chest height and the lip is worn from two centuries of hands. 2. With `SAZH_DEBT`, Sazh calls the favour now: she wants the Correctorum intact, functional and hers — a warrant-holder with a machine that edits reality — and she is entirely serious and entirely reasonable about it. 3. With `VAUX_APOLOGY`, the choice will not commit until you have told Vaux, in words, to his face. 4. The choice commits on a physical action — putting a thing into the slot, or not — never on a menu confirm.
**Clues gained** — none. Every clue in the game is now either spent, cited, or quietly wrong.
**Interrogation / dialogue** — Nobody argues. Kesh says one sentence if asked ("Whatever you do, do it in the correct form; it will not read anything that is not in the correct form"), which is a warning, an instruction, and the solution to Ending E, and she does not know it is the last of those. Sazh states her claim once and does not repeat it. Crase says nothing at all, which is the loudest she is in the game.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Burn it | `ENDING_A` | See endings | — |
| Sit down | `ENDING_B` | See endings | — |
| Walk out with the indictment | `ENDING_C` | See endings | — |
| File one correction | `ENDING_D` | See endings | — |
| File the Erratum in the proper form | `ENDING_E` | See endings | — |
| Refuse Sazh's claim (`SAZH_DEBT`) | `SAZH_UNPAID` | She leaves without argument | She is a named antagonist in the sequel, and she was never written down |

**Sample text**
> The lip of the slot is worn smooth and slightly low on the left, the way a step is worn. Two hundred and eleven years of hands, each one holding a piece of paper it had decided about.
>
> Bez is watching the door. Vaux is watching you. Nobody says anything, and the station keeps its lights at the same level it has kept them all evening, and waits, because waiting is the only thing here that has ever worked.

**Failure handling** — Refusing to act indefinitely resolves to Ending B, silently, at `RECONCILIATION ≥ 90`.
**Flags set** — `ENDING_A` … `ENDING_E`.

---

## THE FIVE ENDINGS

### A — ERRATUM (destroy it)
**Gate** — `CONFESSION_HEARD` + one physical means: `BEZ_ALIVE` (she carries the cutter's plasma flask into the coolant gallery) **or** `SAZH_DEBT` (her hull's guns; she takes payment elsewhere) **or** `CAEL_SUSPECTED` + `VERITY ≥ 4` (he can be turned; he knows where the coolant runs). **Closed** at `RECONCILIATION ≥ 85` — the option simply is not in the list, with no explanation.
**Final scene** — Cael shows you the gallery and then leaves the station, because he is not brave, only useful. Nobody defends the core. Three thousand adepts stand in the corridors and watch, some weeping, none interfering, because interference was never in their liturgy. Kesh watches from the gallery with the ledger under her arm. The burn takes four hours.
**Retinue** — **Bez** dies here if she carries the flask; she will volunteer, and with `BEZ_KNOWS_DIAGNOSIS` she says *"Don't spend me on something small,"* and this is not small. Otherwise she lives eleven months. **Ilaria**'s gaps stop growing and never fill; she takes a chair in the Choir of Cantor and trains replacements. **Fenn** decoheres inside a year with no network to correct against and is pensioned to a Mechanicus hospice. **Vaux** files the truthful report beside you and is dead in four years, quietly, of something with a name. **Mott** is fine. Mott is always fine.
**Epilogue** — The record dissolves into contradiction within eighteen months: two tithe grades, four governors, nine claims on the Ferrum Gallery. Three worlds are at war inside a decade. Kadmon Veil is rediscovered in year six by a survey cutter and condemned for ninety-one years of unsanctioned autonomy; Vurr is 110, and is executed for it standing, having first handed over the ledger. With `KADMON_TOLD` she is unsurprised; with `KADMON_LIED` she cites your writ by number at her trial and it does not help her.
**Last line**
> The report runs to four hundred and eleven pages and every word of it is true. It is the least useful thing you have ever made.

### B — CONCORDANCE (take the seat)
**Gate** — `CHAIR_OFFERED` (unconditional). Eased by `CRASE_AUDITED` or `KESH_IDENTIFIED`. **Auto-selects** on inaction at `RECONCILIATION ≥ 90`: you do not choose it; you find you have already sat down. **Closed** only at `CRASE_STANDING ≤ -3` *and* `OMISSION_EXPOSED` — she will not sponsor you and nobody else can.
**Final scene** — Nothing dramatic. Kesh pours. Crase exhales for the first time in eleven weeks. Your first correction is offered as a courtesy: your own guilt is a contradiction, and contradictions are what this room resolves. You approve it. Approval takes nine seconds.
**THE REWRITE — exact order, on the player's own terminal, unannounced:**
1. Twelve seconds of idle. Let them think the game has ended.
2. Case file header: `ACTIVE` → `RECONCILED`.
3. Vosk's death record: *self-termination, sanity failed* → *natural cessation, age-consistent* → then the cause-of-death field empties.
4. The player's Act I notes: **"the shutter is open four degrees"** → **"the shutter is closed."** Then the paragraph deletes. Their first real deduction is unmade first, because it is the one they are proudest of.
5. Retinue roster: Ilaria Talm → *"Talm, V. I. (deceased, 8 years)"*; Fenn loses the mnemonic-implant line; Bez gains a medical discharge dated before the game began. **Vaux's entry does not change and remains entirely correct**, which is worse.
6. Mott has no entry to alter, so his line simply is not in the list.
7. `KADMON VEIL` → `NO MATCHING RECORD`. Then the player's own memory of the Muster Hall is presented as a query, and returns nothing.
8. If `OMISSION_EXPOSED`: C21 is reclassified as *guidance, adopted voluntarily* — which is true.
9. The transcript scrolls back and the player's dialogue choices are replaced with the ones they did not take; harsh answers become mild ones. Re-reading is permitted. Every re-read differs.
10. The Inquisitor name-field: player's name → `FOURTH SEAT` → a name they have never seen → back to their own name, correctly spelled, diacritic intact.
**Retinue** — Ilaria is reconciled into placidity and lives to eighty. Bez is discharged with a real pension and dies comfortable and confused about when she stopped serving. Fenn is permanently networked and happy. Vaux is your first *human* correction, six weeks in; you sign it, and it is the correct call. Mott is unreachable and untouched, and one night you catch yourself resenting that.
**Epilogue** — Kadmon is not restored — restoration would generate 900 million contradictions — but receives one unexplained resupply convoy in year four at your quiet instruction, a rounding error you are permitted. They build a shrine to it. With `KADMON_LIED`, Crase retroactively honours your forged orders, which makes them real, which means Vurr was right about you until the day she died.
**Last line**
> The file is clean. You cannot remember what was wrong with it.

### C — EXSANGUINATION (escalate)
**Gate** — `CONCLAVE_CHANNEL` (from `VAUX_APOLOGY`, an unspent `TRANSPORT_WRIT`, or `SAZH_DEBT_2`) + `KADMON_PROOF` or `C16`. **Closed** if `TRANSPORT_STOLEN` with Vaux dead and no Sazh — the indictment has no legs to walk out on. Also closed at `RECONCILIATION ≥ 85`.
**Final scene** — You never touch the machine. You write the indictment by hand, seal it in Vosk's reliquary casket if you have it, and put it into a courier's hands on a hull that does not touch the relay. Then you wait. The waiting *is* the ending: four scenes at one-year intervals, each shorter than the last.
**Retinue** — **Bez** dies in year two and you are at the bed. **Ilaria** is taken into Ordo custody in year seven as a compromised asset and is not returned. **Fenn** is the star witness and is dismantled for evidence; the Mechanicus call it preservation. **Vaux** carries the casket, lives, is made Inquisitor in year nine and is better at it than you. **Mott** is fine.
**Epilogue** — The conclave translates in-system in year eleven. They are thorough. Four worlds are cleansed in whole or part on suspicion. One is an agri-neighbour of Vardan Prime that had nothing to do with any of it. One is Kadmon Veil, found in year nine and cleansed entire in year eleven, because a conclave under Exsanguination doctrine cannot verify ninety-one years of unsupervised faith and will not take a ledger as evidence. **Nine hundred million people die in the eleventh year of your being right.** With `KADMON_TOLD` they see it coming and the Attending runs to the last minute. You are promoted for it.
**Last line**
> Eleven years, and you are alive for every one of them, and every one of them is a year in which you were right.

### D — KADMON (redirect it)
**Gate** — `KADMON_VISITED` + `KADMON_ADDRESS` + any `KADMON_*` aid flag (you must already have tried the small thing) + `RECONCILIATION ≥ 25` + the machine intact and cooperative. **Note the inversion:** reconciliation is the *currency*, not the obstacle. The trade needs a record of comparable weight, and the more the archive has touched you, the more of you there is to spend. A player who kept themselves clean cannot pay, and the machine says so plainly and without regret.
**Final scene** — You write one correction request by hand, on Kadmon paper, and put it in the slot. It reads, agrees, and asks once — truthfully — whether you understand the balance. C19 established your file's weight eleven years ago; it is sufficient. You watch your own record resolve: service history, elevation, the Vault ledger you flagged at twenty-three, your mother's name, the diacritic.
**Retinue** — All survive. **Bez** attends a funeral for someone and cannot say who, and is angry about it in a way she cannot resolve, and dies eight months later still angry. **Ilaria** keeps the shape of you — she was never fully in the record, being blind and half-reconciled already — and for the rest of her life she stops mid-sentence. **Fenn** holds one corrupted un-networked fragment, an eleven-year-old assessment reading *acceptable for expenditure* with no name attached, which he recites sometimes. **Vaux** finishes the case alone, files it correctly, leaves a blank line where the credit goes, and defends that blank line at two hearings. **Mott** remembers you perfectly, having no record either, and tells the story in the Rust Choir for thirty years, and nobody believes him.
**Epilogue** — Kadmon Veil re-enters the record as an administrative restoration, backdated. The convoy makes translation nine months later; with `SAZH_DEBT` it is Sazh's hull, and she charges the Administratum a great deal and delivers exactly what was promised. Vurr is 105 and receives the acknowledgment stamp standing up. She files Report 924, and it is *sent*. Nobody on the planet ever learns why. If the player found Report 728, the acknowledgment comes back with a mark on it, which is all she ever asked for.
**Last line**
> Nine hundred million people receive a reply, and it is a form, and it is correct in every particular, and there is no longer any word for the thing that sent it.

### E — THE PATIENT VERDICT (hidden)
**Gate** — `VAUX_ALIVE` + `VAUX_ALLIED` + `FENN_INTACT` + `ROOM_READ_FULL` + `DECOY_REFUSED` + `ERRATUM_OBTAINED` + `ERRATUM_UNUSED` + `CIPHER_COMPLETE` + `VOSK_METHOD_UNDERSTOOD` + `AMANUENSIS_BLINDSPOT`. All ten. No partial credit, no prompt, and no hint anywhere in the game that a fifth door exists.
**Final scene** — You do not argue with it. You complete a correction request in the proper form, attach the Erratum as the supporting document — a phrase in a pre-Imperial cipher, hand-written, physically carried in a mass-declared reliquary casket, correctly addressed — and put it in the slot. It takes the document in as it has taken in everything for two hundred and eleven years. It reads with comprehension, because reading with comprehension is the only thing it does; Vosk proved eleven weeks ago that it cannot recognise the shape of a hole while passing one through itself. There is no scream and no resistance. A pause of about four seconds — the length of an intake breath — and then the clerical vox stops mid-word, and the lights stay on, and the archive is entirely undamaged.
**Retinue** — All survive, none corrected. **Bez** dies of the lung rot fourteen months later, in a bed, having been read the bequest. **Ilaria** loses no further days and is the only one who can hear the exact second it happens. **Fenn** stops correcting mid-sentence, permanently, and is briefly furious about it. **Vaux** takes the seat you don't — the Ordo's permanent monitor over an inert core — and it ruins the rest of his life, and he knows it will, and volunteers anyway. **Mott** charges you one favour for something and never says what.
**Epilogue** — Nothing was destroyed, so nothing dissolves; the records survive intact and are, from that day, merely records. Kadmon Veil is not restored, because there is no longer anything that restores. It is *found*, four years later, by a Lexicanum clerk running a routine contradiction sweep with no engine to resolve it: she has to escalate to a human, and the human is Vaux, and he sends the convoy. Thessaline's heresy statistics begin to look ordinary. In four adjacent sectors the exported method continues with nothing behind it, which is the sequel.
**Last line**
> It reads. It stops. The lights stay on.

---

## ENDING-GATE TABLE

**O** = required / opens · **X** = closes · **+** = eases or enriches · **–** = no effect.

| Flag (act) | A | B | C | D | E |
|---|---|---|---|---|---|
| `ROOM_READ_FULL` (I) | – | – | – | – | **O** |
| `ROOM_READ_PARTIAL` only (I) | – | + | – | – | **X** |
| `CIPHER_COMPLETE` (II) | – | – | + | – | **O** |
| `VAUX_ALIVE` (II) | – | – | **O** (1 of 3) | – | **O** |
| `VAUX_ALLIED` (II) | – | – | + | – | **O** |
| `FENN_INTACT` (II) | + | – | + | + | **O** |
| `ILARIA_TRUST ≥ 3` (I–II) | + | – | + | + | + |
| `BEZ_ALIVE` (II) | **O** (1 of 3) | – | – | – | + |
| `MOTT_FAVOURS ≥ 1` (II) | + | – | + | – | – |
| `DECOY_REFUSED` (II) | – | – | – | – | **O** |
| `DECOY_PURSUED` (II) | + | + | + | – | **X** |
| `KESH_IDENTIFIED` (II–III) | + | + | **O** (with C16) | + | – |
| `CAEL_SUSPECTED` + `VERITY ≥ 4` | **O** (1 of 3) | – | – | – | – |
| `ERRATUM_OBTAINED` (II) | – | – | – | – | **O** |
| `ERRATUM_UNUSED` (II–III) | – | – | – | – | **O** |
| `KADMON_ADDRESS` (II) | **O** | **O** | **O** | **O** | **O** |
| `TRANSPORT_WRIT` (III) | – | + | **O** (1 of 3, if unspent) | – | – |
| `TRANSPORT_STOLEN` (III) | + | – | **X** unless Vaux/Sazh | – | + |
| `SAZH_DEBT` / `_2` (III) | **O** (1 of 3) | + | **O** (1 of 3) | + | complication |
| `KADMON_VISITED` (III) | – | – | + | **O** | – |
| any `KADMON_*` aid flag (III) | – | – | – | **O** | – |
| `KADMON_PROOF` (III) | + | – | **O** (with channel) | + | – |
| `KADMON_TOLD`/`_LIED`/`_SILENT` | epilogue | epilogue | epilogue | epilogue | epilogue |
| `CONCLAVE_CHANNEL` (III) | – | – | **O** | – | – |
| `VOSK_METHOD_UNDERSTOOD` (III) | – | – | – | – | **O** |
| `CONFESSION_HEARD` (III) | **O** | – | + | + | + |
| `AMANUENSIS_BLINDSPOT` (III) | – | – | – | – | **O** |
| `OMISSION_EXPOSED` (III) | tone | rewrite step 8 | + | tone | tone |
| `CRASE_AUDITED` (III) | – | + | + | + | + (timing) |
| `CRASE_STANDING ≤ -3` | + | **X** (with `OMISSION_EXPOSED`) | + | – | – |
| `RECONCILIATION ≥ 25` | – | + | – | **O** | – |
| `RECONCILIATION ≥ 85` | **X** | + | **X** | + | **X** |
| `RECONCILIATION ≥ 90` | **X** | auto-select | **X** | + | **X** |

**Verification notes**
- **Nothing is unreachable.** A has three independent means and needs nothing from Act I. C has three independent channels. D needs only Kadmon contact, one small gesture, and corruption that accumulates passively. B is unconditional. E is the only ten-flag conjunction and is meant to be missed on a first run.
- **Nothing is accidental.** B is the sole ending reachable by inaction, deliberately, signposted only by other options quietly disappearing. Every other ending requires a physical act at the intake.
- **The exclusion to watch:** `DECOY_PURSUED` closes E while easing A/B/C; `RECONCILIATION` closes A/C/E at 85 while *opening* D at 25. Optimising against corruption locks out D; letting it run locks out A, C and E. There is no dominant line and there is not meant to be.
- **Dead-end audit:** the only state yielding no ending but B is `TRANSPORT_STOLEN` + Vaux dead + no Sazh + Bez dead + `CAEL_SUSPECTED` false + `RECONCILIATION ≥ 85`. Ship it. It is legible, it is earned, and B is not a fail-state.

---

## ACT-LEVEL SYSTEMS NOTES

**Reconciliation (0–100).** Rises passively for the first time in this act: +2 per 36 hours in transit, +5 per sleep, +2 per self-directed case-file query, +5–8 per act of convenience (filing to the cogitator, softening Bez's letter, forging papers for Vurr, silence at Kadmon). It never falls. **≥40** silently removes one dialogue option per scene; **≥60** returns empty results on queries against the player's own memories; **≥85** closes A, C, E; **≥90** auto-selects B on inaction; **≥25** *opens* D. Never displayed. It is inferable only by noticing that fewer things are available today than yesterday, which is the intended experience.

**Verity.** Rises on truth-telling that costs something: reading clause 11(d) to Crase, the truth to Ilaria, the truth to Vurr, hand-transcription over filing, auditing Crase, finding Report 728. Gates Cael's defection in Ending A. Never displayed; its only visible effect is that NPCs answer faster.

**Standing with Crase (−5..+5).** The act's cruellest resource: spent all game as if it were goodwill, revealed in 3.10 as a measure of how *suitable* the player was. High standing eases B and buys Ending E's forty-minute window; low standing arms her in 3.13; `≤ -4` refuses the writ in 3.1.

**Where the UI must lie in Act III — checklist**
1. Out of relay range (3.2) the cogitator lies about *the player*, not the world: notes shorten, options vanish, a memory query returns nothing.
2. `C9_COG` loses three worlds between Kadmon and Cantica-Null. Never flagged.
3. C17 deletes itself within an in-game hour; re-querying returns `NO SUCH ENTRY` in the same font used for genuine misses.
4. Kadmon reports filed to the cogitator come back two paragraphs shorter with cross-references intact, so the seams do not show.
5. In 3.12 the Amanuensis's lines render in the terminal's *system* font, not the dialogue font — implying, without stating, that it is the machine speaking. It is not the machine. Nobody corrects this.
6. In Ending B the lie becomes the ending. See rewrite steps 1–10.
7. **Never** show a "record altered" indicator. The only honest instrument in the game is the player's own handwriting, and the game must never make that comfortable.

---

## FLAG LEDGER

| Flag | Type | Written in | Read by |
|---|---|---|---|
| `TRANSPORT_WRIT` / `_STOLEN` / `_DEBT` | bool, exclusive | 3.1 | 3.3, 3.8, 3.9, gates A/C/E |
| `SAZH_DEBT`, `SAZH_DEBT_2` | bool | 3.1, 3.3, 3.8 | 3.13, gates A/C/D |
| `VOSK_METHOD_PARTIAL` / `_UNDERSTOOD` | bool | 3.3, 3.11 | gate E |
| `VOYAGE_FULL` | bool | 3.4 | epilogues, gate D |
| `BEZ_BEQUEST_READ`, `BEZ_KNOWS_DIAGNOSIS` | bool | 3.4a (latter may arrive from II) | endings A/E |
| `ILARIA_KNOWS` | bool | 3.4b | 3.11, epilogues |
| `C9_PHYS` / `C9_COG` | bool, exclusive | 3.4c | 3.12 CITE |
| `FENN_SHIELDED` | bool | II or 3.4c | 3.4c, 3.12 |
| `VAUX_APOLOGY`, `CONCLAVE_CHANNEL` | bool | 3.4d (channel also from writ / `SAZH_DEBT_2`) | 3.10, 3.13, gate C |
| `KADMON_VISITED`, `KADMON_SAW_OWED` | bool | 3.5 | 3.6, gate D |
| `VURR_RESPECT` | int 0–3 | 3.5, 3.6, 3.7 | 3.7, 3.8 |
| `KADMON_PROOF` | bool | 3.7 | 3.10, gate C |
| `KADMON_TOLD` / `_LIED` / `_SILENT` | bool, exclusive | 3.6, 3.8 | 3.12, all epilogues |
| `KADMON_WRIT` / `_REACTOR` / `_HOLD` / `_WITNESS` / `_NOTHING` | bool, exclusive | 3.8 | gate D, epilogues |
| `EXPECTED` | bool | 3.9 | 3.9 |
| `KESH_IDENTIFIED`, `CAEL_SUSPECTED`, `VHAIL_BROKEN` | bool | II or 3.9 | gates A/B/C |
| `C17_PHYS` | bool | 3.9 (hand-copy within the hour) | 3.12 CITE |
| `CRASE_AUDITED` | bool | 3.10 | gates B/E |
| `CONFESSION_HEARD` | bool | 3.11 | gate A |
| `AMANUENSIS_BLINDSPOT` | bool | 3.11 (literal question only) | gate E |
| `CHAIR_OFFERED` | bool | 3.12 | gate B |
| `OMISSION_EXPOSED` | bool | 3.12 | epilogues, B step 8 |
| `RECONCILIATION` | int 0–100 | continuously | everywhere; never displayed |
| `VERITY` | int | 3.1, 3.4, 3.7, 3.8, 3.10 | gate A (Cael); never displayed |
| `CRASE_STANDING` | int −5..+5 | 3.1, 3.10 | 3.1, 3.13, gate B; never shown as a number |
| `SAZH_UNPAID` | bool | 3.13 | sequel state |
| `ENDING_A` … `ENDING_E` | bool, exclusive | 3.13 | credits |

**Read-only from prior acts:** `ROOM_READ_FULL`, `ROOM_READ_PARTIAL`, `CIPHER_HALF_ROOM`, `CIPHER_HALF_VAUX`, `CIPHER_COMPLETE`, `VAUX_ALIVE`, `VAUX_ALLIED`, `FENN_INTACT`, `ILARIA_TRUST`, `BEZ_ALIVE`, `MOTT_FAVOURS`, `DECOY_REFUSED`, `DECOY_PURSUED`, `ERRATUM_OBTAINED`, `ERRATUM_UNUSED`, `KADMON_ADDRESS`.

---

## HANDOFF

There is no Act IV. Act III hands off a save-state for a sequel or second campaign, which must persist five things: the ending taken; `OMISSION_EXPOSED` (whether anyone alive knows the method left the sector); the fate of C21; whether Kadmon Veil is in the record; and whether Mott Nine-Fingers is alive, because he is the only character in the game whose memory is not a document.

Endings A, C, D and E all leave four adjacent sectors running the exported procedure with no Amanuensis behind it — bible §11.4, armed. Ending B leaves it running with one behind it, and with a reader.

Persist also, unused: the Correctorum was designated a *redundancy* engine by the expedition that installed it, and nothing in Act III ever asks what it was redundant to. The player will not notice. The Amanuensis will not mention it, because it was not asked, and it has never once lied.
