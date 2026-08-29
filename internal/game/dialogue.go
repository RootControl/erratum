package game

import (
	"strings"

	"github.com/RootControl/erratum/internal/content"
)

// Conversation works the way examination works, and for the same reason: the
// game never shows the player a list of the things worth asking about. You
// name a subject, the loose matcher finds the nearest thing the person has a
// position on, and what you get back depends on what you are already carrying.
//
// Nothing here signals that a topic was gated, refused or unlocked. A closed
// subject produces a real line of dialogue, not a rejection.

// npcIn resolves the person the player means. With one person in the room,
// naming them is optional — but naming somebody is never an invitation to
// substitute whoever happens to be standing there instead.
func (s *State) npcIn(word string) *content.NPC {
	if strings.TrimSpace(word) != "" {
		return s.Room.FindNPC(word)
	}
	if p := s.Room.People(); len(p) == 1 {
		return p[0]
	}
	return nil
}

// whyNot explains a person the player named who is not going to answer:
// somebody who is not in this room, or somebody who is — Bez against a
// doorframe in 1.1 — and is not taking questions in it.
func (s *State) whyNot(word string) string {
	if strings.TrimSpace(word) == "" {
		return s.S("no_one_here")
	}
	if s.Room.Find(word) != nil {
		return s.S("not_talking")
	}
	return s.S("not_here")
}

// meet renders the approach the first time and nothing afterwards. A person
// does not re-introduce themselves between two questions.
func (s *State) meet(n *content.NPC) []string {
	if s.Met[n.Key()] {
		return nil
	}
	s.Met[n.Key()] = true
	s.Set(strings.ToUpper(n.ID) + "_MET")
	s.Set(strings.ToUpper(n.ID) + "_ALIVE")
	if n.Greeting != "" {
		return []string{n.Greeting}
	}
	return nil
}

func (s *State) talk(arg string) ([]string, int) {
	n := s.npcIn(arg)
	if n == nil {
		return []string{s.whyNot(arg)}, 0
	}
	out := s.meet(n)
	// A beat waiting on its gates rather than on the conversation fires the
	// moment they open, and walking back up to somebody is one of the ways
	// they open. Without this it can only land on a question the player has
	// not already asked, which is not what "the moment" means.
	if beats := s.beats(n); len(out) > 0 || len(beats) > 0 {
		return append(out, beats...), 1
	}
	if n.Reprise != "" {
		return []string{n.Reprise}, 0
	}
	return []string{"You have nothing further to open with."}, 0
}

// ask is the whole verb. ASK VAUX ABOUT THE SEAL, ASK ABOUT THE SEAL, ASK SEAL.
func (s *State) ask(arg string) ([]string, int) {
	if len(s.Room.People()) == 0 {
		return []string{s.S("no_one_here")}, 0
	}
	who, subject := splitAsk(arg)
	if who == "" {
		// "ASK VAUX THE SEAL" — no "about", but the first word is a person.
		if head, rest, ok := strings.Cut(subject, " "); ok && s.Room.FindNPC(head) != nil {
			who, subject = head, strings.TrimSpace(rest)
		}
	}
	n := s.npcIn(who)
	if n == nil {
		return []string{s.whyNot(who)}, 0
	}
	if subject == "" {
		return s.talk(who)
	}

	var out []string
	out = append(out, s.meet(n)...)

	matches := n.MatchTopics(subject)
	// Reconciliation does not close subjects, it removes them. A topic past
	// its threshold stops matching at all, so what comes back is the ordinary
	// "I have nothing on that" and there is no gap to notice.
	kept := matches[:0]
	for _, t := range matches {
		if t.HiddenAbove == 0 || s.Reconciliation() < t.HiddenAbove {
			kept = append(kept, t)
		}
	}
	matches = kept
	if len(matches) == 0 {
		if n.Unknown != "" {
			return append(out, n.Unknown), 1
		}
		return append(out, s.S("no_such_topic")), 0
	}

	var chosen *content.Topic
	for _, t := range matches {
		if s.topicOpen(t) {
			chosen = t
			break
		}
	}
	if chosen == nil {
		// Nothing is open, so the answer is whichever variant actually has a
		// closed line written for it. Several topics share a name on purpose
		// and only one of them is the one he says no with.
		chosen = matches[0]
		for _, t := range matches {
			if t.Closed != "" {
				chosen = t
				break
			}
		}
	}

	// A subject he refused an hour ago is not spent. Accuracy opens him, and
	// it has to be able to open him retroactively or the mechanic is a lock
	// the player can turn the wrong way once and never turn again.
	if !s.topicOpen(chosen) {
		if s.Fired[n.Key()+"/closed/"+chosen.ID] {
			return append(out, s.deflect(n)), 0
		}
		s.Fired[n.Key()+"/closed/"+chosen.ID] = true
		s.Pressure[n.Key()] += chosen.Pressure
		for k, v := range chosen.ClosedCounters {
			s.Add(k, v)
		}
		if chosen.Closed != "" {
			out = append(out, chosen.Closed)
		} else {
			out = append(out, s.deflect(n))
		}
		return append(out, s.beats(n)...), 1
	}

	if s.Topics[n.Key()+"/"+chosen.ID] {
		if chosen.Repeat != "" {
			return append(out, chosen.Repeat), 0
		}
		return append(out, s.deflect(n)), 0
	}
	s.Topics[n.Key()+"/"+chosen.ID] = true
	s.Asked[n.Key()]++
	out = append(out, chosen.Text)
	// Pressure is what pushing on a closed door costs a person. An answered
	// question costs them nothing, however pointed it was.
	out = append(out, s.apply(n, chosen.Sets, chosen.Clears, chosen.Grants, chosen.Counters, 0)...)
	return append(out, s.beats(n)...), 1
}

