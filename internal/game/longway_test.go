package game

import (
	"strings"
	"testing"
)

// lost is the state 2.10 exists to answer: half the cipher off the room, and
// the man who had the other half in the yard.
var lost = []string{"CIPHER_HALF_ROOM", "CIPHER_HALF_LOST"}

// The long way round does not exist while there is a short way round.
func TestLongWayOnlyOpensWhenTheHalfIsLost(t *testing.T) {
	_, log := playFrom(t, "standing_board", []string{"CIPHER_HALF_ROOM"}, "talk to ilaria")
	if strings.Contains(log, "said down") {
		t.Fatalf("she offered the wax while Vaux was still alive:\n%s", log)
	}

	s, log := playFrom(t, "standing_board", lost, "talk to ilaria")
	if !s.Has("ILARIA_OFFERED") {
		t.Fatalf("she should suggest it herself:\n%s", log)
	}
	if !strings.Contains(log, "It isn't written down anywhere. It's *said* down.") {
		t.Fatalf("beat 1 is the premise and she carries it:\n%s", log)
	}
}

// She states the price plainly, and there is no third option to find.
func TestIlariaPricesItHerself(t *testing.T) {
	_, log := playFrom(t, "standing_board", lost,
		"talk to ilaria",
		"ask ilaria about the cost",
		"ask ilaria about forty-one seconds",
		"ask ilaria about another way",
	)
	for _, want := range []string{
		"It ends the part of me you have any use for.",
		"Forty hours is the same thing eight hundred times. That's arithmetic.",
		"It is not arithmetic and you both know it.",
		"There's no third one. There's no clever one.",
	} {
		if !strings.Contains(log, want) {
			t.Fatalf("missing: %q\n%s", want, log)
		}
	}
}

// What she spends is not the player's regard for her. The counter is only how
// the regard is displayed, and it stops where she stops.
func TestSpendingIlariaCapsHerTrust(t *testing.T) {
	s, log := playFrom(t, "standing_board", lost, "talk to ilaria")
	s.Add("ILARIA_TRUST", 5)
	s.Do("send her down")

	if !s.Has("CIPHER_COMPLETE") || !s.Has("ILARIA_BURNED") {
		t.Fatal("six days of wax should finish the cipher and her")
	}
	if _, held := s.Held["B17"]; !held {
		t.Fatal("expected the rote")
	}
	if s.Ints["ILARIA_TRUST"] != 2 {
		t.Fatalf("trust should clamp to 2 on the burn, got %d", s.Ints["ILARIA_TRUST"])
	}
	s.Add("ILARIA_TRUST", 3)
	if s.Ints["ILARIA_TRUST"] != 2 {
		t.Fatalf("and it does not climb back, got %d", s.Ints["ILARIA_TRUST"])
	}
	if strings.Contains(log, "ILARIA_BURNED") || strings.Contains(log, "TRUST") {
		t.Fatalf("nothing announces what it cost:\n%s", log)
	}
}

// playDelta puts the player in Stack Delta with the headers already in hand.
func playDelta(t *testing.T, extra ...string) *State {
	t.Helper()
	s, _ := playFrom(t, "delta", append([]string{"FENN_INTACT", "FENN_ALIVE"}, append(lost, extra...)...), "look")
	s.Grant("B18")
	return s
}

func TestFennNeedsTheHeadersAndTheTerminal(t *testing.T) {
	// The same conversation in the safe house, where there is no live index.
	s, _ := playFrom(t, "safehouse", append([]string{"FENN_INTACT", "FENN_ALIVE"}, lost...), "look")
	s.Grant("B18")
	s.Do("ask fenn about the cipher")
	if s.Has("FENN_OFFERED") {
		t.Fatal("two metres of rockcrete is the one place he cannot do it")
	}

	s = playDelta(t)
	out, _ := s.Do("ask fenn about the cipher")
	if !s.Has("FENN_OFFERED") {
		t.Fatalf("in the stack, with the headers, he works it out:\n%s", strings.Join(out, "\n"))
	}
}

// He says yes immediately, and asks the question nobody can answer.
func TestFennAsksWhatToRemember(t *testing.T) {
	s := playDelta(t)
	s.Do("ask fenn about the cipher")
	out, _ := s.Do("ask fenn about the cost")
	log := strings.Join(out, "\n")
	if !strings.Contains(log, "would you tell me what I ought to try to remember?") {
		t.Fatalf("the question is the whole of the scene:\n%s", log)
	}
	if !strings.Contains(log, "I'd like to say yes now") {
		t.Fatalf("he does not have to be talked into it:\n%s", log)
	}
}

