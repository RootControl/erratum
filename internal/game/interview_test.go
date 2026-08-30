package game

import (
	"strings"
	"testing"
)

// Accuracy is what opens Vaux, and authority is not. A player who pulled the
// seal-integrity register before coming down gets a different man.
func TestAccuracyOpensVaux(t *testing.T) {
	careful, log := play(t,
		"x register", "x register",
		"go vaux",
		"ask vaux about the seal",
		"ask about the caseload",
	)
	if !careful.Has("VAUX_RESPECT") {
		t.Fatal("citing the integrity register should recalibrate him")
	}
	if !strings.Contains(log, "He stopped *filing*") {
		t.Fatal("with respect earned, the caseload answer is the real one")
	}

	blunt, bluntLog := play(t,
		"go vaux",
		"ask vaux about the seal",
		"ask about the caseload",
	)
	if blunt.Has("VAUX_RESPECT") {
		t.Fatal("he does not recalibrate for a player who arrived empty-handed")
	}
	if strings.Contains(bluntLog, "He stopped *filing*") {
		t.Fatal("the real caseload answer must be gated on accuracy")
	}
	if !blunt.Holds("A7") && blunt.Held["A7"] == StoreNone {
		// He concedes the seal either way; only the manner differs.
		if _, held := blunt.Held["A7"]; !held {
			t.Fatal("he admits the seal to everyone who asks")
		}
	}
}

// The boots invariant, in a conversation instead of a room. A9 is a sound a
// man made without knowing he made it, and the cogitator has no field for
// that, so the record can never hold it.
func TestMnemonicCannotBeFiled(t *testing.T) {
	s, log := play(t,
		"go vaux",
		"ask about your investigation",
		"ask about being passed over",
		"file A9",
		"note phrase",
	)
	if _, heard := s.Held["A9"]; !heard {
		t.Fatal("real pressure should surface the mnemonic")
	}
	if s.Held["A9"] == StoreCog || s.Held["A9"] == StoreBoth {
		t.Fatal("A9 must never reach the cogitator")
	}
	if !s.Has("A9_HELD") {
		t.Fatal("writing it down is what makes it evidence")
	}
	if !strings.Contains(log, "no field for a sound") {
		t.Fatal("the filing template should refuse it by having nowhere to put it")
	}
	if len(s.Notebook) != 1 {
		t.Fatalf("the mnemonic belongs in the coat and nowhere else: %v", s.Notebook)
	}
}

// The archive never announces itself and neither does the game. A clue the
// player did not ask for arrives without a bracket telling them it mattered;
// a player who does not notice it walks out without it, which is the design.
func TestMnemonicIsNotAnnounced(t *testing.T) {
	s, log := play(t,
		"go vaux",
		"ask about your investigation",
		"ask about being passed over",
		"go back",
	)
	if _, heard := s.Held["A9"]; !heard {
		t.Fatal("he says it whether or not they are listening")
	}
	if strings.Contains(log, "A9") {
		t.Fatal("nothing may flag the mnemonic as a clue — it is a thing he said")
	}
	if s.Has("A9_HELD") {
		t.Fatal("hearing it is not recording it")
	}
	if len(s.Notebook) != 0 {
		t.Fatal("the game must not write the player's notebook for them")
	}
}

// VAUX_DETAINED | VAUX_SHOWN_ROOM | VAUX_UNCHALLENGED is an enum in the act
// document and three flat bools in the engine, and exactly one of them is
// true when the player walks out — including for the player who did nothing.
func TestDispositionIsOneOfN(t *testing.T) {
	shown, _ := play(t,
		"x boots", "x boots", "note A3",
		"go vaux", "show A3 to vaux", "go back",
	)
	assertOne(t, shown, "VAUX_SHOWN_ROOM")

	detained, _ := play(t,
		"x boots", "x boots", "note A3",
		"go vaux", "show A3 to vaux", "detain vaux", "go back",
	)
	assertOne(t, detained, "VAUX_DETAINED")
	if detained.Ints["CRASE_STANDING"] != 0 { // opens at +1
		t.Fatalf("four hours and nothing gained should cost standing: %d", detained.Ints["CRASE_STANDING"])
	}

	quiet, _ := play(t, "go vaux", "ask about procedure", "go back")
	assertOne(t, quiet, "VAUX_UNCHALLENGED")
}

func assertOne(t *testing.T, s *State, want string) {
	t.Helper()
	group := []string{"VAUX_DETAINED", "VAUX_SHOWN_ROOM", "VAUX_UNCHALLENGED"}
	var set []string
	for _, f := range group {
		if s.Has(f) {
			set = append(set, f)
		}
	}
	if len(set) != 1 || set[0] != want {
		t.Fatalf("expected exactly %s, got %v", want, set)
	}
}

// The lift is not free, and the body does not wait in the room while you are
// in one. A player who goes down to tier four early comes back to furniture.
func TestTravelCostsTheBody(t *testing.T) {
	s, log := play(t,
		"go vaux", // three actions in a lift
		"ask about procedure", "ask about seniority", "ask about vosk",
		"ask about the seal", "ask about crase", "ask about how he died",
		"ask about the verdict", "ask about evidence", "ask about investigation",
		"go back",
		"x boots",
	)
	if s.BodyPresent {
		t.Fatal("the lift and the interview together outlast eight minutes")
	}
	if s.Holds("A3") {
		t.Fatal("the boots do not survive the trip")
	}
	if !strings.Contains(log, "Somewhere above you") {
		t.Fatal("the removal must be narrated even when the player is elsewhere")
	}
	if !strings.Contains(log, "medicae photograph") {
		t.Fatal("coming back to an empty floor should say so")
	}
}

// Evidence travels between scenes even though scenes do not. A3 is picked up
// on tier nine and put on a desk on tier four.
func TestClueTravelsBetweenScenes(t *testing.T) {
	s, log := play(t,
		"x boots", "x boots", "note A3",
		"go vaux",
		"show the boots to vaux",
	)
	if !strings.Contains(log, "I would like to see the room") {
		t.Fatal("showing him the boots is the largest single step toward an alliance")
	}
	if !s.Has("VAUX_SHOWN_ROOM") {
		t.Fatal("expected VAUX_SHOWN_ROOM")
	}
}

// The flag ledger has 1.5 reading SAID_TIRED, and it should: the one true,
// easy thing the player says in 1.1 is the strongest line in the file against
// them, and Vaux has read the file.
func TestVauxHasReadYourStatement(t *testing.T) {
	spoke, spokeLog := playAct1(t,
		"ask about sealed", "ask about the log", "ask about who signed", "tired",
		"go on", "take rosette", "go in",
		"go vaux", "ask about procedure",
	)
	if !spoke.Has("SAID_TIRED") {
		t.Fatal("1.1 writes it")
	}
	if !strings.Contains(spokeLog, "strongest single line in the file") {
		t.Fatal("1.5 reads it")
	}

	_, quietLog := playAct1(t,
		"ask about sealed", "ask about the log", "ask about who signed", "nothing",
		"go on", "take rosette", "go in",
		"go vaux", "ask about procedure",
	)
	if strings.Contains(quietLog, "strongest single line in the file") {
		t.Fatal("said nothing, and there is nothing for him to hold")
	}
}
