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
	moesic data.Moesic
	audioUrl    string
}

type options struct {
	isPlayOne bool // just play one moesic
	moesic    *data.Moesic
}

type playerModel struct {
	currentPlayer *exec.Cmd
	currentMoesic data.Moesic
	next          *next
	loadingNext   bool
	options       options
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
			audioUrl, err := GetAudio(newNextSong.Url)
			if err == nil {
				return fetchNextSongMsg{
					next: &next{
						moesic: newNextSong,
						audioUrl:    audioUrl,
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
		audioUrl string
		err      error
		m        = playerModel{options: options}
		song     data.Moesic
	)

	for {
		var url string
		if options.moesic == nil {
			song = data.GetRandomSong(flatSongs)
			url = song.Url
		} else {
			song = *options.moesic
			url = song.Url
		}

		audioUrl, err = GetAudio(url)
		if err == nil {
			m.currentMoesic = song
			break
		}

	}

	m.currentPlayer = play(audioUrl)
	return m
}

func (m playerModel) Init() tea.Cmd {
	return tickProgress()
}

func (m playerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.currentPlayer.Process.Kill()
			return m, tea.Quit
		case "s":
			if m.options.isPlayOne {
				m.currentPlayer.Process.Kill()
				return m, tea.Quit
			}
			if m.next != nil {
				m.currentPlayer.Process.Kill()
				cmd := play(m.next.audioUrl)
				m.currentPlayer = cmd
				m.currentMoesic = m.next.moesic
				m.next = nil
			} else {
				m.currentPlayer.Process.Kill()
				for {
					song := data.GetRandomSong(flatSongs)
					audioUrl, err := GetAudio(song.Url)
					if err == nil {
						cmd := play(audioUrl)
						m.currentPlayer = cmd
						m.currentMoesic = song
						break
					}
				}
			}
		}

	case progressTickMsg:

		// * song completed
		if globalCurrentTime > globalCurrentDuration {
			m.currentPlayer.Process.Kill()

			switch {
			case m.options.isPlayOne:
				return m, tea.Quit
			default:
				if m.next != nil {
					cmd := play(m.next.audioUrl)
					m.currentPlayer = cmd
					m.currentMoesic = m.next.moesic
					m.next = nil
				} else {
					for {
						song := data.GetRandomSong(flatSongs)
						audioUrl, err := GetAudio(song.Url)
						if err == nil {
							cmd := play(audioUrl)
							m.currentPlayer = cmd
							m.currentMoesic = song
							break
						}
					}
				}
				return m, tickProgress()
			}
		}

		// * prepare next song
		if m.next == nil && !m.loadingNext && !m.options.isPlayOne {
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
	percent := int(float64(globalCurrentTime) / float64(globalCurrentDuration) * 100)
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
		formatTime(globalCurrentTime),
		formatTime(globalCurrentDuration),
		lipgloss.NewStyle().Render(fmt.Sprintf("%skip S%surce %suit",
			lipgloss.NewStyle().Bold(true).Render("[S]"),
			lipgloss.NewStyle().Bold(true).Render("[o]"),
			lipgloss.NewStyle().Bold(true).Render("[Q]"),
		)),
	)

	return boxStyle.Render(info)
}
