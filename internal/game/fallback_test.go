package game

import (
	"strings"
	"testing"
)

// The act documents give every skippable scene a punishing default, and a
// skipped scene that simply does not happen is the one thing the design says
// must not occur. These are those defaults.

// KADMON_PROOF is documented "always obtainable". It was gated on
// VURR_RESPECT >= 1, so a player who pitied the Marchioness could never reach
// the Reading Room at all.
func TestRhodSmugglesTheProof(t *testing.T) {
	s, _ := playFrom(t, "kadmon", landed, "look")
	s.Set("VURR_ASKED")
	var log []string
	for i := 0; i < 8; i++ {
		log = append(log, s.Tick()...)
	}
	if _, held := s.Held["C12"]; !held {
		t.Fatalf("the proof is always obtainable:\n%s", strings.Join(log, "\n"))
	}
	if !s.Has("KADMON_PROOF") || !s.Has("RHOD_SMUGGLED") {
		t.Fatal("expected the smuggled copy")
	}
	flat := flat(strings.Join(log, "\n"))
	if !strings.Contains(flat, "He was not asked. He was not hinted at.") {
		t.Fatalf("he decides it himself:\n%s", flat)
	}
	if !strings.Contains(flat, "you are never told about it") {
		t.Fatalf("what is done to him happens off-screen:\n%s", flat)
	}

	// A player who earned the room gets it the ordinary way and no copy.
	earned, _ := playFrom(t, "kadmon", landed, "look")
	earned.Set("VURR_ASKED")
	earned.Add("VURR_RESPECT", 1)
	for i := 0; i < 8; i++ {
		earned.Tick()
	}
	if earned.Has("RHOD_SMUGGLED") {
		t.Fatal("he does not risk it for somebody she trusted")
	}
}

// Walk out of the reading and it follows you, on printed slips, in order.
// Stopping it is a different decision and is not followed.
func TestTheReadingFollowsYouOut(t *testing.T) {
	s, _ := playFrom(t, "correctorum", docked,
		"ask vox about minutes", "ask vox about kadmon", "go up")
	if s.Has("FILE_HEARD") {
		t.Fatal("setup: the player left before C19")
	}
	var log []string
	for i := 0; i < 6; i++ {
		log = append(log, s.Tick()...)
	}
	if !s.Has("FILE_HEARD") {
		t.Fatalf("it follows:\n%s", strings.Join(log, "\n"))
	}
	if _, held := s.Held["C19"]; !held {
		t.Fatal("expected the requisition")
	}
	if !strings.Contains(flat(strings.Join(log, "\n")), "a slip is what you take away") {
		t.Fatalf("nobody pushed them:\n%s", strings.Join(log, "\n"))
	}

	// Stopping it is the branch that costs C19, and it stays stopped.
	stopped, _ := playFrom(t, "correctorum", docked,
		"ask vox about minutes", "ask vox about kadmon", "stop it", "go up")
	for i := 0; i < 10; i++ {
		stopped.Tick()
	}
	if stopped.Has("FILE_HEARD") {
		t.Fatal("a closed request is not walked up to your cabin")
	}
}

// You do not choose the chair. You find you have already sat down.
func TestInactionSitsYouDown(t *testing.T) {
	s, _ := playFrom(t, "intake", atIntake, "look")
	s.Add("RECONCILIATION", 90)
	var log []string
	for i := 0; i < 3; i++ {
		log = append(log, s.Tick()...)
	}
	if !s.Has("ENDING_B") {
		t.Fatalf("expected the inaction ending:\n%s", strings.Join(log, "\n"))
	}
	if !strings.Contains(flat(strings.Join(log, "\n")), "You do not remember deciding.") {
		t.Fatalf("expected the line:\n%s", strings.Join(log, "\n"))
	}

	// Below ninety it never fires, and nothing says the player was spared it.
	clean, _ := playFrom(t, "intake", atIntake, "look")
	clean.Add("RECONCILIATION", 80)
	for i := 0; i < 20; i++ {
		clean.Tick()
	}
	if clean.Ending() != "" {
		t.Fatal("a player who kept themselves clean still chooses")
	}
}

