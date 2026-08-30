package game

import (
	"strings"
	"testing"
)

// docked is a player who has made the routing address.
var docked = append(append([]string{}, landed...), "EXPECTED")

// No challenge, no search, and the Erratum walks aboard under your coat.
func TestNobodyLooksUp(t *testing.T) {
	s, log := playFrom(t, "cantica", docked, "x placard", "x placard", "x floor", "x floor")
	if _, held := s.Held["C15"]; !held {
		t.Fatalf("expected the placard:\n%s", log)
	}
	if !strings.Contains(flat(log), "It is spelled correctly.") {
		t.Fatalf("the way nobody spells it:\n%s", log)
	}
	if !strings.Contains(log, "REYL, AMMON") {
		t.Fatalf("the record renders the player's own name:\n%s", log)
	}
	if !strings.Contains(flat(log), "It is not new.") {
		t.Fatalf("the polish is worn round the letters:\n%s", log)
	}
	if !strings.Contains(flat(log), "the most frightening room you have ever stood in") {
		t.Fatalf("there is no cabal on that floor:\n%s", log)
	}
}

// She answers everything, threatens nobody, and gives the arithmetic actually
// used rather than the arithmetic people assume.
func TestKeshIndexesRatherThanDefends(t *testing.T) {
	s, log := playFrom(t, "cantica", docked,
		"talk to kesh", "accept the tea", "ask kesh about the numbers", "ask kesh about kadmon")
	if !s.Has("KESH_IDENTIFIED") || s.Ints["VERITY"] != 1 {
		t.Fatalf("the tea is real and taking it is not a mistake:\n%s", log)
	}
	if _, held := s.Held["C16"]; !held {
		t.Fatal("expected the ledger")
	}
	if !strings.Contains(flat(log), "Because a screen is a rumour") {
		t.Fatalf("she keeps it on paper:\n%s", log)
	}
	if !strings.Contains(flat(log), "We were wrong by a factor of about six") {
		t.Fatalf("the admission is the point:\n%s", log)
	}
	if !strings.Contains(flat(log), "I have not asked to have it taken out of me") {
		t.Fatalf("she has not had it reconciled:\n%s", log)
	}
}

// Cael files the only evidence of his own reconciliation into the one store
// the station can reach, helpfully, because he is Cael.
func TestC17DeletesItself(t *testing.T) {
	s, log := playFrom(t, "cantica", docked,
		"talk to cael", "ask cael about the windows", "ask cael about the protocols")
	if _, held := s.Held["C17"]; !held {
		t.Fatalf("expected the protocols:\n%s", log)
	}
	if !s.Has("C17_FILED") {
		t.Fatal("he files it without being asked")
	}
	for i := 0; i < 8; i++ {
		s.Tick()
	}
	if _, held := s.Held["C17"]; held {
		t.Fatal("it is gone within the hour")
	}
	if s.Has("C17_FILED") {
		t.Fatal("and it is not in the case-file either")
	}
	// No klaxon, no notice, and the same three words a genuine miss returns.
	notes := s.caseNotesDump()
	if strings.Contains(notes, "C17") {
		t.Fatalf("no gap is left:\n%s", notes)
	}
	out, _ := s.Do("query C17")
	if !strings.Contains(strings.Join(out, "\n"), "NO SUCH ENTRY") {
		t.Fatalf("expected the ordinary miss: %s", strings.Join(out, "\n"))
	}
}

// Four hours of copying shapes, while a Magos watches a person write.
func TestHandCopyingC17IsTheOnlyWayToKeepIt(t *testing.T) {
	s, log := playFrom(t, "cantica", docked,
		"ask cael about the windows", "ask cael about the protocols", "copy it")
	if !s.Has("C17_PHYS") {
		t.Fatalf("expected the copy:\n%s", log)
	}
	for i := 0; i < 8; i++ {
		s.Tick()
	}
	if _, held := s.Held["C17"]; !held {
		t.Fatal("the hand copy survives")
	}
	if !strings.Contains(flat(log), "the only copy that will still exist tomorrow") {
		t.Fatalf("he has worked it out and will not write it down:\n%s", log)
	}
}

