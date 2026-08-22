package game

import "fmt"

// Archive is the Lexicanum Vault and the player's case-notes: everything the
// enemy can reach. It is a real mutable store, not a scripted cutscene. When
// the Correctorum corrects, it calls Correct() and the data actually changes.
//
// The player's notebook is deliberately NOT in here. That is the whole game.
type Archive struct {
	records   map[string]string
	caseNotes []string
	corrected bool
}

const (
	CodeLive    = "C9-F4-D2" // Choir of Cantor hardcopy dispatch manifests
	CodeForward = "D2-F4-C9" // Munitorum drainage requisition ledgers
)

func NewArchive() *Archive {
	a := &Archive{records: map[string]string{}}
	a.records[CodeLive] = dispatchRecord
	a.records[CodeForward] = drainageRecord
	return a
}

// Query returns what the Vault says today. Not what it said an hour ago.
func (a *Archive) Query(code string) (string, bool) {
	r, ok := a.records[code]
	return r, ok
}

// File appends a line to the case-notes. Anything in here is editable.
func (a *Archive) File(line string) {
	a.caseNotes = append(a.caseNotes, line)
}

func (a *Archive) CaseNotes() []string { return a.caseNotes }

func (a *Archive) Corrected() bool { return a.corrected }

// Correct performs the reconciliation.
//
// It does not delete the player's work. Deletion would be noticeable, and
// noticeable is the one thing the Amanuensis is not. Instead it substitutes
// the plausible mistake a novice would have made — the exact wrong code the
// player typed, or nearly typed, an hour ago. Same timestamp, same session,
// same author. It edits you into being incompetent, in your own handwriting.
func (a *Archive) Correct() {
	if a.corrected {
		return
	}
	a.corrected = true

	delete(a.records, CodeLive)

	for i, line := range a.caseNotes {
		if len(line) > 7 && line[:7] == "REQ 001" {
			a.caseNotes[i] = fmt.Sprintf(
				"REQ 001 — %s — SANITATION REQUISITION LEDGERS (MUNITORUM LIAISON)\n"+
					"         — RETURNED — NO ACTION — CLOSED", CodeForward)
		}
	}
}

const dispatchRecord = `LEXICANUM VAULT :: BLOCK C / LEVEL 9 :: AISLE F4 :: SHELF D2
QUERY: DISPATCH MANIFEST, CHOIR OF CANTOR, CYCLE 0997.4
AUTH: VOSK, H. [REASSIGNED — INQ. ROSETTE 4471-K]
RETURNING 1 RECORD.
--
DISPATCH 0997.4-1188
  ORIGIN ......... RELAY THESSALINE-SECUNDUS
  TRANSIT AUTH ... VERMILION / CONCORDANCE PRIORITY (SCRUB-ON-CARRY)
  RECEIVED ....... 0997.4.211 / TWELFTH HOUR, 11 MIN
  CHOIR-VESSEL ... TALM, V. I. — CELL 4
  LENGTH ......... 41 SEC
  CLASS .......... ROUTINE / TITHE VARIANCE ADVISORY
  ADDRESSEE ...... VOSK, H., INQUISITOR — ORDO: (UNDECLARED)
  DUTY ARCHIVIST . KESH, A., LORD-PREFECT — COUNTERSIGNED IN PERSON
--
END OF RECORD. 1 OF 1.`

const drainageRecord = `LEXICANUM VAULT :: BLOCK D / LEVEL 2 :: AISLE F4 :: SHELF C9
QUERY: SANITATION REQUISITION LEDGERS (MUNITORUM LIAISON)
RETURNING 214,880 RECORDS. NARROW QUERY.`

const dispatchGone = `LEXICANUM VAULT :: BLOCK C / LEVEL 9 :: AISLE F4 :: SHELF D2
QUERY: DISPATCH MANIFEST, CHOIR OF CANTOR, CYCLE 0997.4
RETURNING 0 RECORDS.
NO RECORD MATCHES THIS FORM. CHECK QUERY.`
