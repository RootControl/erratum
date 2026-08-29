package game

import (
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
		"A21_HELD": true, "PRECIS_SIGNED": true, "PRECIS_ANNOTATED": true,
		"PRECIS_REFUSED": true, "ACT2_OPEN": true, "FENN_BURNED_PARTIAL": true,
		"HAS_VOSK_NOTEBOOK": true,
		// Written in Act II scenes that are not built yet.
		"ERRATUM_OBTAINED": true, "ERRATUM_UNUSED": true, "VAUX_ALLIED": true,
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
	for _, x := range gates {
		for _, f := range append(append([]string{}, x.g.Requires...), x.any...) {
			if !written[f] && !engine[f] {
				t.Errorf("%s gates on %s, which nothing writes", x.where, f)
			}
		}
	}
}
