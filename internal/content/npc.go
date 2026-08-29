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

// Gate is every condition the content can put on a line of dialogue. It is
// embedded in topics, acts, shows and beats, so all four gate the same way and
// there is one place to add a condition.
type Gate struct {
	Requires []string `json:"requires"`
	Blocks   []string `json:"blocks"`
	// Has gates on a clue being in the player's hands at all.
	Has []string `json:"has"`
	// AtLeast and Below gate on counters — a savant three points into drift is
	// a different man, and the content should be able to say so.
	AtLeast Counters `json:"at_least"`
	Below   Counters `json:"below"`
}

// Topic is one thing an NPC can be asked about.
//
// Requires/Blocks read flags; Holds reads the player's evidence. A topic whose
// gates are unmet is not hidden and does not error — it produces Closed, which
// is a real line of dialogue and usually costs the NPC something to say.
type Topic struct {
	Gate
	ID    string   `json:"id"`
	Names []string `json:"names"`
	// Holds gates on a clue being recorded somewhere, where Has only asks
	// whether it is in hand. Citing a register you read is knowledge;
	// producing one is a document. They are not the same move.
	Holds []string `json:"holds"`
	// Never marks a subject a person will not be moved on at all. It is not a
	// gate the player can pick; it is a door, and pushing on it is the point.
	Never bool `json:"never"`
	// HiddenAbove is a Reconciliation threshold past which this topic does not
	// exist. Not closed — absent: the name stops matching at all, so the
	// person answers as though they had never heard of it and there is no gap
	// to notice. Zero means never hidden.
	HiddenAbove int    `json:"hidden_above"`
	Text        string `json:"text"`
	Closed      string `json:"closed"`
	// ClosedCounters is what the refusal costs, kept separate from Counters
	// because a topic usually rewards being answered and charges for being
	// pushed on, and those are not the same number.
	ClosedCounters Counters `json:"closed_counters"`
	Repeat         string   `json:"repeat"`
	Sets           []string `json:"sets"`
	Clears         []string `json:"clears"`
	Grants         []string `json:"grants"`
	Pressure       int      `json:"pressure"`
	Counters       Counters `json:"counters"`
}

// Show is what an NPC does when the player puts a held clue in front of him.
type Show struct {
	Gate
	Clue     string   `json:"clue"`
	Refused  string   `json:"refused"` // what they say when they have already had this
	Text     string   `json:"text"`
	Sets     []string `json:"sets"`
	Clears   []string `json:"clears"`
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
	"BEZ_TRUST": true, "ILARIA_TRUST": true, "FENN_DRIFT": true, "MOTT_FAVOURS": true, "QUERY_TIER": true,
	"VAUX_SUSPICION": true,
}

// Act is a verb the player uses on a person rather than a thing: detain,
// dismiss, threaten. Acts are one-way. There is no undo and no confirmation.
type Act struct {
	Gate
	ID       string   `json:"id"`
	Names    []string `json:"names"`
	Refused  string   `json:"refused"` // what they say when it is not yours to do yet
	Text     string   `json:"text"`
	Sets     []string `json:"sets"`
	Clears   []string `json:"clears"`
	Grants   []string `json:"grants"`
	Pressure int      `json:"pressure"`
	Counters Counters `json:"counters"`
	Cost     int      `json:"cost"`
}

// Beat fires on its own once a conversation has gone far enough. Vaux's two
// plants are beats: the player never asks for them, and cannot miss them by
// asking about the wrong thing. Whether they write them down is another matter.
type Beat struct {
	Gate
	ID            string   `json:"id"`
	AfterTopics   int      `json:"after_topics"`
	AfterPressure int      `json:"after_pressure"`
	Text          string   `json:"text"`
	Sets          []string `json:"sets"`
	Clears        []string `json:"clears"`
	Grants        []string `json:"grants"`
	// Counters is how a beat charges for something the player did before the
	// conversation started. Vaux is already two points suspicious of anybody
	// who asked to see him through the Crown's own message system, and the
	// beat that says so is the only place that fact exists.
	Counters Counters `json:"counters"`
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
	Gate
	ID    string   `json:"id"`
	Names []string `json:"names"`
	To    string   `json:"to"`
	// Sets are flags taking this door writes. Which door you go down into the
	// coolant galleries by decides who does not come back up.
	Sets []string `json:"sets"`
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

// MatchActs resolves a player-typed act to every act that answers to it, best
// match first and otherwise in authored order — the same rule topics use.
// Several acts may share a name on purpose: the death in the coolant galleries
// is one verb with four outcomes, and which one happens depends on who the
// player brought.
func (n *NPC) MatchActs(word string) []*Act {
	type hit struct {
		a *Act
		s int
	}
	var hits []hit
	for i := range n.Acts {
		if s := match(word, n.Acts[i].ID, n.Acts[i].Names); s > 0 {
			hits = append(hits, hit{&n.Acts[i], s})
		}
	}
	sort.SliceStable(hits, func(a, b int) bool { return hits[a].s > hits[b].s })
	out := make([]*Act, 0, len(hits))
	for _, h := range hits {
		out = append(out, h.a)
	}
	return out
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