// She has never lied. That is the thing that lands, and it is not that she is
// theirs.
func TestCraseHasBeenRunningTheRuleAllAlong(t *testing.T) {
	s, log := playFrom(t, "fifth_chair", docked,
		"talk to crase", "ask crase about muria", "audit her", "ask crase about my elevation")
	if !s.Has("CRASE_AUDITED") {
		t.Fatalf("expected the audit:\n%s", log)
	}
	if s.Ints["VERITY"] != 3 {
		t.Fatalf("two for the audit and one for the elevation, got %d", s.Ints["VERITY"])
	}
	if !strings.Contains(flat(log), "omission, never falsehood") {
		t.Fatalf("she talks the way it talks:\n%s", log)
	}
	if !strings.Contains(flat(log), "you have been doing it yourself since Kadmon") {
		t.Fatalf("and she taught you:\n%s", log)
	}
	// The only thing the player has that she does not.
	if !strings.Contains(flat(log), "She does not know she was used.") {
		t.Fatalf("expected the one win:\n%s", log)
	}
}

func TestMuriaIsWhatSheDefendsAgainst(t *testing.T) {
	_, log := playFrom(t, "fifth_chair", docked, "ask crase about muria")
	if !strings.Contains(flat(log), "cleansed four worlds on suspicion") {
		t.Fatalf("expected Muria:\n%s", log)
	}
	if !strings.Contains(flat(log), "I have carried it. You have read about it.") {
		t.Fatalf("expected the close:\n%s", log)
	}
	// And she does not defend the Concordance on principle.
	if strings.Contains(flat(log), "necessary evil") {
		t.Fatalf("she never argues it that way:\n%s", log)
	}
}

// Refusing her costs the account you will never hear.
func TestRefusingCraseCostsTheAccount(t *testing.T) {
	s, log := playFrom(t, "fifth_chair", docked, "refuse")
	if !s.Has("CRASE_REFUSED") || s.Has("CRASE_AUDITED") {
		t.Fatalf("expected the refusal:\n%s", log)
	}
	if s.Ints["CRASE_STANDING"] != -2 {
		t.Fatalf("standing opens at +1: expected -2, got %d", s.Ints["CRASE_STANDING"])
	}
	out, _ := s.Do("audit her")
	if s.Has("CRASE_AUDITED") {
		t.Fatalf("it is not on offer twice: %s", strings.Join(out, "\n"))
	}
}

// The worst moment in the game, in order.
func TestTheFiledWords(t *testing.T) {
	s, log := playFrom(t, "correctorum", docked,
		"talk to the vox", "ask vox about minutes", "ask vox about kadmon", "ask vox about my file")
	if !s.Has("CONFESSION_HEARD") || !s.Has("FILE_HEARD") {
		t.Fatalf("expected both readings:\n%s", log)
	}
	if !strings.Contains(flat(log), "I'd take it as read") {
		t.Fatalf("he used that on you weekly for eleven years:\n%s", log)
	}
	if !strings.Contains(flat(log), "I price nine hundred million at less, and I ask the table to price it with me") {
		t.Fatalf("he argues for it, and well:\n%s", log)
	}
	if !strings.Contains(flat(log), "He worked out that they would be fine for a while.") {
		t.Fatalf("the grain marginalia is the knife:\n%s", log)
	}
	if !strings.Contains(flat(log), "*Acceptable for expenditure.*") {
		t.Fatalf("expected the last line:\n%s", log)
	}
	if !strings.Contains(flat(log), "You were vetted by the thing you came to kill.") {
		t.Fatalf("expected the knife:\n%s", log)
	}
}

// The omission that is a gift: it never said it composed the phrase, and the
// player has to ask the precise question.
func TestTheBlindspotMustBeAskedFor(t *testing.T) {
	s, log := playFrom(t, "correctorum", docked,
		"ask vox about minutes", "ask vox about kadmon", "ask vox about my file")
	if s.Has("AMANUENSIS_BLINDSPOT") {
		t.Fatalf("hearing the file is not asking the question:\n%s", log)
	}
	// And nothing in what it said claimed authorship.
	if strings.Contains(flat(log), "I composed") || strings.Contains(flat(log), "we composed") {
		t.Fatalf("it never said that:\n%s", log)
	}

	out, _ := s.Do("ask vox about did you compose it")
	log = strings.Join(out, "\n")
	if !s.Has("AMANUENSIS_BLINDSPOT") {
		t.Fatalf("expected the answer:\n%s", log)
	}
	if !strings.Contains(flat(log), "It has an intake and no reader.") {
		t.Fatalf("expected the exploit:\n%s", log)
	}
	if _, held := s.Held["C20"]; !held {
		t.Fatal("expected C20")
	}
}

