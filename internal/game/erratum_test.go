package game

import (
	"strings"
	"testing"
)

// whole is the state 2.11 opens on: both halves of the cipher, and the room
// read at least once in Act I.
var whole = []string{"CIPHER_COMPLETE", "CIPHER_HALF_ROOM"}

func playWorkshop(t *testing.T, flags []string, cmds ...string) (*State, string) {
	t.Helper()
	return playFrom(t, "workshop", append([]string{"ACT2_OPEN"}, flags...), cmds...)
}

// The imaging is a favour, freely given, on a Tuesday, by a man who does not
// ask what it is.
func TestCaelNeverAsksWhatItIs(t *testing.T) {
	s, log := playWorkshop(t, []string{"A4_HELD"}, "ask cael about the ash", "image it")
	if !s.Has("ERRATUM_OBTAINED") || !s.Has("ERRATUM_UNUSED") {
		t.Fatal("four passes should end with the plate in your hands")
	}
	for _, id := range []string{"B21", "B22"} {
		if _, held := s.Held[id]; !held {
			t.Fatalf("expected %s", id)
		}
	}
	if !strings.Contains(log, "Mind the edges") {
		t.Fatalf("one of the murderers hands it over:\n%s", log)
	}
	if !strings.Contains(log, "Neither of them reads pre-Imperial") {
		t.Fatalf("the reason they live is that they cannot read it:\n%s", log)
	}
}

// If the ash was missed in Act I, Bez collected it. Vosk taught her to.
func TestTheAshIsNeverActuallyMissed(t *testing.T) {
	s, log := playWorkshop(t, nil, "ask cael about the ash", "image it")
	if !s.Has("ERRATUM_OBTAINED") {
		t.Fatalf("a player who never examined the censer still gets there:\n%s", log)
	}
	if !strings.Contains(log, "you take the ash and you don't ask") {
		t.Fatalf("expected the sergeant's provenance:\n%s", log)
	}
}

// Cutting the telemetry cost a fortnight of not trusting each other, and he
// names the price before he does the favour.
func TestDeceivingCaelHasAPrice(t *testing.T) {
	s, log := playWorkshop(t, []string{"A4_HELD", "CAEL_DECEIVED"}, "ask cael about the ash")
	if !s.Has("CAEL_SQUARED") {
		t.Fatalf("he asks, and he asks straight away:\n%s", log)
	}
	if !strings.Contains(log, "It has not sent me a burst since the day you took it") {
		t.Fatalf("he already knew:\n%s", log)
	}

	// And he does not ask a player who never lied to him.
	_, clean := playWorkshop(t, []string{"A4_HELD"}, "ask cael about the ash")
	if strings.Contains(clean, "the shroud.") {
		t.Fatalf("there is nothing to square:\n%s", clean)
	}
}

// The plate is unmissable. A player who never goes to the workshop gets it
// anyway, and loses only the signature on the recovery.
func TestTassBringsItUnbidden(t *testing.T) {
	s, log := playFrom(t, "standing_board", whole,
		"wait", "wait", "wait", "wait", "wait", "wait",
		"wait", "wait", "wait", "wait", "wait", "wait")
	if !s.Has("ERRATUM_OBTAINED") || !s.Has("TASS_BROUGHT_IT") {
		t.Fatalf("the apprentice comes up four tiers:\n%s", log)
	}
	if _, held := s.Held["B22"]; held {
		t.Fatal("nobody asked Cael to do it, so there is no log of him doing it")
	}
	if !strings.Contains(log, "Mind the edges") {
		t.Fatalf("Tass says what his master would have said:\n%s", log)
	}
}

// DECODE needs the chambers for the same reason QUERY needs a cogitator.
func TestDecodeNeedsTheRoomAndBothHalves(t *testing.T) {
	s, log := playFrom(t, "workshop", append([]string{"ERRATUM_OBTAINED"}, whole...), "decode")
	s.Grant("B21")
	if out, _ := s.Do("decode"); !strings.Contains(strings.Join(out, "\n"), "not the room it was made for") {
		t.Fatalf("the sentence is read out of furniture:\n%s\n%s", log, strings.Join(out, "\n"))
	}

	half, _ := playFrom(t, "standing_board", []string{"CIPHER_HALF_ROOM"}, "look")
	half.Grant("B21")
	out, _ := half.Do("decode")
	if !strings.Contains(strings.Join(out, "\n"), "It is a different object") {
		t.Fatalf("half a key reads nothing: %s", strings.Join(out, "\n"))
	}
	if half.Has("ACT2_COMPLETE") {
		t.Fatal("half a key must not finish the act")
	}
}

func TestTheSentence(t *testing.T) {
	s, _ := playFrom(t, "standing_board", whole, "look")
	s.Grant("B21")
	out, _ := s.Do("decode")
	log := strings.Join(out, "\n")

	if !s.Has("ACT2_COMPLETE") || !s.Has("KADMON_ADDRESS") {
		t.Fatal("the sentence is the act's turn")
	}
	for _, id := range []string{"B19", "B20"} {
		if _, held := s.Held[id]; !held {
			t.Fatalf("expected %s", id)
		}
	}
	for _, want := range []string{
		"THE INDEX IS THE CRIME.",
		"THE NINTH OBJECT IS THE WEAPON. IT IS IN THE ASH.",
		"ABOUT WHAT I AM MAKING YOU. — H.V.",
		"*making*, not *made*",
	} {
		if !strings.Contains(log, want) {
			t.Fatalf("missing: %q\n%s", want, log)
		}
	}

	// The canon rule: the message is not the method.
	if s.Has("ROOM_READ_FULL") {
		t.Fatal("decrypting the testimony must not grant the method")
	}
	if again, _ := s.Do("decode"); !strings.Contains(strings.Join(again, "\n"), "rest of your life") {
		t.Fatal("reading it again does not cost an afternoon each time")
	}
}