func TestSpendingFennEndsHim(t *testing.T) {
	s := playDelta(t)
	s.Do("ask fenn about the cipher")
	out, _ := s.Do("let him do it")

	if !s.Has("CIPHER_COMPLETE") || !s.Has("FENN_SPENT") {
		t.Fatal("eleven hours should finish the cipher and him")
	}
	if s.Has("FENN_INTACT") {
		t.Fatal("FENN_INTACT is what he spent")
	}
	if _, held := s.Held["B18b"]; !held {
		t.Fatal("expected the group boundaries")
	}
	if s.Ints["FENN_DRIFT"] != 6 {
		t.Fatalf("all of it, he said: got %d", s.Ints["FENN_DRIFT"])
	}
	s.Set("KADMON_HEARD")
	if st := s.statusDump(); !strings.Contains(strings.Join(st, "\n"), "CONCORDANCE OF RECORD: 0/6") {
		t.Fatalf("the readout is the only thing the player sees: %v", st)
	}
	if log := strings.Join(out, "\n"); strings.Contains(log, "FENN_DRIFT") || strings.Contains(log, "INTACT") {
		t.Fatalf("no counter names itself:\n%s", log)
	}
}

// Cael's shroud thirds incidental accrual. It is not a defence against sitting
// down at a live terminal on purpose for eleven hours, and the act takes it
// off him before the exposure lands.
func TestTheShroudDoesNotSoftenIt(t *testing.T) {
	s := playDelta(t, "FENN_SHIELDED")
	s.Do("ask fenn about the cipher")
	s.Do("let him do it")
	if s.Ints["FENN_DRIFT"] != 6 || s.Has("FENN_INTACT") {
		t.Fatalf("the shroud should not survive this: drift %d, intact %v",
			s.Ints["FENN_DRIFT"], s.Has("FENN_INTACT"))
	}
}

// Mott's elder exists only to stop a soft-lock, and only once there is nobody
// left to spend.
func TestTheThirdRouteWaitsForBothToBeGone(t *testing.T) {
	open := append([]string{"MOTT_TERMS", "ILARIA_BURNED", "FENN_INTACT"}, lost...)
	s, _ := playFrom(t, "rust_choir", open, "the elder")
	if s.Has("CIPHER_COMPLETE") {
		t.Fatal("the savant was still intact and she sold anyway")
	}

	gone := append([]string{"MOTT_TERMS", "ILARIA_BURNED"}, lost...)
	s, log := playFrom(t, "rust_choir", gone, "the elder")
	if !s.Has("CIPHER_COMPLETE") || !s.Has("NAMEBANK_DEBT") {
		t.Fatalf("with both spent, Cyra Onn is the way through:\n%s", log)
	}
	if s.Ints["MOTT_FAVOURS"] != 4 {
		t.Fatalf("four favours, got %d", s.Ints["MOTT_FAVOURS"])
	}
	if s.Ints["RECONCILIATION"] != 15 {
		t.Fatalf("the archive gets a hook in you: got %d", s.Ints["RECONCILIATION"])
	}
}

// The one invariant the scene is for: there is no way to complete the cipher
// that does not spend a named person. Assert it against the content, so a
// later scene cannot quietly add a free one.
func TestNoFreeWayToCompleteTheCipher(t *testing.T) {
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	prices := map[string]bool{"ILARIA_BURNED": true, "FENN_SPENT": true, "NAMEBANK_DEBT": true}
	paid := func(sets []string) bool {
		for _, f := range sets {
			if prices[f] {
				return true
			}
		}
		return false
	}
	// A person authored once and referenced from three scenes is one person.
	seen := map[string]bool{}
	found := 0
	check := func(where string, sets []string) {
		if !contains(sets, "CIPHER_COMPLETE") {
			return
		}
		found++
		if !paid(sets) {
			t.Errorf("%s completes the cipher and spends nobody", where)
		}
	}
	for _, sid := range g.Order {
		r := g.Scene(sid)
		check(sid+"/scene", r.Sets)
		for _, e := range r.Exits {
			check(sid+"/exit "+e.ID, e.Sets)
		}
		for _, c := range r.Clues {
			check(sid+"/clue "+c.ID, c.Sets)
		}
		for _, n := range r.People() {
			if seen[n.Key()] {
				continue
			}
			seen[n.Key()] = true
			for _, x := range n.Topics {
				check(n.Key()+"/topic "+x.ID, x.Sets)
			}
			for _, x := range n.Acts {
				check(n.Key()+"/act "+x.ID, x.Sets)
			}
			for _, x := range n.Beats {
				check(n.Key()+"/beat "+x.ID, x.Sets)
			}
			for _, x := range n.Shows {
				check(n.Key()+"/show "+x.Clue, x.Sets)
			}
		}
	}
	if found != 3 {
		t.Fatalf("expected exactly three ways through, found %d", found)
	}
}
