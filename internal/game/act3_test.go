package game

import (
	"strings"
	"testing"
)

// departed is a player at the top of Act III: the sentence read, the address
// known, the plate in a felt case.
var departed = []string{"ACT2_OPEN", "ACT2_COMPLETE", "KADMON_ADDRESS",
	"ERRATUM_OBTAINED", "ERRATUM_UNUSED", "CIPHER_COMPLETE"}

func playDock(t *testing.T, extra []string, cmds ...string) (*State, string) {
	t.Helper()
	return playFrom(t, "skydock", append(append([]string{}, departed...), extra...), cmds...)
}

// The booking is not denied. Six is more than five, five is the threshold, and
// there is no hand anywhere on it to push against.
func TestThePassageIsPendedNotRefused(t *testing.T) {
	s, log := playDock(t, nil, "talk to rill", "ask rill about the booking")
	if !s.Has("PASSAGE_PENDED") {
		t.Fatalf("expected the pending:\n%s", log)
	}
	if strings.Contains(log, "denied") && !strings.Contains(log, "isn't denied") {
		t.Fatalf("it is not a refusal:\n%s", log)
	}
	if !strings.Contains(log, "Eleven days") {
		t.Fatalf("the cost is the queue:\n%s", log)
	}
}

// Three hulls, and taking one closes the other two.
func TestTheThreeHullsAreExclusive(t *testing.T) {
	routes := []struct {
		name  string
		cmds  []string
		flag  string
		other []string
	}{
		{"writ", []string{"talk to crase", "ask crase about the writ", "take the writ"},
			"TRANSPORT_WRIT", []string{"TRANSPORT_STOLEN", "TRANSPORT_DEBT"}},
		{"theft", []string{"talk to rill", "ask rill about the booking", "lean on him", "break the seal"},
			"TRANSPORT_STOLEN", []string{"TRANSPORT_WRIT", "TRANSPORT_DEBT"}},
		{"debt", []string{"talk to sazh", "ask sazh about the price", "take the favour"},
			"TRANSPORT_DEBT", []string{"TRANSPORT_WRIT", "TRANSPORT_STOLEN"}},
	}
	for _, r := range routes {
		t.Run(r.name, func(t *testing.T) {
			s, log := playDock(t, nil, r.cmds...)
			if !s.Has(r.flag) {
				t.Fatalf("%s did not book:\n%s", r.name, log)
			}
			for _, o := range r.other {
				if s.Has(o) {
					t.Fatalf("%s also set %s", r.name, o)
				}
			}
			if _, ticks := s.Do("go depart"); ticks == 0 || s.Room.ID != "cerement" {
				t.Fatalf("%s cannot leave; in %s", r.name, s.Room.ID)
			}
		})
	}
}

// The debt is the only contract in the sector nobody can edit, and the game
// refuses to file it for the same reason.
func TestSazhsTallyCannotBeFiled(t *testing.T) {
	s, log := playDock(t, nil, "talk to sazh", "ask sazh about the price", "file C6")
	if _, held := s.Held["C6"]; !held {
		t.Fatalf("expected the tally:\n%s", log)
	}
	if !strings.Contains(log, "There is no such class.") {
		t.Fatalf("no template takes it:\n%s", log)
	}
	if s.Has("C6_FILED") {
		t.Fatal("nothing written down, on either side, ever")
	}
}

// The truest thing anybody says in the game, said without a hedge.
func TestCraseHasNeverLiedToYou(t *testing.T) {
	s, log := playDock(t, nil,
		"talk to crase", "ask crase about the session", "ask crase about the truth")
	if !s.Has("CRASE_NEVER_LIED") {
		t.Fatalf("she answers it:\n%s", log)
	}
	if !strings.Contains(log, "\"I have never lied to you,\" she says.") {
		t.Fatalf("expected the sentence:\n%s", log)
	}
	// She says it and does not defend it. The prose names the four hedges an
	// officer of eleven years would reach for; her line contains none of them.
	said := log[strings.Index(log, "You ask her, straight out"):]
	spoken := said[:strings.Index(said, "There is no hedge in it.")]
	for _, hedge := range []string{"knowingly", "recollection", "as far as", "believe"} {
		if strings.Contains(spoken, hedge) {
			t.Fatalf("she hedged it: %s\n%s", hedge, spoken)
		}
	}
	if s.Ints["VERITY"] != 1 {
		t.Fatalf("bank it: expected 1, got %d", s.Ints["VERITY"])
	}
}

