package game

import (
	"strings"
	"testing"
)

// The playthrough tests. Everything else in this package teleports into a
// scene with the flags it needs already set, which is the right way to test
// one guarantee in isolation and says nothing at all about whether the scenes
// join up. These five walk the whole game from New(g) with nothing but typed
// commands, and each one ends on a different door.
//
// Both of the worst bugs this repo has had were in exactly that gap: every
// scene individually correct, and no legal path between two of them.

// walk plays a script from a cold start and returns the state and everything
// the game printed. Any unresolved authored string fails, as everywhere else.
func walk(t *testing.T, script ...string) (*State, string) {
	t.Helper()
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	s := New(g)
	var log strings.Builder
	for _, c := range script {
		out, ticks := s.Do(c)
		// A script is only evidence if every line of it did something. A
		// mistyped command would otherwise sit in here for years looking like
		// coverage.
		for _, dead := range []string{
			"That is not a verb this terminal knows",
			"You hold no such clue",
		} {
			if strings.Contains(strings.Join(out, "\n"), dead) {
				t.Fatalf("dead command %q: %s", c, strings.Join(out, "\n"))
			}
		}
		log.WriteString("> " + c + "\n")
		for _, o := range out {
			log.WriteString(o + "\n")
		}
		for i := 0; i < ticks; i++ {
			for _, e := range s.Tick() {
				log.WriteString(e + "\n")
			}
		}
	}
	out := log.String()
	if strings.Contains(out, "[missing string:") {
		i := strings.Index(out, "[missing string:")
		t.Fatalf("unresolved authored string: %s", out[i:min(i+48, len(out))])
	}
	// The thesis of the whole game, asserted over a whole game rather than
	// over one scene: the archive never announces itself, anywhere, once.
	for _, word := range []string{"DELETED", "ERROR", "CORRUPTED", "TAMPER"} {
		if strings.Contains(out, word) {
			i := strings.Index(out, word)
			t.Fatalf("the archive announced itself (%s): …%s…",
				word, out[max(0, i-120):min(i+120, len(out))])
		}
	}
	return s, out
}

