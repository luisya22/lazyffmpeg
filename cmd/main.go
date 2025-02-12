package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luisya22/lazyffmpeg/tui"
)

func main() {
	// video.Convert("/home/luismatos/Downloads/file_example_MP4_1920_18MG.mp4", "./", "output", "mkv")
	m, _ := tui.InitialModel()
	// m, _ := tui.NewConvertModel()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
