# ACT I — THE LOCKED ROOM

**Dramatic question:** *Was it suicide?* **Answer:** *Yes. And that's worse.*

Ninety minutes arguing with the player's own competence. The verdict is written before the player sees the body; the room says otherwise — not in evidence but in *arrangement*. Five objects placed by a man with six days left and no channel he trusted. Reading them proves Vosk chose the hour of his death, which converts despair into scheduling and scheduling into a message. Then the act takes one reading away: a real record, read once, unmade, with the player's own case-notes retroactively agreeing they never asked for it. Act I does not prove a conspiracy. It proves the player's instruments are not theirs.

## SCENE INDEX

| # | Scene | Location | Type | Gated by | Est. minutes |
|---|---|---|---|---|---|
| 1.1 | SIX HOURS | Vosk's antechamber, Adamant Crown | Cold open | — | 6 |
| 1.2 | THE ROSETTE | Antechamber | Elevation / first branch | 1.1 | 7 |
| 1.3 | THE ROOM | Vosk's chambers | Free-form examination | `ROSETTE_TAKEN` | 22 |
| 1.4 | STANDING WHERE HE STOOD | Vosk's chambers | Synthesis | ≥3 signs held | 8 |
| 1.5 | THE CORRECT PROCEDURE | Vaux's study, tier 4 | Interview | `ROSETTE_TAKEN` | 10 |
| 1.6 | THE DOOR AND THE BODY | Medicae deck, tier 1 | Witness walk + autopsy | `ROSETTE_TAKEN` | 12 |
| 1.7 | NOTHING TO REQUISITION | Adjutant's office | Interview / systems plant | `ROSETTE_TAKEN` | 7 |
| 1.8 | FORTY-ONE SECONDS | Choir Spire, cell 4 | Interview | A2 **or** `CRASE_LEAD_CHOIR` | 12 |
| 1.9 | THE REQUISITION | Terminal | Archive query | `A1_FILED` or `DISPATCH_ID_KNOWN` | 6 |
| 1.10 | ONE OF ONE, NONE OF NONE | Terminal | The Turn / act close | 1.9 + 1h, or fallback | 10 |

---

### 1.1 SIX HOURS
**Location** — Antechamber of Vosk's chambers, Adamant Crown. Three wardens who don't know whether to salute you.
**Present** — Inquisitor Yolanda Crase; Sgt Bez Oduya (on the doorframe, pretending to be fine); wardens; the body, visible through the inner door.
**Player goal** — Get into the room. You cannot yet.
**Entry gate** — Game start.
**Beats**
1. You are summoned by a runner, not a vox. Nobody in the Crown is putting this on the relay.
2. Crase blocks the doorway kindly. She has been here five hours and has already been gentle with Bez.
3. The verdict is read to you before you see anything, countersigned by Crase, timed **two hours before the body was found**. A clerical artefact of a pre-generated form; she explains it truthfully. Included so the player's first suspicion of the record is *wrong*.
4. You see the body from eleven metres. Autopistol in the right hand. From here you cannot see which foot.
5. Crase asks one question — not *what do you think* but *did he seem tired to you* — and writes your answer in a paper notebook. Note the paper.
6. Bez says the door was sealed and she broke it. She says it twice, because she needs to have been heard saying it.

**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| A0 | Verdict timestamp precedes discovery by 2h | Crase's writ, read at distance | Pressure line in 1.2. Lands in whichever store the player picks — the game's first silent notebook test |

**Interrogation / dialogue** — She gives: sealed from inside, 90mm vents, one aperture-log entry in eleven hours (Bez, forced). She won't volunteer the other countersignature; asked, she names the Provost-Adjunct at once. What opens her is grief — she liked him. Say so and the Choir lead arrives a scene early.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Challenge the timestamp | `TIMESTAMP_NOTICED` | Explained; Provost-Adjunct embarrassed; `CRASE_STANDING` +1 | Act III: she cites this as when she decided you were suitable |
| "He seemed tired" | `SAID_TIRED` | Enters the file as corroboration of suicide | Act II: Vaux reads it and uses it |
| Say nothing | — | *Declined to characterise* | Neutral |

**Sample text**
> Five hours ago somebody decided what this was, and the decision went into a form, and the form has been countersigned, and now it is simply a thing that happened.
>
> Crase steps sideways to let you look and not to let you in. From the doorway you can see one hand, the pistol still in it, the fingers loose the way fingers go loose. You have been in this room a hundred times. You cannot remember what colour the floor is.
>
> "Did he seem tired to you," she says. It is not an accusation. That is the trouble with her.

**Failure handling** — Unmissable. Total silence reads as shock and advances after ninety seconds.
**Flags set** — `ACT1_OPEN`, opt. `TIMESTAMP_NOTICED`, `SAID_TIRED`.

---

### 1.2 THE ROSETTE
**Location** — Antechamber. Witnesses assembled: deliberate.
**Present** — Crase; Vaux (late, at the back); Bez; the Crown Provost-Adjunct; two recording adepts.
**Player goal** — Accept. There is no refusing.
**Entry gate** — 1.1 complete.
**Beats**
1. Crase produces the rosette from a cloth, not a case. She has carried it against her ribs for five hours and apologises for that.
2. The elevation takes forty seconds, spoken plainly, witnessed, written down as it happens. Vaux watches the styluses, not you.
3. Three things change hands: the rosette, Vosk's cogitator authority-key, the retinue — *"and I'm sorry about all three."*
4. **Her one real lead:** the Choir of Cantor keeps **wax-cylinder redundancy recordings** of every transcription — physical, un-networked, held on the Spire, never in the Vault. She tells you because she cannot imagine you getting anywhere with them. It is the most useful sentence anyone says in Act I.
5. **Her one warning:** do not requisition Lord-Prefect Kesh's countersignature chain, and do not pull Vosk's personnel file, in your first week. *"Let me ask her. She'll answer me and she'll bury you."* Good advice. Also the exact door she needs closed.
6. Vaux exits with one sentence — *"Congratulations, Interrogator"* — the wrong rank, deliberately, in front of the recording adepts.

**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| A1-lead | Choir keeps physical wax redundancies | Crase, spoken | 1.8 access. Verbal — notebook only unless transcribed. Free proof that some things never enter the cogitator at all |

**Interrogation / dialogue** — She'll discuss his mood, three boring Administratum debts, four enemies of whom two are dead. She won't discuss his last cycle's work: *"nothing, he'd stopped taking cases."* She believes it — she has been reconciled, and from inside her head it is true.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Ask why not Vaux | `ASKED_WHY_NOT_VAUX` | Honest answer — *"because he'd have been very good at it"*; `CRASE_STANDING` +1 | Opens the honest branch in 1.5 |
| Accept the Kesh warning aloud | `KESH_WARNED_OFF` | She relaxes visibly | Act II: breaking it costs `CRASE_STANDING` −2, gains the countersign chain |
| Reject it aloud | `KESH_DEFIED` | `CRASE_STANDING` −1 | She starts logging you three weeks early |