// If the well is never entered, Bez has the slate and the station has the plot.
func TestTheHullsEvidenceSurfacesAnyway(t *testing.T) {
	s, log := playFrom(t, "gunbay", append(aboard, "IN_TRANSIT"), "talk to bez")
	if _, held := s.Held["C4"]; !held {
		t.Fatalf("she had it in the kit all along:\n%s", log)
	}
	if !strings.Contains(flat(log), "That's him about to say something else") {
		t.Fatalf("she finds the fourth second either way:\n%s", log)
	}

	k, klog := playFrom(t, "cantica", docked, "talk to kesh")
	if _, held := k.Held["C2"]; !held {
		t.Fatalf("a courier files a plot or it does not get a berth:\n%s", klog)
	}
	if !strings.Contains(flat(klog), "I would like several things to have come from me") {
		t.Fatalf("expected her line:\n%s", klog)
	}

	// And neither fires for a player who went down the well themselves.
	own, ownLog := playFrom(t, "gunbay",
		append(aboard, "IN_TRANSIT", "SPIRIT_ORDERS_READ"), "talk to bez")
	if strings.Contains(ownLog, "I'd have mentioned it sooner") {
		t.Fatal("she does not hand you a thing you already have")
	}
	_ = own
}

// She follows you to the intake if you never went to the gallery.
func TestCraseFollowsYouToTheSlot(t *testing.T) {
	s, log := playFrom(t, "intake", atIntake, "talk to crase")
	if !strings.Contains(log, "fifth chair") && !strings.Contains(log, "eleven weeks") {
		t.Fatalf("she is there:\n%s", log)
	}
	s.Do("ask crase about muria")
	s.Do("audit her")
	if !s.Has("CRASE_AUDITED") {
		t.Fatal("the account is still available at the slot")
	}
}

// C16 is re-offered unprompted, and she will not argue about a number you have
// not seen.
func TestTheLedgerIsReOffered(t *testing.T) {
	_, log := playFrom(t, "fifth_chair", docked, "talk to crase")
	if !strings.Contains(flat(log), "Has Anneke given you the ledger?") {
		t.Fatalf("she asks first:\n%s", log)
	}
	if !strings.Contains(flat(log), "I would win, and it would be worthless") {
		t.Fatalf("expected her reason:\n%s", log)
	}

	held, hlog := playFrom(t, "fifth_chair", docked, "look")
	held.Grant("C16")
	held.storeClue("C16", StoreBook)
	held.Set("C16_HELD")
	out, _ := held.Do("talk to crase")
	hlog = strings.Join(out, "\n")
	if strings.Contains(hlog, "Has Anneke given you") {
		t.Fatalf("not to somebody who has it:\n%s", hlog)
	}
}

// The Reconciliation symptoms. Nothing below is ever announced, and none of
// them is a message from the archive — they are all the record being ordinary
// about something the player knows is not.

// ≥ 40: the retinue remembers what you don't. She raises it unprompted, which
// is what makes it unanswerable — there is no topic you failed to ask.
func TestBezRemembersSubLevelNine(t *testing.T) {
	s, _ := playFrom(t, "standing_board",
		[]string{"ACT2_OPEN", "SUBLEVEL_NINE_SEEDED"}, "look")
	s.Add("RECONCILIATION", 40)
	var log []string
	for i := 0; i < 10; i++ {
		log = append(log, s.Tick()...)
	}
	joined := flat(strings.Join(log, "\n"))
	if !s.Has("BEZ_REMEMBERS") {
		t.Fatalf("she says it in her own time:\n%s", joined)
	}
	if !strings.Contains(joined, "don't tell me you didn't, I was awake for most of it") {
		t.Fatalf("expected her line:\n%s", joined)
	}
	if !strings.Contains(joined, "You have never said that to her.") {
		t.Fatalf("expected the turn:\n%s", joined)
	}
	// And the Vault answers about an operation that never was, in the words it
	// uses for a shelf that was never filled.
	out, _ := s.Do("query SUBLEVEL NINE")
	answer := flat(strings.Join(out, "\n"))
	if !strings.Contains(answer, "NO SUCH OPERATION.") {
		t.Fatalf("expected the null: %s", answer)
	}
	if !strings.Contains(answer, "the same spacing it uses for a shelf somebody emptied") {
		t.Fatalf("the spacing is the whole of it: %s", answer)
	}

	// Below forty it never happens, and a clean player is never told they were
	// spared it.
	clean, _ := playFrom(t, "standing_board",
		[]string{"ACT2_OPEN", "SUBLEVEL_NINE_SEEDED"}, "look")
	clean.Add("RECONCILIATION", 35)
	for i := 0; i < 20; i++ {
		clean.Tick()
	}
	if clean.Has("BEZ_REMEMBERS") {
		t.Fatal("thirty-five is not forty")
	}
}

