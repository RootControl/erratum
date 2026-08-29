package game

import (
	"strings"
	"testing"
)

// partial is a player who read the room and did not finish it, which is the
// only player 2.0b exists for.
var partial = []string{"ACT2_OPEN", "ROOM_READ_PARTIAL"}

// playTheta puts the player in the cage with the crates already requisitioned.
func playTheta(t *testing.T, extra []string, cmds ...string) (*State, string) {
	t.Helper()
	return playFrom(t, "theta", append(append([]string{}, partial...), extra...), cmds...)
}

// Nothing points at it. A bare GO lists the doors and this one is not on the
// list, for anybody, ever.
func TestThetaIsNeverOffered(t *testing.T) {
	_, log := playFrom(t, "standing_board", append(partial, "THETA_KNOWN"), "go")
	if strings.Contains(log, "THETA") {
		t.Fatalf("the door offered itself:\n%s", log)
	}
	s, _ := playFrom(t, "standing_board", append(partial, "THETA_KNOWN"), "go theta")
	if s.Room.ID != "theta" {
		t.Fatalf("it still opens to its own name; got %s", s.Room.ID)
	}
}

// Invisible to a player who already has the method, and it is the Vault that
// declines rather than the door.
func TestThetaIsInvisibleToAFullReader(t *testing.T) {
	s, log := playFrom(t, "standing_board", []string{"ACT2_OPEN", "ROOM_READ_PARTIAL", "ROOM_READ_FULL"},
		"query THETA EFFECTS")
	if s.Has("THETA_KNOWN") {
		t.Fatalf("a player who finished the room has no gap to fill:\n%s", log)
	}

	s, _ = playFrom(t, "standing_board", partial, "query THETA EFFECTS")
	if !s.Has("THETA_KNOWN") {
		t.Fatal("an ordinary, boring, entirely legitimate requisition")
	}
}

// The clearest statement the game makes about what the enemy can and cannot
// do, and it is made by a discrepancy in a stock manifest.
func TestTheArchiveCanReachThePaperAndNotTheBoots(t *testing.T) {
	_, manifest := playFrom(t, "standing_board", partial, "query EFFECTS")
	if !strings.Contains(manifest, "CRATES: 7") {
		t.Fatalf("the return says seven:\n%s", manifest)
	}
	if !strings.Contains(manifest, "GAMING SET, INCOMPLETE") {
		t.Fatalf("a stock clerk does not call a full board incomplete:\n%s", manifest)
	}
	if strings.Contains(strings.ToUpper(manifest), "FOOTWEAR,") {
		t.Fatalf("there is no line for the boots:\n%s", manifest)
	}

	s, log := playTheta(t, nil, "x crates", "x crates", "ask habrach about the count")
	if !strings.Contains(log, "You count nine.") {
		t.Fatalf("nine crates, seven listed:\n%s", log)
	}
	if !strings.Contains(log, "It cannot reach the boots.") {
		t.Fatalf("expected the statement:\n%s", log)
	}
	// The clerk counts things. It is the job. And he rubs it out.
	if !strings.Contains(log, "rubs it out with his thumb") {
		t.Fatalf("Habrach raises it and then does not:\n%s", log)
	}
	if s.Has("CAGE_OPEN") {
		t.Fatal("counting is free; opening is not")
	}
}

// A crate that has not been unsealed is not a thing you can look at, and it
// reads as absent rather than as refused.
func TestClosedCratesReadAsAbsent(t *testing.T) {
	s, log := playTheta(t, nil, "x boots", "x board")
	if strings.Contains(log, "wrong boots") {
		t.Fatalf("the cage is shut:\n%s", log)
	}
	if !strings.Contains(log, "There is nothing else down here") {
		t.Fatalf("a closed object is absent, not refused:\n%s", log)
	}
	if _, held := s.Held["A3"]; held {
		t.Fatal("nothing should have been granted")
	}
}

// The countersign is standing, ordinary, in advance, and signed by the one
// name that is on everything.
func TestTheCountersignCostsExposureQuietly(t *testing.T) {
	s, log := playTheta(t, nil, "x cage", "x cage", "open the cage")
	if !s.Has("CAGE_OPEN") {
		t.Fatal("the gate opens; it always opens")
	}
	if s.Ints["RECONCILIATION"] != 8 {
		t.Fatalf("expected 8 of exposure, got %d", s.Ints["RECONCILIATION"])
	}
	if !strings.Contains(log, "A. KESH") {
		t.Fatalf("it is the name that is on everything:\n%s", log)
	}
	if strings.Contains(log, "RECONCILIATION") {
		t.Fatalf("the counter named itself:\n%s", log)
	}
}

