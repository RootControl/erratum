package content

import (
	"sort"
	"strings"
)

// NPC is a person who can be asked things.
//
// Dialogue is free-form for the same reason examination is: nothing in this
// game hands the player a list of the things worth caring about. Topics are
// found by naming them, with the same loose matching objects use, so a
// conversation is never lost to vocabulary either.
type NPC struct {
	ID       string   `json:"id"`
	Names    []string `json:"names"`
	Title    string   `json:"title"`
	Portrait string   `json:"portrait"` // what EXAMINE gives you
	Greeting string   `json:"greeting"`
	Reprise  string   `json:"reprise"` // second and later approach
	Deflect  []string `json:"deflect"` // cycled when he has nothing he'll say
	Unknown  string   `json:"unknown"` // the player named something he has no line for

	Topics []Topic `json:"topics"`
	Shows  []Show  `json:"shows"`
	Acts   []Act   `json:"acts"`
	Beats  []Beat  `json:"beats"`

	Disposition *Disposition `json:"disposition"`

	// ListsTopics is false for anyone the player is supposed to have to think
	// about. It exists because 1.10 needs a listed conversation whose list can
	// quietly lose a line at RECONCILIATION 10, and that only works if some
	// conversations are listed and some are not.
	ListsTopics  bool   `json:"lists_topics"`
	TopicRefusal string `json:"topic_refusal"`

	key string
}

// Key identifies one authored conversation. Crase in 1.1 and Crase in 1.2 are
// the same woman and two different conversations — she greets you twice, once
// in a doorway and once with a rosette in a handkerchief — so conversation
// state is keyed by where the person was authored, not by their name. Fenn
// followed into the gallery arch is one definition seen from two rooms, and
// keeps one key and one conversation.
func (n *NPC) Key() string { return n.key }

// Topic is one thing an NPC can be asked about.
//
// Requires/Blocks read flags; Holds reads the player's evidence. A topic whose
// gates are unmet is not hidden and does not error — it produces Closed, which
// is a real line of dialogue and usually costs the NPC something to say.
type Topic struct {
	ID       string   `json:"id"`
	Names    []string `json:"names"`
	Requires []string `json:"requires"`
	Blocks   []string `json:"blocks"`
	// Has gates on a clue being in the player's hands at all; Holds gates on
	// its being recorded somewhere. Citing a register you read is knowledge.
	// Producing one is a document. They are not the same move.
	Has   []string `json:"has"`
	Holds []string `json:"holds"`
	// Never marks a subject a person will not be moved on at all. It is not a
	// gate the player can pick; it is a door, and pushing on it is the point.
	Never  bool   `json:"never"`
	Text   string `json:"text"`
	Closed string `json:"closed"`
	// ClosedCounters is what the refusal costs, kept separate from Counters
	// because a topic usually rewards being answered and charges for being
	// pushed on, and those are not the same number.
	ClosedCounters Counters `json:"closed_counters"`
	Repeat         string   `json:"repeat"`
	Sets           []string `json:"sets"`
	Grants         []string `json:"grants"`
	Pressure       int      `json:"pressure"`
	Counters       Counters `json:"counters"`
}

// Show is what an NPC does when the player puts a held clue in front of him.
type Show struct {
	Clue     string   `json:"clue"`
	Requires []string `json:"requires"`
	Blocks   []string `json:"blocks"`
	Has      []string `json:"has"`
	Refused  string   `json:"refused"` // what they say when they have already had this
	Text     string   `json:"text"`
	Sets     []string `json:"sets"`
	Grants   []string `json:"grants"`
	Pressure int      `json:"pressure"`
	Counters Counters `json:"counters"`
}

// Counters are the standing and trust movements a line of dialogue causes,
// keyed by the same names the act documents use.
type Counters map[string]int

// KnownCounters is every counter the game keeps. A line of dialogue that moves
// anything else has a typo in it, and a typo here is an award that silently
// never lands.
var KnownCounters = map[string]bool{
	"VERITY": true, "CRASE_STANDING": true, "RECONCILIATION": true,
	"BEZ_TRUST": true, "ILARIA_TRUST": true,
}

// Act is a verb the player uses on a person rather than a thing: detain,
// dismiss, threaten. Acts are one-way. There is no undo and no confirmation.
type Act struct {
	ID       string   `json:"id"`
	Names    []string `json:"names"`
	Requires []string `json:"requires"`
	Blocks   []string `json:"blocks"`
	Has      []string `json:"has"`     // clues that must be in hand
	Refused  string   `json:"refused"` // what they say when it is not yours to do yet
	Text     string   `json:"text"`
	Sets     []string `json:"sets"`
	Grants   []string `json:"grants"`
	Pressure int      `json:"pressure"`
	Counters Counters `json:"counters"`
	Cost     int      `json:"cost"`
}