**Sample text**
> The rosette is warmer than it should be. She has had it against her ribs since morning and she hands it over the way you'd hand somebody a bird.
>
> "There's a form," she says. "There's always a form. I've signed the part that matters."
>
> Behind her the adepts are writing. Everything in this room is being written. Vaux stands at the back with his hands folded, watching the styluses, and on his way out he calls you Interrogator, and two of the four people present notice, and one of them is you.

**Failure handling** — Cannot fail. If the player skips all optional lines, Crase's aide delivers the Choir lead as a written note within thirty minutes. The game will not lose it.
**Flags set** — `ROSETTE_TAKEN`, `CRASE_LEAD_CHOIR`, opt. `KESH_WARNED_OFF`/`KESH_DEFIED`, `ASKED_WHY_NOT_VAUX`.

---

### 1.3 THE ROOM
**Location** — Vosk's chambers. Body present for eight minutes, then removed by medicae on a timer the player cannot stop — apologetic, immovable, diegetic. Evidence expires.
**Present** — The player alone, by request; Crase grants it instantly, which is a kindness. Bez outside. Fenn if summoned.
**Player goal** — Find the five signs.
**Entry gate** — `ROSETTE_TAKEN`.

**THE FIVE SIGNS.** Findable in any order, none gating another. Three read-depths each. Filing a sign is what makes it actionable and what makes it editable.

**SIGN 1 — THE REGICIDE BOARD (A1).** *Careless:* an unfinished game. Vosk played; everyone knew. *Careful:* three pieces cannot legally stand where they stand — a **second white Emperor** on **D2**; a **black Lexicanum**, a dark-square-bound piece, on the light square **F4**; a **white Sentinel** on **C9**, the far rank, unpromoted, where promotion is compulsory. *Meaning:* Lexicanum Vault shelf-addresses and Cantor regicide notation are **the same alphabet** — both codified by the same Administratum office, both addressing a nine-by-nine grid over nine levels. A Vault address is three file-rank pairs: block-and-level, aisle-run, shelf. A board with exactly three anomalous pieces is an address written in furniture, which no cogitator will ever index. Read in piece-hierarchy order (Emperor, Lexicanum, Sentinel) it is **D2–F4–C9** — real, and wrong: Munitorum drainage requisition ledgers, two centuries of them. Read **backwards**, as SIGN 3 instructs, it is **C9–F4–D2**: Block C, Level 9, the Choir of Cantor hardcopy dispatch manifests, current cycle. The physical ledger of every message the Choir received, including the one that killed him. *Unlocks:* 1.9. **Store: cogitator required** — you cannot requisition against a code you only wrote by hand. The game forces the exposure.

**SIGN 2 — THE SHUTTER (A2).** *Careless:* a window ajar in a sealed room. It is 90mm; nobody came through it. *Careful:* open **four degrees** and *fixed* — the ratchet shimmed with a folded requisition slip. Vosk never opened it in eleven years. Two readings, both needed. **Light in:** at the twelfth hour the blade crosses the desk and ends at the **third drawer**, which is empty and dusted. That drawer held the regicide set. He is telling you the board is out on the table *on purpose*. **Sightline out:** at the twelfth hour, and only then, the aperture frames one thing across the Crown gap — Transit Gallery 9 of the Choir Spire, and in it the shuttered window of **transcription cell 4**. Ilaria Talm's cell. *Meaning:* a time and a place. Twelfth hour. Cell four. *Unlocks:* 1.8 at usable resolution — without A2 the Choir will only process a full-shift requisition, four hundred cylinders, eleven days, and Act I ends first. Store: either.

**SIGN 3 — THE BOOTS (A3).** *Careless:* wrong feet; the man was distressed. A rushing player files it as such and actively strengthens the case against themselves. *Careful:* both boots **laced evenly, tension matched, tongues straight**; no sock displacement; the bootjack used and returned to the wall in the wrong orientation for a right-handed man in a hurry. He sat down. He took his time. He did this last. *Meaning:* **Sarrow's Reach, eleven years ago.** Your first field posting, still a records-adept, across a table from the Governor's people with a manifest in front of you that you both knew was forged. Vosk could not speak. At the rest-halt he came back with his boots swapped and held your eye exactly as long as it took. *Read it backwards.* You never discussed it. It has never been written down by anyone, which is precisely why it is the only language he had left. *Payoffs:* immediately reverses A1 into the live shelf-code; in Act II it is the insight that Vosk's cipher runs right-to-left, and half of `CIPHER_HALF_ROOM`. **Store: notebook only.** The cogitator offers a filing template whose classification options are *distress / disorder / self-neglect*. There is no field for what it is.

**SIGN 4 — THE THUMB AND THE BLANK PAGE (A4).** *Careless:* ink, a blank sheet, ash in the censer. He wrote something and burned it. Dead end. *Careful:* the ink is **auto-quill** ink, not his pen ink — Vosk wrote longhand out of vanity and used the machine deliberately. The sheet on the blotter is the **second sheet**, the impression platen: illegible, but it carries a strike-count of over nine hundred glyphs on one page where his longhand fits two hundred. He wrote compressed, in a notation that isn't Gothic. And the quill is a Cael-pattern unit whose **pressure-drum tension log survives the burning of the paper**, giving a duration and a stop-time: written, and the censer lit, **four minutes before the estimated hour of death**. *Meaning:* he knew the minute. A man in despair does not schedule. Every other object in the room is retroactively converted from *sad* to *arranged*. *Unlocks:* Act II ash-recovery via Cael's equipment. Store: either — but the drum log extracts only by filing the quill, the first time the player must hand the enemy something to learn something.

**SIGN 5 — SIX DAYS (A5).** *Careless:* a date circled six days out. An appointment. *Careful:* circled with the **auto-quill**, not his pen — same sitting as the burned log; circled **twice**, one ring inside the other, pressed hard; and under the sheet, on the leaves for previous cycles, the ghosts of erased circles, same weekday, **every eighty-four days**, back as far as the pad goes. There is no corresponding entry anywhere — not the diary, the adjutant's book, the cogitator, or the Crown's audience schedule. An eighty-one-year standing appointment that has never once generated a record. *Meaning:* his next reconciliation. **The player cannot learn this in Act I and must not.** What they can learn is the shape of the wrongness: a thing that happens on a schedule and leaves no trace, in a hive built entirely for the making of traces. *Unlocks:* the Fenn beat in 1.7; in Act III it dates Vosk's defection — the second ring is the one he wasn't going to keep. **Store: notebook strongly incentivised.** Filed, the ghost circles are among the first things quietly gone from the scanned image in Act II.

