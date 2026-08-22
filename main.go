// Command erratum is a vertical slice of THE ERRATUM: Act I, scene 1.3-1.4
// and the Act I turn. It exists to prove four things are buildable —
// free-form examination, the two evidence stores, the flag engine, and an
// archive that edits the player's own case-notes without ever saying so.
package main

import (
	"fmt"
	"os"

	"github.com/RootControl/erratum/internal/content"
	"github.com/RootControl/erratum/internal/game"
	"github.com/RootControl/erratum/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	room, err := content.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	st := game.New(room)
	opts := []tea.ProgramOption{tea.WithMouseCellMotion()}
	if os.Getenv("ERRATUM_NOALT") == "" {
		opts = append(opts, tea.WithAltScreen())
	}
	p := tea.NewProgram(ui.New(st), opts...)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