// The coerced half is complete, in the wrong order, and fluent. It holds
// exactly one scene, and what breaks it is counting.
func TestTheCorruptHalfNamesTheWrongMan(t *testing.T) {
	s, _ := playFrom(t, "standing_board", append([]string{"CIPHER_HALF_VAUX_CORRUPT"}, whole...), "look")
	s.Grant("B21")
	s.Add("VERITY", 3)

	first, _ := s.Do("decode")
	log := strings.Join(first, "\n")
	if !strings.Contains(log, "CARDINAL OSTREM VHAIL") {
		t.Fatalf("the false sentence names the man the player already wanted:\n%s", log)
	}
	if s.Has("ACT2_COMPLETE") {
		t.Fatal("a false reading must not turn the act")
	}
	// Nothing on the screen may hedge it. The player is supposed to believe it.
	for _, bad := range []string{"false", "wrong", "corrupt", "mis-order", "does not"} {
		if strings.Contains(strings.ToLower(log), bad) {
			t.Fatalf("the first reading told on itself: %s\n%s", bad, log)
		}
	}

	second, _ := s.Do("decode")
	log = strings.Join(second, "\n")
	if !s.Has("ACT2_COMPLETE") || !s.Has("SENTENCE_REPAIRED") {
		t.Fatalf("counting is what breaks it:\n%s", log)
	}
	if !strings.Contains(log, "Eight.") || !strings.Contains(log, "There are nine objects.") {
		t.Fatalf("nine objects, eight tokens:\n%s", log)
	}
	if s.Ints["VERITY"] != 2 {
		t.Fatalf("the repair costs a point of independent proof, got %d", s.Ints["VERITY"])
	}
}

// The game refuses to file the testimony, and the refusal is the player's.
func TestTheTestimonyCannotBeFiled(t *testing.T) {
	s, _ := playFrom(t, "standing_board", whole, "look")
	s.Grant("B21")
	s.Do("decode")

	out, _ := s.Do("file B19")
	log := strings.Join(out, "\n")
	if !strings.Contains(log, "You close the template.") {
		t.Fatalf("expected the refusal:\n%s", log)
	}
	if s.Has("B19_FILED") {
		t.Fatal("it must not reach the cogitator")
	}

	// And the plate itself goes into neither store, with no explanation.
	for _, verb := range []string{"file B21", "note B21"} {
		out, _ := s.Do(verb)
		log := strings.Join(out, "\n")
		if !strings.Contains(log, "You do not write down why.") {
			t.Fatalf("%s should refuse flatly:\n%s", verb, log)
		}
	}
	if s.Has("B21_FILED") || s.Has("B21_PHYSICAL") {
		t.Fatal("the Erratum has no representation anywhere")
	}
}

// The Vault answers a question about Cantica-Null twice, and only one of the
// answers is yours.
func TestTheVaultConfirmsWhatNobodyAsked(t *testing.T) {
	s, _ := playFrom(t, "standing_board", append([]string{"FENN_INTACT", "FENN_ALIVE"}, whole...), "look")
	s.Grant("B21")
	s.Do("decode")

	out, _ := s.Do("query CANTICA-NULL")
	if !strings.Contains(strings.Join(out, "\n"), "QUERY LOGGED.") {
		t.Fatalf("expected the null and the log line: %s", strings.Join(out, "\n"))
	}
	var echo []string
	for i := 0; i < 3; i++ {
		echo = append(echo, s.Tick()...)
	}
	log := strings.Join(echo, "\n")
	if !strings.Contains(log, "CONFIRMED: NO SUCH ROUTING ADDRESS.") {
		t.Fatalf("nobody confirms a nil return:\n%s", log)
	}
	if !strings.Contains(log, "\"Kadmon Veil.\"") {
		t.Fatalf("an intact savant says it from the next room:\n%s", log)
	}
}

func TestSharingTheTestimonyCostsNothingAndBuysThem(t *testing.T) {
	s, _ := playFrom(t, "standing_board", whole, "look")
	s.Grant("B21")
	s.Do("decode")
	s.Do("note B19")

	out, _ := s.Do("show B19 to ilaria")
	log := strings.Join(out, "\n")
	if !s.Has("TESTIMONY_SHARED") {
		t.Fatalf("expected the retinue to hear it:\n%s", log)
	}
	if s.Ints["ILARIA_TRUST"] != 1 {
		t.Fatalf("expected +1, got %d", s.Ints["ILARIA_TRUST"])
	}
	if !strings.Contains(log, "What's a tithe dispute") {
		t.Fatalf("Bez asks, and it is not a question about vocabulary:\n%s", log)
	}
}

// Every route into Act III carries the plate, and the plate is one use.
func TestTheErratumIsOneUse(t *testing.T) {
	s, _ := playFrom(t, "ferrum", []string{"ACT2_OPEN", "DECOY_PURSUED", "ERRATUM_OBTAINED", "ERRATUM_UNUSED"},
		"the erratum")
	if !s.Has("ERRATUM_SPENT") || s.Has("ERRATUM_UNUSED") {
		t.Fatal("using it on a person spends it")
	}
	// And nothing anywhere says what it just cost.
	_, log := playFrom(t, "ferrum", []string{"ACT2_OPEN", "DECOY_PURSUED", "ERRATUM_OBTAINED", "ERRATUM_UNUSED"},
		"the erratum")
	for _, bad := range []string{"ENDING", "foreclos", "no longer", "you have lost"} {
		if strings.Contains(log, bad) {
			t.Fatalf("the game warned the player: %s\n%s", bad, log)
		}
	}
}
