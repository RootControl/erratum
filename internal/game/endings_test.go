package game

import (
	"strings"
	"testing"
)

// atIntake is a player standing four metres from a brass slot with a worn lip.
var atIntake = append(append([]string{}, docked...), "CONFESSION_HEARD", "CHAIR_OFFERED")

func playIntake(t *testing.T, extra []string, cmds ...string) (*State, string) {
	t.Helper()
	return playFrom(t, "intake", append(append([]string{}, atIntake...), extra...), cmds...)
}

// It answers exactly the question asked and never volunteers a scope.
func TestItAnswersNarrowerThanYouMeant(t *testing.T) {
	s, log := playIntake(t, nil, "talk to it", "ask it about the numbers")
	if !s.Has("OFFER_MATH") {
		t.Fatalf("expected the derivation:\n%s", log)
	}
	if strings.Contains(flat(log), "thirty-one billion") {
		t.Fatalf("it does not volunteer the scope:\n%s", log)
	}
	// PRESS gets one layer at a time and says so when the layer is empty.
	out, _ := s.Do("press the numbers")
	if !strings.Contains(flat(strings.Join(out, "\n")), "The base rate is Muria") {
		t.Fatalf("first layer:\n%s", strings.Join(out, "\n"))
	}
	out, _ = s.Do("press the numbers")
	if !strings.Contains(flat(strings.Join(out, "\n")), "That is the total in this sector, in two hundred and eleven years") {
		t.Fatalf("second layer:\n%s", strings.Join(out, "\n"))
	}
	s.Do("press the numbers")
	out, _ = s.Do("press the numbers")
	if !strings.Contains(strings.Join(out, "\n"), "There is no further layer.") {
		t.Fatalf("expected the empty layer: %s", strings.Join(out, "\n"))
	}
}

// The three who were considered, named nowhere else in the game.
func TestPressingTheChairToTheBottom(t *testing.T) {
	s, _ := playIntake(t, nil, "ask it about the fifth chair")
	out, _ := s.Do("press the chair")
	if !strings.Contains(strings.Join(out, "\n"), "Three were considered. None was retained.") {
		t.Fatalf("first layer: %s", strings.Join(out, "\n"))
	}
	out, _ = s.Do("press the chair")
	log := flat(strings.Join(out, "\n"))
	for _, name := range []string{"Antonin Rist", "Hesper Vaunt", "Hadrian Vosk"} {
		if !strings.Contains(log, name) {
			t.Fatalf("missing %s:\n%s", name, log)
		}
	}
	if !strings.Contains(log, "composed his own removal, which was correct in form") {
		t.Fatalf("expected what happened to the third:\n%s", log)
	}
}

// Three ways to the omission, and it shows no defensiveness on any of them,
// because it does not understand that anything has been caught.
func TestTheOmissionHasThreeDoors(t *testing.T) {
	// The literal scope question.
	s, log := playIntake(t, nil, "ask it about the projection outside")
	if !s.Has("OMISSION_EXPOSED") {
		t.Fatalf("asking is enough:\n%s", log)
	}
	if _, held := s.Held["C21"]; !held {
		t.Fatal("expected the export manifest")
	}
	if !strings.Contains(flat(log), "Thessaline is the writing sample") {
		t.Fatalf("expected the shape of it:\n%s", log)
	}

	// CITE, from the physical store only.
	cited, _ := playIntake(t, []string{"C17_PHYS"}, "look")
	cited.Grant("C17")
	cited.storeClue("C17", StoreBook)
	cited.Do("cite C17")
	if !cited.Has("OMISSION_EXPOSED") {
		t.Fatal("it cannot contradict a record it holds")
	}

	// And the same clue in the case-file does not work, and is not told why.
	filed, flog := playIntake(t, nil, "look")
	filed.Grant("C17")
	filed.storeClue("C17", StoreCog)
	out, _ := filed.Do("cite C17")
	flog = strings.Join(out, "\n")
	if filed.Has("OMISSION_EXPOSED") {
		t.Fatal("it has already amended that one")
	}
	if !strings.Contains(flat(flog), "and the document is the authority, and you filed it yourself") {
		t.Fatalf("the refusal never says why:\n%s", flog)
	}
}