**Beats**
1. Eight-minute body timer. Boots, thumb and the pistol's grip-angle are readable only while it is there.
2. Free movement. No highlights, no glints, no prompts. Objects are examined by naming them at the terminal.
3. The aperture log shows one entry in eleven hours. The **seal-integrity register** shows a second event: a warden's seal broken and re-set forty minutes before Bez arrived (A7). Provosts get this free; others ask Bez.
4. **First filing decision.** The cogitator offers to file, with a suggested classification. Accepting is fast, indexes the clue, makes it requisition-eligible. Declining writes it to the notebook, inert and incorruptible.
5. Vosk's own notebook sits on the desk — eleven years of longhand in a room with a cogitator that would have been faster. The last page is a shopping list. That's the tell.

**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| A1 | Three illegal pieces → Vault address | Side table | 1.9. **Cogitator required** |
| A2 | Shutter at 4°: twelfth hour, cell 4, third drawer | Window / desk | 1.8 at usable resolution. Either |
| A3 | Boots deliberately swapped — Sarrow's Reach code | Body (timed) | Reverses A1; half of `CIPHER_HALF_ROOM`. **Notebook only** |
| A4 | Auto-quill entry written + burned 4 min pre-death | Thumb, platen, censer, drum log | Act II ash-recovery. Either (drum log needs filing) |
| A5 | Six days circled; 84-day ghost cycle; no record | Desk calendar | 1.7; Act III dating. Notebook incentivised |
| A6 | Vosk never trusted the cogitator | His longhand notebook | Two-store tutorial. **Physical item** |
| A7 | Seal broken and re-set 40 min before Bez forced it | Seal-integrity register | Vaux placement in 1.5 |

**Interrogation / dialogue** — None. The only scene in the game with no one to ask.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| File A1 | `A1_FILED` | Requisition possible | It is this requisition the archive edits in 1.10 |
| File A3 as *distress* | `A3_MISFILED` | Strengthens the suicide verdict in your own file | Act II: Vaux quotes your filing at you; Vaux's cipher half costs more |
| Keep A5 in notebook | `A5_PHYSICAL` | Nothing visible | Act II: ghost circles vanish from the filed image; you can prove it |
| Take Vosk's notebook | `HAS_VOSK_NOTEBOOK` | Physical store established | Act III: the only surviving handwriting sample |

**Sample text**
> The shutter is open four degrees. You know this because Vosk never opened it — the man complained about Cantor's light for the eleven years you knew him, said it came off the manufactorum stacks the colour of old teeth.
>
> Four degrees. At the twelfth hour it would put a line across the desk, and the line would end at the third drawer, and the third drawer is empty and has been dusted.
>
> He is telling you something. He is being very patient about it, the way he was patient about everything, and he has been dead for six hours.

**Failure handling** — No sign is individually required; the act proceeds on any subset. Two re-entries are permitted before the Crown reclaims the chambers at act close; after the body goes, A3 degrades to a medicae photograph showing lacing but not tension. No soft-lock, only a shallower Act II.
**Flags set** — Per-clue `A#_HELD` / `A#_FILED`; `HAS_VOSK_NOTEBOOK`.

---

### 1.4 STANDING WHERE HE STOOD
**Location** — The chambers, after the body is gone.
**Present** — Player alone; Bez may be invited and mostly doesn't speak.
**Player goal** — Combine.
**Entry gate** — **Three or more** of A1–A5 held, either store.
**Beats**
1. A terminal verb appears with no tutorial and no explanation: `READ ROOM`. It has existed since the game started, greyed out, unremarked.
2. The player must *order* the held signs. Real input, not a cutscene.
3. Any ordering produces a reading. The **correct** one requires A3 in hand — inversion is the only thing that turns a drainage ledger into the dispatch manifests.
4. Three or four signs → `ROOM_READ_PARTIAL`: the Vault address in its **forward** form, the twelfth hour, the four-minute window. Enough to run 1.8–1.10. Not enough for the cipher.
5. All five plus A3 → `ROOM_READ_FULL`: the room resolves into one sentence of intent — a time, a place, a witness, a shelf, a deadline — and grants `CIPHER_HALF_ROOM`.
6. Either way the room stops being a crime scene and becomes correspondence. Hold on silence before offering the next scene.

**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| A8 | The reading: staged in advance, for a specific reader | Synthesis | `ROOM_READ_PARTIAL`/`FULL`; on FULL, `CIPHER_HALF_ROOM`. Store = wherever composed; filing it is the most dangerous act available in Act I and the only way to get Crase's institutional help in Act II. There is no right answer; log the choice |

**Interrogation / dialogue** — Bez, if present, asks what you're doing. Told the truth, she says: *"He did that with his boots once. On the ship."* She saw it at Sarrow's Reach and never knew what it meant. Confirms A3 without decoding it. `BEZ_TRUST` +1.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Correct order, all five | `ROOM_READ_FULL`, `CIPHER_HALF_ROOM` | Live shelf-code | **Ending E stays open** |
| Any order, 3–4 signs | `ROOM_READ_PARTIAL` | Forward (wrong) code | Ending E closed; Act II cipher needs Vaux alive AND allied AND a Vault break |
| File A8 | `A8_FILED` | Crase can act on it; clerks cooperate | The Concordance knows what you know, three weeks early |

**Sample text**
> Three of them and the room starts to hum. Not literally. It is only that the objects stop being objects.
>
> The board is an address. The shutter is a clock and a finger. The boots are a sentence you learned in a rest-halt on Sarrow's Reach when you were twenty-six and too frightened to eat.
>
> He had six days. No vox he trusted, no page that wasn't read, no confidant the record couldn't map. So he put it in furniture, and he put the key to the furniture in a thing that exists only between you and a dead man.

**Failure handling** — Below three signs at the end of 1.7, Bez volunteers the boots anecdote unprompted at the medicae deck, granting A3 at reduced fidelity and pushing most players over the line. Otherwise 1.9 is entered on Ilaria's dispatch ID and the act completes at `ROOM_READ_NONE`.
**Flags set** — `ROOM_READ_FULL` | `_PARTIAL` | `_NONE`; `CIPHER_HALF_ROOM`; `BEZ_TRUST`.

---

### 1.5 THE CORRECT PROCEDURE
**Location** — Vaux's study, tier 4. Forty years of case-files, indexed, none of them on a cogitator.
**Present** — Interrogator Dmetri Vaux.
**Player goal** — Establish whether he was in the room and what he took.
**Entry gate** — `ROSETTE_TAKEN`.
**Beats**
1. He does not stand. He uses every correct form of address, perfectly, with the emphasis very slightly on *correct*.
2. **His theory:** Crase had Vosk killed and you are her instrument. Wrong about the agent — no human did this — and closer to right about Crase than anyone alive. He won't say why, or what he's been doing for two years.
3. Confronted with A7 he admits the seal at once, with contempt: came to see the body, found it sealed, resealed what he'd broken, took nothing. All true.
4. **Plant one:** he never asks *how*. Forty years an apprentice and he does not ask how his master died. He asks *who signed the verdict*, twice.
5. **Plant two:** under real pressure he produces a six-syllable nonsense phrase in a fixed rhythm, learned by rote, apparently without noticing he said it aloud. It is a segment of Vosk's cipher key. Log as A9 whether or not the player notices. He will undercut you with Crase before the day ends regardless; the note is already drafted.