// She asked to be an instrument rather than a suspect, and here is the invoice.
func TestIlariaStopsAtExpenditure(t *testing.T) {
	flags := append(append([]string{}, docked...), "ILARIA_KNOWS")
	s, log := playFrom(t, "correctorum", flags,
		"ask vox about minutes", "ask vox about kadmon", "ask vox about my file")
	if !s.Has("ILARIA_STOPPED") {
		t.Fatalf("she stops mid-line:\n%s", log)
	}
	if !strings.Contains(flat(log), "I'd like to be asked before I read the next part") {
		t.Fatalf("expected her line:\n%s", log)
	}

	ordered, _ := playFrom(t, "correctorum", flags,
		"ask vox about minutes", "ask vox about kadmon", "ask vox about my file")
	ordered.Add("ILARIA_TRUST", 4)
	ordered.Do("order her")
	if ordered.Ints["ILARIA_TRUST"] != 2 {
		t.Fatalf("expected 4 less 2, got %d", ordered.Ints["ILARIA_TRUST"])
	}

	asked, alog := playFrom(t, "correctorum", flags,
		"ask vox about minutes", "ask vox about kadmon", "ask vox about my file", "please")
	if asked.Ints["ILARIA_TRUST"] != 1 {
		t.Fatalf("expected +1, got %d", asked.Ints["ILARIA_TRUST"])
	}
	if !strings.Contains(flat(alog), "Not this bit. You read this bit.") {
		t.Fatalf("she refuses and hands it over:\n%s", alog)
	}
}

// Stopping the reading is relief that lasts nine seconds and costs a thing you
// will be charged for later without being told what for.
func TestStoppingTheReadingHidesThePrice(t *testing.T) {
	s, log := playFrom(t, "correctorum", docked,
		"ask vox about minutes", "ask vox about kadmon", "stop it")
	if !s.Has("READING_STOPPED") {
		t.Fatalf("expected the stop:\n%s", log)
	}
	if s.Ints["RECONCILIATION"] != 4 {
		t.Fatalf("expected 4, got %d", s.Ints["RECONCILIATION"])
	}
	out, _ := s.Do("ask vox about my file")
	if s.Has("FILE_HEARD") {
		t.Fatalf("C19 is never heard: %s", strings.Join(out, "\n"))
	}
	if s.Has("AMANUENSIS_BLINDSPOT") {
		t.Fatal("and the question cannot be reached")
	}
}

// The Conclave channel has three sources in the act document and the build had
// one, which closed Ending C to two of the three routes with no fiction saying
// why. Spending the writ at Kadmon takes it back — unless something else is
// still holding it up.
func TestTheChannelHasThreeSources(t *testing.T) {
	for _, src := range []string{"VAUX_APOLOGY", "SAZH_DEBT_2", "TRANSPORT_WRIT"} {
		s, _ := playFrom(t, "cantica", docked, "look")
		s.Set(src)
		if !s.Has("CONCLAVE_CHANNEL") {
			t.Errorf("%s should carry the indictment out", src)
		}
	}

	// The writ is one channel and one resupply order, not both.
	spent, _ := playFrom(t, "cantica", docked, "look")
	spent.Set("TRANSPORT_WRIT")
	spent.Set("KADMON_WRIT")
	if spent.Has("CONCLAVE_CHANNEL") {
		t.Fatal("a writ left on Kadmon is not a channel out of the sector")
	}
	// Unless Vaux gave you his, in which case the writ was never the channel.
	kept, _ := playFrom(t, "cantica", docked, "look")
	kept.Set("VAUX_APOLOGY")
	kept.Set("TRANSPORT_WRIT")
	kept.Set("KADMON_WRIT")
	if !kept.Has("CONCLAVE_CHANNEL") {
		t.Fatal("his dead-man writ does not go to Kadmon")
	}
}

// Cael files it to your case-notes without being asked. It has to actually be
// there, or the archive takes back something the player never saw.
func TestCaelsFilingReachesTheCaseNotes(t *testing.T) {
	s, _ := playFrom(t, "cantica", docked,
		"ask cael about the windows", "ask cael about the protocols")
	notes := s.caseNotesDump()
	if !strings.Contains(notes, "C17") {
		t.Fatalf("he said he filed it:\n%s", notes)
	}
	if !s.Holds("C17") {
		t.Fatal("a filed clue is a stored clue")
	}
	for i := 0; i < 8; i++ {
		s.Tick()
	}
	if strings.Contains(s.caseNotesDump(), "C17") {
		t.Fatal("and then it is not one of the lines any more")
	}
}
