// Package content holds the authored material for the game. Everything the
// player reads is data, not code: scenes ship as JSON so writing and
// engineering stay separable.
package content

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

//go:embed strings.json
var sharedStringsJSON []byte

//go:embed arrival.json
var arrivalJSON []byte

//go:embed antechamber.json
var antechamberJSON []byte

//go:embed room.json
var roomJSON []byte

//go:embed vaux_study.json
var vauxStudyJSON []byte

//go:embed medicae.json
var medicaeJSON []byte

//go:embed adjutant.json
var adjutantJSON []byte

//go:embed arch.json
var archJSON []byte

//go:embed choir.json
var choirJSON []byte

//go:embed board.json
var boardJSON []byte

//go:embed safehouse.json
var safehouseJSON []byte

//go:embed workshop.json
var workshopJSON []byte

//go:embed delta.json
var deltaJSON []byte

//go:embed theta.json
var thetaJSON []byte

//go:embed choir_settlement.json
var rustChoirJSON []byte

//go:embed annexe.json
var annexeJSON []byte

//go:embed rows.json
var rowsJSON []byte

//go:embed cathedral.json
var cathedralJSON []byte

//go:embed confessional.json
var confessionalJSON []byte

//go:embed crase_office.json
var craseOfficeJSON []byte

//go:embed ferrum.json
var ferrumJSON []byte

//go:embed cylinder_vault.json
var cylinderVaultJSON []byte

//go:embed skydock.json
var skydockJSON []byte

//go:embed cerement.json
var cerementJSON []byte

//go:embed spirit.json
var spiritJSON []byte

//go:embed gunbay.json
var gunbayJSON []byte

//go:embed cell.json
var cellJSON []byte

//go:embed hold.json
var holdJSON []byte

//go:embed wardroom.json
var wardroomJSON []byte

//go:embed kadmon.json
var kadmonJSON []byte

//go:embed tithe_yard.json
var titheYardJSON []byte

//go:embed muster.json
var musterJSON []byte

//go:embed reading_room.json
var readingRoomJSON []byte

//go:embed cantica.json
var canticaJSON []byte

//go:embed fifth_chair.json
var fifthChairJSON []byte

//go:embed correctorum.json
var correctorumJSON []byte

//go:embed intake.json
var intakeJSON []byte

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
	// StoreNeither is a clue that cannot be recorded at all, in any medium,
	// because the medium is four hundred people saying names out loud and
	// writing it down stops them. Both FILE and NOTE refuse.
	StoreNeither Store = "neither"
)

// Object is a thing in a room that can be named at the terminal.
//
// Objects gate like everything else. A crate that has not been unsealed is not
// a thing the player can look at, and a closed object reads as absent rather
// than as refused, because a refusal is a signpost.
type Object struct {
	Gate
	ID       string   `json:"id"`
	Names    []string `json:"names"`
	Careless string   `json:"careless"`
	Careful  string   `json:"careful"`
	Clue     string   `json:"clue"`
	BodyOnly bool     `json:"body_only"`
	Takeable bool     `json:"takeable"`
	TakeText string   `json:"take_text"`
	// Sets are flags that become true on the careful look, when the player has
	// actually handled the thing. TakeSets are the ones that need it in their
	// coat: reading a dead man's longhand off his desk is not the same as
	// walking out of his chambers with it, and the engine used to know the
	// difference for exactly one object, by name.
	Sets     []string `json:"sets"`
	TakeSets []string `json:"take_sets"`
}

// Clue is a piece of evidence. Sign marks the five signs Vosk arranged.
//
// Names exist because not every clue announces an ID. A thing overheard in a
// conversation has to be writable down by the name the player would use for
// it, or the player who noticed cannot act on having noticed.
type Clue struct {
	ID       string   `json:"id"`
	Names    []string `json:"names"`
	Label    string   `json:"label"`
	Store    Store    `json:"store"`
	Sign     bool     `json:"sign"`
	FileText string   `json:"file_text"`
	// Reconciliation is what filing this costs in exposure. Zero means the
	// default of one; A8 costs two, because A8 is the thing worth reading.
	Reconciliation int    `json:"reconciliation"`
	NoteText       string `json:"note_text"`
	Refusal        string `json:"refusal"`
	// Sets are flags that become true the moment this clue is in hand,
	// whichever way the player came by it. The dispatch number opens the Vault
	// whether it was read off a jacket or off an archivist's index.
	Sets []string `json:"sets"`
}

