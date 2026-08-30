package game

import (
	"strings"
	"testing"
)

// playMedicae waits the body out and follows it down to tier one.
func playMedicae(t *testing.T, cmds ...string) (*State, string) {
	t.Helper()
	pre := []string{}
	for i := 0; i < 12; i++ {
		pre = append(pre, "wait")
	}
	pre = append(pre, "go down")
	return play(t, append(pre, cmds...)...)
}

// A11 and A12 are the act's guaranteed floor: the pistol fired eleven minutes
// post-mortem, and there is no mechanism. Every player gets them, unasked, for
// standing on the deck — Act II builds the Erratum thread on this regardless
// of how badly Act I went.
func TestAutopsyFloorIsUnmissable(t *testing.T) {
	s, log := playMedicae(t)
	if !s.Has("AUTOPSY_DONE") {
		t.Fatal("being in the room is the whole requirement")
	}
	for _, id := range []string{"A11", "A12"} {
		if _, held := s.Held[id]; !held {
			t.Fatalf("%s is delivered unasked", id)
		}
	}
	if !strings.Contains(log, "He died and then he shot himself. In that order.") {
		t.Fatal("Ostro's summary is the same for everybody")
	}
	if !strings.Contains(log, "there is no other box") {
		t.Fatal("and so is his reason for the box he struck")
	}
}

// You cannot autopsy a man who is still on the floor upstairs. Tier one opens
// when the medicae take him, which is what the eight minutes have been for.
func TestTierOneOpensWhenTheyTakeHim(t *testing.T) {
	early, earlyLog := play(t, "go down")
	if early.Room.ID != "vosk_chambers" {
		t.Fatal("there is nothing on tier one yet")
	}
	if !strings.Contains(earlyLog, "still on the floor of this room") {
		t.Fatal("and the refusal should say why")
	}
	late, _ := playMedicae(t)
	if late.Room.ID != "medicae" {
		t.Fatal("once he is down there, so can you be")
	}
}

// A14 is the one clue the cogitator would take. There is a field for it, it is
// the right field, and it would end her career — so the refusal is the
// player's, not the machine's. That is what makes it different from the boots
// and the mnemonic, which the form simply cannot hold.
func TestTheBequestIsNeverFiled(t *testing.T) {
	s, log := playMedicae(t,
		"ask bez about the door", "ask bez about her theory",
		"read roster",
		"ask bez about habits", "ask bez about food",
		"file A14", "note the envelope",
	)
	if _, held := s.Held["A14"]; !held {
		t.Fatal("the cough and the sealed envelope arrive together, unasked")
	}
	if s.Held["A14"] == StoreCog || s.Held["A14"] == StoreBoth {
		t.Fatal("A14 must never reach the cogitator")
	}
	if !strings.Contains(log, "exactly the right field") {
		t.Fatal("the machine would take it, and the refusal has to admit that")
	}
	if !strings.Contains(log, "declined to file for a reason that is not tactical") {
		t.Fatal("the game names it as the betrayal it would be")
	}
	if !s.Has("A14_PHYSICAL") {
		t.Fatal("it belongs in the coat")
	}
}

// The roster is handed to you with a reason that is not true, in the best
// light in the Adamant Crown. Reading it costs nothing and buys everything;
// naming it is taken perfectly and never forgiven.
func TestTheRosterIsTwoAnswers(t *testing.T) {
	kind, kindLog := playMedicae(t,
		"ask bez about the door", "ask bez about her theory", "read roster")
	if !kind.Has("BEZ_ILLITERACY_KNOWN") || kind.Ints["BEZ_TRUST"] != 2 {
		t.Fatalf("reading it aloud is worth two: %d", kind.Ints["BEZ_TRUST"])
	}
	if strings.Contains(kindLog, "cannot read") {
		t.Fatal("and nobody says it out loud, including the game")
	}

	blunt, _ := playMedicae(t,
		"ask bez about the door", "ask bez about her theory", "name it")
	if !blunt.Has("BEZ_ILLITERACY_KNOWN") || blunt.Ints["BEZ_TRUST"] != 0 {
		t.Fatalf("naming it costs two, floored at nought: %d", blunt.Ints["BEZ_TRUST"])
	}

	// One or the other, never both, and neither before she offers the sheet.
	early, earlyLog := playMedicae(t, "read roster")
	if early.Has("BEZ_ILLITERACY_KNOWN") {
		t.Fatal("she has not handed you anything yet")
	}
	if !strings.Contains(earlyLog, "has not asked you to read her anything") {
		t.Fatal("and the refusal should say so")
	}
	both, _ := playMedicae(t,
		"ask bez about the door", "ask bez about her theory", "read roster", "name it")
	if both.Ints["BEZ_TRUST"] != 2 {
		t.Fatalf("the roster is answered once: %d", both.Ints["BEZ_TRUST"])
	}
}

