package main

import (
	"fmt"
	"moesic/data"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func printHelp() {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5F87")).
		PaddingBottom(1)

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Underline(true).
		Foreground(lipgloss.Color("#00D787"))

	optionNameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5FD7FF"))

	descriptionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AFAFAF"))

	// Output
	fmt.Println(titleStyle.Render(` 
 __  __  ___  ___ ___ ___ ___ 
|  \/  |/ _ \| __/ __|_ _/ __|
| |\/| | (_) | _|\__ \| | (__ 
|_|  |_|\___/|___|___/___\___|
                               `))
	fmt.Println("⭐️ Star to support our work!")
	fmt.Print(lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("#4157ba")).Render("   https://github.com/angga7togk/moesic"))
	fmt.Println()
	fmt.Println()
	fmt.Println(sectionStyle.Render("Usage:"))
	fmt.Println("  " + optionNameStyle.Render(fmt.Sprintf("%-30s", "moesic <options>")) + descriptionStyle.Render("Moesic CLI"))

	fmt.Println()
	fmt.Println(sectionStyle.Render("Options:"))
	fmt.Println("  " + optionNameStyle.Render(fmt.Sprintf("%-30s", "--random, --play, -p")) + descriptionStyle.Render("Play random flat moesic"))
	fmt.Println("  " + optionNameStyle.Render(fmt.Sprintf("%-30s", "--random-playlist, -rp")) + descriptionStyle.Render("Play random playlist"))
	fmt.Println("  " + optionNameStyle.Render(fmt.Sprintf("%-30s", "--random-single, -rs")) + descriptionStyle.Render("Play random single moesic"))
	fmt.Println("  " + optionNameStyle.Render(fmt.Sprintf("%-30s", "--help, -h")) + descriptionStyle.Render("Command help"))
	fmt.Println("  " + optionNameStyle.Render(fmt.Sprintf("%-30s", "--info, -i")) + descriptionStyle.Render("Moesic info"))
	fmt.Print("\n\n")
}

var (
	playlists        []data.Playlist = []data.Playlist{}
	flatSongs        []data.Moesic   = []data.Moesic{}
	globalPlayerTime float64         = 0 // * global progres music player
)

func main() {
	InstallDependencies()

	playlists = data.GetMoesic()
	flatSongs = data.FlatSongs(playlists)

	args := os.Args

	if len(args) < 2 {
		printHelp()
		return
	}

	command := args[1]
	switch command {
	case "--random", "--play", "-p":
		p := tea.NewProgram(initialModel(options{
			isRandomFlatMoesic: true,
		}), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	case "--random-single", "-rs":
		p := tea.NewProgram(initialModel(options{
			isRandomSingle: true,
		}), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	case "--random-playlist", "-rp":
		p := tea.NewProgram(initialModel(options{
			isRandomPlaylist: true,
		}), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	default:
		printHelp()
	}
}
