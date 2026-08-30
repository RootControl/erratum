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