func join(parts ...[]string) []string {
	var out []string
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

// Act I, played carefully: the finding, the elevation, all five signs, the
// room read backwards, the requisition, the archive caught in the act, and
// the Case Précis annotated.
func actOne() []string {
	return []string{
		"x writ", "x writ", "note A0",
		// The two things Act I gives an attentive player standing for.
		"ask about the form", "go elevation", "ask about why not vaux",
		"take rosette", "go chambers",
		"x boots", "x boots", "note A3",
		"x board", "x board", "file A1",
		"x shutter", "x shutter", "note A2",
		"x quill", "x quill", "note A4",
		"x calendar", "x calendar", "note A5",
		"read room backward", "note A8",
		"query C9-F4-D2", "note A20",
		"wait", "wait", "wait", "wait", "wait",
		"query C9-F4-D2",
		"wait", "wait", "wait", "wait", "wait",
		"wait", "wait", "wait", "wait", "wait",
		"annotate",
	}
}

// Two of Vosk's four open files, which is what summons Crase.
func twoThreads() []string {
	return []string{
		"note B0",
		"go fenn", "ask fenn about kadmon", "note B1", "go board",
		"go delta", "ask fenn about outgoing", "note B18", "go board",
	}
}

// The Gallery brief refused, a runner bought from Mott, and Vaux allied.
func vauxAllied() []string {
	return []string{
		"go crase", "refuse", "go board",
		"go sump", "ask mott about favours", "surrender", "send a runner", "go board",
		"go vaux",
		"ask vaux about the room", "ask vaux about vosk", "ask vaux about why me",
		"ask vaux about the cipher", "ask vaux about my answer", "ally",
		"note B16", "note B15", "go out",
	}
}

// The ash imaged, the sentence read, the act turned.
func actTwoTurn() []string {
	return []string{
		"go cael", "ask cael about the ash", "image it", "note B22", "go board",
		"decode", "note B19", "note B20",
	}
}

// Off Vardan Prime on a stolen hull, which is the only way nobody has been
// through the cutter before you.
func stealAndSail() []string {
	return []string{
		"go skydock",
		"talk to rill", "ask rill about the booking", "lean on him", "break the seal",
		"go depart",
		"go well", "ask spirit about the orders",
		"x nav", "x nav", "note C2",
		"x manifest", "x manifest", "note C3",
		"x drawer", "x drawer", "note C5",
		"go up",
	}
}

// Kadmon, and then the address.
func kadmonTo(cmds ...string) []string {
	return join([]string{"go translate", "go hall"}, cmds,
		[]string{"go down", "go onward"})
}

// The filed words, in order, and the precise question.
func theReading() []string {
	return []string{
		"go the fifth chair", "go the floor",
		"ask vox about minutes", "ask vox about kadmon", "ask vox about my file",
	}
}

func TestPlaythroughEndingE(t *testing.T) {
	s, log := walk(t, join(
		actOne(), twoThreads(), vauxAllied(), actTwoTurn(), stealAndSail(),
		[]string{"go wardroom", "ask vaux about the dossier", "ask vaux about the apology",
			"agree", "go out"},
		kadmonTo("ask vurr about the report", "ask vurr about the census", "tell her the truth"),
		theReading(),
		[]string{"ask vox about did you compose it", "note C20",
			"go the intake", "talk to it", "tell vaux", "the proper form"},
	)...)

	if !s.Has("ENDING_E") {
		t.Fatalf("the hidden door did not open:\n%s", tail(log))
	}
	if !s.Done {
		t.Fatal("an ending ends the game")
	}
	if !strings.Contains(flat(log), "It reads. It stops. The lights stay on.") {
		t.Fatalf("expected the last line:\n%s", tail(log))
	}
	// The counters land where the documents say a careful player lands.
	// STORY_BIBLE §12 and ACT_II: standing opens at +1 and Act I is built to
	// push an attentive player to +2 or +3; VERITY is contradictions recorded
	// outside the cogitator, not a tally of everything written down, and Act
	// II's honest maximum is six.
	if v := s.Ints["VERITY"]; v < 3 || v > 8 {
		t.Errorf("VERITY %d is outside the range the documents describe", v)
	}
	if r := s.Ints["RECONCILIATION"]; r > 25 {
		t.Errorf("a ledger-keeping player should finish quiet, got %d", r)
	}
	if s.Has("BEZ_REMEMBERS") {
		t.Error("and should never see the forty-point symptom")
	}
	// All ten, reached by playing rather than by being handed them.
	for _, f := range []string{
		"VAUX_ALIVE", "VAUX_ALLIED", "FENN_INTACT", "ROOM_READ_FULL", "DECOY_REFUSED",
		"ERRATUM_OBTAINED", "ERRATUM_UNUSED", "CIPHER_COMPLETE",
		"VOSK_METHOD_UNDERSTOOD", "AMANUENSIS_BLINDSPOT",
	} {
		if !s.Has(f) {
			t.Errorf("reached the door without %s", f)
		}
	}
}

func TestPlaythroughEndingB(t *testing.T) {
	s, log := walk(t, join(
		actOne(), twoThreads(), vauxAllied(), actTwoTurn(), stealAndSail(),
		kadmonTo("ask vurr about the report", "ask vurr about the census", "tell her the truth"),
		theReading(),
		[]string{"go the intake", "talk to it", "ask it about the numbers", "sit down"},
	)...)
	if !s.Has("ENDING_B") {
		t.Fatalf("the chair is always on offer:\n%s", tail(log))
	}
	if !strings.Contains(flat(log), "You cannot remember what was wrong with it.") {
		t.Fatalf("expected the last line:\n%s", tail(log))
	}
	// The rewrite unmakes the first thing the player worked out, first.
	if !strings.Contains(flat(log), "the shutter is closed") {
		t.Fatalf("expected the rewrite:\n%s", tail(log))
	}
}

func TestPlaythroughEndingA(t *testing.T) {
	s, log := walk(t, join(
		actOne(), twoThreads(), vauxAllied(), actTwoTurn(), stealAndSail(),
		kadmonTo("ask vurr about the report", "ask vurr about the census", "tell her the truth"),
		theReading(),
		[]string{"go the intake", "burn it"},
	)...)
	if !s.Has("ENDING_A") {
		t.Fatalf("Bez carries the flask:\n%s", tail(log))
	}
	if !strings.Contains(flat(log), "It is the least useful thing you have ever made.") {
		t.Fatalf("expected the last line:\n%s", tail(log))
	}
}

func TestPlaythroughEndingC(t *testing.T) {
	s, log := walk(t, join(
		actOne(), twoThreads(), vauxAllied(), actTwoTurn(), stealAndSail(),
		[]string{"go wardroom", "ask vaux about the dossier", "ask vaux about the apology",
			"agree", "go out"},
		kadmonTo("ask vurr about the report", "ask vurr about the census", "tell her the truth"),
		[]string{"talk to kesh", "accept the tea", "ask kesh about the numbers", "note C16"},
		theReading(),
		[]string{"go the intake", "tell vaux", "walk out"},
	)...)
	if !s.Has("ENDING_C") {
		t.Fatalf("the indictment needs legs to walk out on:\n%s", tail(log))
	}
	if !s.Has("CONCLAVE_CHANNEL") {
		t.Fatal("the channel is the price of a promise")
	}
	if !strings.Contains(flat(log), "every one of them is a year in which you were right") {
		t.Fatalf("expected the last line:\n%s", tail(log))
	}
}

// D inverts the whole game: exposure is the currency. This path is the
// careless one on purpose — it sleeps through the passage, says nothing to the
// Marchioness, and files everything.
func TestPlaythroughEndingD(t *testing.T) {
	s, log := walk(t, join(
		actOne(), twoThreads(), vauxAllied(), actTwoTurn(),
		// Filing is what the archive is for. A player who uses the case-file
		// the way the Crown intends pays for it in exposure.
		[]string{"file B0", "file B1", "file B15", "file B16", "file B18", "file B22"},
		[]string{
			"go skydock",
			"talk to rill", "ask rill about the booking", "lean on him", "break the seal",
			"go depart", "go sleep",
			"go hall", "ask vurr about the report", "ask vurr about the census",
			"file C14", "say nothing", "take a witness",
			"go down", "go onward",
		},
		[]string{"talk to kesh", "ask kesh about the numbers", "file C16"},
		theReading(),
		[]string{"file C18", "file C19", "stop it",
			"go the intake", "talk to it", "ask it about the projection outside",
			"file C21", "one correction"},
	)...)
	if !s.Has("ENDING_D") {
		t.Fatalf("expected the trade (RECONCILIATION %d):\n%s",
			s.Ints["RECONCILIATION"], tail(log))
	}
	if s.Ints["RECONCILIATION"] < 25 {
		t.Fatalf("the currency is exposure, got %d", s.Ints["RECONCILIATION"])
	}
	// And the symptoms are not theoretical. A careless finish crosses forty,
	// so the retinue remembers something the player has no record of, in a
	// playthrough rather than in a test that set the counter by hand.
	if s.Ints["RECONCILIATION"] < 40 || !s.Has("BEZ_REMEMBERS") {
		t.Fatalf("expected the sub-level nine symptom in play: %d, fired=%v",
			s.Ints["RECONCILIATION"], s.Has("BEZ_REMEMBERS"))
	}
	if !strings.Contains(flat(log), "there is no longer any word for the thing that sent it") {
		t.Fatalf("expected the last line:\n%s", tail(log))
	}
}

// tail is the last of a very long transcript, which is the useful part when a
// walk stops early.
func tail(log string) string {
	if len(log) > 3000 {
		return "…" + log[len(log)-3000:]
	}
	return log
}

// Act I is built to push an attentive player to +2 or +3 with Crase, from an
// opening +1, and to give them nought to four points of independent proof.
// Both were true on paper and neither was checked, and the walks that check
// everything else had been taking none of the standing on offer.
func TestActOneLandsWhereTheDocumentsSayItDoes(t *testing.T) {
	s, _ := walk(t, actOne()...)

	if c := s.Ints["CRASE_STANDING"]; c < 2 || c > 3 {
		t.Errorf("CRASE_STANDING %+d at the act turn; the documents say +2 or +3", c)
	}
	if v := s.Ints["VERITY"]; v < 1 || v > 4 {
		t.Errorf("VERITY %d at the act turn; ACT_I gives the range 0–4", v)
	}
	if !s.Has("A21_HELD") {
		t.Error("a player who transcribed the dispatch can prove the edit")
	}
	if !s.Has("ACT1_COMPLETE") {
		t.Error("the Précis closes the act")
	}
	// The archive has touched a careful player barely at all by the turn.
	if r := s.Ints["RECONCILIATION"]; r > 6 {
		t.Errorf("RECONCILIATION %d after one filing and one requisition", r)
	}
}
