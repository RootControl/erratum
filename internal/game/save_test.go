package game

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestSnapshotCarriesEveryField is the guard the save file needs and cannot
// have any other way. A snapshot that silently drops one of State's timers
// restores a session that is *almost* right — the body still leaves, but not
// when it should — which is the archive's own trick played on the player by
// the engine, with nothing on screen to say so.
//
// It cannot check that a field is carried *correctly*; the round-trip below
// does that. What it can do is fail the moment State grows a field, so that
// adding one is a decision rather than an omission.
func TestSnapshotCarriesEveryField(t *testing.T) {
	// Game is the compiled-in content and is never saved; a session is only
	// ever the player. Room is saved by ID rather than by pointer.
	const notSaved = 1 // Game
	const known = 45   // fields in State when the snapshot was last reviewed

	got := reflect.TypeOf(State{}).NumField()
	if got != known {
		t.Fatalf("State has %d fields and the snapshot was written for %d.\n"+
			"Add the new one to snapshot() and restore() — or, if it genuinely\n"+
			"must not survive a save, say so here and raise the count.", got, known)
	}
	if reflect.TypeOf(snapshot{}).NumField() != known-notSaved+3 {
		// +3: the Archive's three fields travel flattened into the snapshot.
		t.Fatalf("snapshot has %d fields; expected %d",
			reflect.TypeOf(snapshot{}).NumField(), known-notSaved+3)
	}
}

// A session written mid-playthrough and taken back up finishes the same game.
func TestASuspendedSessionFinishesTheSameWay(t *testing.T) {
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "erratum.session")

	// Play Act I and most of Act II, then put it down.
	first := New(g)
	for _, c := range join(actOne(), twoThreads(), vauxAllied()) {
		out, ticks := first.Do(c)
		for _, o := range out {
			if strings.Contains(o, "[missing string:") {
				t.Fatalf("%s: %s", c, o)
			}
		}
		for i := 0; i < ticks; i++ {
			first.Tick()
		}
	}
	if !first.Has("VAUX_ALLIED") {
		t.Fatal("setup: the walk did not reach the alliance")
	}
	if err := first.Save(path); err != nil {
		t.Fatal(err)
	}

	// Take it up again in a fresh process's worth of state.
	second := New(g)
	if err := second.Load(path); err != nil {
		t.Fatal(err)
	}
	if second.Room.ID != first.Room.ID {
		t.Fatalf("resumed in %s, left off in %s", second.Room.ID, first.Room.ID)
	}
	if !reflect.DeepEqual(second.Flags, first.Flags) {
		t.Fatal("the flags did not survive")
	}
	if !reflect.DeepEqual(second.Ints, first.Ints) {
		t.Fatal("the counters did not survive")
	}
	if !reflect.DeepEqual(second.Notebook, first.Notebook) {
		t.Fatal("the notebook did not survive, which is the one thing that must")
	}

	// And it finishes.
	rest := join(actTwoTurn(), stealAndSail(),
		[]string{"go wardroom", "ask vaux about the dossier",
			"ask vaux about the apology", "agree", "go out"},
		kadmonTo("ask vurr about the report", "ask vurr about the census",
			"tell her the truth"),
		theReading(),
		[]string{"ask vox about did you compose it", "note C20",
			"go the intake", "talk to it", "tell vaux", "the proper form"})
	for _, c := range rest {
		_, ticks := second.Do(c)
		for i := 0; i < ticks; i++ {
			second.Tick()
		}
	}
	if !second.Has("ENDING_E") {
		t.Fatal("a resumed session cannot reach the hidden door")
	}
}

