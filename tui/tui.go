package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices []string
	cursor  int
}

func InitialModel() (tea.Model, tea.Cmd) {
	return model{
		choices: []string{
			"Convert",
			"Extract Audio",
			"Trim Video",
			"Resize Video",
			"Generate GIF",
			"Optimize Video",
			"Extract Frames",
		},
	}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			if m.choices[m.cursor] == "Convert" {
				return NewConvertModel()
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder
	s.WriteString("Lazy FFmpeg - Select a Task\n-----------------------\n")

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}

	s.WriteString("\n-----------------------\n")
	s.WriteString("[Press ↑↓ or (k,j) to navigate, 'Enter' or 'Space' to select, 'q' to quit]")

	return s.String()
}
