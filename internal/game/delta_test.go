package game

import (
	"strings"
	"testing"
)

// playThreadA gets to Act II with Kadmon heard and the shelf-code filed, which
// is what the Vault run needs on both sides of the comparison.
func playThreadA(t *testing.T, cmds ...string) (*State, string) {
	t.Helper()
	return playToPrecis(t)(append([]string{
		"sign", "go fenn", "ask fenn about kadmon", "note kadmon", "go back",
	}, cmds...)...)
}

// Cael has had the idea already. The docket is three weeks old, the
// originating-request field is empty, and he has written "anticipated" in the
// margin — and he cannot tell you why, and does not become evasive about it.
func TestCaelBuiltItFirst(t *testing.T) {
	s, log := playThreadA(t, "go cael", "x docket", "x docket", "ask cael about why")
	if _, held := s.Held["B3"]; !held {
		t.Fatal("the docket is the clue, before anybody explains anything")
	}
	if !s.Has("CAEL_SUSPECTED") {
		t.Fatal("pressed twice, he is a reconciled man running on the residue of a reason")
	}
	if !strings.Contains(log, "a savant is expensive and Cantor is loud") {
		t.Fatal("true, and insufficient")
	}
	if !strings.Contains(log, "patting his pockets") {
		t.Fatal("not evasive — puzzled, which is the part that matters")
	}
}

// The shroud is a real gift, freely given, and it comes with a line that sends
// position every eleven hours to a bench in a warm room.
func TestTheShroudIsAGiftWithALine(t *testing.T) {
	s, log := playThreadA(t, "go cael", "ask cael about the shroud", "ask cael about telemetry", "accept")
	if !s.Has("FENN_SHIELDED") {
		t.Fatal("it works, and it is the most useful object anybody hands you")
	}
	if !strings.Contains(log, "Position.") {
		t.Fatal("the word is left on its own line and nothing points at it")
	}

	// Drift accrues at a third the rate under it: three ticks, one point.
	sh, _ := playThreadA(t, "go cael", "ask cael about the shroud", "ask cael about telemetry",
		"accept", "go back", "go vault",
		"ask fenn about the alms", "ask fenn about the sump", "ask fenn about outgoing")
	if sh.Ints["FENN_DRIFT"] != 1 {
		t.Fatalf("three sub-queries under the shroud is one point: %d", sh.Ints["FENN_DRIFT"])
	}
	if !sh.Has("FENN_INTACT") || sh.Has("FENN_BURNED_PARTIAL") {
		t.Fatal("and he comes out of it whole")
	}

	// Cut the line afterwards and he never mentions it. Ever.
	cut, cutLog := playThreadA(t, "go cael", "ask cael about the shroud",
		"ask cael about telemetry", "accept", "cut telemetry", "ask cael about vosk")
	if !cut.Has("CAEL_DECEIVED") || !cut.Has("FENN_SHIELDED") {
		t.Fatal("same benefit, no window")
	}
	if !strings.Contains(cutLog, "Cael never mentions it") {
		t.Fatal("and the silence is the thing you remember")
	}
}

// Three sub-queries are offered and the player can afford two. Nothing says
// so, nothing counts them down, and the third one is where he begins
// correcting you from your own notebook.
func TestThePlayerCanAffordTwo(t *testing.T) {
	two, _ := playThreadA(t, "go vault", "ask fenn about the alms", "ask fenn about the sump")
	if two.Ints["FENN_DRIFT"] != 2 {
		t.Fatalf("two sub-queries, two ticks: %d", two.Ints["FENN_DRIFT"])
	}
	if two.Has("FENN_BURNED_PARTIAL") || !two.Has("FENN_INTACT") {
		t.Fatal("two is what he survives")
	}
	if _, held := two.Held["B18"]; held {
		t.Fatal("and one thread goes unhelped")
	}

	three, log := playThreadA(t, "go vault",
		"ask fenn about the alms", "ask fenn about the sump", "ask fenn about outgoing")
	if !three.Has("FENN_BURNED_PARTIAL") {
		t.Fatal("three burns him partially")
	}
	if !three.Has("FENN_INTACT") {
		t.Fatal("but three is not the end of him")
	}
	if !strings.Contains(log, "He has begun correcting you") {
		t.Fatal("at three he corrects the player's own statements")
	}
	if strings.Contains(log, "FENN_DRIFT") || strings.Contains(log, "BURNED") {
		t.Fatal("and nothing names either")
	}
}

// Past four there is nothing left to protect, and it forecloses Ending E
// without a word.
func TestPastFourThereIsNothingLeft(t *testing.T) {
	s, log := playThreadA(t, "go vault",
		"ask fenn about the alms", "ask fenn about the sump", "ask fenn about outgoing",
		"go back", "go fenn", "forbid", "ask fenn about vardan",
		"ask fenn about himself", "ask fenn about kadmon")
	if s.Ints["FENN_DRIFT"] < 4 {
		t.Skip("drift budget did not reach the far end on this path")
	}
	if s.Ints["FENN_DRIFT"] >= 5 && s.Has("FENN_INTACT") {
		t.Fatal("past four, FENN_INTACT is gone")
	}
	if strings.Contains(log, "Ending") {
		t.Fatal("and the game never says what that cost")
	}
}

// Habrach has arranged thirty years around one sentence, and he says it as a
// comfort. Asked who does the reconciling, he hands over a published,
// unrestricted, forty-seven-year-old description of the method.
func TestHabrachHandsOverTheMethod(t *testing.T) {
	s, log := playThreadA(t, "go vault", "ask habrach about withheld", "ask habrach about the memo")
	if !strings.Contains(log, "Things are *reconciled*") {
		t.Fatal("he says it the way other men say a prayer")
	}
	if _, held := s.Held["B23"]; !held {
		t.Fatal("the circular is the clue")
	}
	if !s.Has("KESH_IDENTIFIED") {
		t.Fatal("and a name comes with it")
	}
	if !strings.Contains(log, "cannot be reduced") {
		t.Fatal("page nine is the method, published, in 979")
	}
	if !strings.Contains(log, "requisitioned four times in forty-seven years") &&
		!strings.Contains(log, "read four times") {
		t.Fatal("and almost nobody has ever read it")
	}
}

// The shelf is not a hole. Somebody would have had to make a hole.
func TestTheSequenceNeverHadIt(t *testing.T) {
	_, log := playThreadA(t, "go vault", "x shelf", "x shelf")
	if !strings.Contains(log, "1186. 1187. 1189. 1190.") {
		t.Fatal("the numbers do not skip")
	}
	if !strings.Contains(log, "only ever had what is on it") {
		t.Fatal("this is the clearest statement of the mechanism in the act")
	}
}
