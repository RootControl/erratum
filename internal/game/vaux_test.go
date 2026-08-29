package game

import (
	"strings"
	"testing"
)

// playVault drops the player into 2.9 the way the board does, with whatever
// approach and history the test is about.
func playVault(t *testing.T, flags []string, cmds ...string) (*State, string) {
	t.Helper()
	return playFrom(t, "cylinder_vault", append([]string{"ACT2_OPEN", "TWO_THREADS", "VAUX_ALIVE"}, flags...), cmds...)
}

// The four exchanges that open the three doors, in the order Vaux insists on.
var toTheDoors = []string{
	"ask vaux about the room",
	"ask vaux about vosk",
	"ask vaux about why me",
	"ask vaux about the cipher",
}

func TestVaultOpensAtTwoThreads(t *testing.T) {
	s, log := playFrom(t, "standing_board", []string{"ACT2_OPEN"}, "go vaux")
	if s.Room.ID != "standing_board" {
		t.Fatal("the Spire opened before two threads were touched")
	}
	if !strings.Contains(log, "has not answered anything") {
		t.Fatalf("expected the blocked line, got:\n%s", log)
	}

	s, _ = playFrom(t, "standing_board", []string{"ACT2_OPEN", "TWO_THREADS"}, "go vaux")
	if s.Room.ID != "cylinder_vault" {
		t.Fatalf("two threads should open the Spire; still in %s", s.Room.ID)
	}
}

// The approach is the whole of I-4, and neither half of it is announced.
func TestApproachSetsSuspicionSilently(t *testing.T) {
	s, log := playVault(t, nil, "ask vaux about the room")
	if s.Ints["VAUX_SUSPICION"] != 2 {
		t.Fatalf("archive approach should reach him at 2, got %d", s.Ints["VAUX_SUSPICION"])
	}
	for _, bad := range []string{"SUSPICION", "suspicion"} {
		if strings.Contains(log, bad) {
			t.Fatalf("the counter named itself:\n%s", log)
		}
	}

	s, _ = playVault(t, []string{"VAUX_RUNNER_SENT"}, "ask vaux about the room")
	if s.Ints["VAUX_SUSPICION"] != 0 {
		t.Fatalf("a child with a spoken message should cost nothing, got %d", s.Ints["VAUX_SUSPICION"])
	}
}

func TestTheDoorsOpenInOrder(t *testing.T) {
	s, log := playVault(t, nil, toTheDoors...)
	for _, f := range []string{"VAUX_FILE_SEEN", "VAUX_INDICTMENT_HEARD", "VAUX_QUESTION_ASKED", "VAUX_DOORS_OPEN"} {
		if !s.Has(f) {
			t.Fatalf("%s never fired:\n%s", f, log)
		}
	}
	// Beat 4: he is wrong about exactly one thing, and the prose says so
	// without anybody in the room saying it.
	if !strings.Contains(log, "twenty-three") && !strings.Contains(log, "pointed") {
		t.Fatalf("the indictment did not arrive:\n%s", log)
	}
}

// The ledger answer is the alliance's real hinge, and a player who filed
// everything cannot give it — not because a gate says no, but because there is
// nothing in the coat to read out.
func TestLedgerAnswerIsTheHinge(t *testing.T) {
	s, _ := playVault(t, []string{"STORE_HABIT_LEDGER"}, append(toTheDoors, "ask vaux about my answer")...)
	if !s.Has("VAUX_PERSUADED") || !s.Has("VAUX_ANSWERED") {
		t.Fatal("the ledger answer should persuade him")
	}
	if s.Ints["VAUX_SUSPICION"] != 1 {
		t.Fatalf("answering should take a point off, got %d", s.Ints["VAUX_SUSPICION"])
	}

	s, log := playVault(t, []string{"STORE_HABIT_FILED"}, append(toTheDoors, "ask vaux about my answer")...)
	if s.Has("VAUX_PERSUADED") {
		t.Fatal("a filed player has no answer to give")
	}
	if !strings.Contains(log, "correctly classed") {
		t.Fatalf("expected the filed variant:\n%s", log)
	}
}

// Physical evidence is the other route in, and it is the only one open to a
// player who has been filing since the chambers reopened.
func TestEvidenceOpensHimToo(t *testing.T) {
	s, _ := playVault(t, []string{"STORE_HABIT_FILED", "VAUX_RUNNER_SENT"})
	s.Grant("B12")
	run(s, append(append([]string{}, toTheDoors...), "show B12 to vaux", "ally")...)
	if !s.Has("VAUX_ALLIED") {
		t.Fatal("B12 shown should open the alliance")
	}
	if !s.Has("VAUX_EVIDENCE_SHOWN") {
		t.Fatal("a thing fetched with your own hands is the other way in")
	}
}

