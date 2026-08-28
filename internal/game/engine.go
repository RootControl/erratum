package game

import (
	"fmt"
	"sort"
	"strings"

	"github.com/RootControl/erratum/internal/content"
)

// Do executes one player command and returns the text to print, plus how many
// timed actions it cost. Free actions (inventory, help, reading your own
// notebook) return zero; the body is on a timer and the player should not be
// punished for checking what they have. Travel and custody cost more than one,
// because they do.
func (s *State) Do(raw string) (out []string, ticks int) {
	line := strings.TrimSpace(raw)
	if line == "" {
		return nil, 0
	}
	fields := strings.Fields(strings.ToLower(line))
	verb := fields[0]
	arg := strings.Join(fields[1:], " ")

	switch verb {
	case "help", "?", "verbs":
		return []string{s.S("help")}, 0

	case "quit", "exit":
		s.Done = true
		return []string{"You leave it where it is."}, 0

	case "look", "l":
		// Look, LookNoBody, Intro and BodyRemoved are fields on the scene, not
		// keys in its string table. Reaching for them through S() is how this
		// room spent its whole life answering LOOK with a missing-string
		// marker. Anything authored as a field is read as a field.
		if !s.BodyPresent && s.Room.LookNoBody != "" {
			return []string{s.Room.LookNoBody}, 0
		}
		return []string{s.Room.Look}, 0

	case "go", "leave", "visit":
		// LEAVE IT is something you do with a cylinder. A person's act wins
		// only when the argument does not name a door.
		if arg != "" && s.Room.FindExit(arg) == nil {
			if out, cost, ok := s.act(verb, arg); ok {
				return out, cost
			}
		}
		return s.travel(arg)

	case "ask", "interview", "question":
		return s.ask(arg)

	case "show", "present":
		return s.show(arg)

	case "talk", "greet", "approach", "see":
		return s.talk(strings.TrimPrefix(arg, "to "))

	case "topics", "subjects":
		return s.topicList(arg)

	case "who":
		return s.whoDump(), 0

	case "examine", "x", "inspect", "read":
		if verb == "read" && (arg == "room" || strings.HasPrefix(arg, "room ")) {
			return s.readRoom(strings.TrimPrefix(arg, "room "))
		}
		// READ ROSTER is something you do for somebody, not something you look
		// at. A person's act wins only when nothing here answers to the name.
		if s.Room.Find(arg) == nil && s.Room.FindNPC(arg) == nil {
			if out, cost, ok := s.act(verb, arg); ok {
				return out, cost
			}
		}
		return s.examine(arg)

	case "take", "get", "pocket":
		// TAKE ROSETTE is a ceremony, not an object. A person's act wins only
		// when there is no actual takeable thing by that name.
		if o := s.Room.Find(arg); o == nil || !o.Takeable {
			if out, cost, ok := s.act(verb, arg); ok {
				return out, cost
			}
		}
		return s.take(arg)

	case "file":
		return s.file(arg)

	case "note", "write":
		return s.note(arg)

	case "clues", "evidence":
		return []string{s.cluesDump()}, 0

	case "notes", "casenotes", "case-notes":
		return []string{s.caseNotesDump()}, 0

	case "notebook", "coat":
		return []string{s.NotebookDump()}, 0

	case "query", "requisition", "req":
		if !s.Room.Terminal {
			return []string{s.S("no_terminal")}, 0
		}
		return s.query(strings.ToUpper(strings.TrimSpace(arg)))

	case "wait", "z":
		return []string{"Time passes. It is the one thing in this hive nobody has to requisition."}, 1

	case "status":
		return []string{fmt.Sprintf("VERITY %d   CRASE_STANDING %+d   RECONCILIATION %d   %s",
			s.Ints["VERITY"], s.Ints["CRASE_STANDING"], s.Ints["RECONCILIATION"], s.TimeString())}, 0
	}

	// Acts are verbs a person defines, not the parser: DETAIN belongs to Vaux,
	// not to the terminal.
	if out, cost, ok := s.act(verb, arg); ok {
		return out, cost
	}
	return []string{"That is not a verb this terminal knows. HELP lists them."}, 0
}