// ≥ 55: the retinue's own files come back shorter, with every cross-reference
// still closing and nothing struck.
func TestTheRetinueFilesThin(t *testing.T) {
	full, _ := playFrom(t, "standing_board", []string{"ACT2_OPEN"}, "look")
	for _, want := range []struct{ q, phrase string }{
		{"query TALM", "does not agree with the transmission"},
		{"query FENN", "un-networked by design"},
		{"query ODUYA", "Literacy: not assessed."},
	} {
		out, _ := full.Do(want.q)
		if !strings.Contains(flat(strings.Join(out, "\n")), want.phrase) {
			t.Errorf("%s is missing %q:\n%s", want.q, want.phrase, strings.Join(out, "\n"))
		}
	}

	thin, _ := playFrom(t, "standing_board", []string{"ACT2_OPEN"}, "look")
	thin.Add("RECONCILIATION", 55)
	// The lines that go are the record's lines. What the player thinks about
	// their absence is their own and stays on the screen.
	for _, gone := range []struct{ q, phrase string }{
		{"query TALM", "does not agree with the transmission"},
		{"query FENN", "un-networked by design"},
		{"query ODUYA", "Literacy: not assessed."},
	} {
		out, _ := thin.Do(gone.q)
		joined := strings.Join(out, "\n")
		if strings.Contains(joined, gone.phrase) {
			t.Errorf("%s still carries %q", gone.q, gone.phrase)
		}
		// Nothing is struck, marked or amended. It is simply a clean file.
		for _, bad := range []string{"DELETED", "STRUCK", "AMENDED", "REDACTED"} {
			if strings.Contains(joined, bad) {
				t.Errorf("%s announced itself: %s", gone.q, bad)
			}
		}
	}
}

// ≥ 70: you go. The record is ordinary about it and there is nobody to raise
// it with who would not ask you to spell the name you say is missing.
func TestAtSeventyYouGo(t *testing.T) {
	s, _ := playFrom(t, "standing_board", []string{"ACT2_OPEN"}, "look")
	out, _ := s.Do("query ROSETTE")
	if !strings.Contains(flat(strings.Join(out, "\n")), "The diacritic is wrong.") {
		t.Fatalf("your name is on it:\n%s", strings.Join(out, "\n"))
	}

	gone, _ := playFrom(t, "standing_board", []string{"ACT2_OPEN"}, "look")
	gone.Add("RECONCILIATION", 70)
	out, _ = gone.Do("query ROSETTE")
	joined := flat(strings.Join(out, "\n"))
	if !strings.Contains(joined, "Where the given name was there is an initial.") {
		t.Fatalf("expected the initial:\n%s", joined)
	}
	if !strings.Contains(joined, "It is the most ordinary thing in the world.") {
		t.Fatalf("the record is ordinary about it:\n%s", joined)
	}
	if strings.Contains(joined, "DELETED") || strings.Contains(joined, "REDACTED") {
		t.Fatal("nothing announces it")
	}
}

// Options are removed, not closed. The person answers as though they had never
// heard of the subject, and there is no gap where the option was.
func TestAtSeventyOptionsAreAbsent(t *testing.T) {
	open, log := playFrom(t, "fifth_chair", docked,
		"talk to crase", "ask crase about muria", "ask crase about the case")
	if !strings.Contains(flat(log), "It was Hadrian.") {
		t.Fatalf("the reframe is there below seventy:\n%s", log)
	}
	_ = open

	s, _ := playFrom(t, "fifth_chair", docked, "ask crase about muria")
	s.Add("RECONCILIATION", 70)
	out, _ := s.Do("ask crase about the case")
	answer := strings.Join(out, "\n")
	if strings.Contains(answer, "It was Hadrian.") {
		t.Fatal("past seventy it is not a subject she has")
	}
	if !strings.Contains(answer, "I don't know that") && !strings.Contains(answer, "Ask me one I do") {
		t.Fatalf("she answers as though she had never heard of it:\n%s", answer)
	}
	for _, bad := range []string{"unavailable", "no longer", "removed", "RECONCILIATION"} {
		if strings.Contains(answer, bad) {
			t.Fatalf("there is no gap to notice: %s", bad)
		}
	}
}

