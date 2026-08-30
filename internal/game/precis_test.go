package game

import (
	"strings"
	"testing"
)

// playToPrecis runs the Act I turn and then waits out the day.
func playToPrecis(t *testing.T, pre ...string) func(...string) (*State, string) {
	t.Helper()
	head := append([]string{"x board", "x board", "file A1"}, pre...)
	head = append(head, "query C9-F4-D2", "wait", "wait", "wait", "wait", "wait",
		"query C9-F4-D2")
	for i := 0; i < 9; i++ {
		head = append(head, "wait")
	}
	return func(tail ...string) (*State, string) {
		return play(t, append(head, tail...)...)
	}
}

// The précis is the hard floor of Act I: unskippable, undismissable, and four
// paragraphs of true sentences describing a day that did not happen. Even a
// player who filed everything ends the act having been shown it.
func TestThePrecisCannotBePutDown(t *testing.T) {
	s, log := playToPrecis(t)("wait", "look", "x board")
	if !s.PrecisPending {
		t.Fatal("the Crown composes it whether or not anybody asked")
	}
	if !s.Has("ACT1_TURN_SEEN") {
		t.Fatal("expected ACT1_TURN_SEEN")
	}
	if _, held := s.Held["A22"]; !held {
		t.Fatal("the précis is itself the last clue of the act")
	}
	if !strings.Contains(log, "requisitioned one Vault shelf in error") {
		t.Fatal("paragraph four is the archive's account of the player")
	}
	if !strings.Contains(log, "You cannot put it down") {
		t.Fatal("everything else waits")
	}
	if s.Has("ACT1_COMPLETE") {
		t.Fatal("and nothing advances until it is answered")
	}

	// Reading your own coat is free, and is the entire point of the pause.
	_, free := playToPrecis(t, "x quill", "x quill", "note A4")("notebook")
	tail := free[strings.LastIndex(free, "> notebook"):]
	if !strings.Contains(tail, "YOUR NOTEBOOK") {
		t.Fatal("the player may consult what they wrote while it waits")
	}
	if strings.Contains(tail, "You cannot put it down") {
		t.Fatal("and consulting it is not an attempt to put the précis down")
	}
}

// Three answers, and the third is only available to a player who wrote
// something down. Filing everything is the fast way, and it costs you the
// ability to disagree.
func TestThePrecisHasThreeAnswers(t *testing.T) {
	signed, sLog := playToPrecis(t)("sign")
	if !signed.Has("PRECIS_SIGNED") || !signed.Has("ACT1_COMPLETE") {
		t.Fatal("signing closes the act")
	}
	if signed.Ints["RECONCILIATION"] < 3 {
		t.Fatalf("and costs exposure: %d", signed.Ints["RECONCILIATION"])
	}
	if !strings.Contains(sLog, "you will not once be able to work out what it cost") {
		t.Fatal("institutional cooperation, and no way to price it")
	}

	refused, _ := playToPrecis(t)("refuse")
	if !refused.Has("PRECIS_REFUSED") || refused.Ints["VERITY"] < 1 {
		t.Fatal("refusing is worth a point of VERITY")
	}
	if refused.Ints["CRASE_STANDING"] != -1 { // opens at +1
		t.Fatalf("and two of standing: %d", refused.Ints["CRASE_STANDING"])
	}

	// A player who filed everything has nothing in their own hand to add.
	empty, eLog := playToPrecis(t)("annotate")
	if empty.Has("PRECIS_ANNOTATED") {
		t.Fatal("an annotation needs something to annotate from")
	}
	if !strings.Contains(eLog, "it is blank") {
		t.Fatal("and the refusal should show them the empty book")
	}
	if !empty.PrecisPending {
		t.Fatal("and the précis is still there")
	}

	kept, kLog := playToPrecis(t, "x quill", "x quill", "note A4")("annotate")
	if !kept.Has("PRECIS_ANNOTATED") || kept.Ints["VERITY"] < 1 {
		t.Fatal("with a notebook, the annotation is available and is worth a point")
	}
	if !strings.Contains(kLog, "disagrees with itself") {
		t.Fatal("a signed document that disagrees with itself is the only honest object of the day")
	}
}

// Act II opens on the far side of the précis, eight days later, in the same
// room, with four of Vosk's case-files still open on the cogitator.
func TestActTwoOpensOnTheFarSide(t *testing.T) {
	s, log := playToPrecis(t)("sign")
	if s.Room.ID != "standing_board" || !s.Has("ACT2_OPEN") {
		t.Fatal("signing it is the door into Act II")
	}
	if _, held := s.Held["B0"]; !held {
		t.Fatal("the four open files are handed to every player who arrives")
	}
	if !strings.Contains(log, "These have been open eleven weeks") {
		t.Fatal("he was pointing for eleven weeks")
	}
	if !strings.HasPrefix(s.TimeString(), "0997.4.219") {
		t.Fatalf("a new act is a new day: %s", s.TimeString())
	}
	if !strings.Contains(s.TimeString(), "05h4") {
		t.Fatalf("and the clock resets with it: %s", s.TimeString())
	}
}

// Ilaria will not cross the threshold, and she is not being poetic: Vosk had
// the vox-plate taken out four years before he used the room, so it has no
// shadow for her to stand on.
func TestIlariaWillNotComeIn(t *testing.T) {
	s, log := playToPrecis(t)("sign", "ask ilaria about the room")
	if !strings.Contains(log, "It's *absent*") {
		t.Fatal("she is describing an instrument reading")
	}
	if !s.Has("VOSK_STRIPPED_ROOM") {
		t.Fatal("four years and two months before he died")
	}
	if !strings.Contains(log, "at his own cost") {
		t.Fatal("the Crown would not pay to remove a working instrument")
	}
}
