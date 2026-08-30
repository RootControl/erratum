// Command erratum is THE ERRATUM: all three acts, thirty-one scenes, and five
// endings. Free-form examination, two evidence stores the game never explains,
// free-form dialogue, a flat flag engine, and an archive that edits the
// player's own case-notes without ever saying so.
//
// ERRATUM_NOALT=1 disables the alternate screen. ERRATUM_PLAIN=1 skips the TUI
// entirely and plays over stdin/stdout, which is what anything scripted should
// use.
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
	g, err := content.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	st := game.New(g)
	if os.Getenv("ERRATUM_PLAIN") != "" {
		// Scripted play: no TUI, no repaints, one block of text per command.
		if err := ui.Plain(st, os.Stdin, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
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
