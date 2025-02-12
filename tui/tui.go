package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luisya22/lazyffmpeg/video"
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
	s := "Lazy FFmpeg - Select a Task\n-----------------------\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\n-----------------------\n"
	s += "[Pres ↑↓ or (k,j) to navigate, 'Enter' or 'Space' to select, 'q' to quit]"

	return s
}

func convertVideo() tea.Msg {
	video.Convert("/home/luismatos/Downloads/file_example_MP4_1920_18MG.mp4", "./", "output", "mkv")

	return "Success"
}