// Room is the authored content for a location.
type Room struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Intro string `json:"intro"`
	// Revisit replaces Intro on every entry after the first. An intro carries
	// the beats of arriving somewhere — Fenn standing up to congratulate you,
	// Ostro reading out the autopsy — and those happen once. Empty falls back
	// to Look.
	Revisit     string `json:"revisit"`
	Look        string `json:"look"`
	LookNoBody  string `json:"look_no_body"`
	BodyRemoved string `json:"body_removed"`
	// ClockAt is the hive time this scene opens at, "06h12". The clock runs
	// forward through a scene and jumps on entering the next one, because the
	// act documents place the scenes hours apart and the player does not play
	// the hours in between.
	ClockAt string `json:"clock_at"`
	// Date is the hive date this scene happens on. Act II runs over weeks.
	// Empty inherits the previous scene's date.
	Date string `json:"date"`
	// Body marks the scene the eight-minute timer belongs to. It starts when
	// the player first arrives and does not stop when they leave.
	Body bool `json:"body"`
	// Terminal marks a scene with a cogitator in it. There is exactly one in
	// the slice, and Vaux's study is pointedly not it.
	Terminal bool `json:"terminal"`
	// Cipher marks the room the objects were arranged in. DECODE needs the
	// chambers the way QUERY needs a cogitator: the sentence is read out of
	// furniture, and the furniture is in one place.
	Cipher bool `json:"cipher"`
	// Sets are the flags arriving in this scene writes, once.
	Sets []string `json:"sets"`
	// Grants are the clues arriving in this scene puts in the player's hands,
	// once and unasked. The autopsy floor is the act's guaranteed minimum and
	// a player cannot stand on the medicae deck and fail to be told it.
	Grants []string `json:"grants"`

	// While are flags that are true only while the player is standing in this
	// scene. Relay shielding is a property of a place, not an event.
	While []string `json:"while"`

	Objects []Object `json:"objects"`
	NPCs    []NPC    `json:"npcs"`
	// NPCRefs name people defined in another scene who are also here. Fenn
	// follows you to the gallery arch and does not become a different man on
	// the way, so both scenes have to reach the same person and the same
	// conversation.
	NPCRefs []string `json:"npc_refs"`
	people  []*NPC

	Exits   []Exit            `json:"exits"`
	Clues   map[string]*Clue  `json:"clues"`
	Strings map[string]string `json:"strings"`
}

// Game is every scene the slice has, plus one index across all of them.
//
// Clues are indexed game-wide on purpose: the player carries A3 out of the
// chambers and into Vaux's study, and the engine has to know what it is when
// they put it on his desk.
type Game struct {
	Scenes map[string]*Room
	Order  []string
	Start  string
	// BodyScene is where Vosk is, for as long as he is anywhere.
	BodyScene string
	Clues     map[string]*Clue
	byUpper   map[string]*Clue
	npcs      map[string][]*NPC
	byKey     map[string]*NPC
}

// Seconds parses "06h12" into seconds since midnight. Zero for an unset field.
func (r *Room) Seconds() int {
	s, _ := r.seconds()
	return s
}

func (r *Room) seconds() (int, error) {
	if r.ClockAt == "" {
		return 0, nil
	}
	var h, m int
	if _, err := fmt.Sscanf(r.ClockAt, "%dh%d", &h, &m); err != nil {
		return 0, fmt.Errorf("clock_at %q: %w", r.ClockAt, err)
	}
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, fmt.Errorf("clock_at %q is not a time", r.ClockAt)
	}
	return h*3600 + m*60, nil
}

// People are everyone in the scene: those authored here, and those who
// followed you in.
func (r *Room) People() []*NPC { return r.people }

// Scene returns a scene by ID, or nil.
func (g *Game) Scene(id string) *Room { return g.Scenes[id] }

