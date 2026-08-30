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
	secondsPerAction    = 40
	bodyLimit           = 8 * 60 // eight minutes, and they mean it
	correctionDelay     = 5      // actions after the requisition before the archive reconciles
	fallbackDelay       = 9      // if the player never looks again, Crase delivers it
	precisDelay         = 14     // and then the Crown wants the day signed off
	precisFloor         = 50     // a player who never requisitions still reaches the end of the day
	lyleConfirmDelay    = 6      // and then the Vault confirms a null result you never asked it to
	canticaConfirmDelay = 1      // forty seconds, and the second one is not yours
	tassDelay           = 12     // and if you never went to the Magos, his apprentice comes to you
	symptomInterval     = 6      // thirty-six hours on the Cerement's edge
	c17Delay            = 5      // an hour, and the intake protocols are not there any more
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
	Layers   map[string]int // how far PRESS has got into each topic
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

	canticaQueried   int
	canticaConfirmed bool

	tassWaited int

	// Out of relay range the cogitator stops lying about the world and starts
	// lying about the player. One symptom every thirty-six hours, in order.
	transitTicks int
	symptomsSeen int

	// C17 goes from the cogitator within the hour, silently, unless the
	// player copied it out by hand while Cael watched.
	c17Waited int
	c17Gone   bool

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
		Layers:   map[string]int{},
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

// Set writes a flag. A few of them have something attached that has to move
// with them, because the flag is the event and the rest is bookkeeping.
func (s *State) Set(f string) {
	s.Flags[f] = true
	if f == "ILARIA_BURNED" && s.Ints["ILARIA_TRUST"] > ilariaBurnedCap {
		s.Ints["ILARIA_TRUST"] = ilariaBurnedCap
	}
	// Having both halves is what a complete cipher is. It used to be a flag
	// three acts in 2.10 remembered to write, which meant the two doors in 2.9
	// that hand the second half over — alliance and coercion — left the player
	// holding all of it and unable to read a word.
	// C9_PHYS / C9_COG is which medium the forward recitation went into, and
	// the act document names them, so the engine translates rather than
	// asking the content to write a second flag beside FILE and NOTE.
	// The Conclave channel has three sources in the act document and only one
	// of them is a person: Vaux's dead-man writ, an unspent transit writ, or
	// the second favour owed to Sazh. Spending the writ at Kadmon takes it
	// back, unless one of the other two is standing.
	if f == "VAUX_APOLOGY" || f == "SAZH_DEBT_2" || f == "TRANSPORT_WRIT" {
		if f != "TRANSPORT_WRIT" || !s.Flags["KADMON_WRIT"] {
			s.Flags["CONCLAVE_CHANNEL"] = true
		}
	}
	if f == "KADMON_WRIT" && !s.Flags["VAUX_APOLOGY"] && !s.Flags["SAZH_DEBT_2"] {
		delete(s.Flags, "CONCLAVE_CHANNEL")
	}
	// A door closes the game. There is no scene after the intake and no
	// prompt saying so; the act simply is the last thing that happens.
	if strings.HasPrefix(f, "ENDING_") {
		s.Done = true
	}
	if id, ok := strings.CutSuffix(f, "_FILED"); ok {
		s.storeFromFlag(id, StoreCog)
	}
	if id, ok := strings.CutSuffix(f, "_PHYSICAL"); ok {
		s.storeFromFlag(id, StoreBook)
	}
	if f == "C9_PHYSICAL" {
		s.Flags["C9_PHYS"] = true
	}
	if f == "C9_FILED" {
		s.Flags["C9_COG"] = true
	}
	// VOYAGE_FULL is every watch the player could have stood, stood. A member
	// who is dead, burned or blank is not a watch they failed to keep.
	if strings.HasPrefix(f, "WATCH_") && !s.Flags["VOYAGE_FULL"] {
		done := (!s.Flags["BEZ_ALIVE"] || s.Flags["WATCH_BEZ"]) &&
			(s.Ints["ILARIA_TRUST"] < 2 || s.Flags["WATCH_ILARIA"]) &&
			(!s.Flags["FENN_INTACT"] || s.Flags["WATCH_FENN"]) &&
			(!s.Flags["VAUX_ALLIED"] || s.Flags["WATCH_VAUX"])
		if done {
			s.Flags["VOYAGE_FULL"] = true
		}
	}
	// The method is what the nav plot, the casket and the uncoupled vox say
	// together: delivery is physical, in a reliquary, correctly addressed.
	if voskMethod[f] {
		if s.Flags["C2_HELD"] && s.Flags["C3_HELD"] {
			s.Flags["VOSK_METHOD_PARTIAL"] = true
			if s.Flags["C5_HELD"] {
				s.Flags["VOSK_METHOD_UNDERSTOOD"] = true
			}
		}
	}
	if cipherHalves[f] && !s.Flags["CIPHER_COMPLETE"] &&
		s.Flags["CIPHER_HALF_ROOM"] &&
		(s.Flags["CIPHER_HALF_VAUX"] || s.Flags["CIPHER_HALF_REBUILT"]) {
		s.Flags["CIPHER_COMPLETE"] = true
	}
}

