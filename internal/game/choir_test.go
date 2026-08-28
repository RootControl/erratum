package game

import (
	"strings"
	"testing"
)

// playSpire takes the shutter across the Crown gap, which is the narrow way in.
func playSpire(t *testing.T, cmds ...string) (*State, string) {
	t.Helper()
	return play(t, append([]string{"x shutter", "x shutter", "note A2", "go spire"}, cmds...)...)
}

// The Spire opens on A2 or on Crase's lead, and either is enough. With
// neither, you would be walking up there to ask about four hundred vessels in
// nine galleries and you know it.
func TestTheSpireOpensTwoWays(t *testing.T) {
	shut, _ := play(t, "x shutter", "x shutter", "go spire")
	if shut.Room.ID != "choir" {
		t.Fatal("the shutter names the cell, and that is enough")
	}

	lead, _ := playFrom(t, "vosk_chambers", []string{"ROSETTE_TAKEN", "CRASE_LEAD_CHOIR"}, "go spire")
	if lead.Room.ID != "choir" {
		t.Fatal("Crase's lead names the cylinders, and that is also enough")
	}

	blind, log := play(t, "go spire")
	if blind.Room.ID != "vosk_chambers" {
		t.Fatal("with neither, there is nothing to go up there and ask")
	}
	if !strings.Contains(log, "You need a cell, or an hour, or a reason") {
		t.Fatal("and the refusal should say what is missing")
	}
}

// Dross is not corrupt and not loyal to anybody. He is a man who has survived
// twenty years by never being the name at the bottom of a decision, so a broad
// request costs eleven days and a specific one costs him nothing at all.
func TestDrossGrantsWhatHeCanSurvive(t *testing.T) {
	narrow, log := playSpire(t, "ask dross about access")
	if !narrow.Has("CHOIR_ACCESS") {
		t.Fatal("twelfth hour, cell four, a forty-minute window is a retrieval")
	}
	if !strings.Contains(log, "isn't a release") {
		t.Fatal("and the difference is the whole of his position")
	}

	broad, broadLog := playFrom(t, "vosk_chambers",
		[]string{"ROSETTE_TAKEN", "CRASE_LEAD_CHOIR"}, "go spire", "ask dross about access")
	if broad.Has("CHOIR_ACCESS") {
		t.Fatal("the lead alone names the cylinders but not the window")
	}
	if !strings.Contains(broadLog, "Eleven days") {
		t.Fatal("four hundred cylinders across nine cells and three watches")
	}
	if _, held := broad.Held["A17"]; held {
		t.Fatal("and Act I ends before eleven days do")
	}
}

// Cylinder six: forty-one seconds of her own voice in a language she does not
// know, and a dispatch number on the jacket in her own hand.
func TestFortyOneSeconds(t *testing.T) {
	s, log := playSpire(t, "ask dross about access", "play cylinder")
	for _, id := range []string{"A17", "A18"} {
		if _, held := s.Held[id]; !held {
			t.Fatalf("%s comes off the lathe", id)
		}
	}
	if !strings.Contains(log, "That is my voice") ||
		!strings.Contains(log, "I would like you to write down that I said that first") {
		t.Fatal("she says it in that order and the order matters")
	}
	if !s.Has("DISPATCH_ID_KNOWN") {
		t.Fatal("the number on the jacket is an address")
	}

	// And not before Dross has found it for you.
	early, earlyLog := playSpire(t, "play cylinder")
	if _, held := early.Held["A17"]; held {
		t.Fatal("there are four hundred cylinders and nobody has told you which")
	}
	if !strings.Contains(earlyLog, "nobody has told you which one") {
		t.Fatal("and the refusal should say so")
	}
}

// The dispatch number is the second door into 1.9. A player who never filed
// the shelf-address can requisition against the number instead — and the
// archive edits that route exactly as silently.
func TestTheDispatchNumberOpensTheVault(t *testing.T) {
	s, log := playSpire(t,
		"ask dross about access", "play cylinder", "go back",
		"query 0997.4-1188", "note A20",
		"wait", "wait", "wait", "wait", "wait",
		"query 0997.4-1188",
	)
	if s.Has("A1_FILED") {
		t.Fatal("this route never files the shelf-address")
	}
	if !strings.Contains(log, "DISPATCH 0997.4-1188") {
		t.Fatal("the number is an address and the Vault answers to it")
	}
	if !s.Archive.Corrected() {
		t.Fatal("the archive reconciles either way")
	}
	if !s.Has("A21_HELD") {
		t.Fatal("and a player who wrote it down can still catch it")
	}
	if strings.Contains(log, "DELETED") || strings.Contains(log, "ERROR") {
		t.Fatal("the archive must never announce itself, by any route")
	}

	// Knowing the number is the requirement; guessing it is not.
	guess, guessLog := play(t, "query 0997.4-1188")
	if guess.Has("KESH_NAMED") {
		t.Fatal("you cannot requisition a dispatch number you have never seen")
	}
	if strings.Contains(guessLog, "DISPATCH 0997.4-1188") {
		t.Fatal("the Vault does not hand it to a guess")
	}
}

