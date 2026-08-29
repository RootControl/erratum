package game

import (
	"strings"
	"testing"
)

func playSump(t *testing.T, cmds ...string) (*State, string) {
	t.Helper()
	return playToPrecis(t)(append([]string{"sign", "go sump"}, cmds...)...)
}

// A settlement fed by the paperwork of its own non-existence: no entitlement,
// so no ration, so it eats what two accounting systems are obliged to write
// off somewhere.
func TestFedByTheirOwnNonExistence(t *testing.T) {
	_, log := playSump(t, "x spill", "x spill")
	if !strings.Contains(log, "write-off docket") {
		t.Fatal("every crate is correct paperwork")
	}
	if !strings.Contains(log, "the destination does not have to exist") {
		t.Fatal("that is the whole mechanism of the place")
	}
}

// Surrender the sidearm or force the gate. Sesh stands aside either way; only
// one of them costs you everything you ask for afterwards, and nobody ever
// tells you that is why.
func TestTheGateIsATrade(t *testing.T) {
	trusted, tLog := playSump(t, "surrender")
	if !trusted.Has("CHOIR_TRUSTED") || trusted.Ints["ILARIA_TRUST"] != 1 {
		t.Fatal("handing it over is worth a point with the astropath")
	}
	if !strings.Contains(tLog, "That was well done, sir") {
		t.Fatal("and Ilaria is the one who says so")
	}

	hostile, hLog := playSump(t, "force", "ask mott about favours", "owe him")
	if !hostile.Has("CHOIR_HOSTILE") {
		t.Fatal("the rosette works and is remembered")
	}
	if hostile.Ints["MOTT_FAVOURS"] != 2 {
		t.Fatalf("everything costs twice, for the rest of the act: %d", hostile.Ints["MOTT_FAVOURS"])
	}
	if strings.Contains(hLog, "twice as") || strings.Contains(hLog, "doubled") {
		t.Fatal("and nobody ever says that is why")
	}
}

// B6 is the one clue that cannot be recorded in any medium at all. Filing it
// and writing it down are the same refusal, and the refusal is four hundred
// people stopping.
func TestTheCensusCannotBeRecorded(t *testing.T) {
	s, log := playSump(t, "surrender",
		"ask mott about the name-bank", "ask mott about himself",
		"note the census", "file B6")
	if _, held := s.Held["B6"]; !held {
		t.Fatal("the Choir recites whether or not you were ready")
	}
	if s.Held["B6"] != StoreNone {
		t.Fatal("it goes in neither store, ever")
	}
	if s.Has("B6_PHYSICAL") || s.Has("B6_FILED") || len(s.Notebook) != 0 {
		t.Fatal("no copy, in any medium")
	}
	if strings.Count(log, "they can't get at a mouth") != 2 {
		t.Fatal("both refusals are the same refusal")
	}
	if !strings.Contains(log, "four hundred people have stopped") {
		t.Fatal("and trying is what stops them")
	}
}

// The purge order removes a record, not a people. Forty thousand were left
// exactly where they were, and that turned out to be sufficient.
func TestThePurgeOrderRemovesRecordOnly(t *testing.T) {
	s, log := playSump(t, "surrender", "ask mott about favours", "owe him",
		"go annexe", "x order", "x order", "note the purge order",
		"x flush", "x flush", "note the flush order")
	if _, held := s.Held["B4"]; !held {
		t.Fatal("the vellum is the clue")
	}
	if !strings.Contains(log, "It does not authorise the removal of anybody") {
		t.Fatal("the document is careful and the game must be as careful")
	}
	if !strings.Contains(log, "left precisely where they were") {
		t.Fatal("nobody was killed and nobody was moved")
	}
	if !s.Has("CLOCK_DAY40_ARMED") {
		t.Fatal("the second sheet arms the clock")
	}
	if !strings.Contains(log, "NO OCCUPANCY OF RECORD") {
		t.Fatal("and the second line of it is true")
	}
}

