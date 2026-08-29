package game

import (
	"fmt"
	"sort"
	"strings"

	"github.com/RootControl/erratum/internal/content"
)

type Store int

const (
	StoreNone Store = iota
	StoreCog
	StoreBook
	StoreBoth
)

func (s Store) String() string {
	switch s {
	case StoreCog:
		return "cogitator"
	case StoreBook:
		return "notebook"
	case StoreBoth:
		return "both"
	}
	return "in hand"
}

const (
	secondsPerAction = 40
	bodyLimit        = 8 * 60 // eight minutes, and they mean it
	correctionDelay  = 5      // actions after the requisition before the archive reconciles
	fallbackDelay    = 9      // if the player never looks again, Crase delivers it
	precisDelay      = 14     // and then the Crown wants the day signed off
	precisFloor      = 50     // a player who never requisitions still reaches the end of the day
	lyleConfirmDelay = 6      // and then the Vault confirms a null result you never asked it to
)

// State is everything the game knows. Flags are strings by design: the act
// documents are written in SCREAMING_SNAKE and the engine should speak the
// same language the writers do.
type State struct {
	Game *content.Game
	Room *content.Room // the scene you are standing in

	// Clock is hive time, in seconds since midnight. It runs forward through a
	// scene and jumps when the player enters the next one.
	Clock int

	// The body has a deadline, not a stopwatch: eight minutes of hive time from
	// the moment the player first reaches the chambers. Walking out does not
	// pause it, and a scene that jumps the clock an hour and a half forward
	// takes him with it.
	BodyPresent  bool
	bodyDeadline int
	Visited      map[string]bool
	date         string   // the date the clock is currently running on
	driftPartial int      // drift accrued under the shroud but not yet whole
	arrived      []string // what the scene handed you on the way in

	Held     map[string]Store
	Examined map[string]int
	Flags    map[string]bool
	Ints     map[string]int
	Carried  map[string]bool

	Notebook []string // append-only. Nothing in the game can rewrite this.
	Archive  *Archive

	// Conversation state, keyed by NPC ID. Pressure is how hard the player
	// has pushed on the things a person will not discuss; it is what buys
	// the plants, and it is never shown.
	Met      map[string]bool
	Asked    map[string]int
	Topics   map[string]bool
	Pressure map[string]int
	Fired    map[string]bool
	deflects map[string]int

	// The Case Précis is the hard floor of Act I: unskippable, undismissable,
	// and four paragraphs of true sentences describing a day that did not
	// happen. While it is pending the player may read their own coat and
	// nothing else.
	PrecisPending bool
	precisDone    bool
	actions       int

	// The Lyle query gets an answer nobody asked for, forty minutes later.
	lyleQueried   int
	lyleConfirmed bool

	requisitioned  bool
	sinceReq       int
	fallbackFired  bool
	seenCorrection bool
	Done           bool
}

func New(g *content.Game) *State {
	s := &State{
		Game:        g,
		Room:        g.Scene(g.Start),
		BodyPresent: true,
		Held:        map[string]Store{},
		Examined:    map[string]int{},
		Flags:       map[string]bool{},
		// CRASE_STANDING starts at +1. She liked him, she liked you, and Act I
		// is built to push an attentive player to +2 or +3 from there.
		Ints:     map[string]int{"VERITY": 0, "CRASE_STANDING": 1, "RECONCILIATION": 0},
		Carried:  map[string]bool{},
		Archive:  NewArchive(),
		Visited:  map[string]bool{},
		Met:      map[string]bool{},
		Asked:    map[string]int{},
		Topics:   map[string]bool{},
		Pressure: map[string]int{},
		Fired:    map[string]bool{},
		deflects: map[string]int{},
	}
	s.date = s.Room.Date
	s.arrived = s.arrive()
	return s
}

// arrive runs the bookkeeping of entering a scene: the clock jump the act
// documents imply, the body timer starting the first time the player is close
// enough to Vosk to read anything off him, and whatever the scene hands you
// for being in it.
func (s *State) arrive() []string {
	var out []string
	if !s.Visited[s.Room.ID] {
		s.Visited[s.Room.ID] = true
		// A scene's date is the earliest it can happen on, not a fact about
		// where the player is. Later advances the calendar and resets the time
		// of day; earlier or equal only ever runs the clock forward.
		// (Dates compare as strings, which holds because days are written
		// three digits wide. Keep them that way.)
		if s.Room.Date > s.date {
			s.date = s.Room.Date
			if at := s.Room.Seconds(); at > 0 {
				s.Clock = at
			}
		} else if at := s.Room.Seconds(); at > s.Clock {
			s.Clock = at
		}
		for _, f := range s.Room.Sets {
			s.Set(f)
		}
		for _, id := range s.Room.Grants {
			if _, had := s.Held[id]; had {
				continue
			}
			s.Grant(id)
			if c := s.Clue(id); c != nil {
				out = append(out, fmt.Sprintf("[ %s HELD — %s ]\n[ FILE %s   or   NOTE %s ]",
					id, c.Label, id, id))
			}
		}
	}
	for _, f := range s.Room.While {
		s.Set(f)
	}
	if s.Room.ID == s.Game.BodyScene && s.bodyDeadline == 0 {
		s.bodyDeadline = s.Clock + bodyLimit
	}
	return append(out, s.removeBody()...)
}