// The leash is standard, and it is a leash, and saying so out loud costs.
func TestReadingTheClauseAloudCosts(t *testing.T) {
	s, _ := playDock(t, nil, "talk to crase", "ask crase about the writ")
	before := s.Ints["CRASE_STANDING"]
	out, _ := s.Do("show C1 to crase")
	log := strings.Join(out, "\n")
	if !strings.Contains(log, "\"Yes,\" she says.") {
		t.Fatalf("she agrees and does not defend it:\n%s", log)
	}
	if s.Ints["CRASE_STANDING"] != before-1 {
		t.Fatalf("expected one layer of warmth: %d -> %d", before, s.Ints["CRASE_STANDING"])
	}
	if s.Ints["VERITY"] != 1 {
		t.Fatalf("expected VERITY +1, got %d", s.Ints["VERITY"])
	}
}

// Out of relay range the cogitator stops lying about the world and starts
// lying about the player, one symptom every thirty-six hours.
func TestTheCerementLiesAboutThePlayer(t *testing.T) {
	s, log := playFrom(t, "cerement", departed,
		"wait", "wait", "wait", "wait", "wait", "wait",
		"wait", "wait", "wait", "wait", "wait", "wait")
	if !s.Has("IN_TRANSIT") || !s.Has("OUT_OF_RELAY") {
		t.Fatal("nine days with no relay")
	}
	if !strings.Contains(log, "one clause shorter than you remember") {
		t.Fatalf("first symptom:\n%s", log)
	}
	if !strings.Contains(log, "you have never in your life mislaid a sentence") {
		t.Fatalf("second symptom:\n%s", log)
	}
	// The third needs the exposure to have got somewhere, and a quiet passage
	// never says it was quiet.
	if strings.Contains(log, "NO SUCH ENTRY") {
		t.Fatalf("the third symptom fired below 60:\n%s", log)
	}

	hot, log := playFrom(t, "cerement", departed,
		"wait", "wait", "wait", "wait", "wait", "wait",
		"wait", "wait", "wait", "wait", "wait", "wait",
		"wait", "wait", "wait", "wait", "wait", "wait")
	hot.Add("RECONCILIATION", 60)
	for i := 0; i < 8; i++ {
		for _, line := range hot.Tick() {
			log += line
		}
	}
	if !strings.Contains(log, "NO SUCH ENTRY") {
		t.Fatalf("at 60 the record forgets a thing you requisitioned yourself:\n%s", log)
	}
}

// Reassurance is not free out here, and the cost is never shown.
func TestReadingYourOwnFileCostsInTransit(t *testing.T) {
	s, log := playFrom(t, "cerement", departed, "notes", "notes", "notes")
	if s.Ints["RECONCILIATION"] != 6 {
		t.Fatalf("two an entry: expected 6, got %d", s.Ints["RECONCILIATION"])
	}
	if strings.Contains(log, "RECONCILIATION") {
		t.Fatalf("nothing says what it cost:\n%s", log)
	}
	// In relay range it is free, as it has been all game.
	cheap, _ := play(t, "notes", "notes")
	if cheap.Ints["RECONCILIATION"] != 0 {
		t.Fatalf("checking your own coat is free: got %d", cheap.Ints["RECONCILIATION"])
	}
}

// The hull is running orders filed six days before he died, and the method is
// what three of them say together.
func TestTheMethodIsAssembledNotAnnounced(t *testing.T) {
	// Stolen: nobody has been through the hull.
	s, log := playFrom(t, "spirit", append(departed, "TRANSPORT_STOLEN"),
		"talk to spirit", "ask spirit about the orders",
		"x nav", "x nav", "note C2",
		"x manifest", "x manifest", "note C3",
		"x drawer", "x drawer", "note C5",
	)
	if !s.Has("SPIRIT_OPEN") {
		t.Fatalf("a stolen hull was never curated:\n%s", log)
	}
	if !s.Has("VOSK_METHOD_UNDERSTOOD") {
		t.Fatal("plot, casket and uncoupled vox is the method")
	}
	if !strings.Contains(log, "a struck one says *I considered it*") {
		t.Fatalf("expected the return leg:\n%s", log)
	}
	if strings.Contains(log, "VOSK_METHOD") {
		t.Fatalf("nothing announces it:\n%s", log)
	}

	// Two of the three is the shape of it and not the thing.
	part, _ := playFrom(t, "spirit", append(departed, "TRANSPORT_STOLEN"),
		"x nav", "x nav", "note C2", "x manifest", "x manifest", "note C3")
	if !part.Has("VOSK_METHOD_PARTIAL") || part.Has("VOSK_METHOD_UNDERSTOOD") {
		t.Fatal("the vox order is the third leg")
	}
}

