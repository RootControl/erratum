package game

import (
	"strings"
	"testing"
)

// landed is a player who has made fast at Kadmon Veil with the retinue intact.
var landed = append(append([]string{}, aboard...), "KADMON_VISITED")

// They drilled it every watch for ninety-one years against the day.
func TestTheyHailYouCorrectly(t *testing.T) {
	_, log := playFrom(t, "kadmon", landed, "talk to rhod", "ask rhod about the hail")
	if !strings.Contains(flat(log), "the drill is run every watch") {
		t.Fatalf("expected the challenge:\n%s", log)
	}
	if !strings.Contains(flat(log), "He is eleven and he did not fumble a syllable") {
		t.Fatalf("somebody had the watch this morning:\n%s", log)
	}
	if !strings.Contains(flat(log), "taught it to their replacement, and their replacement's replacement") {
		t.Fatalf("ninety-one years of it:\n%s", log)
	}
}

// Ninety-one cycles, sealed, on a world that has been hungry since year forty.
func TestTheOwedIsNeverOpened(t *testing.T) {
	s, log := playFrom(t, "tithe_yard", landed, "x stacks", "x stacks", "x seals", "x seals")
	if !s.Has("KADMON_SAW_OWED") {
		t.Fatal("walking it is the branch")
	}
	if _, held := s.Held["C11"]; !held {
		t.Fatalf("expected the Owed:\n%s", log)
	}
	if !strings.Contains(flat(log), "one year and four months") {
		t.Fatalf("she recalculates it every cycle:\n%s", log)
	}
	if !strings.Contains(flat(log), "It is not ours") {
		t.Fatalf("Rhod is genuinely confused by the question:\n%s", log)
	}
}

// She will not falsify, and she has been holding a false number at the head of
// her own ledger for sixty years in order not to.
func TestVurrWillNotCorrectTerra(t *testing.T) {
	s, log := playFrom(t, "muster", landed, "talk to vurr", "ask vurr about the census")
	if _, held := s.Held["C14"]; !held {
		t.Fatalf("expected the census:\n%s", log)
	}
	if !strings.Contains(flat(log), "I am authorised to be *accurate about being unable to*") {
		t.Fatalf("expected her justification:\n%s", log)
	}
	// And the cogitator does in nine seconds the thing she refused for sixty years.
	out, _ := s.Do("file C14")
	if !strings.Contains(flat(strings.Join(out, "\n")), "the acknowledged figure is authoritative pending reassessment") {
		t.Fatalf("the standing rule resolves it:\n%s", strings.Join(out, "\n"))
	}
}

// Pity is what closes her. Competence is what opens her.
func TestPityClosesHerAndCitationOpensHer(t *testing.T) {
	s, log := playFrom(t, "muster", landed, "ask vurr about the census", "ask vurr about pity")
	if s.Ints["VURR_RESPECT"] != 0 {
		t.Fatalf("expected 0 after pity, got %d", s.Ints["VURR_RESPECT"])
	}
	if !strings.Contains(flat(log), "feel it after you have read the ledger") {
		t.Fatalf("cut off cold:\n%s", log)
	}
	if _, ticks := s.Do("go reading room"); ticks != 0 {
		t.Fatal("the Reading Room is not offered to somebody who has not earned it")
	}

	cited, _ := playFrom(t, "muster", landed, "look")
	cited.Grant("C11")
	run(cited, "show C11 to vurr")
	if cited.Ints["VURR_RESPECT"] != 1 {
		t.Fatalf("accurate citation opens her, got %d", cited.Ints["VURR_RESPECT"])
	}
	if _, ticks := cited.Do("go reading room"); ticks == 0 || cited.Room.ID != "reading_room" {
		t.Fatal("she grants it herself")
	}
}

// The question, and the three answers to it.
func TestTheUsefulQuestion(t *testing.T) {
	s, log := playFrom(t, "muster", landed,
		"ask vurr about the report", "ask vurr about the census")
	if !s.Has("VURR_ASKED") {
		t.Fatalf("she asks it:\n%s", log)
	}
	if !strings.Contains(flat(log), "Are we being punished, or are we being forgotten?") {
		t.Fatalf("expected the question:\n%s", log)
	}
	if !strings.Contains(flat(log), "I should like to know whether I prepared for the wrong thing") {
		t.Fatalf("expected the second half:\n%s", log)
	}
}

