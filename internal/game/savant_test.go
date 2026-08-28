package game

import (
	"strings"
	"testing"
)

// playFenn takes the calendar down to tier eight, which is the whole entry
// condition for the scene: you cannot ask a man to corroborate a pattern you
// have not found.
func playFenn(t *testing.T, cmds ...string) (*State, string) {
	t.Helper()
	return play(t, append([]string{"x calendar", "x calendar", "note A5", "go fenn"}, cmds...)...)
}

// The act's most important thirty seconds. Inside the shielding his recall is
// perfect; four metres and one brass strip away it corrects itself
// mid-sentence; walk him back in and it returns, whole, with no memory of ever
// having lost it. The same man, the same conversation, three answers.
func TestTheThreshold(t *testing.T) {
	s, log := playFenn(t,
		"ask about the dates",
		"go arch", "ask about the dates",
		"go back", "ask about the dates",
	)
	if _, held := s.Held["A15"]; !held {
		t.Fatal("inside, he gives you the recurrence")
	}
	if !strings.Contains(log, "going back to —") || !strings.Contains(log, "I don't know why I said that") {
		t.Fatal("outside, it corrects itself mid-sentence")
	}
	if !strings.Contains(log, "He has no idea he has ever said anything else") {
		t.Fatal("and inside again it comes back whole")
	}
	if !s.Has("FENN_TESTED") {
		t.Fatal("expected FENN_TESTED")
	}
	if s.Room.ID != "adjutant" {
		t.Fatal("the test ends where it started")
	}
}

// A16 is an observation, and the game offers no comment on it and no dialogue
// option about it. It arrives with no bracket: a player who was not watching
// walks out with nothing, exactly as with the mnemonic.
func TestTheCorrectionIsNotAnnounced(t *testing.T) {
	s, log := playFenn(t, "ask about the dates", "go arch", "ask about the dates")
	if _, held := s.Held["A16"]; !held {
		t.Fatal("the correction is the clue")
	}
	if strings.Contains(log, "A16") {
		t.Fatal("nothing may flag it — the game does not comment on what just happened")
	}
	if s.Has("A16_HELD") {
		t.Fatal("watching it is not recording it")
	}
}

// Filed, A16 becomes savant fatigue and a recommended rest cycle. The record
// takes it, keeps it, and quietly holds something else — a different failure
// from the boots (no field), the mnemonic (nothing was said) and the bequest
// (the player refuses). A player who only files it has nothing.
func TestFilingTheCorrectionLosesIt(t *testing.T) {
	filed, log := playFenn(t,
		"ask about the dates", "go arch", "ask about the dates", "go back", "file A16")
	if !strings.Contains(log, "SAVANT FATIGUE") || !strings.Contains(log, "RECOMMEND REST CYCLE") {
		t.Fatal("the cogitator classifies it helpfully")
	}
	if strings.Contains(log, "shielding") && strings.Contains(log, "FILED. CLASS: RETINUE") &&
		strings.Contains(log[strings.Index(log, "FILED. CLASS: RETINUE"):], "threshold") {
		t.Fatal("the filed entry must not contain what the player actually observed")
	}
	if len(filed.Notebook) != 1 { // A5 only
		t.Fatal("filing it does not also write it down for you")
	}

	noted, _ := playFenn(t,
		"ask about the dates", "go arch", "ask about the dates", "go back", "note the correction")
	if noted.Held["A16"] != StoreBook || !noted.Has("A16_PHYSICAL") {
		t.Fatal("the coat keeps what he actually said")
	}
	if noted.Ints["VERITY"] <= filed.Ints["VERITY"] {
		t.Fatal("only the written version is worth anything")
	}
}

// You cannot ask a man to corroborate a pattern you have not found. Without
// the calendar he asks which dates you mean, helpfully, and you have nothing.
func TestFennNeedsTheCalendar(t *testing.T) {
	s, log := play(t, "go fenn", "ask about the dates")
	if _, held := s.Held["A15"]; held {
		t.Fatal("a recurrence needs at least three points and you have none")
	}
	if !strings.Contains(log, "Which dates?") {
		t.Fatal("he asks, helpfully, and means it")
	}
	if !s.Has("FENN_INTACT") {
		t.Fatal("he is intact from the moment you walk in, whatever you ask him")
	}
}

