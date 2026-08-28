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
		if at := s.Room.Seconds(); at > s.Clock {
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
	"RECONCILIATION": {0, 100},
}

// Add moves a counter, held inside its range.
func (s *State) Add(k string, n int) {
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

// TimeString renders the clock the way a cogitator would.
func (s *State) TimeString() string {
	return fmt.Sprintf("0997.4.211 / %02dh%02d", s.Clock/3600, (s.Clock%3600)/60)
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
	return out
}

const craseFallback = `Crase puts her head round the door on her way somewhere else, entirely cheerful.

"Why in the Throne's name did you pull two hundred years of drainage ledgers?"

Vaux is behind her in the corridor. He does not look up. He writes something down.

You know you did not pull two hundred years of drainage ledgers. You know it the way
you know your own name, which is to say: you know it, and there is no longer a single
object in this hive that agrees with you.`

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