**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| A7 | Vaux entered the corridor and reset the seal | His admission | Removes him as a locked-room impossibility |
| A9 | Six-syllable mnemonic fragment | Under pressure | Act II `CIPHER_HALF_VAUX`. **Notebook only** — the cogitator has no field for a sound |
| A10 | He believes Crase is responsible | Dialogue | Act II alliance branch |

**Interrogation / dialogue** — He'll discuss procedure, the aperture log, Vosk's declining caseload, his own seniority. Not: his own investigation, being passed over, or Crase's name in a sentence containing *evidence*. What opens him is not authority — he's ignored authority for forty years — but **accuracy**: cite correctly a procedural detail he expects you to fumble and he visibly recalibrates. With `ASKED_WHY_NOT_VAUX`, the only warm thing he says all act: *"He talked about you. Not often. Not much. But he did, and he never talked about me."*
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Detain him for the seal | `VAUX_DETAINED` | Four hours, nothing gained, `CRASE_STANDING` −1 | `VAUX_ALLIED` far more expensive in Act II |
| Show him A3 or A8 | `VAUX_SHOWN_ROOM` | He goes very quiet and asks to see the room | He brings his cipher half unprompted; largest single step to `VAUX_ALLIED` |
| Record the mnemonic | `A9_HELD` | — | Act II decryption survives his death |

**Sample text**
> "You will want to establish that I was in the corridor," Vaux says. "I was in the corridor. I broke a warden's seal at the eighth hour and re-set it, which is an offence under Crown statute, and which I will sign for now if it moves us along."
>
> He has a pen out before you answer. That is the thing about him. He is never the one wasting your time.
>
> "You haven't asked me how he died," you say.
>
> A pause of exactly the wrong length. "No," he says. "Who signed the verdict?"

**Failure handling** — Fully skippable. Skipped, A9 exists only in Vaux's head, making `VAUX_ALIVE` a hard requirement for the Act II cipher. The game does not warn about this.
**Flags set** — `VAUX_MET`, `VAUX_ALIVE`, `VAUX_DETAINED`|`_UNCHALLENGED`|`_SHOWN_ROOM`, `A9_HELD`, `VAUX_RESPECT`.

---

### 1.6 THE DOOR AND THE BODY
**Location** — Chambers corridor, then the medicae deck, tier 1.
**Present** — Sgt Bez Oduya; Medicae-Adept Suryn Ostro (sixty, bored, honest, four hundred Crown autopsies and never one like this).
**Player goal** — Eliminate every mechanism of death, and fail.
**Entry gate** — `ROSETTE_TAKEN`.
**Beats**
1. Bez walks you through the door physically. Eleven hours. One knock at the sixth, no answer, and she waited five more *because he told her to that morning*, and she is holding that like a stone.
2. **Her theory:** somebody got in and out and she missed it. She is hunting a route because the alternative is that Vosk arranged to be alone and used her to do it. She is wrong, and the truth is worse for her than the theory.
3. **Plant one:** she asks you to read the duty-roster to her. *"Light's bad down here."* The light is fine. She cannot read, and twenty-six years of service have been organised downstream of that.
4. **Plant two:** she coughs three times; the third time she turns away, coughs into her hand, and puts the hand in her pocket. A sealed Inquisitorial bequest-envelope sits in her kit, which she calls *"paperwork, sir, I'll get to it."*
5. **Autopsy**, fast and without ceremony. No toxin. No coercion residue. No psychic trace — confirmed by the Choir's own auspex, which is the part Ostro dislikes. No wound but the pistol's, and **the pistol was fired eleven minutes after cardiac arrest**. A three-hundred-year-old heart in the condition of a man of ninety, in excellent order, stopped. Ostro's summary, repeated verbatim to anyone who asks: *"He died and then he shot himself. In that order. I've written self-termination because there is no other box."*

**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| A11 | Pistol discharged 11 min post-mortem | Autopsy | Destroys the physical suicide mechanism |
| A12 | No toxin, coercion or psychic trace; the heart simply stopped | Autopsy | Act II Erratum thread |
| A13 | Vosk ordered Bez away that morning | Bez | With A4, confirms scheduling |
| A14 | Bez cannot read; sealed bequest unopened | Observation, not dialogue | Act II/III Bez arc. **Notebook only — filing it is a betrayal and the game treats it as one** |

**Interrogation / dialogue** — She discusses the door, the timings, Vosk's habits, what he ate. Not her health; asked directly she lies badly and you can watch her decide to. No pressure opens her. Kindness does.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Read the roster aloud, no comment | `BEZ_TRUST` +2, `BEZ_ILLITERACY_KNOWN` | She relaxes | Act III: she brings you the bequest herself |
| Name her illiteracy aloud | `BEZ_TRUST` −2 | She takes it well and is never easy with you again | Bequest stays sealed; every Bez ending worsens |
| Order a medicae exam of the cough | `BEZ_KNOWS_DIAGNOSIS` | Fury; confirmed; fourteen months | Act II: she can be sent knowingly on the errand that kills her |
| Say nothing | — | — | She finds out at the worst possible moment |

**Sample text**
> "He said leave it," Bez says. "That morning. He said, Sergeant, whatever you hear, leave it until the eleventh hour, and I said yes sir, and I left it until the eleventh hour."
>
> She has told this to four people today. It is not getting smoother.
>
> On the slab Ostro peels back a sheet with the boredom of a man doing inventory. "Heart's better than mine. Better than yours. Nothing wrong with it. It stopped." He looks up. "Then eleven minutes later somebody pulled a trigger with his hand. I've written what I've written. There's no box for this."

**Failure handling** — A11 and A12 are delivered unasked; they are the act's guaranteed floor. A13/A14 are missable and cost only depth.
**Flags set** — `BEZ_MET`, `BEZ_ALIVE`, `BEZ_TRUST` (0–5), `BEZ_ILLITERACY_KNOWN`, `BEZ_KNOWS_DIAGNOSIS`, `AUTOPSY_DONE`.

---