// 3.4b's punishing default. Skip the cell and she raises it herself, on a
// street, at the worst possible hour, and does not accept the answer.
func TestIlariaFiresOneSidedAtKadmon(t *testing.T) {
	s, log := playFrom(t, "kadmon", landed, "talk to ilaria")
	if !s.Has("ILARIA_ONE_SIDED") {
		t.Fatalf("she is not prepared to be that patient:\n%s", log)
	}
	if _, held := s.Held["C8"]; !held {
		t.Fatal("expected the tally")
	}
	if !strings.Contains(flat(log), "there's a shrine queue eleven metres away") &&
		!strings.Contains(flat(log), "this is where you are going to have this") {
		t.Fatalf("the place is the punishment:\n%s", log)
	}

	// The clean outcome is not available. Telling her the truth is still
	// right, still costs nothing, and does not land.
	out, _ := s.Do("tell her the truth")
	answer := flat(strings.Join(out, "\n"))
	if s.Has("ILARIA_KNOWS") {
		t.Fatal("she does not accept it")
	}
	if !s.Has("ILARIA_TOLD_LATE") {
		t.Fatalf("the scene still happens:\n%s", answer)
	}
	// But it is not a watch stood. The voyage cannot be completed after it.
	if s.Has("WATCH_ILARIA") {
		t.Fatal("a street on Kadmon is not the cell")
	}
	if !strings.Contains(answer, "No, I don't think I'll take that today, sir.") {
		t.Fatalf("expected her refusal:\n%s", answer)
	}
	if s.Ints["ILARIA_TRUST"] != 0 {
		t.Fatalf("and it buys nothing: trust %d", s.Ints["ILARIA_TRUST"])
	}

	// So she is not the machine's mouth at Cantica-Null.
	if !s.Move("correctorum") {
		t.Fatal("cannot reach the floor")
	}
	run(s, "ask vox about minutes", "ask vox about kadmon", "ask vox about my file")
	if s.Has("ILARIA_STOPPED") {
		t.Fatal("she stands at the back and nobody thinks to ask her")
	}
}

// A player who stood the watch in the cell never meets the fallback.
func TestTheCellClosesTheFallback(t *testing.T) {
	s, log := playFrom(t, "kadmon", append(landed, "WATCH_ILARIA"), "talk to ilaria")
	if s.Has("ILARIA_ONE_SIDED") {
		t.Fatalf("she asked you once and you answered:\n%s", log)
	}
}

// 3.4c's default: Vaux has been writing it down in a corridor, and nine pages
// of a sixty-one-year-old with a pen cannot force anything.
func TestVauxReconstructsItUncitably(t *testing.T) {
	s, log := playFrom(t, "wardroom", aboard, "ask vaux about the projection")
	if _, held := s.Held["C9b"]; !held {
		t.Fatalf("he has been writing it down:\n%s", log)
	}
	if !strings.Contains(flat(log), "there is not one word of it I could swear to") {
		t.Fatalf("he says what it is worth before you take it:\n%s", log)
	}
	if s.Has("C9_TAKEN") || s.Has("C9_PHYS") {
		t.Fatal("a transcription is not the recitation")
	}

	// And it cannot be cited, because there is no record to contradict.
	if !s.Move("intake") {
		t.Fatal("cannot reach the slot")
	}
	s.storeClue("C9b", StoreBook)
	out, _ := s.Do("cite the reconstruction")
	answer := flat(strings.Join(out, "\n"))
	if s.Has("OMISSION_EXPOSED") {
		t.Fatal("secondhand does not force it")
	}
	if !strings.Contains(answer, "There is no record to contradict.") {
		t.Fatalf("expected the refusal:\n%s", answer)
	}
	if !strings.Contains(answer, "You take them back.") {
		t.Fatalf("it offers the intake instead:\n%s", answer)
	}
}

// A player who took it down themselves never gets the reconstruction.
func TestTheHoldClosesTheReconstruction(t *testing.T) {
	s, log := playFrom(t, "wardroom", append(aboard, "C9_TAKEN"),
		"ask vaux about the projection")
	if _, held := s.Held["C9b"]; held {
		t.Fatalf("you have the real one:\n%s", log)
	}
}