func TestAllianceNeedsMoreThanAsking(t *testing.T) {
	s, log := playVault(t, nil, append(toTheDoors, "ally")...)
	if s.Has("VAUX_ALLIED") {
		t.Fatal("he took a hand he was given no reason to take")
	}
	if !strings.Contains(log, "not going to explain it") {
		t.Fatalf("expected the refusal in his own voice:\n%s", log)
	}
	// And the refusal must not tell the player what the missing thing was.
	for _, bad := range []string{"LEDGER", "PERSUAD", "requires", "unavailable"} {
		if strings.Contains(log, bad) {
			t.Fatalf("the refusal explained itself: %s\n%s", bad, log)
		}
	}
}

// I-4's teeth: the standing bar, not the suspicion.
func TestArchiveApproachRaisesTheStandingBar(t *testing.T) {
	cmds := append(append([]string{}, toTheDoors...), "show B4 to vaux", "ally")

	// Standing opens at +1 and one more point is an ordinary month's work.
	archive, _ := playVault(t, []string{"STORE_HABIT_LEDGER"})
	archive.Grant("B4")
	archive.Add("CRASE_STANDING", 1)
	run(archive, cmds...)
	if archive.Ints["VAUX_SUSPICION"] > 1 {
		t.Fatal("the evidence should still have taken a point off")
	}
	if archive.Has("VAUX_ALLIED") {
		t.Fatal("standing +2 should close the archive route")
	}

	runner, _ := playVault(t, []string{"STORE_HABIT_LEDGER", "VAUX_RUNNER_SENT"})
	runner.Grant("B4")
	runner.Add("CRASE_STANDING", 1)
	run(runner, cmds...)
	if !runner.Has("VAUX_ALLIED") {
		t.Fatal("a runner should carry the alliance to standing +2")
	}
}

func TestExecutionTakesTheCipherWithHim(t *testing.T) {
	s, log := playVault(t, nil, append(toTheDoors, "execute")...)
	if !s.Has("VAUX_EXECUTED") || s.Has("VAUX_ALIVE") {
		t.Fatal("execution should end him")
	}
	if !s.Has("CIPHER_HALF_LOST") {
		t.Fatal("execution should lose the half and open the long way round")
	}
	if s.Ints["RECONCILIATION"] != 10 {
		t.Fatalf("expected 10 of exposure, got %d", s.Ints["RECONCILIATION"])
	}
	// Standing opens at +1, per the bible's systems notes.
	if s.Ints["CRASE_STANDING"] != 3 {
		t.Fatalf("expected +2 on the opening +1, got %d", s.Ints["CRASE_STANDING"])
	}
	if !strings.Contains(log, "thanks you for your discretion") {
		t.Fatalf("she has to mean the thanks:\n%s", log)
	}
}

// The corrupt half is the same object as the clean one until 2.11. Nothing
// here may hint otherwise.
func TestBlackmailNeverAnnouncesTheCorruption(t *testing.T) {
	s, _ := playVault(t, nil, toTheDoors...)
	out, _ := s.Do("blackmail")
	if !s.Has("VAUX_COERCED") || !s.Has("CIPHER_HALF_VAUX_CORRUPT") {
		t.Fatal("blackmail should coerce him and corrupt the half")
	}
	if _, held := s.Held["B16"]; !held {
		t.Fatal("blackmail still hands over the cipher")
	}
	// The half he gives under coercion is deliberately mis-ordered and the
	// player does not find that out until 2.11. Nothing on this screen, and
	// nothing the clue itself says when it is written down, may hint at it.
	said := strings.Join(out, "\n") + "\n" + s.NotebookDump()
	s.Do("note B16")
	said += "\n" + strings.Join(mustDo(s, "notebook"), "\n")
	for _, bad := range []string{"corrupt", "mis-order", "misorder", "wrong order", "out of order", "false"} {
		if strings.Contains(strings.ToLower(said), bad) {
			t.Fatalf("the coerced half told on itself: %s\n%s", bad, said)
		}
	}
}

func mustDo(s *State, cmd string) []string {
	out, _ := s.Do(cmd)
	return out
}

// run drives an already-built state, for tests that have to put something in
// the player's hands before the conversation starts.
func run(s *State, cmds ...string) {
	for _, c := range cmds {
		s.Do(c)
	}
}