### 1.7 NOTHING TO REQUISITION
**Location** — Vosk's adjutant's office, one floor down. Ledgers, and a man who is happy.
**Present** — Savant Aleth Fenn.
**Player goal** — Corroborate the calendar; discover corroboration is not a stable substance.
**Entry gate** — `ROSETTE_TAKEN`.
**Beats**
1. Fenn greets you by the correct new title before the elevation has been filed, because he heard it in a corridor, and hearing is filing, to him.
2. **His theory:** it hasn't happened. The death is not in the record yet. He is not callous, he is literal, and mildly annoyed that everyone is behaving as though a thing has occurred which has no reference number. Funny for ninety seconds. Then not.
3. Ask about A5. The Crown's upper chambers are relay-shielded — a Mechanicus courtesy for Inquisitorial privacy, and the game must establish this here. Inside the shielding Fenn is *reliable*: same weekday, every eighty-four days, eighty-one years, no content, no location, no attendees, no travel writ. He has never noticed. He notices now, and it distresses him, and he asks you to stop.
4. Unprompted, he lists four inventory discrepancies in the chambers, one being the regicide set out of its drawer.
5. **The plant — the act's most important thirty seconds.** Walk him to the gallery arch, the shielding threshold, and ask again. Mid-sentence he says *"— every eighty-four days, going back to —"* and stops, restarts, says there were no such dates, says he must have misremembered, apologises. He is not lying. Offer no comment and no dialogue option about it. Let the player walk back inside and get the eighty-four days a third time.
6. He will follow you anywhere and does not understand that he is evidence.

**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| A15 | 84-day recurrence, 81 years, zero records | Fenn, inside shielding | Act III dating of the defection |
| A16 | His recall changes at the shielding threshold | Observation | **Act II Thread A.** Notebook only — filed, it reads *savant fatigue, recommend rest cycle* |

**Interrogation / dialogue** — He'll recite anything he has ever read, at length, on request — usable as a toy for sector minutiae. He won't speculate; no opinions, and he says so cheerfully. Pressure does not apply. He has no guile, which makes him the only person in Act I who endangers the player purely by being helpful.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Test him at the threshold | `FENN_TESTED`, A16 | You have proof and no way to show it | Thread A opens immediately in Act II |
| Confine him to shielded rooms | `FENN_SHIELDED` | Useless for fieldwork | `FENN_INTACT` protected; **Ending E stays viable** |
| Take him into the field | `FENN_FIELDED` | Enormous query advantage | Progressive corruption; `FENN_INTACT` at risk |

**Sample text**
> "There's no file," Fenn says, delighted to help. "So there's nothing to requisition against. When the file exists I'll be able to tell you a great deal."
>
> "He's dead, Fenn."
>
> "Yes," he agrees, "but not *filed*," and goes back to his ledger, and you understand for the first time that this is not a joke he is making.
>
> Later, at the gallery arch, he stops mid-sentence like a needle lifting. "— going back to —" A blink. "Sorry. There aren't any dates like that. I don't know why I said that." He is very sorry. He means it. He always means it.

**Failure handling** — Untested, A16 is re-offered in Act II the first time Fenn goes outside. Act I simply loses its cleanest proof and the reveal lands later and colder.
**Flags set** — `FENN_MET`, `FENN_INTACT`, `FENN_TESTED`, `FENN_SHIELDED`|`FENN_FIELDED`.

---

### 1.8 FORTY-ONE SECONDS
**Location** — Choir Spire, Transit Gallery 9, cell 4. Blood-smell, ozone, the grinding of the redundancy lathes.
**Present** — Vesse Ilaria Talm; Choir-Archivist Hemmel Dross (obstructive, terrified of requisition forms, not corrupt — just bureaucratically load-bearing).
**Player goal** — Establish what she transcribed at the twelfth hour.
**Entry gate** — A2 (names the window) **or** `CRASE_LEAD_CHOIR` (names the cylinders). Both is best.
**Beats**
1. Dross stonewalls. With the lead alone: a full shift, four hundred cylinders, eleven days. With A2: *twelfth hour, cell four, forty-minute window* — nine cylinders in an hour, and instant cooperation, because a specific request is a request he can survive.
2. Ilaria offers her own arrest before you ask anything. She has been waiting six hours to do it. She is not performing.
3. **Her theory:** she killed him. She believes she is the weapon and wants it to be true, because the alternative is the thing she has been hiding. She is wrong: she was a conduit, and the scrub is a property of the channel, not of guilt.
4. Cylinder six, twelfth hour and eleven minutes: **forty-one seconds of her own voice**, unmistakably hers, speaking a language she does not know in a cadence she does not have. She listens standing and does not sit down afterward.
5. Her duty-log for that window is complete, countersigned, and describes forty-one seconds of *routine tithe variance advisory, no action*. She wrote that. She does not remember writing that either.
6. **The plant:** near the end she asks you, casually, to read her own duty-roster back to her — two years, dates only. She listens with her head at an angle, thanks you, and says nothing about what she heard. She is counting gaps. She has been counting gaps for three years and telling no one, because the treatment for a psyker whose memory is failing is not a treatment.

**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| A17 | Cylinder 6 — 41 sec of her voice, unknown language | Redundancy lathe | **Physical object; it can be taken** |
| A18 | Dispatch **0997.4-1188** | Cylinder jacket + Dross's index | Scene 1.9 |
| A19 | Her own log describes what she cannot remember writing | Cell 4 log | Act II reconciliation thread |

**Interrogation / dialogue** — She'll discuss transcription mechanics, priority channels and her guilt, precisely and at length. Not the gaps: named directly in Act I she shuts entirely (`ILARIA_TRUST` −2) — too early, and she has no reason to trust whoever holds her master's rosette. What opens her is **not accusing her**: declining the arrest she offers, +2; letting her keep the cylinder, +1; confiscating it, procedurally correct, −1.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Arrest her as offered | `ILARIA_DETAINED`, trust −3 | Case closes cleanly on paper; `CRASE_STANDING` +1 | Unusable in Act II; no carrier for the unreadable message in Act III |
| Decline, and say why | `ILARIA_TRUST` +2 | Bewildered and useful | She volunteers the gap herself in Act II |
| Take the cylinder | `HAS_CYLINDER_6` | Physical, uneditable evidence | **Critical.** In Act II the Choir index for cylinder 6 becomes *blank stock, reissued*. The wax in your coat is the only counter-proof |

**Sample text**
> The servitor lowers the needle and your astropath's voice comes out of a lathe.
>
> It is her. It is entirely, insultingly her — the flat Vardan vowels, the little catch before a sibilant. It is her voice saying something with no vowels you recognise, forty-one seconds of it, in a rhythm like someone counting down a stair.
>
> Ilaria stands the whole time with her hands at her sides. When it stops she says, "That is my voice," and then, after a while, "I would like you to write down that I said that first."

**Failure handling** — Never reaching the Choir, 1.9 is entered on A1 alone; the act completes without A17/A18, `ILARIA_TRUST` frozen at 0, and Act II opens with her under Crase's detention order rather than the player's authority.
**Flags set** — `ILARIA_MET`, `ILARIA_TRUST` (0–5), `ILARIA_DETAINED`, `HAS_CYLINDER_6`, `DISPATCH_ID_KNOWN`.

---

