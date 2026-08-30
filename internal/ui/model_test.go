package ui

import (
	"strings"
	"testing"

	"github.com/RootControl/erratum/internal/content"
	"github.com/RootControl/erratum/internal/game"
	tea "github.com/charmbracelet/bubbletea"
)

func quits(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	msg := cmd()
	if batch, ok := msg.(tea.BatchMsg); ok {
		for _, c := range batch {
			if quits(c) {
				return true
			}
		}
		return false
	}
	_, ok := msg.(tea.QuitMsg)
	return ok
}

func typed(m Model, s string) (Model, tea.Cmd) {
	for _, r := range s {
		next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = next.(Model)
	}
	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	return next.(Model), cmd
}

func newModel(t *testing.T) (Model, *game.State) {
	t.Helper()
	g, err := content.Load()
	if err != nil {
		t.Fatal(err)
	}
	st := game.New(g)
	return New(st), st
}

// The bug this file exists for. The renderer runs on the alternate screen, so
// quitting the instant the game ended threw away the ending — a thousand words
// the player never got to read, restored out of existence by the terminal.
func TestAnEndingHoldsTheScreen(t *testing.T) {
	m, st := newModel(t)
	if !st.Move("intake") {
		t.Fatal("cannot reach the intake")
	}
	for _, f := range []string{"ACT2_OPEN", "CONFESSION_HEARD", "CHAIR_OFFERED"} {
		st.Set(f)
	}
	m, cmd := typed(m, "sit down")
	if st.Ending() == "" {
		t.Fatal("setup: the chair is always on offer")
	}
	if quits(cmd) {
		t.Fatal("the game closed on top of its own ending")
	}
	if !m.finished {
		t.Fatal("expected the finished state")
	}

	view := m.View()
	if !strings.Contains(view, "ENTER to close the file") {
		t.Fatalf("the footer has to say how to leave:\n%s", view)
	}
	if strings.Contains(view, "HELP for verbs") {
		t.Fatal("there is nothing left to do")
	}
	// And the ending starts at the top of the screen, not the bottom, because
	// it runs past a page and the player is going to scroll.
	if m.vp.YOffset != m.endedAt {
		t.Fatalf("expected the ending at the top: offset %d, ending at %d",
			m.vp.YOffset, m.endedAt)
	}

	// It closes when the player says so, and not before.
	_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !quits(cmd) {
		t.Fatal("ENTER should close the file")
	}
}

// Quitting is a player leaving and still ends the program at once.
func TestQuittingStillQuits(t *testing.T) {
	m, _ := newModel(t)
	m, cmd := typed(m, "quit")
	if !quits(cmd) {
		t.Fatal("QUIT should quit")
	}
	if m.finished {
		t.Fatal("leaving is not an ending")
	}
}

// A finished game takes no more commands, and the keys still scroll.
func TestAFinishedGameTakesNoInput(t *testing.T) {
	m, st := newModel(t)
	if !st.Move("intake") {
		t.Fatal("cannot reach the intake")
	}
	for _, f := range []string{"ACT2_OPEN", "CONFESSION_HEARD", "CHAIR_OFFERED"} {
		st.Set(f)
	}
	m, _ = typed(m, "sit down")

	before := len(m.lines)
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("look")})
	m = next.(Model)
	if len(m.lines) != before {
		t.Fatal("a finished game does not take commands")
	}
	if m.in.Focused() {
		t.Fatal("the prompt is gone")
	}
}