// removeBody takes him when his time is up, wherever the player is standing.
func (s *State) removeBody() []string {
	if !s.BodyPresent || s.bodyDeadline == 0 || s.Clock < s.bodyDeadline {
		return nil
	}
	s.BodyPresent = false
	s.Set("BODY_REMOVED")
	// BodyRemoved is a field on the scene, not a string in its table — which is
	// why this once reached for a key that has never existed and the medicae
	// took him in total silence.
	chambers := s.Game.Scene(s.Game.BodyScene)
	if s.Room == chambers && chambers.BodyRemoved != "" {
		return []string{chambers.BodyRemoved}
	}
	return []string{s.S("body_removed_elsewhere")}
}

// Clue looks a clue up anywhere in the game, not just in the room the player
// is standing in. Evidence travels; scenes do not.
func (s *State) Clue(id string) *content.Clue { return s.Game.Clue(id) }

// S returns an authored string for the current scene.
func (s *State) S(key string) string { return s.Room.S(key) }

// Move puts the player in another scene.
func (s *State) Move(id string) bool {
	sc := s.Game.Scene(id)
	if sc == nil {
		return false
	}
	s.leaving()
	s.Room = sc
	s.arrived = s.arrive()
	return true
}

// leaving settles the disposition of everyone in the room you are walking out
// of. A person you did not challenge is a person the record says you did not
// challenge, and that is a decision whether or not you made it on purpose.
func (s *State) leaving() {
	// Whatever was true of this place stops being true when you walk out of it.
	for _, f := range s.Room.While {
		delete(s.Flags, f)
	}
	for _, n := range s.Room.People() {
		if n.Disposition == nil || !s.Met[n.Key()] {
			continue
		}
		if n.Disposition.Default == "" {
			continue // saying nothing is a real answer here
		}
		settled := false
		for _, f := range n.Disposition.Flags {
			if s.Has(f) {
				settled = true
				break
			}
		}
		if !settled {
			s.Set(n.Disposition.Default)
		}
	}
}

// SetOne writes one flag out of a mutually exclusive group and clears its
// siblings. The act documents call these enums; the engine keeps them as the
// same flat SCREAMING_SNAKE bools as everything else, because a translation
// layer between the writing and the code is how the two stop agreeing.
func (s *State) SetOne(group []string, flag string) {
	for _, f := range group {
		delete(s.Flags, f)
	}
	s.Set(flag)
}

// Reconciliation reports the player's exposure to the archive. Nothing in the
// game displays it and nothing ever names it; it is read only to take things
// away.
func (s *State) Reconciliation() int { return s.Ints["RECONCILIATION"] }

func (s *State) Set(f string)      { s.Flags[f] = true }
func (s *State) Has(f string) bool { return s.Flags[f] }

// key scopes a per-object counter to the scene it was seen in.
func (s *State) key(id string) string { return s.Room.ID + "/" + id }

// bounds are the ranges the bible gives these counters. A counter with no
// entry here is unbounded on purpose.
var bounds = map[string][2]int{
	"CRASE_STANDING": {-5, 5},
	"BEZ_TRUST":      {0, 5},
	"ILARIA_TRUST":   {0, 5},
	"FENN_DRIFT":     {0, 6},
	"RECONCILIATION": {0, 100},
}

// Add moves a counter, held inside its range.
//
// FENN_DRIFT is the one counter with a rule attached: under Cael's shroud it
// accrues at a third the rate, which is a real gift, freely given, by a man
// who built it three weeks before anybody asked him to.
func (s *State) Add(k string, n int) {
	// Force your way into the Rust Choir with a rosette and everything you
	// ever ask for down there costs twice, for the rest of the act.
	if k == "MOTT_FAVOURS" && n > 0 && s.Has("CHOIR_HOSTILE") {
		n *= 2
	}
	if k == "FENN_DRIFT" && n > 0 && s.Has("FENN_SHIELDED") {
		s.driftPartial += n
		n = s.driftPartial / 3
		s.driftPartial -= n * 3
		if n == 0 {
			return
		}
	}
	s.add(k, n)
	if k == "FENN_DRIFT" {
		// Three sub-queries burn him partially; past four there is nothing
		// left to protect. Nothing announces either threshold.
		if s.Ints[k] >= 3 {
			s.Set("FENN_BURNED_PARTIAL")
		}
		if s.Ints[k] >= 5 {
			delete(s.Flags, "FENN_INTACT")
		}
	}
}

