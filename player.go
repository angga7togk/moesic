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

type next struct {
	moesic   data.Moesic
	url      string
	duration int64
}

type options struct {
	isPlayOne bool // just play one moesic
}

type playerModel struct {
	currentPlayer   *exec.Cmd
	currentMoesic   data.Moesic
	currentDuration int64
	next            *next
	loadingNext     bool
	options         options
}

type progressTickMsg struct{}
type fetchNextSongMsg struct {
	next *next
}

/*
* fetching next random song for flat moesic
 */
func fetchNextAsync() tea.Cmd {
	return func() tea.Msg {
		for {
			newNextSong := data.GetRandomSong(flatSongs)
			newNextYtb, err := GetAudio(newNextSong.Url)
			if err == nil {
				return fetchNextSongMsg{
					next: &next{
						moesic:   newNextSong,
						url:      newNextYtb.Url,
						duration: newNextYtb.Duration,
					},
				}
			}
			// sleep 1 detik biar ga spam hehe
			time.Sleep(1 * time.Second)
		}
	}
}

func tickProgress() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return progressTickMsg{}
	})
}

func drawCustomProgressBar(percent int, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := int(float64(percent) / 100 * float64(width))
	empty := width - filled

	return fmt.Sprintf("[%s%s]", strings.Repeat("█", filled), strings.Repeat("░", empty))
}

func initialModel(options options) playerModel {
	var (
		ytb  *Youtube
		err  error
		m    = playerModel{options: options}
		song data.Moesic
	)

	for {
		song = data.GetRandomSong(flatSongs)
		ytb, err = GetAudio(song.Url)
		if err == nil {
			m.currentMoesic = song
			m.currentDuration = ytb.Duration
			break
		}

	}

	m.currentPlayer = play(ytb.Url)
	return m
}

func (m playerModel) Init() tea.Cmd {
	return tickProgress()
}

func (m playerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.currentPlayer.Process.Kill()
			return m, tea.Quit
		case "s":
			if m.options.isPlayOne {
				m.currentPlayer.Process.Kill()
				return m, tea.Quit
			}
			if m.next != nil {
				m.currentPlayer.Process.Kill()
				cmd := play(m.next.url)
				m.currentPlayer = cmd
				m.currentMoesic = m.next.moesic
				m.currentDuration = m.next.duration
				m.next = nil
			} else {
				m.currentPlayer.Process.Kill()
				for {
					song := data.GetRandomSong(flatSongs)
					ytb, err := GetAudio(song.Url)
					if err == nil {
						cmd := play(ytb.Url)
						m.currentPlayer = cmd
						m.currentMoesic = song
						m.currentDuration = ytb.Duration
						break
					}
				}
			}
		}

	case progressTickMsg:

		// * song completed
		if globalPlayerTime > m.currentDuration {
			m.currentPlayer.Process.Kill()

			switch {
			case m.options.isPlayOne:
				return m, tea.Quit
			default:
				if m.next != nil {
					cmd := play(m.next.url)
					m.currentPlayer = cmd
					m.currentMoesic = m.next.moesic
					m.currentDuration = m.next.duration
					m.next = nil
				} else {
					for {
						song := data.GetRandomSong(flatSongs)
						ytb, err := GetAudio(song.Url)
						if err == nil {
							cmd := play(ytb.Url)
							m.currentPlayer = cmd
							m.currentMoesic = song
							m.currentDuration = ytb.Duration
							break
						}
					}
				}
				return m, tickProgress()
			}
		}

		// * prepare next song
		if m.next == nil && globalPlayerTime > 10 && !m.loadingNext && !m.options.isPlayOne {
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

func (m playerModel) View() string {
	percent := int(float64(globalPlayerTime) / float64(m.currentDuration) * 100)
	if percent > 100 {
		percent = 100
	}

	boxStyle := lipgloss.NewStyle().
		Width(40).
		Padding(1).
		BorderStyle(lipgloss.NormalBorder())

	progressBar := drawCustomProgressBar(percent, 15)

	info := fmt.Sprintf(
		"%s\n%s\n\n%s %s/%s\n\n%s",
		lipgloss.NewStyle().Bold(true).Render(m.currentMoesic.Name),
		lipgloss.NewStyle().Italic(true).Render(m.currentMoesic.PlaylistName),
		progressBar,
		formatTime(globalPlayerTime),
		formatTime(m.currentDuration),
		lipgloss.NewStyle().Render(fmt.Sprintf("%skip S%surce %suit",
			lipgloss.NewStyle().Bold(true).Render("[S]"),
			lipgloss.NewStyle().Bold(true).Render("[o]"),
			lipgloss.NewStyle().Bold(true).Render("[Q]"),
		)),
	)

	return boxStyle.Render(info)
}
