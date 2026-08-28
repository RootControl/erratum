// Package ui is the cogitator. It is deliberately plain: the archive's edits
// must never be signalled by a colour change, a glitch effect or a shift in
// spacing, so the renderer has no vocabulary for "something is wrong".
package ui

import (
	"strings"

	"github.com/RootControl/erratum/internal/game"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	amber  = lipgloss.AdaptiveColor{Light: "#8a5a00", Dark: "#e0a53f"}
	dim    = lipgloss.AdaptiveColor{Light: "#6b6157", Dark: "#8a7f72"}
	body   = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#1c1a17", Dark: "#ddd6cc"})
	head   = lipgloss.NewStyle().Foreground(amber).Bold(true)
	faint  = lipgloss.NewStyle().Foreground(dim)
	prompt = lipgloss.NewStyle().Foreground(amber)
)

type Model struct {
	st    *game.State
	vp    viewport.Model
	in    textinput.Model
	lines []string
	ready bool
	w, h  int
}

func New(st *game.State) Model {
	ti := textinput.New()
	ti.Prompt = "> "
	ti.PromptStyle = prompt
	ti.Focus()
	ti.CharLimit = 120

	// Default to a sane geometry so the first frame renders even if the
	// terminal never reports its size.
	m := Model{st: st, in: ti, w: 80, h: 24, ready: true}
	m.vp = viewport.New(80, 20)
	m.push(st.Room.Title)
	m.push(st.Room.Intro)
	m.render()
	return m
}

func (m *Model) push(s string) {
	// The engine emits a scene title as ordinary text when the player travels,
	// so that anything driving the engine gets it. Styling it is this layer's
	// job, and it has to happen wherever the title arrives from.
	if s == m.st.Room.Title {
		s = head.Render(s)
	}
	for _, l := range strings.Split(s, "\n") {
		m.lines = append(m.lines, l)
	}
	m.lines = append(m.lines, "")
}

func (m Model) Init() tea.Cmd { return textinput.Blink }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		vh := msg.Height - 4
		if vh < 5 {
			vh = 5
		}
		if msg.Width > 0 {
			m.vp.Width = msg.Width
		}
		m.vp.Height = vh
		m.ready = true
		m.render()
		m.vp.GotoBottom()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			raw := strings.TrimSpace(m.in.Value())
			if raw == "" {
				break
			}
			m.in.Reset()
			m.push(prompt.Render("> " + raw))

			out, ticks := m.st.Do(raw)
			for _, o := range out {
				m.push(o)
			}
			// Some things cost more than one action. The lift is one of them.
			for i := 0; i < ticks; i++ {
				for _, ev := range m.st.Tick() {
					m.push(ev)
				}
			}
			m.render()
			m.vp.GotoBottom()
			if m.st.Done {
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.in, cmd = m.in.Update(msg)
	cmds = append(cmds, cmd)
	m.vp, cmd = m.vp.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) render() {
	w := m.w
	if w <= 0 {
		w = 80
	}
	wrap := lipgloss.NewStyle().Width(min(w-2, 92))
	var b strings.Builder
	for _, l := range m.lines {
		b.WriteString(wrap.Render(body.Render(l)))
		b.WriteString("\n")
	}
	m.vp.SetContent(b.String())
}

func (m Model) View() string {
	if !m.ready {
		return "…"
	}
	status := faint.Render(strings.Join([]string{
		m.st.TimeString(),
		bodyLine(m.st),
		"HELP for verbs",
	}, "   ·   "))
	return m.vp.View() + "\n" + status + "\n" + m.in.View()
}

func bodyLine(s *game.State) string {
	if s.BodyPresent {
		return "body present"
	}
	return "body removed"
}