// Refuse, walk out, come back: one door, and the two that closed do not
// render as closed.
func TestRefusalReturnsHimAtThree(t *testing.T) {
	s, _ := playVault(t, []string{"STORE_HABIT_LEDGER"},
		append(append([]string{}, toTheDoors...), "refuse", "go out")...)
	if !s.Has("VAUX_LEFT") || s.Has("VAUX_RESOLVED") {
		t.Fatal("refusing is not resolving")
	}
	if !s.Move("cylinder_vault") {
		t.Fatal("cannot get back to the vault")
	}
	if out, _ := s.Do("talk to vaux"); len(out) == 0 {
		t.Fatal("he should be there")
	}
	s.Do("ask vaux about the room")
	if s.Ints["VAUX_SUSPICION"] != 3 || !s.Has("VAUX_FINISHED") {
		t.Fatalf("three weeks in a transit hostel: suspicion %d, finished %v",
			s.Ints["VAUX_SUSPICION"], s.Has("VAUX_FINISHED"))
	}
	s.Do("ally")
	if s.Has("VAUX_ALLIED") {
		t.Fatal("the hand is not on offer any more")
	}
	s.Do("execute")
	if s.Has("VAUX_EXECUTED") {
		t.Fatal("he is out of reach of a warrant")
	}
	s.Do("blackmail")
	if !s.Has("VAUX_COERCED") {
		t.Fatal("the ugly door is the one that stays open")
	}
}

// I-5: the player is punished here for playing the politics well, and the
// deleted options are absent rather than refused.
func TestHighStandingDeletesPersuasion(t *testing.T) {
	s, _ := playFrom(t, "cylinder_vault", []string{"ACT2_OPEN", "TWO_THREADS", "VAUX_ALIVE"}, "status")
	s.Add("CRASE_STANDING", 4)
	var log strings.Builder
	for _, c := range append(append([]string{}, toTheDoors...), "ask vaux about his character", "ask vaux about the indictment") {
		out, _ := s.Do(c)
		log.WriteString(strings.Join(out, "\n") + "\n")
	}
	if s.Ints["VAUX_SUSPICION"] != 2 {
		t.Fatalf("neither persuasion should have landed, got %d", s.Ints["VAUX_SUSPICION"])
	}
	if !strings.Contains(log.String(), "promoted you again last week") {
		t.Fatalf("high standing adds a line to the indictment:\n%s", log.String())
	}
}

// The habit is behaviour, not a menu, and the query tier moves without saying
// that anything moved.
func TestStoreHabitIsNeverChosen(t *testing.T) {
	s, log := playFrom(t, "ferrum", []string{"ACT2_OPEN", "SARN_INTERVIEWED"},
		"ask sarn about historians", "file B12")
	if !s.Has("STORE_HABIT_FILED") || s.Has("STORE_HABIT_LEDGER") {
		t.Fatal("the first filing in Act II is the habit")
	}
	if s.Ints["QUERY_TIER"] != 1 {
		t.Fatalf("filing pleases the Vault, got tier %d", s.Ints["QUERY_TIER"])
	}
	if strings.Contains(log, "HABIT") || strings.Contains(log, "QUERY_TIER") {
		t.Fatalf("the habit announced itself:\n%s", log)
	}

	s, _ = playFrom(t, "ferrum", []string{"ACT2_OPEN", "SARN_INTERVIEWED"},
		"ask sarn about historians", "note B12", "file B13")
	if !s.Has("STORE_HABIT_LEDGER") || s.Has("STORE_HABIT_FILED") {
		t.Fatal("the habit is the first reach, and it does not change its mind")
	}
}

// Every object in the vault, twice, plus the door out — the sweep that catches
// a string authored as a field and read as a key.
func TestVaultRendersEverything(t *testing.T) {
	s, log := playVault(t, nil,
		"look", "x cylinders", "x cylinders", "x cylinders",
		"x reproducer", "x reproducer", "x case", "x case", "x vaux",
		"talk to vaux", "topics", "ask vaux about the weather",
		"status", "go out",
	)
	if s.Room.ID != "standing_board" {
		t.Fatalf("the door out of the vault does not work; still in %s", s.Room.ID)
	}
	for _, want := range []string{
		"A record that cannot be searched cannot be corrected.",
		"recommend continue",
		"re-stitched twice by hand",
		"I have stopped deciding in advance",
	} {
		if !strings.Contains(log, want) {
			t.Fatalf("missing authored prose: %q\n%s", want, log)
		}
	}
	// Nothing down here is a terminal, and the room was not arranged.
	if out, _ := s.Do("query CANTICA-NULL"); len(out) == 0 {
		t.Fatal("expected a refusal, not silence")
	}
}
