package game

import (
	"strings"
	"testing"
)

// playGallery accepts the brief and goes down.
func playGallery(t *testing.T, pre []string, cmds ...string) (*State, string) {
	t.Helper()
	head := []string{"sign", "query cathedral alms", "go fenn", "ask fenn about kadmon"}
	head = append(head, pre...)
	head = append(head, "back", "wait", "go crase", "accept", "back", "go gallery")
	return playToPrecis(t)(append(head, cmds...)...)
}

// Sarn will not discuss the brood, the Patriarch or her own physiology, and no
// amount of standing changes that. She discusses the erasures at length,
// because she has wanted to tell somebody for twenty years.
func TestSarnHasTwoSubjectsAndAThird(t *testing.T) {
	s, log := playGallery(t, nil, "ask sarn about the brood", "ask sarn about historians")
	if !strings.Contains(log, "That's one of the two") {
		t.Fatal("she stops, entirely, in the middle of being co-operative")
	}
	if !s.Has("SARN_INTERVIEWED") || !s.Has("PATTERN_LITERACY") {
		t.Fatal("the third subject is the one nobody else in the game has")
	}
	if _, held := s.Held["B12"]; !held {
		t.Fatal("B12 is the erasure pattern")
	}
	if !strings.Contains(log, "not one of them was leadership") {
		t.Fatal("nineteen in forty years, and it takes the ones who write")
	}
	if !strings.Contains(log, "your own notebook") {
		t.Fatal("and the player is in the pattern too")
	}
}

// The purges are bookkeeping. They follow sector tithe reconciliation, not
// discovery, which means the flush order for the Rows is scheduled and not
// reactive either.
func TestThePurgesAreBookkeeping(t *testing.T) {
	s, log := playGallery(t, nil,
		"ask sarn about historians", "ask sarn about the cycle", "ask sarn about the inspector")
	for _, id := range []string{"B13", "B14"} {
		if _, held := s.Held[id]; !held {
			t.Fatalf("%s comes out of the same conversation", id)
		}
	}
	if !strings.Contains(log, "We're a variance") {
		t.Fatal("a purge is not a response to anything")
	}
	if !strings.Contains(log, "The unwritten inspector's been round") {
		t.Fatal("it is a joke, because the alternative is unsayable on a shift")
	}
	if !strings.Contains(log, "It takes the ones with the pens") {
		t.Fatal("and it never takes the head")
	}
}

// The death is deterministic and ordered so it can be tested: Fenn if brought,
// else Bez, else Vaux if allied, else Ilaria.
func TestTheDeathCascadeIsDeterministic(t *testing.T) {
	fenn, fLog := playGallery(t, []string{"bring him"}, "assault")
	if !fenn.Has("FENN_DEAD") || fenn.Has("BEZ_DEAD") {
		t.Fatal("bring him for the map and he is the one in the sluice")
	}
	if !strings.Contains(fLog, "what he is saying is the flood interval, correctly, twice") {
		t.Fatal("he is not frightened at any point")
	}
	if fenn.Has("FENN_INTACT") {
		t.Fatal("and the implant is not recoverable")
	}

	bez, bLog := playGallery(t, nil, "assault")
	if !bez.Has("BEZ_DEAD") || bez.Has("FENN_DEAD") {
		t.Fatal("leave him behind and the door is held by the sergeant")
	}
	if !strings.Contains(bLog, "the obvious solution to a door") {
		t.Fatal("no speech, no memorable line")
	}

	// Every branch is a real victory and every branch costs somebody.
	for _, s := range []*State{fenn, bez} {
		if !s.Has("GALLERY_CLEARED") || !s.Has("DECOY_RESOLVED") {
			t.Fatal("the victory is real")
		}
		if s.Ints["CRASE_STANDING"] < 4 {
			t.Fatalf("and worth three: %d", s.Ints["CRASE_STANDING"])
		}
	}
}

// Bez knowing her own diagnosis is calm about it, and says so, which is worse.
func TestBezIsCalmAboutItIfSheKnows(t *testing.T) {
	s := setupKnowingBez(t)
	out, _ := s.Do("assault")
	log := strings.Join(out, "\n")
	if !s.Has("BEZ_DEAD") {
		t.Fatal("she still holds the door")
	}
	if !strings.Contains(log, "Fourteen months, sir") {
		t.Fatal("she volunteers before anybody finishes describing the problem")
	}
	if !strings.Contains(log, "you told her the number yourself") {
		t.Fatal("and you are the reason she can be calm about it")
	}
}

func setupKnowingBez(t *testing.T) *State {
	t.Helper()
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	s := New(g)
	for _, f := range []string{"ROSETTE_TAKEN", "ACT1_COMPLETE", "BEZ_ALIVE",
		"BEZ_KNOWS_DIAGNOSIS", "DECOY_PURSUED"} {
		s.Set(f)
	}
	if !s.Move("ferrum") {
		t.Fatal("no ferrum scene")
	}
	return s
}

// Withdrawal is always available, nobody dies, and the cult goes on killing
// people. No later clue depends on it.
func TestWithdrawalIsAlwaysAvailable(t *testing.T) {
	s, log := playGallery(t, nil, "ask sarn about historians", "withdraw")
	if !s.Has("GALLERY_ABORTED") || s.Has("GALLERY_CLEARED") {
		t.Fatal("pulling out is pulling out")
	}
	for _, f := range []string{"FENN_DEAD", "BEZ_DEAD", "VAUX_DEAD", "ILARIA_BURNED"} {
		if s.Has(f) {
			t.Fatalf("nobody dies: %s", f)
		}
	}
	if _, held := s.Held["B12"]; !held {
		t.Fatal("and the information survives, because it was never behind the assault")
	}
	if !strings.Contains(log, "Her regard does not survive it either") {
		t.Fatal("Crase's regard curdles")
	}
}

// A player who refused the brief is not punished with an information hole:
// Mott sells the meeting for two favours, and nobody stands on the wet side of
// a coolant door at all.
func TestTheRefuserCanBuyTheMeeting(t *testing.T) {
	s, log := playToPrecis(t)("sign", "query cathedral alms",
		"go fenn", "ask fenn about kadmon", "back", "wait", "go crase", "refuse", "back",
		"go sump", "surrender", "ask mott about favours", "buy a meeting",
		"back", "go smuggled", "ask sarn about historians")
	if !s.Has("DECOY_REFUSED") {
		t.Fatal("this is the route for a player who said no")
	}
	if s.Ints["MOTT_FAVOURS"] != 2 {
		t.Fatalf("the meeting is two: %d", s.Ints["MOTT_FAVOURS"])
	}
	if _, held := s.Held["B12"]; !held {
		t.Fatal("and it buys the one thing nobody else in the game has")
	}
	if !strings.Contains(log, "nobody at all who is going to be standing on the wet side") {
		t.Fatal("no assault, no enforcers, nobody dies")
	}
}

// The victory is real and then the Arbites turn it into a records purge, and
// every part of that is correct.
func TestTheVictoryBecomesARecordsPurge(t *testing.T) {
	_, log := playGallery(t, nil, "assault", "back")
	if !strings.Contains(log, "forty thousand relocated") {
		t.Fatal("the purge you earned")
	}
	if !strings.Contains(log, "every part of it is correct") {
		t.Fatal("and there is nobody to complain to")
	}
}