// cipherHalves are the flags that can complete the cipher by arriving. The
// room half comes off Vosk's furniture, in Act I or four days late in Theta;
// the other half is Vaux's, or is rebuilt in 2.10 at the cost of a person.
var voskMethod = map[string]bool{"C2_HELD": true, "C3_HELD": true, "C5_HELD": true}

var cipherHalves = map[string]bool{
	"CIPHER_HALF_ROOM":    true,
	"CIPHER_HALF_VAUX":    true,
	"CIPHER_HALF_REBUILT": true,
}

func (s *State) Has(f string) bool { return s.Flags[f] }

// storeFromFlag makes <ID>_FILED and <ID>_PHYSICAL mean the same thing however
// they were written. FILE and NOTE store the clue before they set the flag, so
// this only fires for content that declares the filing itself — Cael putting
// the intake protocols into your case-notes without being asked.
func (s *State) storeFromFlag(id string, st Store) {
	c := s.Clue(id)
	if c == nil {
		return
	}
	if cur, ok := s.Held[id]; ok && (cur == st || cur == StoreBoth) {
		return
	}
	s.Grant(id)
	s.storeClue(id, st)
	if st == StoreCog {
		s.Archive.File(fmt.Sprintf("%s — %s", id, c.Label))
	} else {
		s.Notebook = append(s.Notebook, c.Label)
	}
}

// key scopes a per-object counter to the scene it was seen in.
func (s *State) key(id string) string { return s.Room.ID + "/" + id }

// bounds are the ranges the bible gives these counters. A counter with no
// entry here is unbounded on purpose.
// ilariaBurnedCap is what is left of her after forty hours of raw cylinder:
// enough to carry grain and tithe, and not enough to carry anything that
// matters. The regard is not what she spent.
const ilariaBurnedCap = 2

var bounds = map[string][2]int{
	"CRASE_STANDING": {-5, 5},
	"BEZ_TRUST":      {0, 5},
	"ILARIA_TRUST":   {0, 5},
	"FENN_DRIFT":     {0, 6},
	"RECONCILIATION": {0, 100},
	// Vaux arrives at 0 or 2 depending on how he was asked for, and comes
	// back from three weeks off-world at 3 whatever he left at, which is why
	// the ceiling does the arithmetic instead of the content.
	"VAUX_SUSPICION": {0, 3},
	"VURR_RESPECT":   {0, 3},
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
	// Forty hours of raw cylinder does not cost her the player's regard. It
	// costs her the thing the regard was for, and the number stops at the
	// point where she can still carry grain and tithe.
	if k == "ILARIA_TRUST" && n > 0 && s.Has("ILARIA_BURNED") && s.Ints[k] >= ilariaBurnedCap {
		return
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
	// Forty seconds after the Vault says there is no such routing address, the
	// Vault says it again without being asked. Nothing in the game explains
	// which of the two answers was the query.
	// The one clue in the game the archive takes back. No klaxon, no notice:
	// it is simply not in the case-file any more, and re-querying returns the
	// same NO SUCH ENTRY a shelf that was never filled returns.
	if s.Has("C17_FILED") && !s.c17Gone && !s.Has("C17_PHYS") {
		s.c17Waited++
		if s.c17Waited >= c17Delay {
			s.c17Gone = true
			delete(s.Held, "C17")
			delete(s.Flags, "C17_FILED")
			delete(s.Flags, "C17_HELD")
			s.Archive.Delete("C17")
		}
	}
	if s.Has("IN_TRANSIT") {
		s.transitTicks++
		if s.transitTicks%symptomInterval == 0 {
			if line := s.symptom(); line != "" {
				out = append(out, line)
			}
		}
	}
	// The plate is unmissable. A player who never took the ash to the workshop
	// gets it anyway, brought by an apprentice who was told to, and loses only
	// the Magos's signature on the recovery.
	if s.Has("CIPHER_COMPLETE") && !s.Has("ERRATUM_OBTAINED") {
		s.tassWaited++
		if s.tassWaited == tassDelay {
			s.Grant("B21")
			s.Set("ERRATUM_OBTAINED")
			s.Set("ERRATUM_UNUSED")
			s.Set("TASS_BROUGHT_IT")
			out = append(out, s.S("tass_delivery"))
		}
	}
	if s.Has("CANTICA_QUERIED") && !s.canticaConfirmed {
		s.canticaQueried++
		if s.canticaQueried == canticaConfirmDelay {
			s.canticaConfirmed = true
			out = append(out, s.S("cantica_confirm"))
			if s.Has("FENN_INTACT") && s.Has("FENN_ALIVE") {
				out = append(out, s.S("cantica_fenn"))
			}
		}
	}
	if !s.Has("TWO_THREADS") && s.threadsTouched() >= 2 {
		s.Set("TWO_THREADS")
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

// symptom fires the next thing the record does to the player, in order. The
// third one needs the exposure to have got somewhere; below that the passage
// is simply quiet, and the player has no way to know they were spared it.
func (s *State) symptom() string {
	keys := []string{"symptom_note", "symptom_option", "symptom_memory"}
	if s.symptomsSeen >= len(keys) {
		return ""
	}
	if keys[s.symptomsSeen] == "symptom_memory" && s.Reconciliation() < 60 {
		return ""
	}
	k := keys[s.symptomsSeen]
	s.symptomsSeen++
	return s.S(k)
}

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