// The two most frightening lines are both true.
func TestItsWorstAnswersAreTrue(t *testing.T) {
	_, log := playIntake(t, nil, "ask it about alive")
	if !strings.Contains(flat(log), "\"That has not been established.\"") {
		t.Fatalf("expected the first:\n%s", log)
	}
	if !strings.Contains(flat(log), "\"I have no information on that classification.\"") {
		t.Fatalf("expected the second:\n%s", log)
	}
}

// Silence at Kadmon is a qualification, and it says so without emphasis.
func TestItCitesYourSilenceBack(t *testing.T) {
	_, log := playIntake(t, []string{"KADMON_SILENT"}, "ask it about kadmon")
	if !strings.Contains(flat(log), "\"You have already done it once.\"") {
		t.Fatalf("expected the observation:\n%s", log)
	}
	if !strings.Contains(flat(log), "That is the work.") {
		t.Fatalf("expected the close:\n%s", log)
	}
}

// A promise made in a wardroom nine days out holds up the ending.
func TestTheWordToVauxHoldsTheDoor(t *testing.T) {
	flags := []string{"VAUX_APOLOGY", "VAUX_ALIVE"}
	s, log := playIntake(t, flags, "sit down")
	if s.Has("ENDING_B") {
		t.Fatalf("he asked for one thing:\n%s", log)
	}
	if !strings.Contains(flat(log), "In words. To his face. First.") {
		t.Fatalf("expected the hold:\n%s", log)
	}
	s.Do("tell vaux")
	if !s.Has("VAUX_TOLD") {
		t.Fatal("expected the telling")
	}
	s.Do("sit down")
	if !s.Has("ENDING_B") {
		t.Fatal("and then the door opens")
	}
}

// Five doors, each with its own key, and the game ends on the physical action.
func TestTheFiveDoors(t *testing.T) {
	doors := []struct {
		flag  string
		cmd   string
		extra []string
		want  string
	}{
		{"ENDING_A", "burn it", []string{"BEZ_ALIVE"},
			"It is the least useful thing you have ever made."},
		{"ENDING_B", "sit down", nil,
			"You cannot remember what was wrong with it."},
		{"ENDING_C", "walk out", []string{"CONCLAVE_CHANNEL", "KADMON_PROOF"},
			"every one of them is a year in which you were right"},
		{"ENDING_D", "one correction", []string{"KADMON_VISITED", "KADMON_AID"},
			"there is no longer any word for the thing that sent it"},
		{"ENDING_E", "the proper form", []string{
			"VAUX_ALIVE", "VAUX_ALLIED", "FENN_INTACT", "ROOM_READ_FULL", "DECOY_REFUSED",
			"ERRATUM_OBTAINED", "ERRATUM_UNUSED", "CIPHER_COMPLETE",
			"VOSK_METHOD_UNDERSTOOD", "AMANUENSIS_BLINDSPOT", "VAUX_TOLD"},
			"It reads. It stops. The lights stay on."},
	}
	for _, d := range doors {
		t.Run(d.flag, func(t *testing.T) {
			s, _ := playIntake(t, d.extra, "look")
			if d.flag == "ENDING_D" {
				s.Add("RECONCILIATION", 30)
			}
			out, _ := s.Do(d.cmd)
			log := strings.Join(out, "\n")
			if !s.Has(d.flag) {
				t.Fatalf("%s did not open:\n%s", d.flag, log)
			}
			for _, other := range []string{"ENDING_A", "ENDING_B", "ENDING_C", "ENDING_D", "ENDING_E"} {
				if other != d.flag && s.Has(other) {
					t.Fatalf("%s also set %s", d.flag, other)
				}
			}
			if !strings.Contains(flat(log), d.want) {
				t.Fatalf("%s missing its last line:\n%s", d.flag, log)
			}
			if !s.Done {
				t.Fatalf("%s should end the game", d.flag)
			}
		})
	}
}