// He will follow you anywhere and does not understand that he is evidence, so
// the arch is not somewhere you go without him — and he is one man in two
// rooms, not two men. His greeting happens once.
func TestFennIsOnePerson(t *testing.T) {
	alone, aloneLog := play(t, "go fenn", "go arch")
	if alone.Room.ID == "gallery_arch" {
		t.Fatal("there is no reason to walk to an archway on your own")
	}
	if !strings.Contains(aloneLog, "no reason at all to walk to it") {
		t.Fatal("and the refusal should say so")
	}

	_, log := playFenn(t, "ask about vosk", "go arch", "ask about vosk", "go back", "ask about vosk")
	if strings.Count(log, "Congratulations") != 1 {
		t.Fatalf("he greets you once, not once per room: %d", strings.Count(log, "Congratulations"))
	}
	if strings.Count(log, "but not *filed*") != 1 {
		t.Fatal("and a topic he has answered stays answered across the threshold")
	}
}

// Shielded or fielded, never both, and the third answer is not choosing.
func TestFennIsShieldedOrFielded(t *testing.T) {
	shielded, _ := playFenn(t, "shield him", "go back")
	assertFenn(t, shielded, "FENN_SHIELDED")

	fielded, _ := playFenn(t, "take him", "go back")
	assertFenn(t, fielded, "FENN_FIELDED")

	both, _ := playFenn(t, "shield him", "take him", "go back")
	assertFenn(t, both, "FENN_SHIELDED")

	quiet, _ := playFenn(t, "ask about vosk", "go back")
	assertFenn(t, quiet, "")
	if !quiet.Has("FENN_INTACT") {
		t.Fatal("undecided is not the same as lost")
	}
}

func assertFenn(t *testing.T, s *State, want string) {
	t.Helper()
	var set []string
	for _, f := range []string{"FENN_SHIELDED", "FENN_FIELDED"} {
		if s.Has(f) {
			set = append(set, f)
		}
	}
	if want == "" {
		if len(set) != 0 {
			t.Fatalf("not choosing must set nothing, got %v", set)
		}
		return
	}
	if len(set) != 1 || set[0] != want {
		t.Fatalf("expected exactly %s, got %v", want, set)
	}
}

// An intro carries the beats of arriving somewhere and those happen once.
func TestIntrosHappenOnce(t *testing.T) {
	_, log := playFenn(t, "ask about vosk", "go arch", "go back")
	if strings.Count(log, "Hearing is filing, to him") != 1 {
		t.Fatal("the opening beat of a scene does not replay when you walk back in")
	}
	if !strings.Contains(log, "pleased to see you again") {
		t.Fatal("a revisited room says so in its own words")
	}
}

// A15 arrives once, whichever way the player gets to it — asked for, or shown
// the pad — and the other route knows it has already happened.
func TestTheRecurrenceArrivesOnce(t *testing.T) {
	shown, log := playFenn(t, "show the calendar to fenn", "ask about the dates", "show A5 to fenn")
	if _, held := shown.Held["A15"]; !held {
		t.Fatal("showing him the pad is a way to get there")
	}
	if n := strings.Count(log, "A15 HELD"); n != 1 {
		t.Fatalf("granted once, announced once: %d", n)
	}
	if !strings.Contains(log, "haven't stopped doing it since you asked me") {
		t.Fatal("asked afterwards, he reports what he has been doing about it")
	}
	if !strings.Contains(log, "You've shown me") {
		t.Fatal("shown twice, he says so")
	}

	// And the threshold still works from the shown route.
	tested, _ := playFenn(t, "show the calendar to fenn", "go arch", "ask about the dates")
	if !tested.Has("FENN_TESTED") {
		t.Fatal("the correction does not care how you learned the dates")
	}
}
