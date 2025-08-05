package main

import (
	"fmt"
	"moesic/data"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	player      *exec.Cmd
	moesic      data.Moesic
	duration    float64
	currentTime float64
	startTime   time.Time
}


type progressTickMsg struct{}

func tickProgress() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return progressTickMsg{}
	})
}

func drawCustomProgressBar(percent float64, width int) string {
	filled := int(percent * float64(width))
	empty := width - filled
	return fmt.Sprintf("[%s%s]", strings.Repeat("█", filled), strings.Repeat("░", empty))
}

func initialModel(player *exec.Cmd, moesic data.Moesic, duration float64) model {
	m := model{
		player:      player,
		moesic:      moesic,
		duration:    duration,
		currentTime: 0,
		startTime:   time.Now(),
	}

	return m

}

func (m model) Init() tea.Cmd {
	return tickProgress()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}

	case progressTickMsg:
		m.currentTime = currentTime
		return m, tickProgress()
	}

	return m, nil
}

func (m model) View() string {
	percent := m.currentTime / m.duration
	if percent > 1 {
		percent = 1
	}

	progressBar := drawCustomProgressBar(percent, 30)

	info := fmt.Sprintf("%s\n\n%s %s / %s\n\n[S]kip  [P]ause  S[o]urce  [Q]uit",
		m.moesic.Name,
		progressBar,
		formatTime(m.currentTime),
		formatTime(m.duration),
	)

	boxStyle := lipgloss.NewStyle().
		Width(50).
		Padding(1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))

	return boxStyle.Render(info)
}