// Her health is the one subject no pressure opens. Asked directly she lies,
// badly, and the game shows the player the decision rather than the tell.
func TestSheLiesAboutTheCough(t *testing.T) {
	s, log := playMedicae(t, "ask bez about your cough")
	if s.Has("BEZ_KNOWS_DIAGNOSIS") {
		t.Fatal("asking is not ordering, and asking gets you nothing")
	}
	if !strings.Contains(log, "Dust, sir") {
		t.Fatal("she lies")
	}
	if !strings.Contains(log, "watched her decide to do it") {
		t.Fatal("and the player is shown the decision, not the diagnosis")
	}
}

// Ordering the examination is available, correct, and one of the cruellest
// things the player can do to somebody they like. Fourteen months, confirmed,
// and she thanks the medicae.
func TestOrderingTheExamCosts(t *testing.T) {
	s, log := playMedicae(t,
		"ask bez about the door", "ask bez about her theory",
		"ask bez about habits", "ask bez about food",
		"order exam",
	)
	if !s.Has("BEZ_KNOWS_DIAGNOSIS") {
		t.Fatal("an Inquisitor can order it, and that is the whole mechanism")
	}
	if !strings.Contains(log, "Fourteen months") {
		t.Fatal("Ostro tells her, not you")
	}
	if !strings.Contains(log, "Doesn't change the watch") {
		t.Fatal("and she goes back to the wall")
	}

	// And not before she has coughed in front of you.
	early, _ := playMedicae(t, "order exam")
	if early.Has("BEZ_KNOWS_DIAGNOSIS") {
		t.Fatal("there is nothing to order an examination of yet")
	}
}

// Bez saw the boots eleven years ago on a ship and has thought he was tired
// about it ever since. She confirms A3 without any idea what she is
// confirming, which is the most useful thing anybody says all day.
func TestBezConfirmsTheBootsWithoutKnowing(t *testing.T) {
	s, log := play(t, append(bootsThenDown(), "show the boots to bez")...)
	if !strings.Contains(log, "He did that once") {
		t.Fatal("she has seen it before")
	}
	if !strings.Contains(log, "Sarrow's Reach") {
		t.Fatal("on the ship, the year you came on")
	}
	if s.Ints["BEZ_TRUST"] != 1 {
		t.Fatalf("BEZ_TRUST +1 for telling her the truth: %d", s.Ints["BEZ_TRUST"])
	}
	if s.Has("CIPHER_HALF_ROOM") {
		t.Fatal("she confirms the fact and decodes nothing")
	}
}

func bootsThenDown() []string {
	cmds := []string{"x boots", "x boots", "note A3"}
	for i := 0; i < 11; i++ {
		cmds = append(cmds, "wait")
	}
	return append(cmds, "go down")
}

// Naming somebody who is not answering must not quietly hand the question to
// whoever else is standing there. Bez is against a doorframe in 1.1 and is not
// taking questions; asking her something used to get you Crase's answer.
func TestAskingTheWrongPersonSaysSo(t *testing.T) {
	_, log := playAct1(t, "ask bez about the door")
	if strings.Contains(log, "Crase says") {
		t.Fatal("a question put to Bez must not be answered by Crase")
	}
	if !strings.Contains(log, "not taking questions") {
		t.Fatal("she is in the room, so the refusal should say that")
	}

	_, absent := play(t, "go vaux", "ask ostro about the heart")
	if !strings.Contains(absent, "nobody here by that name") {
		t.Fatal("Ostro is three tiers down and the game should say so")
	}
	// Two people in a room, and each answers only for themselves.
	s, both := playMedicae(t, "ask ostro about the heart", "ask bez about the door")
	if !strings.Contains(both, "Two hundred and eighty-seven, and I have opened up men") {
		t.Fatal("Ostro answers on the heart")
	}
	if !strings.Contains(both, "Eleven hours on that door") {
		t.Fatal("Bez answers on the door")
	}
	if !s.Has("BEZ_MET") || !s.Has("OSTRO_MET") {
		t.Fatal("both are met, separately")
	}
}
