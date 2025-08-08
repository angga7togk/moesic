package main

import (
	"fmt"
	"log"
	"moesic/data"
	"strings"
	"time"

	vlc "github.com/adrg/libvlc-go/v3"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type playing struct {
	player      *vlc.Player
	moesic      data.Moesic
	duration    float64
	currentTime float64
}

type options struct {
	isRandomFlatMoesic bool
	isRandomSingle     bool
	isRandomPlaylist   bool
}

type model struct {
	current       playing
	options       options
	playlistIndex int
	playlists     []data.Moesic
}

type progressTickMsg struct{}

// Inisialisasi libVLC sekali saja
func init() {
	if err := vlc.Init("--no-xlib", "--quiet"); err != nil {
		log.Fatalf("libVLC init error: %v", err)
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

// Main fungsi untuk play lagu + handle event selesai
func (m *model) playSong(song data.Moesic) {
	ytb, err := GetAudio(song.Url)
	if err != nil {
		log.Println("Error get audio:", err)
		return
	}

	// buat player baru
	p, err := vlc.NewPlayer()
	if err != nil {
		log.Fatalf("new player error: %v", err)
	}

	media, err := p.LoadMediaFromURL(ytb.Url)
	if err != nil {
		log.Fatalf("load media error: %v", err)
	}
	defer media.Release()

	// attach event ketika lagu selesai
	em, _ := p.EventManager()
	_, _ = em.Attach(vlc.MediaPlayerEndReached, func(ev vlc.Event, ud interface{}) {
		m.onSongFinished()
	}, em)

	// update current song
	m.current = playing{
		player:      p,
		moesic:      song,
		duration:    ytb.Duration,
		currentTime: 0,
	}

	// mainkan
	if err := p.Play(); err != nil {
		log.Fatalf("play error: %v", err)
	}
}

// Fungsi yang dipanggil ketika lagu selesai
func (m *model) onSongFinished() {
	switch {
	case m.options.isRandomSingle:
		m.current.player.Stop()
		tea.Quit()

	case m.options.isRandomFlatMoesic:
		nextSong := data.GetRandomSong(flatSongs)
		m.playSong(nextSong)

	case m.options.isRandomPlaylist:
		m.playlistIndex++
		if m.playlistIndex >= len(m.playlists) {
			tea.Quit()
		} else {
			m.playSong(m.playlists[m.playlistIndex])
		}
	}
}

func initialModel(opts options) model {
	m := model{options: opts}

	if opts.isRandomSingle || opts.isRandomFlatMoesic {
		song := data.GetRandomSong(flatSongs)
		m.playSong(song)
	} else if opts.isRandomPlaylist {
		playlist := data.GetRandomPlaylist(playlists)
		m.playlists = playlist.Songs
		m.playlistIndex = 0
		m.playSong(m.playlists[m.playlistIndex])
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
			m.current.player.Stop()
			return m, tea.Quit
		case "p":
			m.current.player.TogglePause()
		case "s":
			m.current.player.Stop()
			m.onSongFinished()
		}

	case progressTickMsg:
		pos, err := m.current.player.MediaTime()
		if err == nil {
			m.current.currentTime = float64(pos) / 1000
		}
		return m, tickProgress()
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
		lipgloss.NewStyle().Render(fmt.Sprintf("%skip  %soause  %suit",
			lipgloss.NewStyle().Bold(true).Render("[S]"),
			lipgloss.NewStyle().Bold(true).Render("[P]"),
			lipgloss.NewStyle().Bold(true).Render("[Q]"),
		)),
	)

	return boxStyle.Render(info)
}