### 1.9 THE REQUISITION
**Location** — Terminal. Vosk's cogitator, now yours. Cael-pattern, three years old, beautifully maintained.
**Present** — The player and the interface. Crase, in the room, on other business.
**Player goal** — Read the dispatch.
**Entry gate** — `A1_FILED` **or** `DISPATCH_ID_KNOWN`.
**Beats**
1. **The embarrassment, and Crase saving the player from it.** Submit the forward code `D2–F4–C9` — which is what a player lacking A3 will do — and the terminal accepts it and begins composing a formal Inquisitorial requisition, seal and priority and all, against two centuries of Munitorum drainage ledgers, addressed to the Lord-Prefect's office. Crase, reading over a chair-back, puts one finger on the arrest-key: *"Not into Kesh's clerks, love. They'll frame it."* She fixes it herself, quietly, and never mentions it again. `CRASE_STANDING` +1. It is the warmest thing anyone does for the player in Act I, and in Act III it means she knew which of the two codes was live.
2. With `ROOM_READ_FULL` the player submits `C9–F4–D2` and nobody has to save them.
3. The Vault returns one record. It is mundane. That is the point: the horror is that it is a form. Render exactly:

```
LEXICANUM VAULT :: BLOCK C / LEVEL 9 :: AISLE F4 :: SHELF D2
QUERY: DISPATCH MANIFEST, CHOIR OF CANTOR, CYCLE 0997.4
AUTH: VOSK, H. [REASSIGNED — INQ. ROSETTE 4471-K]
RETURNING 1 RECORD.
--
DISPATCH 0997.4-1188
  ORIGIN ......... RELAY THESSALINE-SECUNDUS
  TRANSIT AUTH ... VERMILION / CONCORDANCE PRIORITY (SCRUB-ON-CARRY)
  RECEIVED ....... 0997.4.211 / TWELFTH HOUR, 11 MIN
  CHOIR-VESSEL ... TALM, V. I. — CELL 4
  LENGTH ......... 41 SEC
  CLASS .......... ROUTINE / TITHE VARIANCE ADVISORY
  ADDRESSEE ...... VOSK, H., INQUISITOR — ORDO: (UNDECLARED)
  DUTY ARCHIVIST . KESH, A., LORD-PREFECT — COUNTERSIGNED IN PERSON
--
END OF RECORD. 1 OF 1.
```

4. Five load-bearing details, none flagged: **CONCORDANCE PRIORITY** — a word that means nothing to the player, appears exactly once in Act I, and is gone an hour later; **SCRUB-ON-CARRY** — the channel erases the transcriber as a routine property, i.e. Ilaria's exoneration in plain text; **41 SEC**, matching the cylinder exactly; **ROUTINE / TITHE VARIANCE ADVISORY**, matching Ilaria's own duty-log, meaning she summarised something she never knew; and a **Lord-Prefect countersigning a routine advisory in person**, which is a sector governor initialling a laundry list.
5. Filing to case-notes is offered and defaulted to *yes*. Copying the ID and the four details into the notebook takes four separate player actions and no prompt suggests it. Everything in 1.10 depends on which they did.

**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| A20 | Dispatch 0997.4-1188; Vermilion scrub; Kesh countersign | Vault query | 1.10; Act II Kesh thread. **Cogitator store by default — this is the trap** |

**Interrogation / dialogue** — Only Crase, and only if the player triggers beat 1. She will not read the returned record over your shoulder and pointedly turns away while it renders — good manners, and the last thing she needs is to be seen having seen it. Asked about Kesh's countersignature she repeats her warning, once, less warmly.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Transcribe to notebook | `A20_PHYSICAL` | Slower | 1.10 lands as *proof*; `VERITY` +1 |
| File only | `A20_FILED` | Fast, indexed | 1.10 lands as *doubt*; the player may conclude they misremembered |

**Sample text**
> One record. Not sealed, not restricted, not interesting. A routine tithe advisory that arrived at eleven minutes past the twelfth hour and took forty-one seconds to say.
>
> Priority Vermilion. Scrub-on-carry. You had to look it up: the channel wipes the transcriber. Standard practice. Nobody's fault. Ilaria has been sitting in a cell for six hours waiting to be executed for a property of the mail.
>
> At the bottom, in the duty-archivist field, a name you have heard twice today, both times from Crase.

**Failure handling** — Never running the query routes the Turn through Fenn (below). The act cannot end without it.
**Flags set** — `A20_FILED`|`A20_PHYSICAL`, `KESH_NAMED`, `CRASE_SAVED_ME`.

---

### 1.10 ONE OF ONE, NONE OF NONE
**Location** — Terminal, an hour later. The player should be doing something else when this becomes available.
**Present** — Whoever they have with them. Ideally Fenn.
**Player goal** — None assigned. The game does not ask for this.
**Entry gate** — 1.9 + 60 minutes + any further Vault query. Fallbacks below.
**Beats**
1. Any Vault query behaves normally. Then the player re-runs 1188, or opens their case-notes, on their own initiative.
2. The re-query returns:

```
LEXICANUM VAULT :: BLOCK C / LEVEL 9 :: AISLE F4 :: SHELF D2
QUERY: DISPATCH MANIFEST, CHOIR OF CANTOR, CYCLE 0997.4
RETURNING 0 RECORDS.
NO RECORD MATCHES THIS FORM. CHECK QUERY.
```

3. The case-notes do not say *deleted*. Where the requisition log read `REQ 001 — C9/F4/D2 — DISPATCH MANIFEST — 1 RECORD RETURNED — READ`, it now reads:

```
REQ 001 — D2/F4/C9 — SANITATION REQUISITION LEDGERS (MUNITORUM LIAISON)
         — RETURNED — NO ACTION — CLOSED
```

Same timestamp. Same session ID. Same author. The archive has not erased the player's work; it has **substituted the plausible mistake a novice would have made** — the exact wrong code they typed or nearly typed an hour ago. It edits you into being incompetent, in your own handwriting.
4. **No announcement, ever.** No klaxon, no glitch, no colour shift, no "something is wrong" line, no dialogue option naming it. Same font, same spacing, same everything. Scroll past it and you have scrolled past it.
5. **The only channels by which the player can know.** *Primary — the notebook:* with `A20_PHYSICAL`, the page in their coat says 1188 / Vermilion / 41 sec / Kesh, and the cogitator says no such dispatch has existed and that they pulled drainage records. The contradiction is between two things the player owns, and nothing in the game adjudicates it. `VERITY` +1. *Secondary — the cylinder:* forty-one seconds of wax, and forty-one is a number that now exists nowhere in the record. *Tertiary — Fenn:* he recites the manifest including 1188 from un-networked memory and, near a live relay, corrects himself mid-recitation and apologises. Doing it twice is the demonstration. *Nothing else.* A player who filed everything, took no cylinder and left Fenn in his office has no evidence and a strong feeling. That is a designed and legitimate outcome.
6. **If the player never re-queries**, three escalating fallbacks. **F1 (~20 min):** Crase, warmly, in passing — *"Why in the Throne's name did you pull two hundred years of drainage ledgers?"* — in front of Vaux. The player knows they didn't. Vaux says nothing and writes it down. Her kindness is the delivery vehicle for the edit for the second time, which on a replay is unbearable. **F2 (~40 min):** Dross sends a courteous physical note; the Vault has queried the Choir index and the Choir's own copy now reads *blank stock, reissued*. He assumes clerical error and apologises for the inconvenience — telling the attentive player that the mechanism reaches paper too, via people, without the game ever saying so. **F3 (hard floor, act close, unskippable):** the cogitator auto-composes the **Case Précis** for the Ordo and asks the player to countersign. Four paragraphs: Vosk self-terminated; the scene was examined and found consistent; the astropath Talm is cleared of transcribing any material dispatch, *no such dispatch being on record*; and the investigating officer requisitioned one shelf in error and closed it. Every sentence is true according to the archive, and every sentence is a small lie.

