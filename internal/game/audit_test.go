package game

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/RootControl/erratum/internal/content"
)

// Reconciliation removes rather than closes, and the act documents name the
// option it removes: the direct accusation, gone at 25, with the answer coming
// back as an ordinary "I have no idea" and no gap where it used to be.
func TestTheAccusationVanishesAtTwentyFive(t *testing.T) {
	pre := []string{"ask ruse about the crypt", "ask ivrit about the crypt", "go vhail"}

	clean, cLog := playCathedral(t, append(append([]string{}, pre...),
		"ask vhail about you wired them")...)
	if !clean.Has("VHAIL_ACCUSED") {
		t.Fatal("below the threshold the option is there")
	}
	if !strings.Contains(cLog, "let me tell you what it was for") {
		t.Fatal("and he answers it in one syllable")
	}

	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	s := New(g)
	for _, f := range []string{"ROSETTE_TAKEN", "ACT1_COMPLETE"} {
		s.Set(f)
	}
	s.Add("RECONCILIATION", 30)
	s.Move("cathedral")
	for _, c := range pre {
		s.Do(c)
	}
	out, _ := s.Do("ask vhail about you wired them")
	joined := strings.Join(out, "\n")
	if s.Has("VHAIL_ACCUSED") {
		t.Fatal("past 25 the option does not exist")
	}
	if strings.Contains(joined, "you wired") || strings.Contains(joined, "let me tell you what") {
		t.Fatal("nothing of it survives")
	}
	if !strings.Contains(joined, "I have no idea") {
		t.Fatalf("and what comes back is ordinary, with no gap in it: %q", joined)
	}
}

// A player smuggled in for the interview has no cordon, no writ and no force,
// and the game says so rather than letting them order an assault they cannot
// mount.
func TestTheSmuggledVisitorHasNoForce(t *testing.T) {
	s, log := playToPrecis(t)("sign", "query cathedral alms",
		"go fenn", "ask fenn about kadmon", "back", "wait", "go crase", "refuse", "back",
		"go sump", "surrender", "ask mott about favours", "owe him", "buy a meeting",
		"back", "go smuggled", "assault")
	if s.Has("GALLERY_CLEARED") {
		t.Fatal("there is nothing to clear it with")
	}
	if !strings.Contains(log, "There is no assault") {
		t.Fatal("and the refusal should say why")
	}
	if !strings.Contains(log, "when you said no") {
		t.Fatal("and name where the decision was actually made")
	}
}

// The purge order is what makes Casmir Lyle askable, and asking Mott about him
// is what sends the player to the query that has no answer.
func TestThePurgeOrderOpensLyle(t *testing.T) {
	blind, blindLog := playToPrecis(t)("sign", "go sump", "surrender",
		"ask mott about favours", "owe him", "ask mott about lyle")
	if strings.Contains(blindLog, "Seat four-one-one") {
		t.Fatal("there is nothing to ask about until the vellum is out of the coat")
	}
	_ = blind

	s, log := playToPrecis(t)("sign", "go sump", "surrender",
		"ask mott about favours", "owe him", "go annexe", "x order", "x order",
		"out", "ask mott about lyle")
	if !s.Has("PURGE_ORDER_SEEN") {
		t.Fatal("holding the vellum is what opens him")
	}
	if !strings.Contains(log, "I'd genuinely like to know") {
		t.Fatal("he has wanted somebody with clearance to ask this for nineteen years")
	}
}

// Two clues answering to one word is a bug the finder cannot see, so no two
// clues in the game may share a name.
func TestNoTwoCluesShareAName(t *testing.T) {
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]string{}
	for _, sid := range g.Order {
		for _, c := range g.Scene(sid).Clues {
			if len(c.Names) == 0 {
				t.Errorf("clue %s answers to nothing but its id", c.ID)
			}
			for _, n := range c.Names {
				n = strings.ToLower(n)
				if prev, dup := seen[n]; dup && prev != c.ID {
					t.Errorf("%q resolves to both %s and %s", n, prev, c.ID)
				}
				seen[n] = c.ID
			}
		}
	}
}