// travel moves between scenes. The lift is not free.
func (s *State) travel(arg string) ([]string, int) {
	arg = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(arg, "to "), "back to "))
	if arg == "" {
		// Doors are visible. Naming them is not the same as listing clues, and
		// guessing on the player's behalf costs three actions and a corpse.
		if len(s.Room.Exits) == 0 {
			return []string{s.S("no_such_exit")}, 0
		}
		var names []string
		for _, e := range s.Room.Exits {
			names = append(names, strings.ToUpper(e.ID))
		}
		return []string{"GO where? " + strings.Join(names, ", ") + "."}, 0
	}
	e := s.Room.FindExit(arg)
	if e == nil {
		return []string{s.S("no_such_exit")}, 0
	}
	if !s.exitOpen(e) {
		if e.Blocked != "" {
			return []string{e.Blocked}, 0
		}
		return []string{s.S("no_such_exit")}, 0
	}
	cost := e.Cost
	if cost < 1 {
		cost = 1
	}
	out := []string{}
	if e.Text != "" {
		out = append(out, e.Text)
	}
	first := !s.Visited[e.To]
	if !s.Move(e.To) {
		return []string{s.S("no_such_exit")}, 0
	}
	out = append(out, s.Room.Title, s.entering(first))
	out = append(out, s.arrived...)
	return out, cost
}

// exitOpen reads an exit's gates. Requires is a conjunction; RequiresAny and
// HasAny are one disjunction between them, because the Choir Spire opens on
// the shutter or on Crase's lead and the act document says "or".
func (s *State) exitOpen(e *content.Exit) bool {
	for _, f := range e.Requires {
		if !s.Has(f) {
			return false
		}
	}
	if len(e.RequiresAny) == 0 && len(e.HasAny) == 0 {
		return true
	}
	for _, f := range e.RequiresAny {
		if s.Has(f) {
			return true
		}
	}
	for _, id := range e.HasAny {
		if _, held := s.Held[id]; held {
			return true
		}
	}
	return false
}

// entering renders arriving somewhere. The intro happens once; after that a
// room is a room.
func (s *State) entering(first bool) string {
	switch {
	case first:
		return s.Room.Intro
	case s.Room.Revisit != "":
		return s.Room.Revisit
	default:
		return s.Room.Look
	}
}

// topicList exists so that the one scene which needs a listed conversation can
// have one. Vaux is not that scene, and says so.
func (s *State) topicList(arg string) ([]string, int) {
	n := s.npcIn(arg)
	if n == nil {
		return []string{s.whyNot(arg)}, 0
	}
	if !n.ListsTopics {
		if n.TopicRefusal != "" {
			return []string{n.TopicRefusal}, 0
		}
		return []string{"You would have to know what you were asking for."}, 0
	}
	var b strings.Builder
	b.WriteString("YOU COULD ASK ABOUT\n\n")
	for i := range n.Topics {
		t := &n.Topics[i]
		if !s.topicOpen(t) {
			continue
		}
		b.WriteString("  · " + strings.ReplaceAll(t.ID, "_", " ") + "\n")
	}
	return []string{b.String()}, 0
}

func (s *State) whoDump() []string {
	if len(s.Room.People()) == 0 {
		return []string{s.S("no_one_here")}
	}
	var out []string
	for _, n := range s.Room.People() {
		out = append(out, n.Title)
	}
	return out
}

func (s *State) examine(arg string) ([]string, int) {
	if arg == "" {
		return []string{"Examine what?"}, 0
	}
	o := s.Room.Find(arg)
	if o == nil {
		// People are things you can look at, and looking at one costs what
		// looking at anything else costs — once.
		if n := s.Room.FindNPC(arg); n != nil && n.Portrait != "" {
			k := s.key("npc/" + n.ID)
			s.Examined[k]++
			if s.Examined[k] > 1 {
				return []string{s.S("already_read")}, 0
			}
			return []string{n.Portrait}, 1
		}
		return []string{s.S("no_such_thing")}, 0
	}
	if o.BodyOnly && !s.BodyPresent {
		return []string{s.S("body_gone")}, 0
	}

	k := s.key(o.ID)
	s.Examined[k]++
	switch s.Examined[k] {
	case 1:
		return []string{o.Careless}, 1
	case 2:
		var out []string
		out = append(out, o.Careful)
		if o.Clue != "" {
			s.Grant(o.Clue)
			c := s.Clue(o.Clue)
			out = append(out, fmt.Sprintf("[ %s HELD — %s ]\n[ FILE %s   or   NOTE %s ]", o.Clue, c.Label, o.Clue, o.Clue))
		}
		return out, 1
	default:
		return []string{s.S("already_read")}, 0
	}
}

func (s *State) take(arg string) ([]string, int) {
	o := s.Room.Find(arg)
	if o == nil {
		return []string{s.S("no_such_thing")}, 0
	}
	if !o.Takeable {
		return []string{"You leave it where he put it."}, 0
	}
	if s.Carried[o.ID] {
		return []string{"You already have it."}, 0
	}
	s.Carried[o.ID] = true
	s.Set("HAS_VOSK_NOTEBOOK")
	return []string{o.TakeText}, 1
}