**Clues gained**

| ID | Evidence | Where found | What it unlocks |
|---|---|---|---|
| A21 | 1188 was real, was read, and is now not | The player's two stores in conflict | The whole game. `VERITY` +1. **Notebook only, by definition** |
| A22 | The Case Précis — the record's own account of Act I | Act close | What every other institution believes you did |

**Interrogation / dialogue** — Anyone the player asks will be reasonable about it. Dross apologises for a clerical error. Crase suggests fatigue, kindly, and means it. Vaux says *"then write down that you saw it, in ink, tonight"* and refuses to elaborate. Fenn is the only one who can be made to contradict the archive, and he cannot be made to do it twice in the same room.
**Branches**

| Player choice | Flag set | Consequence (near) | Consequence (far) |
|---|---|---|---|
| Sign | `PRECIS_SIGNED`, `RECONCILIATION` +3 | Institutional cooperation all through Act II | With `ROOM_READ_PARTIAL`, Ending E is now doubly foreclosed; Crase's confidence is total |
| Annotate with notebook contradictions | `PRECIS_ANNOTATED`, `VERITY` +1, `CRASE_STANDING` −1 | Kesh's office opens a file on you | Act III: your annotation is the only surviving contemporaneous account, and Vurr reads it aloud on Kadmon Veil |
| Refuse | `PRECIS_REFUSED`, `VERITY` +1, `CRASE_STANDING` −2 | Act II opens with a formal competency query | Vaux respects you; `VAUX_ALLIED` cheaper |

**Sample text**
> Returning zero records.
>
> You read it twice, which is what he trained you to do, and then you open your own case-notes, which are the one place you know what you did, because you are the one who did it.
>
> `REQ 001 — D2/F4/C9 — SANITATION REQUISITION LEDGERS — NO ACTION — CLOSED.`
>
> Same hour. Same session. Your authority code. Your handwriting, if a machine can be said to have your handwriting.
>
> In your coat there is a notebook with a number in it. The number is 1188. Nothing in the sector agrees with your notebook, and the notebook does not care.

**Failure handling** — F3 is unskippable and undismissable; the précis terminates the act. Even a player who filed everything ends Act I having been shown a document that misdescribes their own day.
**Flags set** — `ACT1_TURN_SEEN`, `PRECIS_SIGNED`|`_ANNOTATED`|`_REFUSED`, `VERITY`, `RECONCILIATION`.

---

## PLAYER BACKGROUNDS IN ACT I

**LEXMECHANIC.** *1.3* — recognises Cantor regicide notation as Vault shelf-addressing on sight; A1's encoding is stated, not deduced, though A3 is still needed to know the direction. Cost: the boots default to a *disorder* classification, and A3 requires an explicit second examination of the body before the timer expires — many Lexmechanics will lose it. *1.10* — the terminal exposes a query-checksum line no other background sees. The first render's checksum is in the session log; the second differs. This player can *prove* the edit rather than believe it: `VERITY` +2, and Act II opens with a Vault-forensics thread nobody else gets.

**PROVOST.** *1.6* — A7 comes free; Arbites death-scene protocol makes the seal-integrity register the second thing you read, always. The player therefore walks into 1.5 already knowing Vaux was there, reversing the interview's power entirely. *1.5* — procedural compulsion is available: a colder Vaux who answers everything asked and volunteers nothing, so A9 surfaces only if the player *stops* compelling and asks about Vosk as a person. Cost: Dross has been fighting Arbites requisitions for a decade, so cylinder access in 1.8 costs a favour or a point of `CRASE_STANDING`.

**PENITENT.** *1.3* — the absence registers as a presence. Where others get "no psychic trace" from Ostro in 1.6, the Penitent stands in the chambers and feels a hole shaped like an event: A12 arrives early, plus a Penitent-only A12b — the hole has an edge, and the edge is a *shape*, which in Act II is the first evidence the Erratum is structural rather than psychic. *1.8* — Ilaria knows what you are within ten seconds; she is blind and sanctioned and you are neither. Disclose voluntarily → `ILARIA_TRUST` +2 and she becomes the only retinue member who will lie to protect you. Conceal and be caught (automatic, when the cylinder plays) → −1, and she never fully settles. Cost: `RECONCILIATION` starts at 5. The archive has been in you once already.

## THE WRONG ACT I

A player can finish Act I in under thirty minutes: take the rosette, glance at the room, skip Vaux, accept the autopsy, sign the précis. The game neither stops nor scolds them, and gives them a real Act II — Crase, helpfully, hands them the Ferrum Gallery genestealer thread directly, and it is a genuine, well-built investigation with a real cult at the end of it.

What they permanently lose. **The cipher:** fewer than five signs → `ROOM_READ_PARTIAL` → no `CIPHER_HALF_ROOM`, so Act II decryption requires `CIPHER_HALF_VAUX` *and* a full Vault break *and* Vaux alive and cooperative; fail any and the death scene never resolves into a sentence, and they finish the game never knowing what Vosk wrote. **Ilaria:** no A2 → no narrow cylinder window → no forty-one seconds in Act I → `ILARIA_TRUST` ≤1 and she enters Act II under Crase's detention order rather than the player's authority, which removes her Act III function as the carrier of a message the enemy cannot read. **Ending E:** `ROOM_READ_FULL` is a hard gate. It is set or unset inside ninety minutes, and there is exactly one late recovery — Act II scene 2.0b, the holding cage, which is unsignposted, costs four days and eight points of `RECONCILIATION`, and which most players will never find. The chambers are reclaimed at act close; the crate manifest is amended; the objects themselves are not, because the archive edits paper and a boot is not paper. The player is never told. They simply arrive at Act III with four endings instead of five and no sense that a door closed.

Deliberate. Do not soften it. The thesis is that inattention is how the thing wins, and the price of demonstrating that is that some players pay it.

