// Package content holds the authored material for the game. Everything the
// player reads is data, not code: scenes ship as JSON so writing and
// engineering stay separable.
package content

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed room.json
var roomJSON []byte

// Store says which evidence store a clue may live in.
//
// This is the game's central mechanic, so it is a first-class field on every
// clue rather than a special case in the engine. "cogitator" means filing is
// mandatory to make the clue actionable, and filing is what exposes it to the
// archive. "notebook" means the cogitator's filing template has no honest
// classification available. "either" means the player chooses, and the choice
// is the point.
type Store string

const (
	StoreEither    Store = "either"
	StoreCogitator Store = "cogitator"
	StoreNotebook  Store = "notebook"
)

// Object is a thing in a room that can be named at the terminal.
type Object struct {
	ID       string   `json:"id"`
	Names    []string `json:"names"`
	Careless string   `json:"careless"`
	Careful  string   `json:"careful"`
	Clue     string   `json:"clue"`
	BodyOnly bool     `json:"body_only"`
	Takeable bool     `json:"takeable"`
	TakeText string   `json:"take_text"`
}

// Clue is a piece of evidence. Sign marks the five signs Vosk arranged.
type Clue struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Store    Store  `json:"store"`
	Sign     bool   `json:"sign"`
	FileText string `json:"file_text"`
	NoteText string `json:"note_text"`
	Refusal  string `json:"refusal"`
}

// Room is the authored content for a location.
type Room struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Intro       string            `json:"intro"`
	Look        string            `json:"look"`
	LookNoBody  string            `json:"look_no_body"`
	BodyRemoved string            `json:"body_removed"`
	Objects     []Object          `json:"objects"`
	Clues       map[string]*Clue  `json:"clues"`
	Strings     map[string]string `json:"strings"`
}

// Load returns the embedded Act I room.
func Load() (*Room, error) {
	var r Room
	if err := json.Unmarshal(roomJSON, &r); err != nil {
		return nil, fmt.Errorf("content: %w", err)
	}
	return &r, nil
}

// Find resolves a player-typed noun to an object. Matching is loose on
// purpose: the player should never lose a clue to vocabulary.
func (r *Room) Find(word string) *Object {
	word = strings.ToLower(strings.TrimSpace(word))
	if word == "" {
		return nil
	}
	for i := range r.Objects {
		o := &r.Objects[i]
		if o.ID == word {
			return o
		}
		for _, n := range o.Names {
			if n == word {
				return o
			}
		}
	}
	// Fall back to prefix matching so "regi" finds the regicide board.
	for i := range r.Objects {
		o := &r.Objects[i]
		for _, n := range append([]string{o.ID}, o.Names...) {
			if strings.HasPrefix(n, word) && len(word) >= 3 {
				return o
			}
		}
	}
	return nil
}

// S returns an authored string by key, loudly if it is missing.
func (r *Room) S(key string) string {
	if v, ok := r.Strings[key]; ok {
		return v
	}
	return "[missing string: " + key + "]"
}
