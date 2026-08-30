package game

import (
	"strings"
	"testing"
)

// The verdict is timed two hours before the body was found, and the
// explanation is true. This is the only place in the game where the record
// looks wrong and is not: the player's first suspicion of the archive has to
// be mistaken, so that they learn to distrust their paranoia before they learn
// to trust it.
func TestFirstSuspicionIsWrong(t *testing.T) {
	s, log := playAct1(t, "x writ", "x writ", "ask about the form")
	if _, held := s.Held["A0"]; !held {
		t.Fatal("reading the writ carefully is what surfaces the two hours")
	}
	if !s.Has("TIMESTAMP_NOTICED") {
		t.Fatal("challenging it is the branch")
	}
	if s.Ints["CRASE_STANDING"] != 2 { // opens at +1
		t.Fatalf("she is pleased you asked: %d", s.Ints["CRASE_STANDING"])
	}
	if !strings.Contains(log, "pre-generated form") {
		t.Fatal("she explains it truthfully and the explanation exonerates the record")
	}

	// And you cannot challenge a document you did not read.
	blind, blindLog := playAct1(t, "ask about the form")
	if blind.Has("TIMESTAMP_NOTICED") {
		t.Fatal("noticing is a thing the player does, not a thing the topic does")
	}
	if !strings.Contains(blindLog, "Look at it.") {
		t.Fatal("she should decline to do the noticing for them")
	}
}

// A0 is the game's first silent notebook test: it lands in whichever store the
// player picks, and the game expresses no preference between them.
func TestA0LandsWhereverYouPutIt(t *testing.T) {
	filed, _ := playAct1(t, "x writ", "x writ", "file A0")
	if filed.Held["A0"] != StoreCog || !filed.Has("A0_FILED") {
		t.Fatal("A0 may be filed")
	}
	noted, _ := playAct1(t, "x writ", "x writ", "note timestamp")
	if noted.Held["A0"] != StoreBook || !noted.Has("A0_PHYSICAL") {
		t.Fatal("A0 may be noted, by name and not only by ID")
	}
	if len(noted.Notebook) != 1 {
		t.Fatal("and it goes in the coat")
	}
}

// She does not ask what you think. She asks whether he seemed tired, and
// writes the answer down in ink either way — including the answer that you
// gave her nothing.
func TestTheQuestionIsRecordedEitherWay(t *testing.T) {
	said, saidLog := playAct1(t, "ask about sealed", "ask about the log", "ask about who signed", "tired")
	if !said.Has("SAID_TIRED") {
		t.Fatal("the easiest true thing available")
	}
	if !strings.Contains(saidLog, "characterised him as tired") {
		t.Fatal("it enters the file as corroboration and the game should say so plainly")
	}

	quiet, quietLog := playAct1(t, "ask about sealed", "ask about the log", "ask about who signed", "nothing")
	if quiet.Has("SAID_TIRED") {
		t.Fatal("silence is an answer and it is not that one")
	}
	if !strings.Contains(quietLog, "Declined to characterise") {
		t.Fatal("saying nothing is still written down")
	}

	// Neither answer is available before she has asked.
	early, earlyLog := playAct1(t, "tired")
	if early.Has("SAID_TIRED") {
		t.Fatal("you cannot answer a question nobody has put to you")
	}
	if !strings.Contains(earlyLog, "has not asked you anything") {
		t.Fatal("and the refusal should say why")
	}
}

// "What opens her is grief — she liked him. Say so and the Choir lead arrives
// a scene early." And having arrived early, it must not arrive again.
func TestGriefBuysTheLeadASceneEarly(t *testing.T) {
	s, log := playAct1(t, "ask about grief")
	if !s.Has("CRASE_LEAD_CHOIR") {
		t.Fatal("kindness at the fourth hour is worth the most useful sentence in Act I")
	}
	if _, held := s.Held["A1-lead"]; !held {
		t.Fatal("and the clue with it")
	}

	// Play on into 1.2 and take the rosette: she must not tell you twice.
	for _, c := range []string{"ask about sealed", "ask about the log", "ask about who signed",
		"nothing", "go on", "take rosette"} {
		out, ticks := s.Do(c)
		for _, o := range out {
			log += o + "\n"
		}
		for i := 0; i < ticks; i++ {
			s.Tick()
		}
	}
	if strings.Count(log, "six hundred years") != 1 {
		t.Fatalf("the lead is delivered exactly once, got %d", strings.Count(log, "six hundred years"))
	}
	if !s.Has("ROSETTE_TAKEN") {
		t.Fatal("the ceremony still happens")
	}
}

// Bez puts the sentence into the room twice, unprompted, because there is
// going to be a record and she is not the one writing it. She is not met:
// meeting her belongs to 1.6.
func TestBezSaysItTwice(t *testing.T) {
	s, log := playAct1(t, "ask about sealed", "ask about the log")
	if strings.Count(log, "Door was sealed") != 2 {
		t.Fatalf("twice, in the same voice: got %d", strings.Count(log, "Door was sealed"))
	}
	if s.Has("BEZ_MET") {
		t.Fatal("1.6 owns BEZ_MET; standing near her is not meeting her")
	}
}

// Act I opens in order, at the right hours, and hands 1.3 the two flags the
// ledger says are always true by then.
func TestActOneOpensInOrder(t *testing.T) {
	s, _ := playAct1(t)
	if s.Room.ID != "antechamber_arrival" || !s.Has("ACT1_OPEN") {
		t.Fatal("the game starts in the doorway")
	}
	if !strings.Contains(s.TimeString(), "/ 04h") {
		t.Fatalf("1.1 runs at the fourth hour: %s", s.TimeString())
	}

	step := func(cmds ...string) {
		for _, c := range cmds {
			_, ticks := s.Do(c)
			for i := 0; i < ticks; i++ {
				s.Tick()
			}
		}
	}
	step("ask about sealed", "ask about the log", "ask about who signed", "nothing", "go on")
	if s.Room.ID != "antechamber" {
		t.Fatal("1.1 leads to 1.2")
	}
	if !strings.Contains(s.TimeString(), "04h42") {
		t.Fatalf("1.2 opens at forty-two past the fourth: %s", s.TimeString())
	}

	step("take rosette", "go in")
	if s.Room.ID != "vosk_chambers" {
		t.Fatal("1.2 leads to 1.3")
	}
	if !strings.Contains(s.TimeString(), "06h12") {
		t.Fatalf("1.3 opens at twelve past the sixth: %s", s.TimeString())
	}
	for _, f := range []string{"ACT1_OPEN", "ROSETTE_TAKEN", "CRASE_LEAD_CHOIR"} {
		if !s.Has(f) {
			t.Fatalf("%s is true for every player who reaches 1.3", f)
		}
	}
	if !s.BodyPresent {
		t.Fatal("two scenes and an hour and a half later, the eight minutes have not started")
	}
}