// file is the dangerous verb. It makes a clue actionable and hands it to the
// archive in the same keystroke, and the game never says so.
func (s *State) file(arg string) ([]string, int) {
	c := s.namedClue(arg)
	if c == nil {
		return []string{"You hold no such clue."}, 0
	}
	id := c.ID
	if c.Store == content.StoreNotebook {
		return []string{c.Refusal}, 1
	}
	if s.Held[id] == StoreCog || s.Held[id] == StoreBoth {
		return []string{"Already filed."}, 0
	}
	s.storeClue(id, StoreCog)
	s.Set(id + "_FILED")
	s.Set(id + "_HELD")
	cost := c.Reconciliation
	if cost < 1 {
		cost = 1
	}
	s.Add("RECONCILIATION", cost)
	s.Archive.File(fmt.Sprintf("%s — %s", id, c.Label))
	out := []string{c.FileText}
	if c.Store == content.StoreCogitator && c.Refusal != "" {
		out = append([]string{c.Refusal}, out...)
	}
	return out, 1
}

func (s *State) note(arg string) ([]string, int) {
	c := s.namedClue(arg)
	if c == nil {
		return []string{"You hold no such clue."}, 0
	}
	id := c.ID
	if c.Store == content.StoreCogitator {
		s.storeClue(id, StoreBook)
		s.Notebook = append(s.Notebook, c.Label)
		return []string{c.Refusal}, 1
	}
	if s.Held[id] == StoreBook || s.Held[id] == StoreBoth {
		return []string{"Already in your notebook."}, 0
	}
	s.storeClue(id, StoreBook)
	s.Set(id + "_PHYSICAL")
	s.Set(id + "_HELD")
	s.Add("VERITY", 1)
	s.Notebook = append(s.Notebook, c.Label)
	if c.NoteText != "" {
		return []string{c.NoteText}, 1
	}
	return []string{"You write it down. It is now a thing in your coat and nothing else."}, 1
}

// namedClue resolves what the player typed to a clue they actually hold. IDs
// work; so does the word a person would use out loud, because not every clue
// arrives with a label attached to it.
func (s *State) namedClue(arg string) *content.Clue {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return nil
	}
	if c := s.Clue(strings.ToUpper(arg)); c != nil {
		if _, held := s.Held[c.ID]; held {
			return c
		}
		return nil
	}
	c := s.Game.FindClue(arg)
	if c == nil {
		return nil
	}
	if _, held := s.Held[c.ID]; !held {
		return nil
	}
	return c
}

func (s *State) cluesDump() string {
	if len(s.Held) == 0 {
		return "You hold nothing. It is " + s.TimeString() + "."
	}
	var b strings.Builder
	b.WriteString("HELD\n\n")
	for _, id := range sortedKeys(s.Held) {
		c := s.Clue(id)
		label := id
		if c != nil {
			label = c.Label
		}
		b.WriteString(fmt.Sprintf("  %-4s [%-9s] %s\n", id, s.Held[id].String(), label))
	}
	b.WriteString("\nUnfiled clues cannot be requisitioned against. Filed clues can be read by whatever reads the archive.")
	return b.String()
}

func (s *State) caseNotesDump() string {
	n := s.Archive.CaseNotes()
	if len(n) == 0 {
		return "CASE-NOTES :: VOSK, H. :: EMPTY"
	}
	var b strings.Builder
	b.WriteString("CASE-NOTES :: VOSK, H. :: INQ. ROSETTE 4471-K\n\n")
	for _, l := range n {
		b.WriteString("  " + l + "\n")
	}
	if s.Archive.Corrected() {
		s.seenCorrection = true
	}
	return b.String()
}

func (s *State) readRoom(dir string) ([]string, int) {
	signs := s.Signs()
	if len(signs) < 3 {
		return []string{s.S("read_room_locked")}, 0
	}
	dir = strings.TrimSpace(dir)
	if dir != "forward" && dir != "backward" && dir != "backwards" {
		return []string{s.S("read_room_prompt")}, 0
	}

	// The reading is itself a clue, and where the player puts it is the most
	// consequential filing decision in the act.
	grant := func(out []string) []string {
		if _, had := s.Held["A8"]; had {
			return out
		}
		s.Grant("A8")
		c := s.Clue("A8")
		return append(out, fmt.Sprintf("[ A8 HELD — %s ]\n[ FILE A8   or   NOTE A8 ]", c.Label))
	}

	full := len(signs) >= 5 && s.Holds("A3") && dir != "forward"
	if full {
		s.Set("ROOM_READ_FULL")
		s.Set("CIPHER_HALF_ROOM")
		return grant([]string{readingFull}), 1
	}
	if dir != "forward" && !s.Holds("A3") {
		return []string{readingBackwardNoBoots}, 1
	}
	s.Set("ROOM_READ_PARTIAL")
	return grant([]string{readingPartial}), 1
}

