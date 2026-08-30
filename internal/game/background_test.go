package game

import (
	"strings"
	"testing"
)

// The three openings the bible offers: same hook, different training. Each
// buys the player something in Act I and charges for it somewhere else, and
// the game never says which is which.

func inChambers(t *testing.T, bg string, cmds ...string) (*State, string) {
	t.Helper()
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	s := New(g)
	s.SetBackground(bg)
	s.Set("ROSETTE_TAKEN")
	if !s.Move("vosk_chambers") {
		t.Fatal("no chambers")
	}
	var log strings.Builder
	for _, c := range cmds {
		out, ticks := s.Do(c)
		for _, o := range out {
			log.WriteString(o + "\n")
		}
		for i := 0; i < ticks; i++ {
			for _, e := range s.Tick() {
				log.WriteString(e + "\n")
			}
		}
	}
	if strings.Contains(log.String(), "[missing string:") {
		t.Fatalf("unresolved authored string in the %s chambers", bg)
	}
	return s, log.String()
}

// A Lexmechanic reads the board on sight and misreads the boots with total
// confidence, which is the trade.
func TestLexmechanicReadsTheBoardAndMisreadsTheBoots(t *testing.T) {
	s, log := inChambers(t, "lexmechanic", "x board", "x boots")
	if _, held := s.Held["A1"]; !held {
		t.Fatalf("the encoding is stated, not deduced:\n%s", log)
	}
	if !strings.Contains(flat(log), "It is Vault shelf-addressing") {
		t.Fatalf("expected the recognition:\n%s", log)
	}
	if !strings.Contains(flat(log), "What you do not have is the direction") {
		t.Fatalf("A3 is still needed to know which way round:\n%s", log)
	}
	// And the boots read as an ordinary distress classification, correctly,
	// which is the whole of the problem.
	if _, held := s.Held["A3"]; held {
		t.Fatal("the first look must not give A3")
	}
	if !strings.Contains(flat(log), "you are entirely correct about the classification") {
		t.Fatalf("expected the confident misread:\n%s", log)
	}
	// The second look still gets there, for a player who goes back.
	s.Do("x boots")
	if _, held := s.Held["A3"]; !held {
		t.Fatal("a second examination still finds it")
	}

	// Nobody else sees any of that.
	plain, plog := inChambers(t, "", "x board")
	if _, held := plain.Held["A1"]; held {
		t.Fatalf("without the training it is deduced, not stated:\n%s", plog)
	}
}

// And can prove the edit rather than believe it.
func TestLexmechanicProvesTheEdit(t *testing.T) {
	s, log := inChambers(t, "lexmechanic",
		"x board", "x board", "file A1",
		"query C9-F4-D2",
		"wait", "wait", "wait", "wait", "wait",
		"query C9-F4-D2",
	)
	if !strings.Contains(flat(log), "does not match the checksum in the session log") {
		t.Fatalf("expected the checksum:\n%s", tail(log))
	}
	if s.Ints["VERITY"] < 2 {
		t.Fatalf("two numbers, both on the screen: got VERITY %d", s.Ints["VERITY"])
	}
	if !s.Has("VAULT_FORENSICS") {
		t.Fatal("Act II opens a thread nobody else gets")
	}

	// Anybody else has a strong feeling and knows what that is worth.
	other, olog := inChambers(t, "",
		"x board", "x board", "file A1",
		"query C9-F4-D2",
		"wait", "wait", "wait", "wait", "wait",
		"query C9-F4-D2",
	)
	if strings.Contains(olog, "checksum") {
		t.Fatal("no other grade of adept is issued a terminal that shows it")
	}
	if other.Ints["VERITY"] != 0 {
		t.Fatalf("a strong feeling is not evidence: got %d", other.Ints["VERITY"])
	}
}

// A Provost reads the seal register second, always, and walks into the Vaux
// interview already knowing he was there.
func TestProvostGetsTheSealFree(t *testing.T) {
	s, log := inChambers(t, "provost", "x register")
	if _, held := s.Held["A7"]; !held {
		t.Fatalf("second item, always:\n%s", log)
	}
	if !strings.Contains(flat(log), "it has been second since you were nineteen") {
		t.Fatalf("expected the protocol:\n%s", log)
	}

	plain, _ := inChambers(t, "", "x register")
	if _, held := plain.Held["A7"]; held {
		t.Fatal("everybody else has to look twice")
	}
}

// Compulsion buys everything that was asked for and nothing that was not.
// A9 is a thing a man says without knowing he says it, and he will not.
func TestCompulsionCostsTheMnemonic(t *testing.T) {
	s, log := playFrom(t, "vaux_study", []string{"ROSETTE_TAKEN", "BG_PROVOST"},
		"compel", "ask vaux about the seal", "ask vaux about procedure",
		"ask vaux about his seniority", "ask vaux about crase")
	if !s.Has("VAUX_COMPELLED") {
		t.Fatalf("expected the compulsion:\n%s", log)
	}
	if strings.Contains(log, "A9") {
		t.Fatalf("nothing surfaces under compulsion that was not asked for:\n%s", log)
	}
	// A9 arrives on pressure, and pressure is what pushing on a closed door
	// costs a person. Under compulsion nothing is closed, so nothing is
	// pushed on, so the mnemonic has nowhere to come from.
	if _, held := s.Held["A9"]; held {
		t.Fatal("a sound a man makes without knowing he makes it is not an answer")
	}
	// Asking about him as a person is what ends the compulsion.
	s.Do("ask vaux about hadrian")
	if s.Has("VAUX_COMPELLED") {
		t.Fatal("asking about him as a person ends it")
	}
}