// Clue returns a clue from anywhere in the game. Lookup is case-insensitive
// because one clue in Act I is called A1-lead in the act document and the
// document wins.
func (g *Game) Clue(id string) *Clue {
	if c, ok := g.Clues[id]; ok {
		return c
	}
	return g.byUpper[strings.ToUpper(id)]
}

// MatchClues resolves a player-typed reference to every clue that answers to
// it, best match first. Aliases legitimately overlap across an act — two
// different things are both "the pattern" — so the caller decides which one
// the player meant, usually by asking which one they are holding.
func (g *Game) MatchClues(word string) []*Clue {
	type hit struct {
		c *Clue
		s int
	}
	var hits []hit
	for _, id := range g.Order {
		for _, c := range g.Scenes[id].Clues {
			if s := match(word, c.ID, c.Names); s > 0 {
				hits = append(hits, hit{c, s})
			}
		}
	}
	sort.SliceStable(hits, func(a, b int) bool { return hits[a].s > hits[b].s })
	out := make([]*Clue, 0, len(hits))
	for _, h := range hits {
		out = append(out, h.c)
	}
	return out
}

// Load returns every embedded scene.
//
// Scenes after the first inherit any strings they do not define, so a new
// scene only has to author the lines that are actually different. The room is
// first because the room is where the game starts.
func Load() (*Game, error) {
	var shared map[string]string
	if err := json.Unmarshal(sharedStringsJSON, &shared); err != nil {
		return nil, fmt.Errorf("content: strings: %w", err)
	}
	sources := []struct {
		id  string
		raw []byte
	}{
		{"antechamber_arrival", arrivalJSON},
		{"antechamber", antechamberJSON},
		{"chambers", roomJSON},
		{"vaux_study", vauxStudyJSON},
		{"medicae", medicaeJSON},
		{"adjutant", adjutantJSON},
		{"gallery_arch", archJSON},
		{"choir", choirJSON},
		{"standing_board", boardJSON},
		{"safehouse", safehouseJSON},
		{"workshop", workshopJSON},
		{"delta", deltaJSON},
		{"theta", thetaJSON},
		{"rust_choir", rustChoirJSON},
		{"annexe", annexeJSON},
		{"the_rows", rowsJSON},
		{"cathedral", cathedralJSON},
		{"confessional", confessionalJSON},
		{"crase_office", craseOfficeJSON},
		{"ferrum", ferrumJSON},
		{"cylinder_vault", cylinderVaultJSON},
		{"skydock", skydockJSON},
		{"cerement", cerementJSON},
		{"spirit", spiritJSON},
		{"gunbay", gunbayJSON},
		{"cell", cellJSON},
		{"hold", holdJSON},
		{"wardroom", wardroomJSON},
		{"kadmon", kadmonJSON},
		{"tithe_yard", titheYardJSON},
		{"muster", musterJSON},
		{"reading_room", readingRoomJSON},
		{"cantica", canticaJSON},
		{"fifth_chair", fifthChairJSON},
		{"correctorum", correctorumJSON},
		{"intake", intakeJSON},
	}
	g := &Game{Scenes: map[string]*Room{}, Clues: map[string]*Clue{},
		byUpper: map[string]*Clue{}, npcs: map[string][]*NPC{}, byKey: map[string]*NPC{}}
	lastDate := "0997.4.211"
	for _, src := range sources {
		var r Room
		if err := json.Unmarshal(src.raw, &r); err != nil {
			return nil, fmt.Errorf("content: %s: %w", src.id, err)
		}
		if r.ID == "" {
			r.ID = src.id
		}
		if r.Strings == nil {
			r.Strings = map[string]string{}
		}
		if g.Start == "" {
			g.Start = r.ID
		}
		if r.Body {
			g.BodyScene = r.ID
		}
		if r.Date == "" {
			// Scenes are authored in chronological order; one that does not
			// name a date is on the same day as the one before it.
			r.Date = lastDate
		}
		lastDate = r.Date
		for k, v := range shared {
			if _, ok := r.Strings[k]; !ok {
				r.Strings[k] = v
			}
		}
		for id, c := range r.Clues {
			if c.ID == "" {
				c.ID = id
			}
			if _, dup := g.Clues[id]; dup {
				return nil, fmt.Errorf("content: clue %s defined twice", id)
			}
			g.Clues[id] = c
			g.byUpper[strings.ToUpper(id)] = c
		}
		for i := range r.NPCs {
			n := &r.NPCs[i]
			n.key = r.ID + "/" + n.ID
			g.npcs[n.ID] = append(g.npcs[n.ID], n)
			g.byKey[n.key] = n
			r.people = append(r.people, n)
		}
		g.Scenes[r.ID] = &r
		g.Order = append(g.Order, r.ID)
	}
	// Resolved after every scene is loaded, so a scene may reference somebody
	// authored later in the list.
	for _, id := range g.Order {
		r := g.Scenes[id]
		for _, ref := range r.NPCRefs {
			if scene, who, scoped := strings.Cut(ref, "/"); scoped {
				n := g.byKey[scene+"/"+who]
				if n == nil {
					return nil, fmt.Errorf("content: %s: no person %q in scene %q", id, who, scene)
				}
				r.people = append(r.people, n)
				continue
			}
			switch found := g.npcs[ref]; len(found) {
			case 0:
				return nil, fmt.Errorf("content: %s: no person named %q to follow you here", id, ref)
			case 1:
				r.people = append(r.people, found[0])
			default:
				return nil, fmt.Errorf("content: %s: %q is authored in more than one scene — reference it as \"scene/%s\"", id, ref, ref)
			}
		}
	}
	if g.BodyScene == "" {
		return nil, fmt.Errorf("content: no scene claims the body")
	}
	if err := g.validate(); err != nil {
		return nil, err
	}
	return g, nil
}

