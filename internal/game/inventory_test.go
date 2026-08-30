package game

import (
	"strings"
	"testing"
)

func playSafehouse(t *testing.T, cmds ...string) (*State, string) {
	t.Helper()
	return playToPrecis(t)(append([]string{"sign", "go fenn"}, cmds...)...)
}

// The recitation and the retraction, in one breath, in a room with two metres
// of rockcrete and no plate in it. Nine hundred million four hundred and
// twelve thousand souls, and then no such world.
func TestTheKadmonRecitation(t *testing.T) {
	s, log := playSafehouse(t, "ask fenn about kadmon")
	if !s.Has("KADMON_HEARD") {
		t.Fatal("the name is in play")
	}
	if _, held := s.Held["B1"]; !held {
		t.Fatal("B1 is the recitation and its retraction together")
	}
	if !strings.Contains(log, "nine hundred million four hundred and twelve thousand") {
		t.Fatal("too exact to be an estimate")
	}
	if !strings.Contains(log, "there is no Kadmon Veil") {
		t.Fatal("and too calm to be a lie")
	}
	if strings.Contains(log, "B1 HELD") {
		t.Fatal("nothing marks the moment as evidence — he thinks he made an error")
	}
}

// At the gallery arch it looked like the shielding. In a box with no vox in it
// at all, he does it again, in the same words at the same syllable — so it is
// not where he is standing.
func TestItIsNotTheRoom(t *testing.T) {
	s, log := playSafehouse(t, "ask fenn about kadmon", "ask fenn about again")
	if _, held := s.Held["B2"]; !held {
		t.Fatal("the second pass is the clue")
	}
	if !strings.Contains(log, "It is not the room") {
		t.Fatal("and the second pass is what says so")
	}
	if !strings.Contains(log, "He is not distressed") {
		t.Fatal("he believes he made an error and cannot be shown otherwise")
	}
	if !s.Has("FENN_INTACT") {
		t.Fatal("he is still intact; this is not damage he can feel")
	}
}

// FENN_DRIFT is shown as equipment condition, because that is what the Crown's
// inventory system calls a savant with a mnemonic augmetic. Never relabelled.
func TestDriftIsEquipmentCondition(t *testing.T) {
	before, beforeLog := playSafehouse(t, "status")
	if strings.Contains(beforeLog, "CONCORDANCE OF RECORD") {
		t.Fatal("the readout appears once he is an instrument you are using")
	}
	if before.Ints["FENN_DRIFT"] != 0 {
		t.Fatal("drift opens at nought")
	}

	clean, cleanLog := playSafehouse(t, "ask fenn about kadmon", "status")
	if !strings.Contains(cleanLog, "CONCORDANCE OF RECORD: 6/6") {
		t.Fatal("and it reads as condition, not as a man")
	}
	if strings.Contains(cleanLog, "DRIFT") {
		t.Fatal("the label is never the thing it measures")
	}
	_ = clean

	// Forbidding the word costs a point at once, and he forgets he was told.
	gagged, gagLog := playSafehouse(t, "ask fenn about kadmon", "forbid", "status")
	if !gagged.Has("FENN_SUPPRESSED") || gagged.Ints["FENN_DRIFT"] != 1 {
		t.Fatalf("suppression costs a point immediately: %d", gagged.Ints["FENN_DRIFT"])
	}
	if !strings.Contains(gagLog, "CONCORDANCE OF RECORD: 5/6") {
		t.Fatal("and the readout moves")
	}
	if !strings.Contains(gagLog, "no memory whatever of being told") {
		t.Fatal("he complies, and then forgets he was told")
	}
}

// Isolating him preserves the instrument by making it almost useless, which is
// the whole of Thread A stated once, in an order the player gives.
func TestIsolationIsATrade(t *testing.T) {
	s, log := playSafehouse(t, "kill vox")
	if !s.Has("FENN_ISOLATED") || !s.Has("FENN_INTACT") {
		t.Fatal("dark room, intact savant")
	}
	if !strings.Contains(log, "almost useless") {
		t.Fatal("and the cost is said plainly, once")
	}
}

// Reconciliation does not close a subject, it removes it: the name stops
// matching, the person answers as though they never heard of it, and there is
// no gap for the player to notice.
func TestReconciliationRemovesRatherThanCloses(t *testing.T) {
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	fenn := g.Scene("safehouse").FindNPC("fenn")
	var kadmon *struct{}
	_ = kadmon
	found := false
	for i := range fenn.Topics {
		if fenn.Topics[i].ID == "kadmon" {
			fenn.Topics[i].HiddenAbove = 25
			found = true
		}
	}
	if !found {
		t.Fatal("no kadmon topic to hide")
	}

	s := New(g)
	s.Set("ROSETTE_TAKEN")
	s.Move("safehouse")
	s.Add("RECONCILIATION", 30)
	out, _ := s.Do("ask fenn about kadmon")
	joined := strings.Join(out, "\n")
	if strings.Contains(joined, "Kadmon Veil. Vardan subsector") {
		t.Fatal("past its threshold the topic does not exist")
	}
	if !strings.Contains(joined, "I don't have that") {
		t.Fatal("and what comes back is an ordinary answer with no gap in it")
	}
	if s.Has("KADMON_HEARD") {
		t.Fatal("nothing fires")
	}
}
