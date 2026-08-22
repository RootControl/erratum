package game

import (
	"fmt"
	"strings"
)

// Do executes one player command and returns the text to print, plus whether
// the command consumed time. Free actions (inventory, help, reading your own
// notebook) never advance the clock; the body is on a timer and the player
// should not be punished for checking what they have.
func (s *State) Do(raw string) (out []string, timed bool) {
	line := strings.TrimSpace(raw)
	if line == "" {
		return nil, false
	}
	fields := strings.Fields(strings.ToLower(line))
	verb := fields[0]
	arg := strings.Join(fields[1:], " ")

	switch verb {
	case "help", "?", "verbs":
		return []string{s.Room.S("help")}, false

	case "quit", "exit":
		s.Done = true
		return []string{"You leave it where it is."}, false

	case "look", "l":
		if s.BodyPresent {
			return []string{s.Room.S("look")}, false
		}
		return []string{s.Room.S("look_no_body")}, false

	case "examine", "x", "inspect", "read":
		if verb == "read" && (arg == "room" || strings.HasPrefix(arg, "room ")) {
			return s.readRoom(strings.TrimPrefix(arg, "room "))
		}
		return s.examine(arg)

	case "take", "get", "pocket":
		return s.take(arg)

	case "file":
		return s.file(strings.ToUpper(arg))

	case "note", "write":
		return s.note(strings.ToUpper(arg))

	case "clues", "evidence":
		return []string{s.cluesDump()}, false

	case "notes", "casenotes", "case-notes":
		return []string{s.caseNotesDump()}, false

	case "notebook", "coat":
		return []string{s.NotebookDump()}, false

	case "query", "requisition", "req":
		return s.query(strings.ToUpper(strings.TrimSpace(arg)))

	case "wait", "z":
		return []string{"Time passes. It is the one thing in this hive nobody has to requisition."}, true

	case "status":
		return []string{fmt.Sprintf("VERITY %d   CRASE_STANDING %+d   RECONCILIATION %d   %s",
			s.Ints["VERITY"], s.Ints["CRASE_STANDING"], s.Ints["RECONCILIATION"], s.TimeString())}, false
	}
	return []string{"That is not a verb this terminal knows. HELP lists them."}, false
}

func (s *State) examine(arg string) ([]string, bool) {
	if arg == "" {
		return []string{"Examine what?"}, false
	}
	o := s.Room.Find(arg)
	if o == nil {
		return []string{s.Room.S("no_such_thing")}, false
	}
	if o.BodyOnly && !s.BodyPresent {
		return []string{s.Room.S("body_gone")}, false
	}

	s.Examined[o.ID]++
	switch s.Examined[o.ID] {
	case 1:
		return []string{o.Careless}, true
	case 2:
		var out []string
		out = append(out, o.Careful)
		if o.Clue != "" {
			s.Grant(o.Clue)
			c := s.Room.Clues[o.Clue]
			out = append(out, fmt.Sprintf("[ %s HELD — %s ]\n[ FILE %s   or   NOTE %s ]", o.Clue, c.Label, o.Clue, o.Clue))
		}
		return out, true
	default:
		return []string{s.Room.S("already_read")}, false
	}
}

func (s *State) take(arg string) ([]string, bool) {
	o := s.Room.Find(arg)
	if o == nil {
		return []string{s.Room.S("no_such_thing")}, false
	}
	if !o.Takeable {
		return []string{"You leave it where he put it."}, false
	}
	if s.Carried[o.ID] {
		return []string{"You already have it."}, false
	}
	s.Carried[o.ID] = true
	s.Set("HAS_VOSK_NOTEBOOK")
	return []string{o.TakeText}, true
}