// validate refuses to ship content that will fail quietly at play time. Every
// mistake it catches has already happened once in this repo: a clue reference
// that resolves to nothing, an exit to a scene that is not there, prose
// authored under a key nothing reads. None of them raise an error at runtime —
// they just make a line of the game not happen.
func (g *Game) validate() error {
	fail := func(f string, a ...any) error { return fmt.Errorf("content: "+f, a...) }
	for _, id := range g.Order {
		r := g.Scenes[id]
		if _, err := r.seconds(); err != nil {
			return fail("%s: %w", id, err)
		}
		for _, e := range r.Exits {
			if g.Scenes[e.To] == nil {
				return fail("%s: exit %q leads to unknown scene %q", id, e.ID, e.To)
			}
		}
		for _, c := range r.Clues {
			switch c.Store {
			case StoreEither, StoreCogitator, StoreNotebook, StoreNeither:
			default:
				return fail("%s: clue %s has store %q", id, c.ID, c.Store)
			}
		}
		for i := range r.Objects {
			o := &r.Objects[i]
			if o.Clue != "" && g.Clue(o.Clue) == nil {
				return fail("%s: object %s grants unknown clue %q", id, o.ID, o.Clue)
			}
			for _, c := range o.Has {
				if g.Clue(c) == nil {
					return fail("%s: object %s gates on unknown clue %q", id, o.ID, c)
				}
			}
			if err := checkGate(id, "", "object "+o.ID, o.Gate, fail); err != nil {
				return err
			}
		}
		for _, e := range r.Exits {
			if err := checkGate(id, "", "exit "+e.ID, e.Gate, fail); err != nil {
				return err
			}
			if err := checkCounters(id, "", "exit "+e.ID, e.Counters, fail); err != nil {
				return err
			}
		}
		for _, c := range r.Grants {
			if g.Clue(c) == nil {
				return fail("%s: scene grants unknown clue %q", id, c)
			}
		}
		for _, n := range r.people {
			if err := g.validateNPC(id, n, fail); err != nil {
				return err
			}
		}
	}
	return nil
}

// checkGate catches a counter threshold naming a counter that does not exist,
// which would otherwise be a gate that silently never opens.
func checkGate(scene, npc, what string, g Gate, fail func(string, ...any) error) error {
	if err := checkCounters(scene, npc, what+" at_least", g.AtLeast, fail); err != nil {
		return err
	}
	return checkCounters(scene, npc, what+" below", g.Below, fail)
}

