package tui

import (
	"fmt"
	"log"

	textInput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func PromptUser(key string) string {
	ti := textInput.New()
	ti.Placeholder = key
	ti.Focus()

	p := tea.NewProgram(model{input: ti, key: key})

	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	return m.(model).value
}

type model struct {
	input textInput.Model
	key   string
	value string
}

func (m model) Init() tea.Cmd {
	return textInput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			m.value = m.input.Value()
			return m, tea.Quit
		}

		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("🔑 Enter value for %s:\n\n%s\n\n[Enter] → Next  |  [Ctrl+C] → Exit", m.key, m.input.View())
}
