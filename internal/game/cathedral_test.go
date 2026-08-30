package game

import (
	"strings"
	"testing"
)

func playCathedral(t *testing.T, cmds ...string) (*State, string) {
	t.Helper()
	return playToPrecis(t)(append([]string{"sign", "go cathedral"}, cmds...)...)
}

// A wall that is the wrong thickness, in a plan that has never been amended,
// because an amendment is a variance and a variance leaves a name.
func TestTheWallsAreTheWrongThickness(t *testing.T) {
	_, log := playCathedral(t, "x confessionals", "x confessionals", "x plan", "x plan")
	if !strings.Contains(log, "Four hundred millimetres") {
		t.Fatal("the wall is measured, not asserted")
	}
	if !strings.Contains(log, "the two facts have never met") {
		t.Fatal("and nobody measures a wall")
	}
}

// Ruse is not corrupt. She was told it is the fire main and has believed it for
// thirty-one years, and shown the schematic she says it twice more and then
// asks a question the game has no answer to.
func TestRuseSaysItThreeTimes(t *testing.T) {
	s, log := playCathedral(t, "ask ruse about the crypt",
		"ask ivrit about the crypt", "ask ruse about the crypt")
	if strings.Count(log, "It's the fire main") != 2 {
		t.Fatalf("twice, and then the pause: %d", strings.Count(log, "It's the fire main"))
	}
	if !strings.Contains(log, "What do I do") {
		t.Fatal("and then she asks, and there is no good answer")
	}
	if !s.Has("RUSE_ASKED") {
		t.Fatal("expected RUSE_ASKED")
	}

	warned, wLog := playCathedral(t, "ask ruse about the crypt",
		"ask ivrit about the crypt", "ask ruse about the crypt", "warn ruse")
	if !warned.Has("RUSE_WARNED") {
		t.Fatal("warning her is available and costs you the keys")
	}
	if !strings.Contains(wLog, "you making sure I'm not in the way") {
		t.Fatal("and she names it correctly, and is not bitter")
	}
}

// Ivrit sells cheaply, without anguish, and his reason for choosing you is
// that you are the first person he can safely be useful to.
func TestIvritSellsCheaply(t *testing.T) {
	s, log := playCathedral(t, "ask ivrit about the crypt", "ask ivrit about his price")
	if _, held := s.Held["B8"]; !held {
		t.Fatal("the works schematic, off the bills")
	}
	if !s.Has("VHAIL_EXPOSED") || !s.Has("CATHEDRAL_MAPPED") {
		t.Fatal("expected VHAIL_EXPOSED and CATHEDRAL_MAPPED")
	}
	if !strings.Contains(log, "safely be useful to") {
		t.Fatal("you cannot have him killed, and that is the whole of your appeal")
	}
	if !strings.Contains(log, "a mother in Threnody") {
		t.Fatal("his price is one sentence in writing")
	}
}

// He does not move to admission. He goes from denial straight to
// justification, as though admitting were a formality for men who are ashamed.
func TestVhailJustifiesWithoutAdmitting(t *testing.T) {
	s, log := playCathedral(t,
		"ask ruse about the crypt", "ask ruse about the relief",
		"ask ivrit about the crypt", "go vhail",
		"ask vhail about the crypt", "ask vhail about the relief")
	if !s.Has("VHAIL_DENIED") {
		t.Fatal("ten minutes of denial first, and he is good at it")
	}
	if !strings.Contains(log, "arranges eleven true things in an order") {
		t.Fatal("nothing he says in the denial is untrue")
	}
	if !s.Has("VHAIL_BROKEN") {
		t.Fatal("the relief is a lever")
	}
	if _, held := s.Held["B9"]; !held {
		t.Fatal("and the confession is the clue")
	}
	if !strings.Contains(log, "I was the only man in this hive who was *counting*") {
		t.Fatal("ship the confession verbatim")
	}
	if s.Ints["RECONCILIATION"] < 4 {
		t.Fatalf("breaking him costs four: %d", s.Ints["RECONCILIATION"])
	}

	// Without a lever he stonewalls, warmly, and offers you the accounts again.
	stuck, sLog := playCathedral(t, "ask ivrit about the crypt", "go vhail",
		"ask vhail about the relief")
	if stuck.Has("VHAIL_BROKEN") {
		t.Fatal("nothing required is gated behind him, and he does not break for free")
	}
	if !strings.Contains(sLog, "exactly what you walked in with") {
		t.Fatal("and eleven minutes of it gets you nowhere")
	}
}

// The confession is taken inside the instrument, reaches his handlers while he
// is speaking, and nothing happens, because he is not the author.
func TestTheConfessionalIsTheTrap(t *testing.T) {
	_, log := playCathedral(t,
		"ask ruse about the crypt", "ask ruse about the relief",
		"ask ivrit about the crypt", "go vhail", "x wall", "x wall",
		"ask vhail about the relief", "ask vhail about this room",
		"ask vhail about where does it go")
	if !strings.Contains(log, "not marked out of use") {
		t.Fatal("a chamber in service has a line in it")
	}
	if !strings.Contains(log, "there is not one") {
		t.Fatal("you wait for a consequence and there is none")
	}
	if !strings.Contains(log, "I assumed that meant it was yours") {
		t.Fatal("the first true thing he does is be bewildered")
	}
}