const readingPartial = `Three of them and the room starts to hum. Not literally. It is only that the
objects stop being objects.

The board is an address. The shutter is a clock and a finger.

You have a shelf-code, forward, as he set it down: D2-F4-C9. You have the
twelfth hour. You have four minutes on a drum log.

Something is missing and you cannot name it, which is a thing you will think
about later, at length, for years.

[ ROOM_READ_PARTIAL ]`

const readingBackwardNoBoots = `You try it the other way and nothing arrives.

Reversal is a decision and you have no authority for it — nothing in this room
tells you to read it backwards, and you are not going to invert a dead
Inquisitor's last message on a hunch.

[ the boots are on the wrong feet ]`

const readingFull = `You stop moving and let the room assemble itself.

The board is an address. The shutter is a clock and a finger. The boots are a
sentence you learned in a rest-halt on Sarrow's Reach when you were twenty-six
and too frightened to eat.

Read it backwards.

D2-F4-C9 is Munitorum drainage ledgers, two centuries of them, and it is the
mistake he expected somebody else to make. C9-F4-D2 is Block C, Level 9: the
Choir of Cantor hardcopy dispatch manifests, current cycle. The physical ledger
of every message the Choir received.

Including the one that killed him.

He had six days. No vox he trusted, no page that wasn't read, no confidant the
record couldn't map. So he put it in furniture, and he put the key to the
furniture in a thing that exists only between you and a dead man.

[ ROOM_READ_FULL — CIPHER_HALF_ROOM ]
[ the live shelf-code is C9-F4-D2 ]`

func (s *State) query(code string) ([]string, int) {
	if code == "" {
		return []string{"QUERY <shelf-address>. The Vault does not take guesses politely."}, 0
	}
	code = strings.ReplaceAll(code, "/", "-")
	code = strings.ReplaceAll(code, " ", "-")

	if code == CodeDispatch {
		if !s.Has("DISPATCH_ID_KNOWN") {
			return []string{s.S("vault_none")}, 1
		}
		rec, ok := s.Archive.Query(code)
		if !ok {
			s.seenCorrection = true
			return []string{dispatchGone, s.contradiction()}, 1
		}
		return s.afterDispatch(rec), 1
	}
	if !s.Has("A1_FILED") && (code == CodeLive || code == CodeForward) {
		return []string{"The Vault indexes filings. You have not filed the shelf-address,\n" +
			"so as far as this terminal is concerned you do not have one.\n\n[ FILE A1 ]"}, 0
	}

	if code == CodeForward {
		s.Add("CRASE_STANDING", 1)
		s.Set("CRASE_SAVED_ME")
		code = CodeLive
		rec, ok := s.Archive.Query(code)
		if !ok {
			return []string{s.S("crase_rescue"), dispatchGone}, 1
		}
		return append([]string{s.S("crase_rescue")}, s.afterDispatch(rec)...), 1
	}

	if code == CodeLive {
		rec, ok := s.Archive.Query(code)
		if !ok {
			s.seenCorrection = true
			return []string{dispatchGone, s.contradiction()}, 1
		}
		return s.afterDispatch(rec), 1
	}

	return []string{s.S("vault_none")}, 1
}

func (s *State) afterDispatch(rec string) []string {
	out := []string{rec}
	if !s.requisitioned {
		s.requisitioned = true
		s.Set("KESH_NAMED")
		s.Grant("A20")
		s.storeClue("A20", StoreCog)
		s.Set("A20_FILED")
		s.Archive.File("REQ 001 — " + CodeLive + " — DISPATCH MANIFEST — 1 RECORD RETURNED — READ")
		out = append(out, "[ FILED TO CASE-NOTES. ]\n[ A20 HELD ]")
	}
	return out
}

// contradiction is the only place the game acknowledges the edit, and it only
// does so if the player earned the right by writing something down.
func (s *State) contradiction() string {
	if s.Held["A20"] == StoreBook || s.Held["A20"] == StoreBoth {
		s.Add("VERITY", 1)
		s.Set("A21_HELD")
		return `Your notebook says 1188. Vermilion. Forty-one seconds. Kesh, countersigned in person.

The Vault says no such dispatch has ever existed.

Nothing in this room adjudicates between them. There is no third thing to ask.
There is a page in your coat and a machine on the desk and they do not agree,
and one of them was written by your own hand an hour ago.

[ A21 — VERITY +1 ]`
	}
	return `You read it an hour ago. You are sure of that.

You check the case-notes. They say you requisitioned two centuries of drainage
ledgers, in error, and closed it without action. Same timestamp. Same session.
Your name on it.

You have no copy. You have a strong feeling, and a strong feeling is not
evidence, and you of all people know exactly what an Ordo conclave does with a
strong feeling.

[ no A21 — you filed everything and wrote nothing down ]`
}

func sortedKeys(m map[string]Store) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