func checkCounters(scene, npc, what string, c Counters, fail func(string, ...any) error) error {
	for k := range c {
		if !KnownCounters[k] {
			return fail("%s/%s: %s moves unknown counter %q", scene, npc, what, k)
		}
	}
	return nil
}

func (g *Game) validateNPC(scene string, n *NPC, fail func(string, ...any) error) error {
	seen := map[string]bool{}
	grants := func(kind, id string, clues []string) error {
		if key := kind + "/" + id; seen[key] {
			return fail("%s/%s: duplicate %s id %q", scene, n.ID, kind, id)
		} else {
			seen[key] = true
		}
		for _, c := range clues {
			if g.Clue(c) == nil {
				return fail("%s/%s: %s %q grants unknown clue %q", scene, n.ID, kind, id, c)
			}
		}
		return nil
	}
	for _, t := range n.Topics {
		if err := grants("topic", t.ID, t.Grants); err != nil {
			return err
		}
		for _, c := range append(append([]string{}, t.Has...), t.Holds...) {
			if g.Clue(c) == nil {
				return fail("%s/%s: topic %q gates on unknown clue %q", scene, n.ID, t.ID, c)
			}
		}
	}
	for _, a := range n.Acts {
		if err := grants("act", a.ID, a.Grants); err != nil {
			return err
		}
	}
	for _, b := range n.Beats {
		if err := grants("beat", b.ID, b.Grants); err != nil {
			return err
		}
	}
	for _, sh := range n.Shows {
		if g.Clue(sh.Clue) == nil {
			return fail("%s/%s: show reacts to unknown clue %q", scene, n.ID, sh.Clue)
		}
	}
	for _, t := range n.Topics {
		if err := checkCounters(scene, n.ID, "topic "+t.ID, t.Counters, fail); err != nil {
			return err
		}
		if err := checkCounters(scene, n.ID, "topic "+t.ID, t.ClosedCounters, fail); err != nil {
			return err
		}
		if err := checkGate(scene, n.ID, "topic "+t.ID, t.Gate, fail); err != nil {
			return err
		}
	}
	for _, a := range n.Acts {
		if err := checkCounters(scene, n.ID, "act "+a.ID, a.Counters, fail); err != nil {
			return err
		}
		if err := checkGate(scene, n.ID, "act "+a.ID, a.Gate, fail); err != nil {
			return err
		}
	}
	for _, sh := range n.Shows {
		if err := checkCounters(scene, n.ID, "show "+sh.Clue, sh.Counters, fail); err != nil {
			return err
		}
		if err := checkGate(scene, n.ID, "show "+sh.Clue, sh.Gate, fail); err != nil {
			return err
		}
	}
	for _, b := range n.Beats {
		if err := checkCounters(scene, n.ID, "beat "+b.ID, b.Counters, fail); err != nil {
			return err
		}
		if err := checkGate(scene, n.ID, "beat "+b.ID, b.Gate, fail); err != nil {
			return err
		}
	}
	if d := n.Disposition; d != nil && d.Default != "" {
		found := false
		for _, f := range d.Flags {
			found = found || f == d.Default
		}
		if !found {
			return fail("%s/%s: disposition default %q is not one of its flags", scene, n.ID, d.Default)
		}
	}
	return nil
}

// Find resolves a player-typed noun to an object. Matching is loose on
// purpose: the player should never lose a clue to vocabulary.
func (r *Room) Find(word string) *Object {
	var best *Object
	bestScore := 0
	for i := range r.Objects {
		o := &r.Objects[i]
		if s := match(word, o.ID, o.Names); s > bestScore {
			best, bestScore = o, s
		}
	}
	return best
}

// Score is how well a typed word names this object, zero for not at all. It
// is exported because the engine has to pick the best match among the objects
// whose gates are open, and a closed object is not in the room.
func (o *Object) Score(word string) int { return match(word, o.ID, o.Names) }

// S returns an authored string by key, loudly if it is missing.
func (r *Room) S(key string) string {
	if v, ok := r.Strings[key]; ok {
		return v
	}
	return "[missing string: " + key + "]"
}
