package game

import (
	"strings"
	"testing"
)

// Player goal: accept. There is no refusing, and the game enforces it with a
// door rather than with a refusal message, because the door is diegetic and
// the refusal message would be the game talking about itself.
func TestThereIsNoRefusing(t *testing.T) {
	stalling, log := playRosette(t, "go in")
	if stalling.Room.ID != "antechamber" {
		t.Fatal("nobody walks past Bez without a rosette")
	}
	if !strings.Contains(log, "records-adept with eleven years in") {
		t.Fatal("the closed door should say why in her voice, not the parser's")
	}

	elevated, _ := playRosette(t, "take rosette", "go in")
	if elevated.Room.ID != "vosk_chambers" {
		t.Fatal("accepting opens the room")
	}
	if !elevated.Has("ROSETTE_TAKEN") {
		t.Fatal("expected ROSETTE_TAKEN")
	}
}

// "Cannot fail. The game will not lose it." The Choir lead is the most useful
// sentence anyone says in Act I, so it is part of the ceremony and not part of
// the conversation — there is no way to be handed the rosette without it.
func TestChoirLeadCannotBeMissed(t *testing.T) {
	silent, _ := playRosette(t, "take rosette", "go in")
	if !silent.Has("CRASE_LEAD_CHOIR") {
		t.Fatal("a player who asked nothing still leaves with the lead")
	}
	if _, held := silent.Held["A1-lead"]; !held {
		t.Fatal("the lead is a clue, in hand, whether or not they wanted it")
	}

	// It is verbal, so it starts in neither store, and writing it down is the
	// player's own move — the free demonstration that some things never enter
	// the cogitator at all.
	if silent.Has("A1-lead_HELD") {
		t.Fatal("hearing it is not recording it")
	}
	noted, _ := playRosette(t, "take rosette", "note cylinders")
	if noted.Held["A1-lead"] != StoreBook || len(noted.Notebook) != 1 {
		t.Fatal("NOTE CYLINDERS should put it in the coat")
	}
}

// Her warning about Kesh has three answers and the third one is saying
// nothing. Unlike Vaux's disposition there is no default: declining to answer
// a friend is not the same as refusing her, and the record should not pretend
// it is.
func TestKeshWarningIsThreeAnswers(t *testing.T) {
	agreed, _ := playRosette(t, "take rosette", "ask about kesh", "promise", "go in")
	assertKesh(t, agreed, "KESH_WARNED_OFF")
	if agreed.Ints["CRASE_STANDING"] != 1 {
		t.Fatalf("agreeing costs and gains nothing: %d", agreed.Ints["CRASE_STANDING"])
	}

	defied, _ := playRosette(t, "take rosette", "ask about kesh", "defy", "go in")
	assertKesh(t, defied, "KESH_DEFIED")
	if defied.Ints["CRASE_STANDING"] != 0 {
		t.Fatalf("she starts logging you three weeks early: %d", defied.Ints["CRASE_STANDING"])
	}

	quiet, _ := playRosette(t, "take rosette", "ask about kesh", "go in")
	assertKesh(t, quiet, "")
}

func assertKesh(t *testing.T, s *State, want string) {
	t.Helper()
	var set []string
	for _, f := range []string{"KESH_WARNED_OFF", "KESH_DEFIED"} {
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

// The first thing anyone tells you that is not true, told by somebody who is
// not lying. Crase has been reconciled and from inside her head it is true, so
// the game must render it as a plain answer with nothing marking it.
func TestCraseBelievesTheLie(t *testing.T) {
	s, log := playRosette(t, "ask about his last cycle")
	if !strings.Contains(log, "He'd stopped taking cases") {
		t.Fatal("she gives the answer she believes")
	}
	if !s.Has("CRASE_SAYS_NO_CASES") {
		t.Fatal("Act II needs to know she said this")
	}
	for _, tell := range []string{"[", "reconcil", "Reconcil", "RECONCIL", "false", "lie"} {
		if strings.Contains(log, tell) {
			t.Fatalf("nothing may mark the answer as wrong: found %q", tell)
		}
	}
}

// A question asked in 1.2 buys the only warm thing Vaux says in 1.5. Flags
// cross scenes by name, with no translation layer, so this is one string
// written in the antechamber and read on tier four.
func TestWhyNotVauxCrossesScenes(t *testing.T) {
	warm, warmLog := playRosette(t,
		"ask why not vaux",
		"take rosette", "go in",
		"go vaux", "ask about why not you",
	)
	if !warm.Has("ASKED_WHY_NOT_VAUX") {
		t.Fatal("1.2 writes the flag")
	}
	if !strings.Contains(warmLog, "he never talked about me") {
		t.Fatal("1.5 reads it")
	}
	if warm.Ints["CRASE_STANDING"] != 2 {
		t.Fatalf("the honest question is worth a point: %d", warm.Ints["CRASE_STANDING"])
	}

	cold, coldLog := playRosette(t,
		"take rosette", "go in",
		"go vaux", "ask about why not you",
	)
	if cold.Has("ASKED_WHY_NOT_VAUX") {
		t.Fatal("1.5 must never write a flag that belongs to 1.2")
	}
	if strings.Contains(coldLog, "he never talked about me") {
		t.Fatal("unasked in 1.2, the line does not exist in 1.5")
	}
}

// The scenes are hours apart and the player does not play the hours between
// them. The eight minutes belong to the body, not to the clock: they start
// when the player is close enough to read something off him.
func TestClockJumpsAndTheBodyWaits(t *testing.T) {
	cmds := []string{"take rosette"}
	for i := 0; i < 15; i++ { // ten minutes of dithering in front of the adepts
		cmds = append(cmds, "wait")
	}
	s, _ := playRosette(t, cmds...)
	if !strings.Contains(s.TimeString(), "/ 04h") {
		t.Fatalf("1.2 runs at the fourth hour: %s", s.TimeString())
	}
	if !s.BodyPresent {
		t.Fatal("the medicae are not taking him while the player is in a ceremony")
	}

	out, ticks := s.Do("go in")
	_, _ = out, ticks
	for i := 0; i < ticks; i++ {
		s.Tick()
	}
	if !strings.HasSuffix(s.TimeString(), "06h12") {
		t.Fatalf("1.3 opens at twelve past the sixth: %s", s.TimeString())
	}
	if !s.BodyPresent {
		t.Fatal("the eight minutes start now, not an hour and a half ago")
	}

	for i := 0; i < 12; i++ {
		s.Do("wait")
		s.Tick()
	}
	if s.BodyPresent {
		t.Fatal("and they are eight minutes")
	}
}