// Beat fires on its own once a conversation has gone far enough. Vaux's two
// plants are beats: the player never asks for them, and cannot miss them by
// asking about the wrong thing. Whether they write them down is another matter.
type Beat struct {
	ID            string   `json:"id"`
	AfterTopics   int      `json:"after_topics"`
	AfterPressure int      `json:"after_pressure"`
	Requires      []string `json:"requires"`
	Blocks        []string `json:"blocks"`
	Has           []string `json:"has"`
	Text          string   `json:"text"`
	Sets          []string `json:"sets"`
	Grants        []string `json:"grants"`
}

// Disposition is a one-of-N flag group. The act documents write these as an
// enum; the engine writes them as the same flat SCREAMING_SNAKE bools as
// everything else. This struct is the only thing that knows they are mutually
// exclusive, and Default is what the record says about you if you did nothing
// — empty where saying nothing is a real third answer, as it is with Crase's
// warning about Kesh.
type Disposition struct {
	Flags   []string `json:"flags"`
	Default string   `json:"default"`
}

// Exit is somewhere else. Cost is in timed actions: the clock does not stop
// because you are in a lift.
type Exit struct {
	ID       string   `json:"id"`
	Names    []string `json:"names"`
	To       string   `json:"to"`
	Requires []string `json:"requires"`
	// RequiresAny and HasAny are one disjunction between them: when either is
	// non-empty, the exit needs at least one item across both. The Choir Spire
	// opens on A2 or on Crase's lead, and either is enough.
	RequiresAny []string `json:"requires_any"`
	HasAny      []string `json:"has_any"`
	Text        string   `json:"text"`
	Blocked     string   `json:"blocked"`
	Cost        int      `json:"cost"`
}

// match is the game's one matching rule, used by everything the player can
// name: exact id, then exact alias, then a prefix of at least three
// characters. 2 is exact, 1 is a prefix, 0 is no.
func match(word, id string, names []string) int {
	word = trimArticle(word)
	if word == "" {
		return 0
	}
	if strings.ToLower(id) == word {
		return 2
	}
	for _, n := range names {
		// Articles are trimmed from both sides, so an alias authored as "the
		// form" answers to "form" and vice versa.
		if strings.ToLower(n) == word || trimArticle(n) == word {
			return 2
		}
	}
	if len(word) < 3 {
		return 0
	}
	for _, n := range append([]string{id}, names...) {
		n = trimArticle(n)
		// A prefix, so "regi" finds the regicide board; or a suffix, so
		// "being passed over" finds "passed over" without anyone having to
		// author every way an English sentence can start.
		if strings.HasPrefix(n, word) || strings.HasSuffix(word, n) && len(n) >= 3 {
			return 1
		}
	}
	return 0
}

// trimArticle drops the words people put in front of nouns. "ask about the
// seal" and "ask about seal" are the same question and the game should not be
// the kind of thing that has an opinion about which one you typed.
func trimArticle(word string) string {
	word = strings.ToLower(strings.TrimSpace(word))
	for _, a := range []string{"the ", "a ", "an ", "his ", "her ", "their ", "your ", "my "} {
		if strings.HasPrefix(word, a) {
			return strings.TrimSpace(word[len(a):])
		}
	}
	return word
}

// FindNPC resolves a player-typed name to someone in the room.
func (r *Room) FindNPC(word string) *NPC {
	var best *NPC
	bestScore := 0
	for _, n := range r.people {
		if s := match(word, n.ID, n.Names); s > bestScore {
			best, bestScore = n, s
		}
	}
	return best
}

// MatchTopics resolves a player-typed subject to every topic that answers to
// it, best match first and otherwise in authored order. More than one topic
// may share a name on purpose: the engine takes the first whose gates the
// player has actually met, so "the seal" is a different conversation with a
// player holding the integrity register than without one.
func (n *NPC) MatchTopics(word string) []*Topic {
	type hit struct {
		t *Topic
		s int
	}
	var hits []hit
	for i := range n.Topics {
		if s := match(word, n.Topics[i].ID, n.Topics[i].Names); s > 0 {
			hits = append(hits, hit{&n.Topics[i], s})
		}
	}
	sort.SliceStable(hits, func(a, b int) bool { return hits[a].s > hits[b].s })
	out := make([]*Topic, 0, len(hits))
	for _, h := range hits {
		out = append(out, h.t)
	}
	return out
}

// FindAct resolves DETAIN and its neighbours.
func (n *NPC) FindAct(word string) *Act {
	var best *Act
	bestScore := 0
	for i := range n.Acts {
		if s := match(word, n.Acts[i].ID, n.Acts[i].Names); s > bestScore {
			best, bestScore = &n.Acts[i], s
		}
	}
	return best
}

// ShowFor returns what he does with a particular clue, if anything.
func (n *NPC) ShowFor(clue string) *Show {
	for i := range n.Shows {
		if strings.EqualFold(n.Shows[i].Clue, clue) {
			return &n.Shows[i]
		}
	}
	return nil
}

// FindExit resolves a player-typed destination.
func (r *Room) FindExit(word string) *Exit {
	var best *Exit
	bestScore := 0
	for i := range r.Exits {
		if s := match(word, r.Exits[i].ID, r.Exits[i].Names); s > bestScore {
			best, bestScore = &r.Exits[i], s
		}
	}
	return best
}
