package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	player *exec.Cmd
	songName  string
	duration  float64
	startTime time.Time
	quitting  bool
}

type tickMsg time.Time

func drawCustomProgressBar(percent float64, width int) string {
	filled := int(percent * float64(width))
	empty := width - filled
	return fmt.Sprintf("[%s%s]", strings.Repeat("█", filled), strings.Repeat("░", empty))
}

func initialModel(player *exec.Cmd, songName string, duration float64) model {
	return model{
		player: player,
		songName:  songName,
		duration:  duration,
		startTime: time.Now(),
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.quitting {
			return m, tea.Quit
		}
		elapsed := time.Since(m.startTime).Seconds()
		if elapsed >= m.duration {
			m.quitting = true
			return m, tea.Quit
		}
		return m, tick()

	case tea.KeyMsg:
		if msg.String() == "q" {
			m.quitting = true
			m.player.Cancel()
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	elapsed := time.Since(m.startTime).Seconds()
	percent := elapsed / m.duration
	if percent > 1 {
		percent = 1
	}

	progressBar := drawCustomProgressBar(percent, 20)

	info := fmt.Sprintf("%s\n\n%s %s / %s\n\n[S]kip  [P]ause  S[o]urce  [Q]uite",
		m.songName,
		progressBar,
		formatTime(elapsed),
		formatTime(m.duration),
	)

	boxStyle := lipgloss.NewStyle().
		Width(50).
		Padding(1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))

	return boxStyle.Render(info)
}

func main() {
	song := GetRandomSong(FlatSongs(getMoesic()))
	duration, err := getDuration(song.Url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	cmd := play(song.Url)

	p := tea.NewProgram(initialModel(cmd, song.Name, duration))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	cmd.Wait()
	fmt.Println("\nDone.")
}