func TestTheThreeAnswers(t *testing.T) {
	answers := []struct {
		cmd     string
		flag    string
		verity  int
		recon   int
		wants   string
		forbids []string
	}{
		{"tell her the truth", "KADMON_TOLD", 2, 0,
			"A thing said is weather", []string{"KADMON_LIED", "KADMON_SILENT"}},
		{"lie to her", "KADMON_LIED", -3, 0,
			"May I see them?", []string{"KADMON_TOLD", "KADMON_SILENT"}},
		{"say nothing", "KADMON_SILENT", 0, 5,
			"Very well.", []string{"KADMON_TOLD", "KADMON_LIED"}},
	}
	for _, a := range answers {
		t.Run(a.flag, func(t *testing.T) {
			s, log := playFrom(t, "muster", landed,
				"ask vurr about the report", "ask vurr about the census", a.cmd)
			if !s.Has(a.flag) {
				t.Fatalf("expected %s:\n%s", a.flag, log)
			}
			for _, f := range a.forbids {
				if s.Has(f) {
					t.Fatalf("%s also set %s", a.flag, f)
				}
			}
			if s.Ints["VERITY"] != a.verity {
				t.Fatalf("%s: expected VERITY %d, got %d", a.flag, a.verity, s.Ints["VERITY"])
			}
			if s.Ints["RECONCILIATION"] != a.recon {
				t.Fatalf("%s: expected %d exposure, got %d", a.flag, a.recon, s.Ints["RECONCILIATION"])
			}
			if !strings.Contains(flat(log), a.wants) {
				t.Fatalf("%s missing %q:\n%s", a.flag, a.wants, log)
			}
		})
	}
}

// One is a gift. Two is a policy, and you are not staying.
func TestTheAidIsInsufficientAndExclusive(t *testing.T) {
	s, log := playFrom(t, "muster", landed,
		"ask vurr about the report", "ask vurr about the census",
		"tell her the truth", "relight the reactor")
	if !s.Has("KADMON_REACTOR") || !s.Has("KADMON_AID") {
		t.Fatalf("expected the reactor:\n%s", log)
	}
	// The ratio is said out loud, by her, on purpose.
	if !strings.Contains(flat(log), "two and four-tenths per cent of my grid") {
		t.Fatalf("she says the number:\n%s", log)
	}
	if !strings.Contains(flat(log), "neither of us has to pretend this afternoon was a rescue") {
		t.Fatalf("expected the knife:\n%s", log)
	}
	// And a shrine to it inside a week.
	if !strings.Contains(flat(log), "feast day") {
		t.Fatalf("the shrine is the knife:\n%s", log)
	}
	s.Do("take a witness")
	if s.Has("KADMON_WITNESS") {
		t.Fatal("one is a gift; two is a policy")
	}
}

// Do nothing at all and she still thanks you for coming.
func TestSheThanksYouForComingAnyway(t *testing.T) {
	s, log := playFrom(t, "muster", landed,
		"ask vurr about the report", "ask vurr about the census",
		"ask vurr about yourself", "ask vurr about the vigil")
	if s.Has("KADMON_AID") {
		t.Fatal("setup: nothing was given")
	}
	if !s.Has("VURR_THANKED") {
		t.Fatalf("the departure fires anyway:\n%s", log)
	}
	if !strings.Contains(flat(log), "\"Thank you for coming,\" the Marchioness says.") {
		t.Fatalf("expected the line:\n%s", log)
	}
	if !strings.Contains(flat(log), "She does not add a condition to it") {
		t.Fatalf("she means it entirely:\n%s", log)
	}
}

// Sixty years of reporting correctly into a wall.
func TestTheUnsentReportsAreTheProof(t *testing.T) {
	s, log := playFrom(t, "reading_room", append(landed, "VURR_CITED"),
		"look", "x reports", "x reports", "x index", "x index")
	if _, held := s.Held["C12"]; !held {
		t.Fatalf("expected the three:\n%s", log)
	}
	if !s.Has("KADMON_PROOF") {
		t.Fatal("the reports are the documentary proof")
	}
	if !strings.Contains(flat(log), "I record these not as complaint but as evidence, in the event that evidence is what is wanted") {
		t.Fatalf("415 is the one that hurts:\n%s", log)
	}
	if !strings.Contains(flat(log), "for the Emperor's attention or another's") {
		t.Fatalf("833 is the worst one:\n%s", log)
	}
	// Filing them says, in the filing, what will happen to them.
	out, _ := s.Do("file C12")
	if !strings.Contains(flat(strings.Join(out, "\n")), "two of the three will be shorter") {
		t.Fatalf("expected the seam:\n%s", strings.Join(out, "\n"))
	}
}

// Report 728 is a spine unless the player goes looking for it.
func TestReport728MustBeLookedFor(t *testing.T) {
	s, log := playFrom(t, "reading_room", append(landed, "VURR_CITED"), "x 728")
	if _, held := s.Held["C728"]; held {
		t.Fatal("it is one of nine hundred and twenty-two spines")
	}
	if !strings.Contains(log, "Shelves, a desk") {
		t.Fatalf("it is not in the room until you know the year:\n%s", log)
	}

	// Rhod knows the year. He has read all of them twice.
	told, _ := playFrom(t, "kadmon", landed, "ask rhod about her husband")
	if !told.Has("VURR_WIDOWED") {
		t.Fatal("he gives you the cycle")
	}
	if !told.Move("reading_room") {
		t.Fatal("cannot reach the room")
	}
	run(told, "x 728", "x 728")
	if _, held := told.Held["C728"]; !held {
		t.Fatal("expected report 728")
	}
	out, _ := told.Do("file C728")
	log = strings.Join(out, "\n")
	if told.Has("C728_FILED") {
		t.Fatal("the disposition options are all lies")
	}
	if !strings.Contains(flat(log), "you go back and you make the mark") {
		t.Fatalf("the refusal is the mark:\n%s", log)
	}
}