// The timers are the part a careless save loses, and losing them would be
// invisible. Check the ones that are pure clockwork by hand.
func TestTheTimersSurvive(t *testing.T) {
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	s := New(g)
	// Wind every delayed system partway.
	s.Set("LYLE_QUERIED")
	s.Set("C17_FILED")
	s.Set("IN_TRANSIT")
	s.Set("CIPHER_COMPLETE")
	s.Set("CANTICA_QUERIED")
	s.Set("SUBLEVEL_NINE_SEEDED")
	s.Set("BEZ_ALIVE")
	s.Add("RECONCILIATION", 40)
	for i := 0; i < 3; i++ {
		s.Tick()
	}

	path := filepath.Join(t.TempDir(), "s.json")
	if err := s.Save(path); err != nil {
		t.Fatal(err)
	}
	back := New(g)
	if err := back.Load(path); err != nil {
		t.Fatal(err)
	}

	for _, f := range []struct {
		name string
		a, b int
	}{
		{"bodyDeadline", s.bodyDeadline, back.bodyDeadline},
		{"actions", s.actions, back.actions},
		{"lyleQueried", s.lyleQueried, back.lyleQueried},
		{"canticaQueried", s.canticaQueried, back.canticaQueried},
		{"tassWaited", s.tassWaited, back.tassWaited},
		{"bezWaited", s.bezWaited, back.bezWaited},
		{"transitTicks", s.transitTicks, back.transitTicks},
		{"symptomsSeen", s.symptomsSeen, back.symptomsSeen},
		{"c17Waited", s.c17Waited, back.c17Waited},
		{"sinceReq", s.sinceReq, back.sinceReq},
		{"driftPartial", s.driftPartial, back.driftPartial},
		{"clock", s.Clock, back.Clock},
	} {
		if f.a != f.b {
			t.Errorf("%s: %d on the way in, %d on the way out", f.name, f.a, f.b)
		}
	}
	if s.date != back.date {
		t.Errorf("date: %q -> %q", s.date, back.date)
	}
}

// Whatever the archive has already done to the case-notes travels with the
// save, because the case-file is the enemy's store and a session that came
// back clean would be a lie.
func TestTheArchiveTravelsAsItStands(t *testing.T) {
	s, _ := play(t,
		"x board", "x board", "file A1",
		"query C9-F4-D2",
		"wait", "wait", "wait", "wait", "wait", "wait",
	)
	if !s.Archive.Corrected() {
		t.Fatal("setup: the archive has not reconciled yet")
	}
	before := s.caseNotesDump()

	path := filepath.Join(t.TempDir(), "s.json")
	if err := s.Save(path); err != nil {
		t.Fatal(err)
	}
	back := New(s.Game)
	if err := back.Load(path); err != nil {
		t.Fatal(err)
	}
	if !back.Archive.Corrected() {
		t.Fatal("the correction did not survive")
	}
	if back.caseNotesDump() != before {
		t.Fatalf("the case-notes came back different:\n%s\n---\n%s",
			before, back.caseNotesDump())
	}
	if _, ok := back.Archive.Query(CodeLive); ok {
		t.Fatal("a record the archive deleted came back")
	}
}

// A session from another build is refused rather than half-applied.
func TestAForeignSessionIsRefused(t *testing.T) {
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	s := New(g)
	path := filepath.Join(t.TempDir(), "s.json")
	if err := s.Save(path); err != nil {
		t.Fatal(err)
	}
	raw, err := readFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var snap map[string]any
	if err := json.Unmarshal(raw, &snap); err != nil {
		t.Fatal(err)
	}
	snap["version"] = 99
	writeJSON(t, path, snap)

	fresh := New(g)
	if err := fresh.Load(path); err == nil {
		t.Fatal("a session from a different build should be refused")
	}
	if fresh.Room.ID != g.Start {
		t.Fatal("and a refused load must not have half-applied itself")
	}

	// So should one that names a scene this build does not have.
	if err := json.Unmarshal(raw, &snap); err != nil {
		t.Fatal(err)
	}
	snap["room"] = "a_room_that_was_cut"
	writeJSON(t, path, snap)
	if err := fresh.Load(path); err == nil || !strings.Contains(err.Error(), "does not have") {
		t.Fatalf("expected a scene refusal, got %v", err)
	}
}

func readFile(p string) ([]byte, error) { return osReadFile(p) }

func writeJSON(t *testing.T, p string, v any) {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := osWriteFile(p, b, 0o644); err != nil {
		t.Fatal(err)
	}
}