// Every flag the content gates on is written by something. A gate on a flag
// nothing sets is a line of the game that silently never happens.
func TestEveryGatedFlagHasAWriter(t *testing.T) {
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	// Written by the engine rather than by content.
	engine := map[string]bool{
		"ROSETTE_TAKEN": true, "ACT1_COMPLETE": true, "ACT1_TURN_SEEN": true,
		"BODY_REMOVED": true, "LYLE_QUERIED": true, "CRASE_ASKED_FOR_YOU": true,
		"DISPATCH_ID_KNOWN": true, "KESH_NAMED": true, "CRASE_SAVED_ME": true,
		"ROOM_READ_FULL": true, "ROOM_READ_PARTIAL": true, "CIPHER_HALF_ROOM": true,
		"PRECIS_SIGNED": true, "PRECIS_ANNOTATED": true,
		"PRECIS_REFUSED": true, "ACT2_OPEN": true, "FENN_BURNED_PARTIAL": true,
		"HAS_VOSK_NOTEBOOK": true,
		// TWO_THREADS is the same count that summons Crase; the store habit is
		// whichever of FILE and NOTE the player reached for first in Act II,
		// and there is no scene in which either is chosen.
		"TWO_THREADS": true, "STORE_HABIT_FILED": true, "STORE_HABIT_LEDGER": true,
		// The plate is unmissable: Cael's imaging act writes these, and so does
		// the fallback in Tick where Tass brings it up four tiers unbidden.
		"ERRATUM_OBTAINED": true, "ERRATUM_UNUSED": true,
		// CIPHER_COMPLETE is derived in Set from holding both halves, so it has
		// no content writer by design. KADMON_ADDRESS is not here: B20's own
		// sets write it.
		"CIPHER_COMPLETE": true,
		// DECODE, in the chambers, with both halves and the plate.
		"ACT2_COMPLETE": true, "FALSE_SENTENCE_SEEN": true, "SENTENCE_REPAIRED": true,
		"CANTICA_QUERIED": true, "TASS_BROUGHT_IT": true,
		// The effects manifest, requisitioned from the Vault by a player who
		// thought to ask where a dead man's crates went.
		"THETA_KNOWN": true,
		// The three openings the bible offers, chosen before the first
		// command by SetBackground and never mentioned again.
		"BG_LEXMECHANIC": true, "BG_PROVOST": true, "BG_PENITENT": true,
		// Derived in Set from the clues that make up the method, the way
		// CIPHER_COMPLETE is derived from the two halves.
		"VOSK_METHOD_PARTIAL": true, "VOSK_METHOD_UNDERSTOOD": true,
		// Out of relay range, and the endings.
		"IN_TRANSIT": true, "VOYAGE_FULL": true, "C9_PHYS": true, "C9_COG": true,
	}
	written := map[string]bool{}
	note := func(fs []string) {
		for _, f := range fs {
			written[f] = true
		}
	}
	type gated struct {
		where string
		g     content.Gate
		any   []string
	}
	var gates []gated
	for _, sid := range g.Order {
		r := g.Scene(sid)
		note(r.Sets)
		note(r.While)
		for _, e := range r.Exits {
			note(e.Sets)
			gates = append(gates, gated{sid + "/exit " + e.ID, e.Gate, e.RequiresAny})
		}
		for _, c := range r.Clues {
			note(c.Sets)
		}
		for _, n := range r.People() {
			written[strings.ToUpper(n.ID)+"_MET"] = true
			written[strings.ToUpper(n.ID)+"_ALIVE"] = true
			if n.Disposition != nil {
				note(n.Disposition.Flags)
			}
			for _, t := range n.Topics {
				note(t.Sets)
				gates = append(gates, gated{sid + "/" + n.ID + "/topic " + t.ID, t.Gate, nil})
			}
			for _, a := range n.Acts {
				note(a.Sets)
				gates = append(gates, gated{sid + "/" + n.ID + "/act " + a.ID, a.Gate, nil})
			}
			for _, b := range n.Beats {
				note(b.Sets)
				gates = append(gates, gated{sid + "/" + n.ID + "/beat " + b.ID, b.Gate, nil})
			}
			for _, sh := range n.Shows {
				note(sh.Sets)
			}
		}
	}
	// FILE and NOTE write <ID>_HELD, _FILED and _PHYSICAL for every clue in
	// the game, in the engine, so a gate on one of those is written by
	// definition and does not need listing above.
	for id := range g.Clues {
		for _, suffix := range []string{"_HELD", "_FILED", "_PHYSICAL"} {
			engine[strings.ToUpper(id)+suffix] = true
		}
	}
	for _, x := range gates {
		for _, f := range append(append([]string{}, x.g.Requires...), x.any...) {
			if !written[f] && !engine[f] {
				t.Errorf("%s gates on %s, which nothing writes", x.where, f)
			}
		}
	}
}

// parserVerbs are the words Do() handles itself. An act named with one of them
// at the head is unreachable, because the verb is consumed before act() is
// ever called — except for the few that deliberately fall through when their
// argument names no door, object or layer.
var parserVerbs = map[string]bool{
	"help": true, "?": true, "verbs": true, "quit": true, "exit": true,
	"look": true, "l": true, "out": true, "in": true, "back": true, "up": true,
	"down": true, "enter": true, "ask": true, "interview": true, "question": true,
	"show": true, "present": true, "talk": true, "greet": true, "approach": true,
	"see": true, "topics": true, "subjects": true, "who": true, "examine": true,
	"x": true, "inspect": true, "file": true, "note": true, "write": true,
	"clues": true, "evidence": true, "notes": true, "casenotes": true,
	"case-notes": true, "notebook": true, "coat": true, "query": true,
	"requisition": true, "req": true, "wait": true, "z": true, "status": true,
	"decode": true, "decipher": true, "translate": true, "cite": true,
	"produce": true,
}