// The Penitent stands in the shape of it, and pays for having been in the
// record once already.
func TestPenitentFeelsTheEdge(t *testing.T) {
	s, log := inChambers(t, "penitent", "x the absence", "x the absence")
	if _, held := s.Held["A12"]; !held {
		t.Fatalf("A12 arrives early:\n%s", log)
	}
	if _, held := s.Held["A12b"]; !held {
		t.Fatal("and the edge with it")
	}
	if !strings.Contains(flat(log), "It is not a wound in the room. It is a *cut*.") {
		t.Fatalf("expected the distinction:\n%s", log)
	}
	if s.Ints["RECONCILIATION"] != 5 {
		t.Fatalf("the archive has been in you once already: got %d",
			s.Ints["RECONCILIATION"])
	}
	// A12b is the player's own and the template has no field for what they are.
	out, _ := s.Do("file A12b")
	if !strings.Contains(strings.Join(out, "\n"), "neither of them has a field for what you are") {
		t.Fatalf("expected the refusal: %s", strings.Join(out, "\n"))
	}

	// Nobody else finds anything in that room to look at.
	plain, plog := inChambers(t, "", "x the absence")
	if _, held := plain.Held["A12"]; held {
		t.Fatalf("there is nothing there for anybody else:\n%s", plog)
	}
	if plain.Ints["RECONCILIATION"] != 0 {
		t.Fatal("and nobody else starts owing anything")
	}
}

// Disclosing costs nothing and buys the only person who will lie for you.
// Being caught costs a little and cannot be undone.
func TestIlariaKnowsEitherWay(t *testing.T) {
	told, log := playFrom(t, "choir", []string{"ROSETTE_TAKEN", "BG_PENITENT"},
		"talk to ilaria", "disclose")
	if !told.Has("PENITENT_DISCLOSED") || told.Ints["ILARIA_TRUST"] < 2 {
		t.Fatalf("expected the disclosure:\n%s", log)
	}
	if !strings.Contains(flat(log), "she will lie for you") {
		t.Fatalf("expected what it buys:\n%s", log)
	}

	caught, clog := playFrom(t, "choir",
		[]string{"ROSETTE_TAKEN", "BG_PENITENT", "HAS_CYLINDER_6"}, "talk to ilaria")
	if !caught.Has("PENITENT_CAUGHT") {
		t.Fatalf("she is blind:\n%s", clog)
	}
	if !strings.Contains(flat(clog), "That's twice now.") {
		t.Fatalf("expected her line:\n%s", clog)
	}
	if caught.Ints["ILARIA_TRUST"] != 0 {
		t.Fatalf("expected the cost, got %d", caught.Ints["ILARIA_TRUST"])
	}
}

// A background is a choice made before the first command and never mentioned
// again, and none of them is required to finish.
func TestNoBackgroundIsTheDefault(t *testing.T) {
	for _, bg := range []string{"", "lexmechanic", "provost", "penitent", "nonsense"} {
		g, err := loadGame(t)
		if err != nil {
			t.Fatal(err)
		}
		s := New(g)
		s.SetBackground(bg)
		set := 0
		for _, f := range Backgrounds {
			if s.Has(f) {
				set++
			}
		}
		if bg == "" || bg == "nonsense" {
			if set != 0 {
				t.Errorf("%q set a background", bg)
			}
			continue
		}
		if set != 1 {
			t.Errorf("%q set %d backgrounds", bg, set)
		}
	}
}

// STORY_BIBLE §5: the name is player-set and defaults to Ammon Reyl. The
// prose is second person throughout and never says it. The record does, and
// what the record does with it is the point.
func TestTheRecordKnowsYourName(t *testing.T) {
	s, _ := playFrom(t, "standing_board", []string{"ACT2_OPEN"}, "look")

	out, _ := s.Do("query ROSETTE")
	full := strings.Join(out, "\n")
	// The Vault has had one vowel of it wrong since the player was nineteen.
	if !strings.Contains(full, "RAYL, AMMON") {
		t.Fatalf("expected the Vault's version:\n%s", full)
	}
	if strings.Contains(full, "REYL, AMMON") {
		t.Fatal("the Vault has never held it correctly")
	}
	// Cantica has it right, which is the only place that ever does.
	c, clog := playFrom(t, "cantica", docked, "x placard", "x placard")
	if !strings.Contains(clog, "REYL, AMMON") {
		t.Fatalf("the placard is spelled correctly:\n%s", clog)
	}
	_ = c

	// A player-set name propagates everywhere the tokens are.
	named, _ := playFrom(t, "standing_board", []string{"ACT2_OPEN"}, "look")
	named.SetName("Ottiline Vurr")
	out, _ = named.Do("query ROSETTE")
	full = strings.Join(out, "\n")
	if !strings.Contains(full, "VARR, OTTILINE") {
		t.Fatalf("the misspelling follows the name:\n%s", full)
	}
	// And past seventy the record keeps its own spelling and loses the rest.
	named.Add("RECONCILIATION", 70)
	out, _ = named.Do("query ROSETTE")
	if !strings.Contains(strings.Join(out, "\n"), "VARR, O.") {
		t.Fatalf("expected the initial:\n%s", strings.Join(out, "\n"))
	}
}

// The name survives a suspended session, because a session that came back
// under a different name would be the whole game happening by accident.
func TestTheNameSurvivesASession(t *testing.T) {
	g, err := loadGame(t)
	if err != nil {
		t.Fatal(err)
	}
	s := New(g)
	s.SetName("Casmir Lyle")
	path := t.TempDir() + "/s.json"
	if err := s.Save(path); err != nil {
		t.Fatal(err)
	}
	back := New(g)
	if err := back.Load(path); err != nil {
		t.Fatal(err)
	}
	if back.Given != "Casmir" || back.Family != "Lyle" {
		t.Fatalf("came back as %s %s", back.Given, back.Family)
	}
}