## ACT-LEVEL SYSTEMS NOTES

**RECONCILIATION** (0–100). Act I range 0–12; Penitent 5–17. +1 per filed clue naming Vosk's death, +3 `PRECIS_SIGNED`, +2 `A8_FILED`. Never displayed, never named. Its only Act I expression: at ≥10 the dialogue list in 1.10 loses *"ask Crase who else can amend a Vault record."* Do not grey it out. Remove the line.

**VERITY** (int). Defined as contradictions the player has personally recorded *outside* the cogitator. Act I range 0–4. Sources: A21 (+1), Lexmechanic checksum (+2), `PRECIS_ANNOTATED`/`_REFUSED` (+1), `HAS_CYLINDER_6` surviving to act close (+1). Never displayed. Expression: NPCs who are themselves being edited — Ilaria, Fenn, later Crase — respond differently at VERITY ≥2. They stop apologising for their gaps, because someone who has seen the mechanism does not need managing.

**CRASE_STANDING** (−5..+5). Starts +1; Act I is built to push an attentive player to +2 or +3. Intended. She should be the player's favourite NPC at act close.

**Where the UI lies in Act I** — exhaustive. The budget is small on purpose.
1. *1.1* — the verdict timestamp preceding discovery. Not an edit; a real bureaucratic artefact, so that the player's first suspicion of the record is wrong and they learn to distrust their paranoia before learning to trust it.
2. *1.3* — A3's filing template offers only *distress / disorder / self-neglect*. The interface lies by having no field for the truth. The whole game in a dropdown.
3. *1.7* — filed, A16 reads *savant fatigue, recommend rest cycle*. The player wrote something else; the cogitator classified it helpfully.
4. *1.9→1.10* — the Turn: record deleted, requisition log substituted with the plausible novice error, the surrounding session left perfectly intact so nothing looks broken.
5. *1.10 F2* — the Choir's physical index amended to *blank stock, reissued*.
6. *1.10 F3* — the Case Précis: four paragraphs of true sentences describing a day that did not happen.

**Implementation rule:** the game never renders a "corrupted" aesthetic. No flicker, no glyph corruption, no red text. The archive is not damaged. It is tidy. Tidiness is the horror.

## FLAG LEDGER

| Flag | Type | Written by | Read by |
|---|---|---|---|
| `ACT1_OPEN` | bool | 1.1 | — |
| `TIMESTAMP_NOTICED`, `SAID_TIRED` | bool | 1.1 | 1.2, 1.5, Act II file, Act III |
| `ROSETTE_TAKEN` | bool | 1.2 | 1.3–1.7 gates |
| `CRASE_LEAD_CHOIR` | bool | 1.2 | 1.8 gate |
| `KESH_WARNED_OFF` / `KESH_DEFIED` | enum | 1.2 | Act II Kesh thread |
| `ASKED_WHY_NOT_VAUX` | bool | 1.2 | 1.5 |
| `A0`–`A22` `_HELD`/`_FILED` | bool | 1.1–1.10 | 1.4, 1.9, 1.10, Act II |
| `A3_MISFILED`, `A5_PHYSICAL`, `A9_HELD`, `A20_PHYSICAL`, `A8_FILED` | bool | 1.3–1.9 | 1.10, Act II cipher |
| `HAS_VOSK_NOTEBOOK` | bool | 1.3 | Act III |
| `ROOM_READ_FULL`/`_PARTIAL`/`_NONE` | enum | 1.4 | Act II cipher, **Ending E gate** |
| `CIPHER_HALF_ROOM` | bool | 1.4 | Act II turn |
| `VAUX_MET`, `VAUX_ALIVE`, `VAUX_RESPECT` | bool | 1.5 | Act II |
| `VAUX_DETAINED`/`_UNCHALLENGED`/`_SHOWN_ROOM` | enum | 1.5 | `VAUX_ALLIED` cost |
| `BEZ_MET`, `BEZ_ALIVE`, `BEZ_ILLITERACY_KNOWN`, `BEZ_KNOWS_DIAGNOSIS`, `AUTOPSY_DONE` | bool | 1.6 | Act II errand branch, Act III bequest |
| `BEZ_TRUST` | 0–5 | 1.4, 1.6 | Act III |
| `FENN_MET`, `FENN_INTACT`, `FENN_TESTED` | bool | 1.7 | Act II Thread A, **Ending E gate** |
| `FENN_SHIELDED`/`FENN_FIELDED` | enum | 1.7 | Act II corruption clock |
| `ILARIA_MET`, `ILARIA_DETAINED`, `HAS_CYLINDER_6`, `DISPATCH_ID_KNOWN` | bool | 1.8 | 1.9, 1.10, Act II/III |
| `ILARIA_TRUST` | 0–5 | 1.8, 1.10 | Act III message-carrier |
| `KESH_NAMED`, `CRASE_SAVED_ME`, `ACT1_TURN_SEEN` | bool | 1.9, 1.10 | Act II open, Act III callback |
| `PRECIS_SIGNED`/`_ANNOTATED`/`_REFUSED` | enum | 1.10 | Act II standing, Act III Vurr |
| `RECONCILIATION` | 0–100 | 1.3–1.10 | all acts |
| `VERITY` | int | 1.7–1.10 | NPC response tier |
| `CRASE_STANDING` | −5..+5 | 1.1–1.10 | all acts |

Act I never writes `CIPHER_HALF_VAUX`, `VAUX_ALLIED`, `DECOY_*`, `VHAIL_BROKEN`, `KESH_IDENTIFIED`, `CAEL_SUSPECTED`, `ERRATUM_*`, `KADMON_ADDRESS`, `TRANSPORT_*`; all false/zero at handoff. `MOTT_FAVOURS` may be seeded to 1 by an optional find: an unredeemed favour-chit in Vosk's coat lining, no name on it, which is how Mott does business.

## HANDOFF TO ACT II

1. **The player is an Inquisitor** — elevated, witnessed, on record. `ROSETTE_TAKEN` always true.
2. **The Turn has been seen.** `ACT1_TURN_SEEN` always true. Whether they can *prove* it varies: `VERITY` 0 means they suspect; ≥2 means they hold something physical.
3. **All four retinue members alive, met, unaligned.** No Act I path kills anyone. Trust values vary widely.
4. **Autopsy floor guaranteed:** A11 and A12 for every player — the pistol fired eleven minutes post-mortem, and there is no mechanism. Act II can build the Erratum thread on this regardless of skill.
5. **One binary the endgame hangs on:** `ROOM_READ_FULL` or not. Never revisitable.
6. **Crase is trusted.** Act II must spend its first two hours cashing that in.
7. Kesh's name is in play. Cael's, Vhail's and Cantica-Null are not. The word *Concordance* has appeared exactly once — in a routing field, in a document that no longer exists. Act II must assume the player either wrote it down or did not, and must not remind them.