// She offers her own arrest before you ask her anything, and it is the correct
// procedural outcome, and the case closes cleanly on paper. Declining costs
// you that and buys the only person who can carry a message in Act III.
func TestTheArrestSheOffers(t *testing.T) {
	taken, log := playSpire(t, "arrest her", "go back")
	assertIlaria(t, taken, "ILARIA_DETAINED")
	if taken.Ints["ILARIA_TRUST"] != 0 {
		t.Fatalf("trust floors at nought: %d", taken.Ints["ILARIA_TRUST"])
	}
	if taken.Ints["CRASE_STANDING"] != 2 {
		t.Fatalf("the Crown is pleased: %d", taken.Ints["CRASE_STANDING"])
	}
	if !strings.Contains(log, "She hasn't done anything") {
		t.Fatal("Dross puts his name to something for the first time in twenty years")
	}

	declined, dLog := playSpire(t, "decline", "go back")
	assertIlaria(t, declined, "ILARIA_CLEARED")
	if declined.Ints["ILARIA_TRUST"] != 2 {
		t.Fatalf("declining, and saying why, is worth two: %d", declined.Ints["ILARIA_TRUST"])
	}
	if !strings.Contains(dLog, "it is a *route*") {
		t.Fatal("she wants a reason she can repeat to the Choir-Prefect")
	}

	quiet, _ := playSpire(t, "ask ilaria about her theory", "go back")
	assertIlaria(t, quiet, "")
}

func assertIlaria(t *testing.T, s *State, want string) {
	t.Helper()
	var set []string
	for _, f := range []string{"ILARIA_DETAINED", "ILARIA_CLEARED"} {
		if s.Has(f) {
			set = append(set, f)
		}
	}
	if want == "" {
		if len(set) != 0 {
			t.Fatalf("saying nothing must set nothing, got %v", set)
		}
		return
	}
	if len(set) != 1 || set[0] != want {
		t.Fatalf("expected exactly %s, got %v", want, set)
	}
}

// The wax is the only object in Act I that no cogitator can amend, and taking
// it costs the trust of the one person who will carry a message in Act III.
// There is no version where the player gets both.
func TestTheCylinderIsATrade(t *testing.T) {
	kept, keptLog := playSpire(t, "ask dross about access", "play cylinder", "decline", "leave it")
	if kept.Has("HAS_CYLINDER_6") {
		t.Fatal("left on the rack is left on the rack")
	}
	if kept.Ints["ILARIA_TRUST"] != 3 {
		t.Fatalf("declined and left: %d", kept.Ints["ILARIA_TRUST"])
	}
	if !strings.Contains(keptLog, "in a building the Lexicanum Vault has no index of") {
		t.Fatal("and it is still somewhere nothing can reach")
	}

	took, _ := playSpire(t, "ask dross about access", "play cylinder", "decline", "take cylinder")
	if !took.Has("HAS_CYLINDER_6") {
		t.Fatal("impounded is impounded")
	}
	if took.Ints["ILARIA_TRUST"] != 1 {
		t.Fatalf("procedurally correct, and she will not forget it: %d", took.Ints["ILARIA_TRUST"])
	}
}

// Named directly in Act I, the gaps shut her entirely. No pressure opens her;
// she has no reason to trust whoever holds her master's rosette.
func TestTheGapsAreNotForAskingAbout(t *testing.T) {
	s, log := playSpire(t, "decline", "ask ilaria about the gaps")
	if s.Ints["ILARIA_TRUST"] != 0 {
		t.Fatalf("two earned, two spent: %d", s.Ints["ILARIA_TRUST"])
	}
	if !strings.Contains(log, "is not a treatment") {
		t.Fatal("she explains exactly what asking costs her, and then stops")
	}
	if !strings.Contains(log, "Shall we get on") {
		t.Fatal("and offers her hands again")
	}

	// She raises it herself, on her own terms, and says nothing about it.
	plant, plantLog := playSpire(t,
		"ask dross about access", "play cylinder",
		"ask ilaria about her theory", "ask ilaria about transcription", "ask ilaria about vosk",
	)
	if !strings.Contains(plantLog, "counting for three years") {
		t.Fatal("the roster is the plant and it fires on its own")
	}
	if strings.Contains(plantLog, "She has been counting.\n\nShe tells you") {
		t.Fatal("she does not tell you what she heard")
	}
	if plant.Ints["ILARIA_TRUST"] != 0 {
		t.Fatal("and it costs and earns nothing, because she is not doing it for you")
	}
}
