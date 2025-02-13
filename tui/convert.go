package tui

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luisya22/lazyffmpeg/video"
)

type convertModel struct {
	stage              int
	filepicker         filepicker.Model
	quitting           bool
	selectedFile       string
	selectedOutputPath string
	formatChoices      []string
	selectedFormat     string
	cursor             int
	textInput          textinput.Model
	spinner            spinner.Model
	err                error
}

type errMsg struct{ err error }

type videoProcessed struct{}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func NewConvertModel() (tea.Model, tea.Cmd) {
	fp := filepicker.New()
	fp.AllowedTypes = []string{
		".3gp", ".3g2", ".asf", ".avi", ".flv", ".swf", ".mp4", ".m4v", ".mkv",
		".webm", ".mxf", ".gxf", ".ogg", ".ogv", ".mov", ".rm", ".rmvb", ".vob",
		".mpeg", ".mpg", ".mpe", ".ts", ".m2ts", ".dv", ".nut", ".f4v", ".yuv",
		".264", ".265", ".h264", ".h265", ".mp4v", ".wmv", ".mts", ".m2v",
	}
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.Height = 20

	ti := textinput.New()
	ti.Placeholder = "Filename"
	ti.Width = 30

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := convertModel{
		filepicker: fp,
		stage:      0, formatChoices: []string{"mp4", "webm", "mkv", "mov", "avi", "flv", "wmv",
			"mpg", "mpeg", "ts",
		},
		textInput: ti,
		spinner:   s,
	}

	return m, m.Init()
}

func (m convertModel) Init() tea.Cmd {

	return tea.Batch(
		m.filepicker.Init(),
	)
}
func (m convertModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return InitialModel()
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	case videoProcessed:
		m.stage = 5
		return m, tea.Tick(1*time.Second, func(time.Time) tea.Msg {
			return tea.Quit()
		})
	}

	var cmd tea.Cmd
	if m.stage == 0 {

		m.filepicker, cmd = m.filepicker.Update(msg)

		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
			m.selectedFile = path
			m.stage = 1
			m.filepicker.DirAllowed = true
			m.filepicker.FileAllowed = false
			m.filepicker.CurrentDirectory, _ = os.UserHomeDir()
			m.filepicker.AllowedTypes = nil

			return m, m.filepicker.Init()
		}

		if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
			m.err = errors.New(path + " is not valid.")
			m.selectedFile = ""

			return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
		}
	} else if m.stage == 1 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeySpace:
				m.stage = 2
				return m, nil
			}
		}
		m.filepicker, cmd = m.filepicker.Update(msg)

		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
			m.selectedOutputPath = path
		}

		if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
			m.err = errors.New(path + " is not valid.")
			m.selectedFile = ""

			return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
		}

	} else if m.stage == 2 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.formatChoices)-1 {
					m.cursor++
				}
			case tea.KeyEnter.String(), tea.KeySpace.String():
				if m.stage == 2 {
					m.stage = 3
					m.selectedFormat = m.formatChoices[m.cursor]
					m.textInput.Cursor.Blink = true

					return m, m.textInput.Focus()

				}
			}

		}
	} else if m.stage == 3 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
				filename := m.textInput.Value()

				m.stage = 4
				return m, tea.Batch(m.spinner.Tick, convertVideo(m.selectedFile, m.selectedOutputPath, filename, m.selectedFormat))
			}

			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	} else if m.stage == 4 {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func convertVideo(selectedFile string, selectedOutputPath string, filename string, selectedFormat string) tea.Cmd {

	return func() tea.Msg {
		err := video.Convert(selectedFile, selectedOutputPath, filename, selectedFormat)
		if err != nil {
			return errMsg{err}
		}

		return videoProcessed{}
	}
}

func (m convertModel) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder
	s.WriteString("\n")

	if m.stage == 0 {
		if m.err != nil {
			s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
		} else if m.selectedFile == "" {
			s.WriteString("Pick a file:")
		} else {
			s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
		}

		s.WriteString("\n\n" + m.filepicker.View() + "\n")
		s.WriteString("[Press ↑↓ or (k,j) to navigate, 'Enter' to select, 'q' to quit]")
	} else if m.stage == 1 {
		if m.err != nil {
			s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
		} else if m.selectedFile == "" {
			s.WriteString("Pick an output folder:")
		} else {
			s.WriteString("Selected Folder: " + m.filepicker.Styles.Selected.Render(m.selectedOutputPath))
		}
		m.stage = 1

		s.WriteString("\n\n" + m.filepicker.View() + "\n")
		s.WriteString("[Press ↑↓ or (k,j) to navigate, 'Enter' to select, 'q' to quit, after selecting folder press 'Space' to confirm and continue]")
	} else if m.stage == 2 {
		s.WriteString("Select an output format:\n\n")
		for i, c := range m.formatChoices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			s.WriteString(fmt.Sprintf("%s %s\n", cursor, c))
		}

		s.WriteString("\n-----------------------\n")
		s.WriteString("[Pres ↑↓ or (k,j) to navigate, 'Enter' or 'Space' to select, 'q' to quit]")
	} else if m.stage == 3 {
		s.WriteString("Enter file name (output) ")
		s.WriteString(m.textInput.View())
	} else if m.stage == 4 {
		s.WriteString(fmt.Sprintf("%s Working on your video\n", m.spinner.View()))
	} else if m.stage == 5 {
		s.WriteString("✅ Your video was successfully processed!\n")
	}

	return s.String()
}

// TODO: Styles
// TODO: progressbar or waiting spinner at the end. Controls when selecting format
// TODO: Handle and print errors correctly

// TODO: Design all the pages

/**
State manages what it shows
S1 - Pick a file (file picker) - DONE
S2 - Pick output Folder
S3 - Input output filename or use default
S4 - Pick format
**/
