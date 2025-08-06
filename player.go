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

type playing struct {
	player      *exec.Cmd
	moesic      data.Moesic
	duration    float64
	currentTime float64
}

type next struct {
	moesic   data.Moesic
	url      string
	duration float64
}

type model struct {
	current     playing
	next        *next
	loadingNext bool
	single      bool
}

type progressTickMsg struct{}
type fetchNextSongMsg struct {
	next *next
}

func fetchNextAsync() tea.Cmd {
	return func() tea.Msg {
		newNextSong := data.GetRandomSong(flatSongs)
		newNextYtb, err := GetYoutube(newNextSong.Url)
		if err != nil {
			return nil
		}

		return fetchNextSongMsg{
			next: &next{
				moesic:   newNextSong,
				url:      newNextYtb.Url,
				duration: newNextYtb.Duration,
			},
		}
	}
}

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

func initialModel(single bool) model {
	randomSong := data.GetRandomSong(flatSongs)
	ytb, err := GetYoutube(randomSong.Url)
	if err != nil {
		fmt.Println(err)
	}
	cmd := play(ytb.Url)

	m := model{
		current: playing{
			player:      cmd,
			moesic:      randomSong,
			duration:    ytb.Duration,
			currentTime: 0,
		},
		single: single,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return tickProgress()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.current.player.Process.Kill()
			return m, tea.Quit
		case "s":
			if m.next != nil {
				m.current.player.Process.Kill()
				cmd := play(m.next.url)

				m.current = playing{
					player:      cmd,
					moesic:      m.next.moesic,
					duration:    m.next.duration,
					currentTime: 0,
				}
				m.next = nil
			}

		}

	case progressTickMsg:
		m.current.currentTime = globalPlayerTime

		if m.current.currentTime >= m.current.duration {
			if m.single {
				m.current.player.Process.Kill()
				return m, tea.Quit
			} else {
				if m.next != nil {
					cmd := play(m.next.url)

					m.current = playing{
						player:      cmd,
						moesic:      m.next.moesic,
						duration:    m.next.duration,
						currentTime: 0,
					}
					m.next = nil
				}
			}
		}

		if m.next == nil && int(m.current.currentTime) > 10 && !m.loadingNext && !m.single {
			m.loadingNext = true
			return m, tea.Batch(tickProgress(), fetchNextAsync())
		}

		return m, tickProgress()

	case fetchNextSongMsg:
		m.next = msg.next
		m.loadingNext = false
		return m, nil

	}

	return m, nil
}

func (m model) View() string {
	percent := m.current.currentTime / m.current.duration
	if percent > 1 {
		percent = 1
	}

	boxStyle := lipgloss.NewStyle().
		Width(40).
		Padding(1).
		BorderStyle(lipgloss.NormalBorder())

	// Components
	progressBar := drawCustomProgressBar(percent, 15)
	nextSongInfo := "Next: loading...\n\n"
	if m.next != nil {
		nextSongInfo = fmt.Sprintf("Next: %s\n\n", m.next.moesic.Name)
	}
	if m.single {
		nextSongInfo = ""
	}
	info := fmt.Sprintf(
		"%s\n%s\n\n%s %s/%s\n\n%s%s",
		lipgloss.NewStyle().Bold(true).Render(m.current.moesic.Name),
		lipgloss.NewStyle().Italic(true).Render(m.current.moesic.PlaylistName),
		progressBar,
		formatTime(m.current.currentTime),
		formatTime(m.current.duration),
		lipgloss.NewStyle().Italic(true).Render(nextSongInfo),
		lipgloss.NewStyle().Width(38).Align(lipgloss.Right).Render(fmt.Sprintf("%skip S%surce %suit",
			lipgloss.NewStyle().Bold(true).Render("[S]"),
			lipgloss.NewStyle().Bold(true).Render("[o]"),
			lipgloss.NewStyle().Bold(true).Render("[Q]"),
		)),
	)

	return boxStyle.Render(info)
}
