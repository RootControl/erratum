package game

import (
	"strings"
	"testing"
)

// QUERY needs a cogitator. Exactly one scene has one, and Vaux's study — forty
// years of casework in paper, with the socket capped and painted over — is
// pointedly not it. The verb used to reach for chambers-only prose from
// anywhere and render "[missing string: vault_none]" at the player.
func TestQueryNeedsATerminal(t *testing.T) {
	study, studyLog := play(t,
		"x board", "x board", "file A1",
		"go vaux", "query C9-F4-D2",
	)
	if study.Has("KESH_NAMED") {
		t.Fatal("you cannot requisition from a room with no cogitator in it")
	}
	if !strings.Contains(studyLog, "no socket for one") {
		t.Fatal("and the refusal should be in the scene's own voice")
	}

	here, hereLog := play(t, "x board", "x board", "file A1", "query C9-F4-D2")
	if !here.Has("KESH_NAMED") {
		t.Fatal("in the chambers it works")
	}
	if !strings.Contains(hereLog, "1188") {
		t.Fatal("the dispatch renders")
	}
}

// READ ROOM synthesises the room you are standing in. Carrying five signs out
// of the chambers does not make Vaux's study readable.
func TestReadRoomIsAboutThisRoom(t *testing.T) {
	s, log := play(t,
		"x boots", "x boots", "note A3",
		"x board", "x board", "file A1",
		"x shutter", "x shutter", "note A2",
		"go vaux",
		"read room backward",
	)
	if s.Has("ROOM_READ_PARTIAL") || s.Has("ROOM_READ_FULL") {
		t.Fatal("the signs are facts about the chambers, not about wherever you are")
	}
	if !strings.Contains(log, "not have enough of this room") {
		t.Fatal("it should say so in the shared voice, not a missing-string marker")
	}
}

// A bare GO used to guess when a scene had one exit, which silently spent
// three actions and cost the player the boots.
func TestBareGoNamesTheDoors(t *testing.T) {
	s, log := play(t, "go")
	if s.Room.ID != "vosk_chambers" {
		t.Fatal("a bare GO must not move anybody")
	}
	if s.Clock != 6*3600+12*60 {
		t.Fatal("and must not cost anything")
	}
	if !strings.Contains(log, "GO where?") || !strings.Contains(log, "STUDY") {
		t.Fatalf("it should name the ways out: %q", log)
	}
}

// Looking at a person twice reprinted the whole portrait and charged a timed
// action for it, which objects have never done.
func TestPeopleReadLikeObjects(t *testing.T) {
	s, log := play(t, "go vaux", "x vaux", "x vaux")
	if strings.Count(log, "Sixty-one, dry, immaculate") != 1 {
		t.Fatal("the second look is not the first look again")
	}
	if !strings.Contains(log, "look like an adept") {
		t.Fatal("and it should say so the way the room does")
	}
	// go(3) + one look. The second look is free.
	if s.Clock != 6*3600+12*60+4*40 {
		t.Fatalf("the free look is free: %d", s.Clock)
	}
}