// gatesMet reads every condition the content can put on a line of dialogue.
func (s *State) gatesMet(g content.Gate) bool {
	for _, f := range g.Requires {
		if !s.Has(f) {
			return false
		}
	}
	for _, f := range g.Blocks {
		if s.Has(f) {
			return false
		}
	}
	for _, id := range g.Has {
		if _, held := s.Held[id]; !held {
			return false
		}
	}
	for k, v := range g.AtLeast {
		if s.Ints[k] < v {
			return false
		}
	}
	for k, v := range g.Below {
		if s.Ints[k] >= v {
			return false
		}
	}
	return true
}

// topicOpen reports whether a person will take this subject from this player.
func (s *State) topicOpen(t *content.Topic) bool {
	if t.Never {
		return false
	}
	if !s.gatesMet(t.Gate) {
		return false
	}
	for _, id := range t.Holds {
		if !s.Holds(id) {
			return false
		}
	}
	return true
}

// show puts a held clue in front of someone. What they do with it is authored
// per clue; most people, for most evidence, do nothing at all.
func (s *State) show(arg string) ([]string, int) {
	if len(s.Room.People()) == 0 {
		return []string{s.S("no_one_here")}, 0
	}
	what, who := arg, ""
	if i := strings.LastIndex(strings.ToLower(arg), " to "); i >= 0 {
		what, who = arg[:i], arg[i+4:]
	}
	n := s.npcIn(strings.TrimSpace(who))
	if n == nil {
		return []string{s.whyNot(who)}, 0
	}
	var c *content.Clue
	for _, m := range s.Game.MatchClues(strings.TrimSpace(what)) {
		if _, held := s.Held[m.ID]; held {
			c = m
			break
		}
	}
	if c == nil {
		return []string{"You hold no such clue."}, 0
	}

	out := s.meet(n)
	sh := n.ShowFor(c.ID)
	if sh == nil {
		return append(out, "You put it in front of them. They look at it, and it is not the thing they are afraid of."), 1
	}
	if !s.gatesMet(sh.Gate) {
		if sh.Refused != "" {
			return append(out, sh.Refused), 0
		}
		return append(out, s.deflect(n)), 0
	}
	if s.Fired[n.Key()+"/show/"+c.ID] {
		return append(out, s.deflect(n)), 0
	}
	s.Fired[n.Key()+"/show/"+c.ID] = true
	out = append(out, sh.Text)
	out = append(out, s.apply(n, sh.Sets, sh.Clears, sh.Grants, sh.Counters, sh.Pressure)...)
	return append(out, s.beats(n)...), 1
}

// act runs DETAIN, READ ROSTER and their neighbours: one-way, unconfirmed,
// and often expensive.
//
// The person is usually implied. With two people in a room only one of them is
// holding out a duty-roster, so if the player did not name anybody, the act
// itself names them.
func (s *State) act(verb, arg string) ([]string, int, bool) {
	if len(s.Room.People()) == 0 {
		return nil, 0, false
	}
	full := strings.TrimSpace(verb + " " + arg)
	if n := s.npcIn(arg); n != nil {
		for _, word := range []string{full, verb} {
			if a := s.pickAct(n, word); a != nil {
				return s.doAct(n, a)
			}
		}
	}
	for _, word := range []string{full, verb} {
		for _, n := range s.Room.People() {
			if a := s.pickAct(n, word); a != nil {
				return s.doAct(n, a)
			}
		}
	}
	return nil, 0, false
}