// file is the dangerous verb. It makes a clue actionable and hands it to the
// archive in the same keystroke, and the game never says so.
func (s *State) file(id string) ([]string, bool) {
	c, ok := s.Room.Clues[id]
	if !ok || !s.Holds(id) && s.Held[id] != StoreNone {
		return []string{"You hold no such clue."}, false
	}
	if _, held := s.Held[id]; !held {
		return []string{"You hold no such clue."}, false
	}
	if c.Store == "notebook" {
		return []string{c.Refusal}, true
	}
	if s.Held[id] == StoreCog || s.Held[id] == StoreBoth {
		return []string{"Already filed."}, false
	}
	s.storeClue(id, StoreCog)
	s.Set(id + "_FILED")
	s.Add("RECONCILIATION", 1)
	s.Archive.File(fmt.Sprintf("%s — %s", id, c.Label))
	out := []string{c.FileText}
	if c.Store == "cogitator" && c.Refusal != "" {
		out = append([]string{c.Refusal}, out...)
	}
	return out, true
}

func (s *State) note(id string) ([]string, bool) {
	c, ok := s.Room.Clues[id]
	if !ok {
		return []string{"You hold no such clue."}, false
	}
	if _, held := s.Held[id]; !held {
		return []string{"You hold no such clue."}, false
	}
	if c.Store == "cogitator" {
		s.storeClue(id, StoreBook)
		s.Notebook = append(s.Notebook, c.Label)
		return []string{c.Refusal}, true
	}
	if s.Held[id] == StoreBook || s.Held[id] == StoreBoth {
		return []string{"Already in your notebook."}, false
	}
	s.storeClue(id, StoreBook)
	s.Set(id + "_PHYSICAL")
	s.Add("VERITY", 1)
	s.Notebook = append(s.Notebook, c.Label)
	return []string{c.NoteText}, true
}

func (s *State) cluesDump() string {
	if len(s.Held) == 0 {
		return "You hold nothing. You have been in the room for " + s.TimeString() + "."
	}
	var b strings.Builder
	b.WriteString("HELD\n\n")
	for _, id := range sortedKeys(s.Held) {
		c := s.Room.Clues[id]
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

func (s *State) readRoom(dir string) ([]string, bool) {
	signs := s.Signs()
	if len(signs) < 3 {
		return []string{s.Room.S("read_room_locked")}, false
	}
	dir = strings.TrimSpace(dir)
	if dir != "forward" && dir != "backward" && dir != "backwards" {
		return []string{s.Room.S("read_room_prompt")}, false
	}

	full := len(signs) >= 5 && s.Holds("A3") && dir != "forward"
	if full {
		s.Set("ROOM_READ_FULL")
		s.Set("CIPHER_HALF_ROOM")
		s.Set("DISPATCH_ID_KNOWN")
		return []string{readingFull}, true
	}
	if dir != "forward" && !s.Holds("A3") {
		return []string{readingBackwardNoBoots}, true
	}
	s.Set("ROOM_READ_PARTIAL")
	return []string{readingPartial}, true
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

func (s *State) query(code string) ([]string, bool) {
	if code == "" {
		return []string{"QUERY <shelf-address>. The Vault does not take guesses politely."}, false
	}
	code = strings.ReplaceAll(code, "/", "-")
	code = strings.ReplaceAll(code, " ", "-")

	if !s.Has("A1_FILED") && (code == CodeLive || code == CodeForward) {
		return []string{"The Vault indexes filings. You have not filed the shelf-address,\n" +
			"so as far as this terminal is concerned you do not have one.\n\n[ FILE A1 ]"}, false
	}

	if code == CodeForward {
		s.Add("CRASE_STANDING", 1)
		s.Set("CRASE_SAVED_ME")
		code = CodeLive
		rec, ok := s.Archive.Query(code)
		if !ok {
			return []string{s.Room.S("crase_rescue"), dispatchGone}, true
		}
		return append([]string{s.Room.S("crase_rescue")}, s.afterDispatch(rec)...), true
	}

	if code == CodeLive {
		rec, ok := s.Archive.Query(code)
		if !ok {
			s.seenCorrection = true
			return []string{dispatchGone, s.contradiction()}, true
		}
		return s.afterDispatch(rec), true
	}

	return []string{s.Room.S("vault_none")}, true
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
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}