// Everything he gives is readable off an unmarked requisition slip in twenty
// minutes. The query is ungated, available from the first hour, and nothing
// anywhere points at it.
func TestTheQuerySlipMakesHimRedundant(t *testing.T) {
	s, log := playToPrecis(t)("sign", "query cathedral alms")
	if _, held := s.Held["B10"]; !held {
		t.Fatal("second tier, unrestricted, first morning")
	}
	if !strings.Contains(log, "CN-∅ / PRIORITY-BLACK") {
		t.Fatal("the prefix is the thing")
	}
	if !strings.Contains(log, "0.4% UNDER ASSESSED") {
		t.Fatal("and the relief with it")
	}
	if !strings.Contains(log, "could have run on the first") {
		t.Fatal("the slip says so plainly, once, and never again")
	}
}

// Break and denounce, break and hold, or kill him. The trunk keeps running
// under all three, because stopping it was never in his gift.
func TestTheTrunkKeepsRunning(t *testing.T) {
	pre := []string{"ask ruse about the crypt", "ask ruse about the relief",
		"ask ivrit about the crypt", "go vhail", "ask vhail about the relief"}

	den, dLog := playCathedral(t, append(append([]string{}, pre...), "denounce")...)
	assertVhail(t, den, "VHAIL_DENOUNCED")
	if !strings.Contains(dLog, "supplies exactly the same data without any of the counting") {
		t.Fatal("replaced in a month by a quieter man")
	}
	if !strings.Contains(dLog, "never turned off") {
		t.Fatal("the trunk keeps running")
	}

	held, _ := playCathedral(t, append(append([]string{}, pre...), "hold")...)
	assertVhail(t, held, "VHAIL_LEVERAGED")

	dead, kLog := playCathedral(t, append(append([]string{}, pre...), "kill")...)
	assertVhail(t, dead, "VHAIL_DEAD")
	if dead.Ints["RECONCILIATION"] < 16 { // 4 to break, 12 to kill
		t.Fatalf("killing him is twelve: %d", dead.Ints["RECONCILIATION"])
	}
	if !strings.Contains(kLog, "it is satisfying") {
		t.Fatal("and the game says so plainly")
	}
	if !strings.Contains(kLog, "It runs that evening") {
		t.Fatal("and the trunk keeps running anyway")
	}
}

func assertVhail(t *testing.T, s *State, want string) {
	t.Helper()
	var set []string
	for _, f := range []string{"VHAIL_DENOUNCED", "VHAIL_LEVERAGED", "VHAIL_DEAD"} {
		if s.Has(f) {
			set = append(set, f)
		}
	}
	if len(set) != 1 || set[0] != want {
		t.Fatalf("expected exactly %s, got %v", want, set)
	}
}

// Crase asks at two threads, having been useful for a fortnight, and the
// dossier is real. Refusing costs two points and starts the surveillance.
func TestTheGalleryBriefIsRealAndTimed(t *testing.T) {
	early, _ := playToPrecis(t)("sign", "go crase")
	if early.Room.ID == "crase_office" {
		t.Fatal("she asks when there is something to ask about")
	}

	s, log := playToPrecis(t)("sign", "query cathedral alms",
		"go fenn", "ask fenn about kadmon", "back", "wait", "go crase",
		"x dossier", "x dossier", "ask crase about the evidence")
	if !s.Has("CRASE_ASKED_FOR_YOU") {
		t.Fatal("two threads touched")
	}
	if !strings.Contains(log, "Nothing is wrong. — Y.C.") {
		t.Fatal("a runner, not a vox")
	}
	if !s.Has("GALLERY_VERIFIED") {
		t.Fatal("she wants you to satisfy yourself of it rather than take it from her")
	}
	if !strings.Contains(log, "Everything in the folder is true") {
		t.Fatal("the evidence is real, thorough and correct")
	}

	refused, rLog := playToPrecis(t)("sign", "query cathedral alms",
		"go fenn", "ask fenn about kadmon", "back", "wait", "go crase", "refuse")
	if !refused.Has("DECOY_REFUSED") || !refused.Has("CRASE_SURVEILLED") {
		t.Fatal("refusing starts the surveillance")
	}
	if !refused.Has("SUBLEVEL_NINE_SEEDED") {
		t.Fatal("and seeds the string that pays off at Reconciliation 40")
	}
	if !strings.Contains(rLog, "There are eight. You have counted them") {
		t.Fatal("there is no sub-level nine")
	}
	if !strings.Contains(rLog, "asks whether you have eaten") {
		t.Fatal("she stays warm throughout, and no dialogue is attached to the refusal")
	}

	taken, tLog := playToPrecis(t)("sign", "query cathedral alms",
		"go fenn", "ask fenn about kadmon", "back", "wait", "go crase", "accept")
	if !taken.Has("DECOY_PURSUED") || taken.Ints["CRASE_STANDING"] < 4 {
		t.Fatalf("accepting is worth three: %d", taken.Ints["CRASE_STANDING"])
	}
	if !strings.Contains(tLog, "day forty was eleven days ago") {
		t.Fatal("fourteen days, and the clock does not wait")
	}
}
