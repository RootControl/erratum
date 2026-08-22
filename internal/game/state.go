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
	return "unheld"
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
	Room *content.Room

	Clock       int
	BodyPresent bool
	bodyWarned  bool

	Held     map[string]Store
	Examined map[string]int
	Flags    map[string]bool
	Ints     map[string]int
	Carried  map[string]bool

	Notebook []string // append-only. Nothing in the game can rewrite this.
	Archive  *Archive

	requisitioned  bool
	sinceReq       int
	fallbackFired  bool
	seenCorrection bool
	Done           bool
}

func New(r *content.Room) *State {
	return &State{
		Room:        r,
		BodyPresent: true,
		Held:        map[string]Store{},
		Examined:    map[string]int{},
		Flags:       map[string]bool{},
		Ints:        map[string]int{"VERITY": 0, "CRASE_STANDING": 0, "RECONCILIATION": 0},
		Carried:     map[string]bool{},
		Archive:     NewArchive(),
	}
}

func (s *State) Set(f string)        { s.Flags[f] = true }
func (s *State) Has(f string) bool   { return s.Flags[f] }
func (s *State) Add(k string, n int) { s.Ints[k] += n }

// TimeString renders the clock the way a cogitator would.
func (s *State) TimeString() string {
	base := 6*3600 + 12*60 // sixth hour, twelve minutes
	t := base + s.Clock
	return fmt.Sprintf("0997.4.211 / %02dh%02d", t/3600, (t%3600)/60)
}

// Tick advances the scene clock and returns any events that fired.
func (s *State) Tick() []string {
	var out []string
	s.Clock += secondsPerAction

	if s.BodyPresent && s.Clock >= bodyLimit {
		s.BodyPresent = false
		out = append(out, s.Room.S("body_removed"))
	}
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

// Signs returns the held signs (A1-A5) in canonical order.
func (s *State) Signs() []string {
	var out []string
	for id, st := range s.Held {
		if st == StoreNone {
			continue
		}
		if c, ok := s.Room.Clues[id]; ok && c.Sign {
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