// Route A had his effects forwarded, curated.
func TestTheWritRouteGetsACuratedHull(t *testing.T) {
	s, log := playFrom(t, "spirit", append(departed, "TRANSPORT_WRIT"),
		"x slate", "x slate", "x drawer")
	if !strings.Contains(log, "DAMAGED IN TRANSIT") {
		t.Fatalf("the slate came back curated:\n%s", log)
	}
	if !strings.Contains(log, "There is not much down here") {
		t.Fatalf("the drawer is not in the room at all:\n%s", log)
	}
	if _, held := s.Held["C4"]; held {
		t.Fatal("nothing to recover off a re-cut slate")
	}
	if _, held := s.Held["C5"]; held {
		t.Fatal("C5 is absent on the writ route")
	}

	// An intact savant gets it back out of the parity blocks, at a cost.
	fixed, log := playFrom(t, "spirit",
		append(departed, "TRANSPORT_WRIT", "FENN_INTACT", "FENN_ALIVE"),
		"recover it", "x slate", "x slate", "note C4")
	if _, held := fixed.Held["C4"]; !held {
		t.Fatalf("Fenn finds the second copy:\n%s", log)
	}
	if fixed.Ints["FENN_DRIFT"] != 1 {
		t.Fatalf("four hours in the parity blocks costs him: %d", fixed.Ints["FENN_DRIFT"])
	}
}

// Buying the spirit-core back is the second favour, and she makes you hear the
// number.
func TestTheDebtRouteBuysTheCoreBack(t *testing.T) {
	s, log := playFrom(t, "spirit", append(departed, "TRANSPORT_DEBT", "SAZH_DEBT"),
		"x drawer")
	if !strings.Contains(log, "There is not much down here") {
		t.Fatalf("no core, no drawer:\n%s", log)
	}
	s.Do("buy the core")
	if !s.Has("SPIRIT_OPEN") || !s.Has("SAZH_DEBT_2") {
		t.Fatal("one favour, which makes two")
	}
	out, _ := s.Do("x drawer")
	if strings.Contains(strings.Join(out, "\n"), "There is not much") {
		t.Fatal("the drawer is there once the core is")
	}
}

// Wiping the hull removes her foreknowledge and the only corroboration at once.
func TestWipingTheSpiritCutsBothWays(t *testing.T) {
	s, log := playFrom(t, "spirit", append(departed, "TRANSPORT_STOLEN"),
		"ask spirit about the orders", "wipe it")
	if !s.Has("SPIRIT_WIPED") {
		t.Fatalf("expected the purge:\n%s", log)
	}
	if s.Ints["VERITY"] != 1 {
		t.Fatalf("expected VERITY +1, got %d", s.Ints["VERITY"])
	}
	if !strings.Contains(log, "Nobody will ever be able to check it either.") {
		t.Fatalf("the cost is in the same breath:\n%s", log)
	}
}

// aboard is a player nine days out with the whole retinue intact.
var aboard = append(append([]string{}, departed...),
	"TRANSPORT_STOLEN", "BEZ_ALIVE", "FENN_ALIVE", "FENN_INTACT", "VAUX_ALIVE", "VAUX_ALLIED")

// Nobody raises their own subject. The player has to go to them, and going
// costs ship-time.
func TestTheWatchesMustBeVisited(t *testing.T) {
	s, log := playFrom(t, "cerement", aboard, "look", "x bez", "x bez", "x fenn", "x fenn")
	for _, f := range []string{"WATCH_BEZ", "WATCH_FENN"} {
		if s.Has(f) {
			t.Fatal("looking at somebody is not standing a watch with them")
		}
	}
	if strings.Contains(log, "bequest") || strings.Contains(log, "four hundred and six") {
		t.Fatalf("nobody raises their own subject unprompted:\n%s", log)
	}

	s.Add("ILARIA_TRUST", 2)
	for _, door := range []string{"gun bay", "cell", "hold", "wardroom"} {
		if _, ticks := s.Do("go " + door); ticks == 0 {
			t.Fatalf("%s should be open and should cost", door)
		}
		s.Do("go out")
	}
}