// pickAct takes the first act answering to the word whose gates the player has
// met, falling back to the best match so that a refusal still has a voice.
func (s *State) pickAct(n *content.NPC, word string) *content.Act {
	matches := n.MatchActs(word)
	if len(matches) == 0 {
		return nil
	}
	for _, a := range matches {
		if s.gatesMet(a.Gate) {
			return a
		}
	}
	for _, a := range matches {
		if a.Refused != "" {
			return a
		}
	}
	return matches[0]
}

func (s *State) doAct(n *content.NPC, a *content.Act) ([]string, int, bool) {
	if !s.gatesMet(a.Gate) {
		if a.Refused != "" {
			return []string{a.Refused}, 0, true
		}
		return []string{s.deflect(n)}, 0, true
	}
	if s.Fired[n.Key()+"/act/"+a.ID] {
		return []string{s.deflect(n)}, 0, true
	}
	s.Fired[n.Key()+"/act/"+a.ID] = true
	out := s.meet(n)
	out = append(out, a.Text)
	out = append(out, s.apply(n, a.Sets, a.Clears, a.Grants, a.Counters, a.Pressure)...)
	cost := a.Cost
	if cost < 1 {
		cost = 1
	}
	return append(out, s.beats(n)...), cost, true
}

// apply writes the consequences of a line of dialogue. Flags belonging to a
// person's disposition group displace their siblings; everything else is an
// ordinary flag.
func (s *State) apply(n *content.NPC, sets, clears, grants []string, counters content.Counters, pressure int) []string {
	for _, f := range sets {
		if n.Disposition != nil && contains(n.Disposition.Flags, f) {
			s.SetOne(n.Disposition.Flags, f)
			continue
		}
		s.Set(f)
	}
	// A line of dialogue that ends a man has to be able to say so.
	for _, f := range clears {
		delete(s.Flags, f)
	}
	for k, v := range counters {
		s.Add(k, v)
	}
	s.Pressure[n.Key()] += pressure
	for _, id := range grants {
		s.Grant(id)
	}
	return nil
}

// beats are what a person gives up without being asked. Vaux's two plants live
// here: the player cannot ask for them and cannot miss them by asking about
// the wrong thing. Whether they write them down is another matter entirely.
func (s *State) beats(n *content.NPC) []string {
	var out []string
	for i := range n.Beats {
		b := &n.Beats[i]
		if s.Fired[n.Key()+"/beat/"+b.ID] {
			continue
		}
		ready := false
		if b.AfterTopics > 0 && s.Asked[n.Key()] >= b.AfterTopics {
			ready = true
		}
		if b.AfterPressure > 0 && s.Pressure[n.Key()] >= b.AfterPressure {
			ready = true
		}
		// A beat with no counter at all is waiting on its gates rather than on
		// the conversation: it fires the moment they open.
		if b.AfterTopics == 0 && b.AfterPressure == 0 && len(b.Requires) > 0 {
			ready = true
		}
		if !ready || !s.gatesMet(b.Gate) {
			continue
		}
		s.Fired[n.Key()+"/beat/"+b.ID] = true
		out = append(out, b.Text)
		for k, v := range b.Counters {
			s.Add(k, v)
		}
		for _, f := range b.Sets {
			s.Set(f)
		}
		for _, f := range b.Clears {
			delete(s.Flags, f)
		}
		// Granted silently. A clue the player did not ask for does not get a
		// bracket announcing that it matters; that is the entire test.
		for _, id := range b.Grants {
			s.Grant(id)
		}
	}
	return out
}

func (s *State) deflect(n *content.NPC) string {
	if len(n.Deflect) == 0 {
		return "Nothing."
	}
	i := s.deflects[n.Key()] % len(n.Deflect)
	s.deflects[n.Key()]++
	return n.Deflect[i]
}

// splitAsk pulls "vaux about the seal" apart into a person and a subject.
func splitAsk(arg string) (who, subject string) {
	low := strings.ToLower(arg)
	if i := strings.Index(low, " about "); i >= 0 {
		return strings.TrimSpace(arg[:i]), strings.TrimSpace(arg[i+7:])
	}
	if strings.HasPrefix(low, "about ") {
		return "", strings.TrimSpace(arg[6:])
	}
	return "", strings.TrimSpace(arg)
}

func contains(hay []string, needle string) bool {
	for _, h := range hay {
		if h == needle {
			return true
		}
	}
	return false
}
