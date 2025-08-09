package main

import (
	"fmt"
	"moesic/data"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type playlistSongNext struct {
	moesic   data.Moesic
	audioUrl *string // url player
}

type playlistModel struct {
	currentPlayer *exec.Cmd // player saat ini
	loadingNext   bool      // untuk loading fetching next biar ga looping trus
	playlistIndex int       // playlist index yang sedang di play
	playlists     []playlistSongNext
}

type fetchNextPlaylistSongMsg struct {
	index int
	next  *playlistSongNext
}

/*
* fetching next playlist song
 */
func fetchNextPlaylistSongAsync(nextIndex int, nextSong playlistSongNext) tea.Cmd {
	return func() tea.Msg {
		for {
			audioUrl, err := GetAudio(nextSong.moesic.Url)
			if err == nil {
				return fetchNextPlaylistSongMsg{
					index: nextIndex,
					next: &playlistSongNext{
						moesic:   nextSong.moesic,
						audioUrl: &audioUrl,
					},
				}
			}
			// sleep 1 detik biar ga spam hehe
			time.Sleep(1 * time.Second)
		}
	}
}

func initialPlayerPlaylistplaylistModel(pl *data.Playlist) playlistModel {
	var (
		audioUrl string
		err      error
		m        = playlistModel{}
		playlist data.Playlist
	)

	// * get one random playlist
	if pl == nil {
		playlist = data.GetRandomPlaylist(playlists)
	} else {
		playlist = *pl
	}
	m.playlistIndex = 0
	m.playlists = []playlistSongNext{}

	for _, pl := range playlist.Songs {
		m.playlists = append(m.playlists, playlistSongNext{
			moesic: pl,
		})
	}
	// * loop get audio ke yt-dlp
	for {
		playlistSongNext := m.playlists[m.playlistIndex]
		audioUrl, err = GetAudio(playlistSongNext.moesic.Url)
		if err == nil {
			m.playlists[m.playlistIndex].audioUrl = &audioUrl
			break
		}

	}

	m.currentPlayer = play(audioUrl)
	return m
}

func (m playlistModel) Init() tea.Cmd {
	return tickProgress()
}

func (m playlistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.currentPlayer.Process.Kill()
			return m, tea.Quit
		case "s":
			m.currentPlayer.Process.Kill()
			m.playlistIndex++
			if m.playlistIndex >= len(m.playlists) {
				return m, tea.Quit
			}

			// * kalau next song blom di fetching fetching dulu
			nextIndex := m.playlistIndex
			nextSong := m.playlists[nextIndex]
			if nextSong.audioUrl == nil {
				for {
					audioUrl, err := GetAudio(nextSong.moesic.Url)
					if err == nil {
						m.playlists[nextIndex].audioUrl = &audioUrl
						break
					}
				}
			}

			m.currentPlayer = play(*m.playlists[nextIndex].audioUrl)
			return m, tickProgress()
		}

	case progressTickMsg:
		// * song completed
		if globalCurrentTime >= globalCurrentDuration{
			m.currentPlayer.Process.Kill()
			m.playlistIndex++
			if m.playlistIndex >= len(m.playlists) {
				return m, tea.Quit
			}

			// * kalau next song blom di fetching fetching dulu
			nextIndex := m.playlistIndex
			nextSong := m.playlists[nextIndex]
			if nextSong.audioUrl == nil {
				for {
					audioUrl, err := GetAudio(nextSong.moesic.Url)
					if err == nil {
						m.playlists[nextIndex].audioUrl = &audioUrl
						break
					}
				}
			}

			m.currentPlayer = play(*m.playlists[nextIndex].audioUrl)
			return m, tickProgress()
		}

		// * prepare next song
		if globalCurrentTime > 5 && !m.loadingNext {
			m.loadingNext = true
			nextSong := m.playlists[m.playlistIndex+1]
			if nextSong.audioUrl == nil {
				return m, tea.Batch(tickProgress(), fetchNextPlaylistSongAsync(m.playlistIndex+1, nextSong))
			}
		}

		return m, tickProgress()

	case fetchNextPlaylistSongMsg:
		// * set model next song
		nextIndex := msg.index
		m.playlists[nextIndex] = *msg.next

		m.loadingNext = false
		return m, nil
	}

	return m, nil
}

func (m playlistModel) View() string {
	current := m.playlists[m.playlistIndex]
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
		lipgloss.NewStyle().Bold(true).Render(current.moesic.Name),
		lipgloss.NewStyle().Italic(true).Render(current.moesic.PlaylistName),
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
