package game

import (
	"strings"
	"testing"

	"github.com/RootControl/erratum/internal/content"
)

func loadGame(t *testing.T) (*content.Game, error) {
	t.Helper()
	return content.Load()
}

// play drops the player into the chambers at the top of 1.3, with the flag its
// entry gate assumes. The guarantees below are about 1.3-1.5, and re-playing
// the elevation in front of every one of them would only make them harder to
// read. Scene 1.2 has its own tests.
func play(t *testing.T, cmds ...string) (*State, string) {
	t.Helper()
	return playFrom(t, "vosk_chambers", []string{"ROSETTE_TAKEN"}, cmds...)
}

// playAct1 starts where the game starts: 1.1, in the doorway, with nothing.
func playAct1(t *testing.T, cmds ...string) (*State, string) {
	t.Helper()
	return playFrom(t, "", nil, cmds...)
}

// playRosette drops the player into 1.2 with the elevation about to happen.
func playRosette(t *testing.T, cmds ...string) (*State, string) {
	t.Helper()
	return playFrom(t, "antechamber", []string{"ACT1_OPEN"}, cmds...)
}

func playFrom(t *testing.T, scene string, flags []string, cmds ...string) (*State, string) {
	t.Helper()
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	s := New(g)
	for _, f := range flags {
		s.Set(f)
	}
	if scene != "" && !s.Move(scene) {
		t.Fatalf("no scene %q", scene)
	}
	var log strings.Builder
	for _, c := range cmds {
		out, ticks := s.Do(c)
		log.WriteString("> " + c + "\n")
		for _, o := range out {
			log.WriteString(o + "\n")
		}
		for i := 0; i < ticks; i++ {
			for _, e := range s.Tick() {
				log.WriteString(e + "\n")
			}
		}
	}
	// A key that is not in a scene's table renders its own name and keeps
	// going, which is how LOOK and the body's removal both shipped broken.
	// Every playthrough in this package guards against it.
	if out := log.String(); strings.Contains(out, "[missing string:") {
		i := strings.Index(out, "[missing string:")
		t.Fatalf("unresolved authored string: %s", out[i:min(i+48, len(out))])
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
	// ACT_I gives VERITY four sources — A21, the Lexmechanic checksum, the
	// Précis annotated or refused, and cylinder six surviving to act close.
	// This run reaches only the first, and that is the whole of what a
	// careful reader has by the end of 1.9: one contradiction, recorded
	// outside the cogitator, in their own hand.
	if s.Ints["VERITY"] != 1 {
		t.Fatalf("expected the one contradiction, got %d", s.Ints["VERITY"])
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
	// Standing opens at +1 per the bible; the requisition rescue is the second.
	if s.Ints["CRASE_STANDING"] != 2 {
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