// TestNoActIsNamedForAParserVerb is a permanent guard on a mistake this repo
// has made four times: an act named "ask fenn", "the corpus", "file a
// correction" or "press him", which the parser eats before any person sees it.
// GO, READ, TAKE and PRESS fall through on purpose and are not listed above.
func TestNoActIsNamedForAParserVerb(t *testing.T) {
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, sid := range g.Order {
		for _, n := range g.Scene(sid).People() {
			if seen[n.Key()] {
				continue
			}
			seen[n.Key()] = true
			for _, a := range n.Acts {
				for _, name := range a.Names {
					head := strings.ToLower(strings.Fields(name)[0])
					if parserVerbs[head] {
						t.Errorf("%s/act %s: %q is unreachable — %q is a parser verb",
							n.Key(), a.ID, name, head)
					}
				}
			}
		}
	}
}

// TestNoTwoExitsShareAName guards the other half of the same class. FindExit
// picks the best-scoring door regardless of gates, so two exits answering to
// one word means the second is unreachable by that word however the gates fall.
func TestNoTwoExitsShareAName(t *testing.T) {
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	for _, sid := range g.Order {
		owner := map[string]string{}
		for _, e := range g.Scene(sid).Exits {
			for _, name := range append([]string{e.ID}, e.Names...) {
				name = strings.ToLower(name)
				if prev, dup := owner[name]; dup && prev != e.ID {
					t.Errorf("%s: %q names both %s and %s; the second is unreachable by it",
						sid, name, prev, e.ID)
				}
				owner[name] = e.ID
			}
		}
	}
}

// TestDelayedTextResolvesEverywhere is the guard for the bug class that has
// now bitten four times: a string authored in the scene that introduces it and
// emitted from Tick, which reads whichever scene the player happens to be
// standing in. vault_none did it, the Cerement symptoms did it, and both were
// found by hand. Arm every delayed emitter in the game, stand in every scene
// in turn, and let the clock run.
func TestDelayedTextResolvesEverywhere(t *testing.T) {
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	arm := []string{
		"LYLE_QUERIED",    // the Vault confirms a nil return nobody asked it to
		"C17_FILED",       // and takes the intake protocols back within the hour
		"IN_TRANSIT",      // the Cerement's symptoms, one every thirty-six hours
		"CIPHER_COMPLETE", // Tass brings the plate up four tiers unbidden
		"CANTICA_QUERIED", // and the routing index answers a second time
		"FENN_INTACT", "FENN_ALIVE",
	}
	for _, sid := range g.Order {
		s := New(g)
		if !s.Move(sid) {
			t.Fatalf("cannot enter %s", sid)
		}
		for _, f := range arm {
			s.Set(f)
		}
		s.Add("RECONCILIATION", 60)
		for i := 0; i < 60; i++ {
			for _, line := range s.Tick() {
				if strings.Contains(line, "[missing string:") {
					j := strings.Index(line, "[missing string:")
					t.Errorf("%s: %s", sid, line[j:min(j+40, len(line))])
				}
			}
		}
	}
}

// TestNoDateAfterThePresent is the guard for a class that got away from me for
// the whole of Act III. STORY_BIBLE.md §12 fixes the present at 826.M41 and
// says any new date must be derived from Kadmon's 735, the Ferrum cult's 763,
// Vhail's 786 and Mott's district at 807. Three dozen authored years had
// drifted into the 900s, which in a game that trains the player to read every
// number as evidence is worse than a typo — it reads as a clue.
//
// The lint reads the scene files directly, because it is a lint on the prose
// rather than on anything the loader exposes.
func TestNoDateAfterThePresent(t *testing.T) {
	// Numbers above 826 that are not years. Kadmon's reports are numbered,
	// the acknowledged census is 900,412,000, and the game's own date stamps
	// are cycles — 0997.4.211 — not Imperial years.
	notYears := map[string]bool{
		"997": true, "996": true, "993": true, // cycle stamps
		"833": true, "922": true, "923": true, "924": true, // Kadmon's reports
		"900": true, // nine hundred million four hundred and twelve thousand
		"087": true, // cycle 087.M42, the far end of the sum
	}
	files, err := filepath.Glob("../content/*.json")
	if err != nil || len(files) == 0 {
		t.Fatalf("no scene files to lint: %v", err)
	}
	year := regexp.MustCompile(`\b0?(\d{3})\b`)
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range year.FindAllStringSubmatch(string(raw), -1) {
			if notYears[m[1]] {
				continue
			}
			if n, _ := strconv.Atoi(m[1]); n > 826 {
				t.Errorf("%s: %s is after the present day (826.M41)",
					filepath.Base(f), m[1])
			}
		}
	}
}
