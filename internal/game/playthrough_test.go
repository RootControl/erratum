package game

import (
	"strings"
	"testing"

	"github.com/RootControl/erratum/internal/content"
)

func play(t *testing.T, cmds ...string) (*State, string) {
	t.Helper()
	r, err := content.Load()
	if err != nil {
		t.Fatal(err)
	}
	s := New(r)
	var log strings.Builder
	for _, c := range cmds {
		out, timed := s.Do(c)
		log.WriteString("> " + c + "\n")
		for _, o := range out {
			log.WriteString(o + "\n")
		}
		if timed {
			for _, e := range s.Tick() {
				log.WriteString(e + "\n")
			}
		}
	}
	return s, log.String()
}

// The careful path: read everything twice, keep the boots in your coat,
// transcribe the dispatch, and catch the archive in the act.
func TestCarefulPlaythrough(t *testing.T) {
	s, log := play(t,
		"x boots", "x boots", "note A3",
		"x board", "x board", "file A1",
		"x shutter", "x shutter", "note A2",
		"x quill", "x quill", "file A4",
		"x calendar", "x calendar", "note A5",
		"read room", "read room backward",
		"query C9-F4-D2",
		"note A20",
		"wait", "wait", "wait", "wait", "wait",
		"query C9-F4-D2",
	)
	if !s.Has("ROOM_READ_FULL") {
		t.Fatal("expected ROOM_READ_FULL on the careful path")
	}
	if !s.Has("CIPHER_HALF_ROOM") {
		t.Fatal("expected CIPHER_HALF_ROOM")
	}
	if !s.Archive.Corrected() {
		t.Fatal("archive should have reconciled by now")
	}
	if !s.Has("A21_HELD") {
		t.Fatal("player transcribed the dispatch and should be able to prove the edit")
	}
	if s.Ints["VERITY"] < 4 {
		t.Fatalf("VERITY too low on the careful path: %d", s.Ints["VERITY"])
	}
	if !strings.Contains(log, "1188") {
		t.Fatal("dispatch never rendered")
	}
	t.Logf("VERITY=%d RECONCILIATION=%d CRASE=%+d",
		s.Ints["VERITY"], s.Ints["RECONCILIATION"], s.Ints["CRASE_STANDING"])
}

// The trusting path: file everything, write nothing, and end the scene with a
// strong feeling and no evidence. This is a designed outcome, not a failure.
func TestTrustingPlaythrough(t *testing.T) {
	s, log := play(t,
		"x board", "x board", "file A1",
		"x shutter", "x shutter", "file A2",
		"x quill", "x quill", "file A4",
		"read room", "read room forward",
		"query D2-F4-C9",
		"wait", "wait", "wait", "wait", "wait",
		"query C9-F4-D2", "notes",
	)
	if s.Has("ROOM_READ_FULL") {
		t.Fatal("no boots, no full reading")
	}
	if !s.Has("ROOM_READ_PARTIAL") {
		t.Fatal("expected a partial reading")
	}
	if s.Has("A21_HELD") {
		t.Fatal("player filed everything; they must NOT be able to prove the edit")
	}
	if s.Ints["CRASE_STANDING"] != 1 {
		t.Fatalf("Crase should have saved them once: %d", s.Ints["CRASE_STANDING"])
	}
	if !strings.Contains(log, "SANITATION REQUISITION LEDGERS") {
		t.Fatal("case-notes were not substituted")
	}
	if strings.Contains(log, "DELETED") || strings.Contains(log, "ERROR") {
		t.Fatal("the archive must never announce itself")
	}
}

// The boots cannot be filed. There is no field on the form for what they mean.
func TestBootsCannotBeFiled(t *testing.T) {
	s, log := play(t, "x boots", "x boots", "file A3")
	if s.Held["A3"] == StoreCog || s.Held["A3"] == StoreBoth {
		t.Fatal("A3 must never reach the cogitator")
	}
	if !strings.Contains(log, "SELF-NEGLECT") {
		t.Fatal("the filing template should offer only the three wrong answers")
	}
}

// The body leaves on a timer the player cannot stop, and takes A3 with it.
func TestBodyTimerClosesTheBoots(t *testing.T) {
	cmds := []string{}
	for i := 0; i < 13; i++ {
		cmds = append(cmds, "wait")
	}
	cmds = append(cmds, "x boots")
	s, log := play(t, cmds...)
	if s.BodyPresent {
		t.Fatal("body should be gone after eight minutes")
	}
	if s.Holds("A3") {
		t.Fatal("A3 must be unreachable once the body is removed")
	}
	if !strings.Contains(log, "medicae") {
		t.Fatal("removal should be narrated")
	}
}

// You cannot requisition against a shelf-code you only told your coat about.
func TestVaultRequiresFiling(t *testing.T) {
	_, log := play(t, "x board", "x board", "query C9-F4-D2")
	if !strings.Contains(log, "FILE A1") {
		t.Fatal("the Vault should refuse an unfiled address and say so")
	}
}