// She cannot read, and the ask costs her more than the lung rot does.
func TestBezAsksTheOnlyWaySheCan(t *testing.T) {
	s, log := playFrom(t, "gunbay", aboard,
		"talk to bez", "ask bez about the seal", "read it aloud")
	if !s.Has("BEZ_BEQUEST_READ") || !s.Has("WATCH_BEZ") {
		t.Fatalf("expected the reading:\n%s", log)
	}
	if _, held := s.Held["C7"]; !held {
		t.Fatal("expected the bequest")
	}
	if !strings.Contains(log, "the only person in my service I never had to check") {
		t.Fatalf("the fourth line is why he wrote the other three:\n%s", log)
	}
	if s.Ints["VERITY"] != 1 {
		t.Fatalf("expected VERITY +1, got %d", s.Ints["VERITY"])
	}
	// The game never says the word.
	if strings.Contains(log, "illiterate") || strings.Contains(log, "cannot read") {
		t.Fatalf("it is never named:\n%s", log)
	}
}

// Two small kindnesses, both good reasons, and the game lets you notice.
func TestSofteningTheBequestIsWhatTheyDo(t *testing.T) {
	s, log := playFrom(t, "gunbay", aboard,
		"ask bez about the seal", "soften it")
	if !s.Has("BEQUEST_EDITED") {
		t.Fatalf("expected the edit:\n%s", log)
	}
	if s.Ints["RECONCILIATION"] != 5 {
		t.Fatalf("expected 5, got %d", s.Ints["RECONCILIATION"])
	}
	if !strings.Contains(log, "both of the reasons were good ones") {
		t.Fatalf("the point is that it is not malice:\n%s", log)
	}
	if s.Ints["VERITY"] != 0 {
		t.Fatal("an edited reading is not independent proof of anything")
	}
}

// Refusing is legitimate and she never asks again.
func TestRefusingTheBequestIsFinal(t *testing.T) {
	s, _ := playFrom(t, "gunbay", aboard, "ask bez about the seal", "give it back")
	if !s.Has("BEZ_BEQUEST_REFUSED") || s.Has("BEZ_BEQUEST_READ") {
		t.Fatal("expected the refusal")
	}
	out, _ := s.Do("read it aloud")
	if s.Has("BEZ_BEQUEST_READ") {
		t.Fatalf("she does not ask again: %s", strings.Join(out, "\n"))
	}
}

// The honest answer is worse than insanity and she does not break.
func TestIlariaAsksToBeToldSheIsMad(t *testing.T) {
	s, log := playFrom(t, "cell", aboard, "talk to ilaria", "ask ilaria about the tally",
		"ask ilaria about mad", "tell her the truth")
	if !s.Has("ILARIA_KNOWS") || !s.Has("WATCH_ILARIA") {
		t.Fatalf("expected the truth:\n%s", log)
	}
	if !strings.Contains(log, "I don't count the days I lost. I count the days that lost me.") {
		t.Fatalf("expected her line:\n%s", log)
	}
	if !strings.Contains(log, "the only one aboard who's been in the room with it and lived") {
		t.Fatalf("she volunteers:\n%s", log)
	}
	if s.Ints["VERITY"] != 2 {
		t.Fatalf("expected VERITY +2, got %d", s.Ints["VERITY"])
	}
}

// The kind answer is the machine's answer, applied by hand, to one person.
func TestHandlingIlariaIsTheSameInstrument(t *testing.T) {
	s, log := playFrom(t, "cell", aboard,
		"ask ilaria about the tally", "ask ilaria about mad", "you're ill")
	if !s.Has("ILARIA_HANDLED") || s.Has("ILARIA_KNOWS") {
		t.Fatal("expected the kind answer")
	}
	if s.Ints["RECONCILIATION"] != 6 {
		t.Fatalf("expected 6, got %d", s.Ints["RECONCILIATION"])
	}
	if !strings.Contains(log, "using the same instrument, for the same reason") {
		t.Fatalf("the game says it once, plainly:\n%s", log)
	}
}

