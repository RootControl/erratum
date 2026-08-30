package ui

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/RootControl/erratum/internal/game"
)

// Plain drives the game over plain stdin/stdout, with no Bubble Tea, no ANSI
// and no repainting.
//
// It exists because a TUI driven through a pty emits a new frame for every
// keystroke, and a transcript captured that way is every frame concatenated —
// which reads like commands being swallowed and is not. Anything scripted
// should use this instead.
func Plain(st *game.State, in io.Reader, out io.Writer) error {
	w := bufio.NewWriter(out)
	defer w.Flush()

	emit := func(s string) {
		for _, l := range strings.Split(s, "\n") {
			fmt.Fprintln(w, l)
		}
		fmt.Fprintln(w)
	}
	emit(st.Room.Title)
	emit(st.Room.Intro)

	sc := bufio.NewScanner(in)
	sc.Buffer(make([]byte, 0, 64*1024), 64*1024)
	for {
		fmt.Fprint(w, "> ")
		w.Flush()
		if !sc.Scan() {
			break
		}
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		fmt.Fprintln(w, line)

		lines, ticks := st.Do(line)
		for _, l := range lines {
			emit(l)
		}
		for i := 0; i < ticks; i++ {
			for _, ev := range st.Tick() {
				emit(ev)
			}
		}
		if b := bodyLine(st); b != "" {
			fmt.Fprintf(w, "[ %s · %s ]\n\n", st.TimeString(), b)
		} else {
			fmt.Fprintf(w, "[ %s ]\n\n", st.TimeString())
		}
		if st.Done {
			break
		}
	}
	return sc.Err()
}
