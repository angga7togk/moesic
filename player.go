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

type options struct {
	isRandomFlatMoesic bool // []{playlist, songs[]} -> songs[]
	isRandomSingle     bool // just play one moesic
	isRandomPlaylist   bool // play random playlist
}

type model struct {
	current     playing
	next        *next
	loadingNext bool
	options     options

	playlistIndex int
	playlists     []data.Moesic // if options isRandomPlaylist
}

type progressTickMsg struct{}
type fetchNextSongMsg struct {
	next *next
}

/*
* fetching next playlist song
*/
func fetchNextPlaylistSongAsync(nextSong data.Moesic) tea.Cmd {
	return func() tea.Msg {
		ytb, err := GetAudio(nextSong.Url)
		if err != nil {
			return nil
		}

		return fetchNextSongMsg{
			next: &next{
				moesic:   nextSong,
				url:      ytb.Url,
				duration: ytb.Duration,
			},
		}
	}
}

/*
* fetching next random song for flat moesic
 */
func fetchNextAsync() tea.Cmd {
	return func() tea.Msg {
		newNextSong := data.GetRandomSong(flatSongs)
		newNextYtb, err := GetAudio(newNextSong.Url)
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

func initialModel(options options) model {
	var (
		ytb  *Youtube
		err  error
		m    = model{options: options}
		song data.Moesic
	)

	if options.isRandomSingle || options.isRandomFlatMoesic {
		// * get one random moesic from all playlists
		song = data.GetRandomSong(flatSongs)
		ytb, err = GetAudio(song.Url)
		if err != nil {
			panic("Please try again :)")
		}
		m.current = playing{
			moesic:      song,
			duration:    ytb.Duration,
			currentTime: 0,
		}
	} else if options.isRandomPlaylist {
		// * get one random playlist
		playlist := data.GetRandomPlaylist(playlists)
		m.playlistIndex = 0
		m.playlists = playlist.Songs

		song = m.playlists[m.playlistIndex]
		ytb, err = GetAudio(song.Url)
		if err != nil {
			panic("Please try again :)")
		}
		m.current = playing{
			moesic:      song,
			duration:    ytb.Duration,
			currentTime: 0,
		}
	}

	m.current.player = play(ytb.Url)
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

		// * song completed
		if m.current.currentTime >= m.current.duration {
			m.current.player.Process.Kill()

			switch {
			case m.options.isRandomSingle:
				return m, tea.Quit

			case m.options.isRandomFlatMoesic:
				song := data.GetRandomSong(flatSongs)
				ytb, err := GetAudio(song.Url)
				if err != nil {
					fmt.Println("Error:", err)
					return m, tea.Quit
				}
				cmd := play(ytb.Url)
				m.current = playing{
					player:      cmd,
					moesic:      song,
					duration:    ytb.Duration,
					currentTime: 0,
				}
				return m, tickProgress()

			case m.options.isRandomPlaylist:
				m.playlistIndex++
				if m.playlistIndex >= len(m.playlists) {
					return m, tea.Quit
				}

				var (
					song data.Moesic
					url  string
					dur  float64
					cmd  *exec.Cmd
				)

				if m.next != nil {
					song = m.next.moesic
					url = m.next.url
					dur = m.next.duration
					m.next = nil
				} else {
					song = m.playlists[m.playlistIndex]
					ytb, err := GetAudio(song.Url)
					if err != nil {
						fmt.Println("Error:", err)
						return m, tea.Quit
					}
					url = ytb.Url
					dur = ytb.Duration
				}

				cmd = play(url)
				m.current = playing{
					player:      cmd,
					moesic:      song,
					duration:    dur,
					currentTime: 0,
				}
				return m, tickProgress()

			}
		}

		// * prepare next song
		if m.next == nil && int(m.current.currentTime) > 10 && !m.loadingNext {
			m.loadingNext = true

			if m.options.isRandomFlatMoesic {
				return m, tea.Batch(tickProgress(), fetchNextAsync())
			}

			if m.options.isRandomPlaylist && m.playlistIndex+1 < len(m.playlists) {
				nextSong := m.playlists[m.playlistIndex+1]
				return m, tea.Batch(tickProgress(), fetchNextPlaylistSongAsync(nextSong))
			}
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

	progressBar := drawCustomProgressBar(percent, 15)

	info := fmt.Sprintf(
		"%s\n%s\n\n%s %s/%s\n\n%s",
		lipgloss.NewStyle().Bold(true).Render(m.current.moesic.Name),
		lipgloss.NewStyle().Italic(true).Render(m.current.moesic.PlaylistName),
		progressBar,
		formatTime(m.current.currentTime),
		formatTime(m.current.duration),
		lipgloss.NewStyle().Render(fmt.Sprintf("%skip S%surce %suit",
			lipgloss.NewStyle().Bold(true).Render("[S]"),
			lipgloss.NewStyle().Bold(true).Render("[o]"),
			lipgloss.NewStyle().Bold(true).Render("[Q]"),
		)),
	)

	return boxStyle.Render(info)
}
