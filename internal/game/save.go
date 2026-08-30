package game

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/RootControl/erratum/internal/content"
)

// Saving is the one place in this game where losing a field would be
// indistinguishable from the thing the game is about.
//
// State carries about twenty-five unexported counters that are pure timing —
// the body's deadline, the actions since the requisition, how far into the
// Cerement's symptoms the passage has got, how long C17 has left. A snapshot
// that dropped any one of them would restore a session that is *almost*
// right: the body still leaves, but not when it should. That is the
// archive's own trick, played on the player by the engine, and nothing on
// screen would say so.
//
// So the snapshot lists every field by hand — no reflection, no struct tags
// deciding what survives — and TestSnapshotCarriesEveryField fails the moment
// State grows one this file has not been told about.
type snapshot struct {
	Version int `json:"version"`

	Room  string `json:"room"`
	Clock int    `json:"clock"`
	Date  string `json:"date"`

	Given  string `json:"given"`
	Family string `json:"family"`

	BodyPresent  bool `json:"body_present"`
	BodyDeadline int  `json:"body_deadline"`
	DriftPartial int  `json:"drift_partial"`

	Visited  map[string]bool  `json:"visited"`
	Held     map[string]Store `json:"held"`
	Examined map[string]int   `json:"examined"`
	Flags    map[string]bool  `json:"flags"`
	Ints     map[string]int   `json:"ints"`
	Carried  map[string]bool  `json:"carried"`
	Notebook []string         `json:"notebook"`
	Arrived  []string         `json:"arrived"`

	Met      map[string]bool `json:"met"`
	Asked    map[string]int  `json:"asked"`
	Topics   map[string]bool `json:"topics"`
	Layers   map[string]int  `json:"layers"`
	Pressure map[string]int  `json:"pressure"`
	Fired    map[string]bool `json:"fired"`
	Deflects map[string]int  `json:"deflects"`

	PrecisPending bool `json:"precis_pending"`
	PrecisDone    bool `json:"precis_done"`
	Actions       int  `json:"actions"`
	Done          bool `json:"done"`

	LyleQueried      int  `json:"lyle_queried"`
	LyleConfirmed    bool `json:"lyle_confirmed"`
	CanticaQueried   int  `json:"cantica_queried"`
	CanticaConfirmed bool `json:"cantica_confirmed"`
	TassWaited       int  `json:"tass_waited"`
	RhodWaited       int  `json:"rhod_waited"`
	SlipsWaited      int  `json:"slips_waited"`
	BezWaited        int  `json:"bez_waited"`
	TransitTicks     int  `json:"transit_ticks"`
	SymptomsSeen     int  `json:"symptoms_seen"`
	C17Waited        int  `json:"c17_waited"`
	C17Gone          bool `json:"c17_gone"`

	Requisitioned  bool `json:"requisitioned"`
	SinceReq       int  `json:"since_req"`
	FallbackFired  bool `json:"fallback_fired"`
	SeenCorrection bool `json:"seen_correction"`

	// The archive is the enemy's own store and it is mutable, so it travels
	// with the save exactly as it stands — including whatever it has already
	// done to the player's case-notes.
	Records   map[string]string `json:"records"`
	CaseNotes []string          `json:"case_notes"`
	Corrected bool              `json:"corrected"`
}

const saveVersion = 1

func (s *State) snapshot() snapshot {
	return snapshot{
		Version: saveVersion,

		Room:  s.Room.ID,
		Clock: s.Clock,
		Date:  s.date,

		Given:  s.Given,
		Family: s.Family,

		BodyPresent:  s.BodyPresent,
		BodyDeadline: s.bodyDeadline,
		DriftPartial: s.driftPartial,

		Visited:  s.Visited,
		Held:     s.Held,
		Examined: s.Examined,
		Flags:    s.Flags,
		Ints:     s.Ints,
		Carried:  s.Carried,
		Notebook: s.Notebook,
		Arrived:  s.arrived,

		Met:      s.Met,
		Asked:    s.Asked,
		Topics:   s.Topics,
		Layers:   s.Layers,
		Pressure: s.Pressure,
		Fired:    s.Fired,
		Deflects: s.deflects,

		PrecisPending: s.PrecisPending,
		PrecisDone:    s.precisDone,
		Actions:       s.actions,
		Done:          s.Done,

		LyleQueried:      s.lyleQueried,
		LyleConfirmed:    s.lyleConfirmed,
		CanticaQueried:   s.canticaQueried,
		CanticaConfirmed: s.canticaConfirmed,
		TassWaited:       s.tassWaited,
		RhodWaited:       s.rhodWaited,
		SlipsWaited:      s.slipsWaited,
		BezWaited:        s.bezWaited,
		TransitTicks:     s.transitTicks,
		SymptomsSeen:     s.symptomsSeen,
		C17Waited:        s.c17Waited,
		C17Gone:          s.c17Gone,

		Requisitioned:  s.requisitioned,
		SinceReq:       s.sinceReq,
		FallbackFired:  s.fallbackFired,
		SeenCorrection: s.seenCorrection,

		Records:   s.Archive.records,
		CaseNotes: s.Archive.caseNotes,
		Corrected: s.Archive.corrected,
	}
}