// The annexe is where Mott does business, and you have to have done business.
func TestTheAnnexeCostsAFavour(t *testing.T) {
	closed, log := playSump(t, "surrender", "go annexe")
	if closed.Room.ID == "annexe" {
		t.Fatal("he is perfectly friendly about it and does not move")
	}
	if !strings.Contains(log, "have not done any business") {
		t.Fatal("and the refusal should say so")
	}
	open_, _ := playSump(t, "surrender", "ask mott about favours", "owe him", "go annexe")
	if open_.Room.ID != "annexe" {
		t.Fatal("owing him is the whole of the price")
	}
}

// Query a man who does not exist and forty minutes later something confirms
// the nil return, courteously, unasked, over the signature of a clerk who does
// not exist either. Nobody confirms a nil return.
func TestTheArchiveAnswersAQuestionNobodyRepeated(t *testing.T) {
	s, log := playSump(t, "surrender", "ask mott about favours", "owe him",
		"go annexe", "x order", "x order", "out", "back",
		"query casmir lyle",
		"wait", "wait", "wait", "wait", "wait", "wait")
	if !s.Has("LYLE_QUERIED") {
		t.Fatal("the query goes into the log the way every query does")
	}
	if !strings.Contains(log, "There has never been a seat 411") {
		t.Fatal("not struck, not dead — never a seat")
	}
	if _, held := s.Held["B5"]; !held {
		t.Fatal("the confirmation is the clue")
	}
	if !strings.Contains(log, "Nobody confirms a nil return") {
		t.Fatal("and the point of it is that there is no procedure for it")
	}
	if !strings.Contains(log, "There is no Clerk-Ordinate M. Vestin") {
		t.Fatal("nor for the clerk who sent it")
	}
}

// The Rows are the act's spine and nothing is attached to them. Four ways to
// help, all real, all costly, and not one of them gives the player anything.
func TestHelpingBuysNothing(t *testing.T) {
	base, _ := playSump(t, "surrender", "ask mott about favours", "owe him",
		"go annexe", "x order", "x order", "out", "go rows")
	before := len(base.Held)

	s, log := playSump(t, "surrender", "ask mott about favours", "owe him",
		"go annexe", "x order", "x order", "out", "go rows",
		"ask wick about day 40", "empty the account")
	if !s.Has("CHOIR_FED") || !s.Has("TRANSPORT_DEBT") {
		t.Fatal("emptying the account feeds them and mortgages Act III")
	}
	if len(s.Held) != before {
		t.Fatal("no clue is attached to any of this")
	}
	if s.Ints["VERITY"] != base.Ints["VERITY"] {
		t.Fatal("and no stat either")
	}
	if !strings.Contains(log, "That was stupid") || !strings.Contains(log, "Do it again in forty days") {
		t.Fatal("Mott never thanks anybody")
	}
	if !s.Has("WICK_KNOWN") {
		t.Fatal("she asks your name and files you, which is not a reward")
	}

	// Deferring costs standing, and Crase does not ask a single question.
	deferred, dLog := playSump(t, "surrender", "ask mott about favours", "owe him",
		"go annexe", "x order", "x order", "x flush", "x flush", "out", "go rows", "defer")
	if !deferred.Has("CHOIR_DEFERRED") {
		t.Fatal("sixty days")
	}
	if deferred.Ints["CRASE_STANDING"] != -1 { // opens at +1; signing the précis does not move it
		t.Fatalf("two points of standing: %d", deferred.Ints["CRASE_STANDING"])
	}
	if !strings.Contains(dLog, "it is going to be a receipt") {
		t.Fatal("it feels wonderful now and is a receipt later")
	}
}

// A scene's date is the earliest it can happen on, not a fact about where the
// player is standing. Walking back to the hub does not rewind the calendar.
func TestTheCalendarOnlyMovesForward(t *testing.T) {
	s, _ := playSump(t, "surrender", "back")
	if s.Room.ID != "standing_board" {
		t.Fatal("back up the pipe")
	}
	if !strings.HasPrefix(s.TimeString(), "0997.4.224") {
		t.Fatalf("the day the sump happened on is still the day: %s", s.TimeString())
	}
}