// He is not prophesying. He is finishing.
func TestFennIsFinishingNotGuessing(t *testing.T) {
	s, log := playFrom(t, "hold", aboard,
		"talk to fenn", "ask fenn about the recitation", "ask fenn about finishing",
		"copy it", "note C9")
	if !s.Has("C9_TAKEN") || !s.Has("WATCH_FENN") {
		t.Fatalf("expected the copy:\n%s", log)
	}
	if !s.Has("C9_PHYS") || s.Has("C9_COG") {
		t.Fatal("hand-copied is C9_PHYS")
	}
	if !strings.Contains(log, "\"I'm *finishing*.\"") {
		t.Fatalf("expected the answer:\n%s", log)
	}
	if !strings.Contains(log, "I assumed I was.") {
		t.Fatalf("the one hesitation in four days:\n%s", log)
	}
}

// Filing it is fast and editable, and the list arrives three worlds shorter.
func TestFilingTheRecitationIsTheTrap(t *testing.T) {
	s, log := playFrom(t, "hold", aboard,
		"ask fenn about the recitation", "copy it", "file C9")
	if !s.Has("C9_COG") || s.Has("C9_PHYS") {
		t.Fatal("filed is C9_COG")
	}
	// Authored prose is hard-wrapped, so match against a flattened copy.
	if !strings.Contains(flat(log), "the list will be three worlds shorter than the one he gave you") {
		t.Fatalf("the filing says what it will do, in the filing:\n%s", log)
	}
	if strings.Contains(log, "WARNING") || strings.Contains(log, "DELETED") {
		t.Fatalf("and nothing warns:\n%s", log)
	}
}

// The shroud is the kind thing and it costs the back half of the recitation.
func TestShieldingFennTruncatesIt(t *testing.T) {
	s, log := playFrom(t, "hold", aboard, "ask fenn about the recitation", "shield him")
	if !s.Has("FENN_SHIELDED") || !s.Has("C9_TRUNCATED") {
		t.Fatalf("expected the shroud:\n%s", log)
	}
	if !strings.Contains(log, "What were we doing?") {
		t.Fatalf("he stops mid-word:\n%s", log)
	}
}

// The apology is precise about what it is for, and the price of a promise is
// the only sealed channel out of the sector.
func TestVauxsApologyBuysTheChannel(t *testing.T) {
	s, log := playFrom(t, "wardroom", aboard,
		"talk to vaux", "ask vaux about the dossier", "ask vaux about the apology", "agree")
	if !s.Has("VAUX_APOLOGY") || !s.Has("CONCLAVE_CHANNEL") || !s.Has("WATCH_VAUX") {
		t.Fatalf("expected the writ:\n%s", log)
	}
	if !strings.Contains(log, "He chose you because you were not a good clerk.") {
		t.Fatalf("the apology is for forty years of being right about small things:\n%s", log)
	}
	if !strings.Contains(log, "Do not spend it on being right.") {
		t.Fatalf("expected the handover:\n%s", log)
	}

	no, _ := playFrom(t, "wardroom", aboard,
		"ask vaux about the dossier", "ask vaux about the apology", "decline")
	if no.Has("CONCLAVE_CHANNEL") {
		t.Fatal("he keeps the writ, and does not mention it")
	}
}

// VOYAGE_FULL is every watch the player could have stood. A member who is
// dead, burned or blank is not a watch they failed to keep.
func TestVoyageFullCountsWhatWasPossible(t *testing.T) {
	// Only Bez aboard: one watch is the whole voyage.
	only := append(append([]string{}, departed...), "TRANSPORT_STOLEN", "BEZ_ALIVE")
	s, _ := playFrom(t, "gunbay", only, "ask bez about the seal", "read it aloud")
	if !s.Has("VOYAGE_FULL") {
		t.Fatal("a lost retinue is not a skipped scene")
	}

	// Everyone aboard: one watch is not enough.
	all, _ := playFrom(t, "gunbay", aboard, "ask bez about the seal", "read it aloud")
	if all.Has("VOYAGE_FULL") {
		t.Fatal("three still to stand")
	}
}

// flat collapses the hard wrapping in authored prose so an assertion can quote
// a whole sentence without knowing where the line breaks fell.
func flat(s string) string { return strings.Join(strings.Fields(s), " ") }

// PRESS HIM was Rill's act before PRESS was a verb in 3.12. A person's act
// wins when there is no layer to open, the same way READ ROSTER does.
func TestPressStillReachesRillsAct(t *testing.T) {
	for _, cmd := range []string{"press him", "push him", "lean on him"} {
		s, log := playDock(t, nil, "talk to rill", "ask rill about the booking", cmd)
		if !s.Has("RILL_DELAYED") {
			t.Errorf("%q did not reach the act:\n%s", cmd, log)
		}
	}
}