func (s *State) add(k string, n int) {
	s.Ints[k] += n
	if b, ok := bounds[k]; ok {
		if s.Ints[k] < b[0] {
			s.Ints[k] = b[0]
		}
		if s.Ints[k] > b[1] {
			s.Ints[k] = b[1]
		}
	}
}

// TimeString renders the clock the way a cogitator would. Act II runs over
// weeks, so the date belongs to the scene.
func (s *State) TimeString() string {
	d := s.date
	if d == "" {
		d = "0997.4.211"
	}
	return fmt.Sprintf("%s / %02dh%02d", d, s.Clock/3600, (s.Clock%3600)/60)
}

// Tick advances the scene clock and returns any events that fired.
func (s *State) Tick() []string {
	var out []string
	s.Clock += secondsPerAction

	out = append(out, s.removeBody()...)
	if s.requisitioned {
		s.sinceReq++
		if s.sinceReq == correctionDelay && !s.Archive.Corrected() {
			// Silent. No klaxon, no glitch, no colour shift. This is the
			// single most important line in the prototype.
			s.Archive.Correct()
		}
		if s.sinceReq == fallbackDelay && s.Archive.Corrected() && !s.seenCorrection && !s.fallbackFired {
			s.fallbackFired = true
			out = append(out, craseFallback)
		}
	}
	if s.Has("LYLE_QUERIED") && !s.lyleConfirmed {
		s.lyleQueried++
		if s.lyleQueried == lyleConfirmDelay {
			s.lyleConfirmed = true
			s.Grant("B5")
			out = append(out, s.S("lyle_confirmation"))
		}
	}
	if !s.Has("CRASE_ASKED_FOR_YOU") && s.threadsTouched() >= 2 {
		s.Set("CRASE_ASKED_FOR_YOU")
		out = append(out, s.S("crase_summons"))
	}
	s.actions++
	if !s.precisDone && !s.PrecisPending {
		due := s.actions >= precisFloor
		if s.requisitioned && s.sinceReq >= precisDelay {
			due = true
		}
		if due {
			out = append(out, s.openPrecis()...)
		}
	}
	return out
}

const craseFallback = `Crase puts her head round the door on her way somewhere else, entirely cheerful.

"Why in the Throne's name did you pull two hundred years of drainage ledgers?"

Vaux is behind her in the corridor. He does not look up. He writes something down.

You know you did not pull two hundred years of drainage ledgers. You know it the way
you know your own name, which is to say: you know it, and there is no longer a single
object in this hive that agrees with you.`

// threadsTouched counts how many of Vosk's four open files the player has put
// a hand on. Crase asks for them at two, which is the point at which an
// investigation stops looking like grief and starts looking like a direction.
func (s *State) threadsTouched() int {
	threads := [][]string{
		{"B1", "B2"},                // FENN/CUSTODY
		{"B4", "B4b", "B6"},         // SUMP-VARIANCE
		{"B8", "B11", "B10"},        // ECCL/ALMS
		{"B23", "B7", "B6a", "B18"}, // the Vault run
	}
	n := 0
	for _, t := range threads {
		for _, id := range t {
			if _, held := s.Held[id]; held {
				n++
				break
			}
		}
	}
	return n
}

// openPrecis composes the Crown's account of the player's day and asks them to
// countersign it. Every sentence in it is true according to the archive.
func (s *State) openPrecis() []string {
	s.PrecisPending = true
	s.Set("ACT1_TURN_SEEN")
	s.Grant("A22")
	return []string{s.S("precis"), s.S("precis_prompt")}
}

// Signs returns the held signs of the scene the player is standing in, in
// canonical order. A sign is a fact about one room; carrying it into another
// room does not make that room readable.
func (s *State) Signs() []string {
	var out []string
	for id, st := range s.Held {
		if st == StoreNone {
			continue
		}
		if c := s.Room.Clues[id]; c != nil && c.Sign {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}

func (s *State) Holds(id string) bool {
	st, ok := s.Held[id]
	return ok && st != StoreNone
}

// Grant records a clue as held, without yet choosing a store.
func (s *State) Grant(id string) {
	if id == "" {
		return
	}
	if _, ok := s.Held[id]; !ok {
		s.Held[id] = StoreNone
		if c := s.Clue(id); c != nil {
			for _, f := range c.Sets {
				s.Set(f)
			}
		}
	}
}

func (s *State) storeClue(id string, st Store) {
	cur := s.Held[id]
	if cur == StoreNone || cur == st {
		s.Held[id] = st
		return
	}
	s.Held[id] = StoreBoth
}

// NotebookDump renders the physical notebook. Immutable by construction:
// nothing in the engine has a method that edits an existing entry.
func (s *State) NotebookDump() string {
	if len(s.Notebook) == 0 {
		return "Blank. You have written nothing down.\n\nThe cogitator has been very helpful."
	}
	var b strings.Builder
	b.WriteString("YOUR NOTEBOOK — in your own hand, in ink\n\n")
	for _, l := range s.Notebook {
		b.WriteString("  · " + l + "\n")
	}
	return b.String()
}