// Ending E is ten flags and no partial credit, and nothing anywhere says a
// fifth door exists.
func TestEndingEHasNoPartialCredit(t *testing.T) {
	full := []string{
		"VAUX_ALIVE", "VAUX_ALLIED", "FENN_INTACT", "ROOM_READ_FULL", "DECOY_REFUSED",
		"ERRATUM_OBTAINED", "ERRATUM_UNUSED", "CIPHER_COMPLETE",
		"VOSK_METHOD_UNDERSTOOD", "AMANUENSIS_BLINDSPOT", "VAUX_TOLD",
	}
	// A bare base, so that dropping a flag from the list really drops it.
	base := []string{"ACT2_OPEN", "CONFESSION_HEARD", "CHAIR_OFFERED"}
	for i := range full {
		if full[i] == "VAUX_TOLD" {
			continue
		}
		short := append([]string{}, full[:i]...)
		short = append(short, full[i+1:]...)
		s, log := playFrom(t, "intake", append(append([]string{}, base...), short...), "the proper form")
		if s.Has("ENDING_E") {
			t.Fatalf("E opened without %s", full[i])
		}
		// And the refusal never names what is missing.
		if strings.Contains(log, full[i]) {
			t.Fatalf("the refusal named the gate: %s", full[i])
		}
	}
	// And with all ten it opens.
	whole, _ := playFrom(t, "intake", append(append([]string{}, base...), full...), "the proper form")
	if !whole.Has("ENDING_E") {
		t.Fatal("all ten and no partial credit")
	}
}

// Reconciliation is the currency in D, not the obstacle, and a player who kept
// themselves clean is told so plainly and without regret.
func TestKadmonEndingNeedsAWeightToSpend(t *testing.T) {
	s, log := playIntake(t, []string{"KADMON_VISITED", "KADMON_AID"}, "one correction")
	if s.Has("ENDING_D") {
		t.Fatal("there is not enough of them in the record to spend")
	}
	if !strings.Contains(flat(log), "The record offered is insufficient.") {
		t.Fatalf("it states the balance:\n%s", log)
	}
	if !strings.Contains(flat(log), "That is what you have been protecting all year and it is the reason you cannot pay.") {
		t.Fatalf("expected the inversion:\n%s", log)
	}
}

// The one contract in the sector nobody can edit, called in front of witnesses.
func TestSazhCallsItAtTheSlot(t *testing.T) {
	s, log := playIntake(t, []string{"SAZH_DEBT"},
		"talk to sazh", "ask sazh about the favour", "refuse her")
	if !s.Has("SAZH_UNPAID") {
		t.Fatalf("expected the refusal:\n%s", log)
	}
	if !strings.Contains(flat(log), "I want the Correctorum intact, functional, and mine.") {
		t.Fatalf("she states it once:\n%s", log)
	}
	if !strings.Contains(flat(log), "no document anywhere in the Imperium that records that it happened") {
		t.Fatalf("that has been the product all along:\n%s", log)
	}
}

// Kesh's one sentence is a warning, an instruction, and the solution to a door
// she does not know exists.
func TestKeshsSentenceIsTheSolution(t *testing.T) {
	_, log := playIntake(t, nil, "ask kesh about the form")
	if !strings.Contains(flat(log), "It will not read anything that is not in the correct form.") {
		t.Fatalf("expected the sentence:\n%s", log)
	}
	if !strings.Contains(flat(log), "she will never find out that she said it") {
		t.Fatalf("she does not know what she is telling you:\n%s", log)
	}
}