// A clue read off the original four days late is worth less than one read off
// the floor with the body still in the room, and nothing says so.
func TestSalvagedCluesAreWorthLess(t *testing.T) {
	s, log := playTheta(t, nil, "open the cage", "x pad", "x pad", "note A5")
	if _, held := s.Held["A5"]; !held {
		t.Fatalf("the desk pad is in crate one:\n%s", log)
	}
	if !s.Has("A5_SALVAGED") {
		t.Fatal("expected the salvage mark")
	}
	if s.Ints["VERITY"] != 0 {
		t.Fatalf("one VERITY less than its Act I value, got %d", s.Ints["VERITY"])
	}
	if strings.Contains(log, "SALVAGED") || strings.Contains(log, "VERITY") {
		t.Fatalf("nothing says it was worth less:\n%s", log)
	}

	// Read in the chambers with the body still there, the same clue pays.
	fresh, _ := play(t, "x calendar", "x calendar", "note A5")
	if fresh.Ints["VERITY"] != 1 {
		t.Fatalf("the Act I value is one, got %d", fresh.Ints["VERITY"])
	}
}

// Five signs, recorded, reconstructed off paper and a pair of boots.
func TestReconstructionReopensEndingE(t *testing.T) {
	s, log := playTheta(t, nil,
		"open the cage",
		"x board", "x board", "file A1",
		"x captures", "x captures", "file A2",
		"x quill", "x quill", "file A4",
		"x pad", "x pad", "file A5",
		"x boots", "x boots", "note A3",
		"reconstruct",
	)
	if !s.Has("ROOM_READ_FULL") || !s.Has("CIPHER_HALF_ROOM") {
		t.Fatalf("all five and it closes:\n%s", log)
	}
	if s.Ints["CRASE_STANDING"] != 0 {
		t.Fatalf("four days gone: expected +1 less one, got %d", s.Ints["CRASE_STANDING"])
	}
	if !strings.Contains(log, "you read it backwards") {
		t.Fatalf("the boots are still the instruction:\n%s", log)
	}
}

// One short and it does not close, because he did not leave a message that
// could be partly read.
func TestFourSignsIsNotFive(t *testing.T) {
	s, log := playTheta(t, nil,
		"open the cage",
		"x board", "x board", "file A1",
		"x captures", "x captures", "file A2",
		"x quill", "x quill", "file A4",
		"x pad", "x pad", "file A5",
		"reconstruct",
	)
	if s.Has("ROOM_READ_FULL") {
		t.Fatal("four signs must not close it")
	}
	if !strings.Contains(log, "It does not close.") {
		t.Fatalf("expected the incomplete reading:\n%s", log)
	}
	// And it does not name the thing that is missing.
	for _, bad := range []string{"A3", "boots are", "you need"} {
		if strings.Contains(log[strings.Index(log, "> reconstruct"):], bad) {
			t.Fatalf("the refusal explained itself: %s", bad)
		}
	}
}

// The reconstruction is only as good as what the player filed, and what they
// filed is editable. This is the Reconciliation-55 symptom, delivered where it
// costs something.
func TestTheCaseFileComesBackShort(t *testing.T) {
	s, _ := playTheta(t, nil, "open the cage")
	s.Add("RECONCILIATION", 60)
	s.Grant("A1")
	out, _ := s.Do("reconstruct")
	log := strings.Join(out, "\n")

	if s.Has("ROOM_READ_FULL") {
		t.Fatal("a thin file cannot rebuild the room")
	}
	if !strings.Contains(log, "no residue was detected") {
		t.Fatalf("the autopsy entry loses its clauses:\n%s", log)
	}
	if !strings.Contains(log, "Nothing has been struck.") {
		t.Fatalf("there must be no gap to notice:\n%s", log)
	}
	for _, bad := range []string{"DELETED", "ERROR", "AMENDED", "REVISED"} {
		if strings.Contains(log, bad) {
			t.Fatalf("the archive announced itself: %s\n%s", bad, log)
		}
	}
}

// He is not on the manifest and never was, so there is nothing to sign out.
func TestTakingTheBootsLeavesNoTrace(t *testing.T) {
	s, log := playTheta(t, nil, "open the cage", "x boots", "x boots", "take boots")
	if !s.Has("VOSK_BOOTS") {
		t.Fatalf("expected the boots:\n%s", log)
	}
	if !strings.Contains(log, "you have removed nothing from Stack Theta at all") {
		t.Fatalf("the record cannot see it:\n%s", log)
	}
	if s.Ints["RECONCILIATION"] != 8 {
		t.Fatalf("taking them costs nothing, got %d", s.Ints["RECONCILIATION"])
	}
}

// The notebook used to get its flag from a line in the engine that named it.
func TestObjectsWriteTheirOwnFlags(t *testing.T) {
	s, _ := play(t, "take notebook")
	if !s.Has("HAS_VOSK_NOTEBOOK") {
		t.Fatal("taking a dead man's longhand is still how you get it")
	}
}

// The Vault's null answer was authored in the chambers only, so every Act II
// query for a code it does not know rendered a marker instead of a refusal.
// Found by 2.0b's own query, which is the first one Act II makes it easy to
// get wrong.
func TestVaultRefusesEverywhere(t *testing.T) {
	for _, scene := range []string{"vosk_chambers", "standing_board"} {
		_, log := playFrom(t, scene, []string{"ROSETTE_TAKEN", "ACT2_OPEN"}, "query NONSENSE-1")
		if strings.Contains(log, "missing string") {
			t.Fatalf("%s: %s", scene, log)
		}
	}
}