// restore rebuilds a session against a freshly loaded Game. The content is
// never saved — it is compiled in — so a save is only ever the player.
func restore(g *content.Game, snap snapshot) (*State, error) {
	if snap.Version != saveVersion {
		return nil, fmt.Errorf("session was written by a different build (version %d, this is %d)",
			snap.Version, saveVersion)
	}
	room := g.Scene(snap.Room)
	if room == nil {
		return nil, fmt.Errorf("session refers to a scene this build does not have: %q", snap.Room)
	}
	s := New(g)
	s.Room = room

	s.Clock = snap.Clock
	s.date = snap.Date
	if snap.Given != "" {
		s.Given, s.Family = snap.Given, snap.Family
	}

	s.BodyPresent = snap.BodyPresent
	s.bodyDeadline = snap.BodyDeadline
	s.driftPartial = snap.DriftPartial

	s.Visited = orEmptyBool(snap.Visited)
	s.Held = orEmptyStore(snap.Held)
	s.Examined = orEmptyInt(snap.Examined)
	s.Flags = orEmptyBool(snap.Flags)
	s.Ints = orEmptyInt(snap.Ints)
	s.Carried = orEmptyBool(snap.Carried)
	s.Notebook = snap.Notebook
	s.arrived = snap.Arrived

	s.Met = orEmptyBool(snap.Met)
	s.Asked = orEmptyInt(snap.Asked)
	s.Topics = orEmptyBool(snap.Topics)
	s.Layers = orEmptyInt(snap.Layers)
	s.Pressure = orEmptyInt(snap.Pressure)
	s.Fired = orEmptyBool(snap.Fired)
	s.deflects = orEmptyInt(snap.Deflects)

	s.PrecisPending = snap.PrecisPending
	s.precisDone = snap.PrecisDone
	s.actions = snap.Actions
	s.Done = snap.Done

	s.lyleQueried = snap.LyleQueried
	s.lyleConfirmed = snap.LyleConfirmed
	s.canticaQueried = snap.CanticaQueried
	s.canticaConfirmed = snap.CanticaConfirmed
	s.tassWaited = snap.TassWaited
	s.rhodWaited = snap.RhodWaited
	s.slipsWaited = snap.SlipsWaited
	s.bezWaited = snap.BezWaited
	s.transitTicks = snap.TransitTicks
	s.symptomsSeen = snap.SymptomsSeen
	s.c17Waited = snap.C17Waited
	s.c17Gone = snap.C17Gone

	s.requisitioned = snap.Requisitioned
	s.sinceReq = snap.SinceReq
	s.fallbackFired = snap.FallbackFired
	s.seenCorrection = snap.SeenCorrection

	s.Archive.records = orEmptyString(snap.Records)
	s.Archive.caseNotes = snap.CaseNotes
	s.Archive.corrected = snap.Corrected
	return s, nil
}

func orEmptyBool(m map[string]bool) map[string]bool {
	if m == nil {
		return map[string]bool{}
	}
	return m
}

func orEmptyInt(m map[string]int) map[string]int {
	if m == nil {
		return map[string]int{}
	}
	return m
}

func orEmptyStore(m map[string]Store) map[string]Store {
	if m == nil {
		return map[string]Store{}
	}
	return m
}

func orEmptyString(m map[string]string) map[string]string {
	if m == nil {
		return map[string]string{}
	}
	return m
}

// Thin aliases so the tests can read and write a session file without
// importing os alongside everything else.
var (
	osReadFile  = os.ReadFile
	osWriteFile = os.WriteFile
)

// SavePath is where a suspended session lives. ERRATUM_SAVE overrides it.
func SavePath() string {
	if p := os.Getenv("ERRATUM_SAVE"); p != "" {
		return p
	}
	return "erratum.session"
}

// Save writes the session. It is the player's own file and nothing in the
// game reads it back except LOAD.
func (s *State) Save(path string) error {
	b, err := json.MarshalIndent(s.snapshot(), "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

// Load reads a session over the top of this one. The Game is not replaced:
// the content is compiled in, so a save from this build always matches it.
func (s *State) Load(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var snap snapshot
	if err := json.Unmarshal(b, &snap); err != nil {
		return fmt.Errorf("session is not readable: %w", err)
	}
	next, err := restore(s.Game, snap)
	if err != nil {
		return err
	}
	*s = *next
	return nil
}
