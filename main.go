package main

import (
	"fmt"
	"moesic/data"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func headerString() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5F87")).
		PaddingBottom(1)

	return fmt.Sprintf("%s\n%s\n",
		titleStyle.Render(` 
 __  __  ___  ___ ___ ___ ___ 
|  \/  |/ _ \| __/ __|_ _/ __|
| |\/| | (_) | _|\__ \| | (__ 
|_|  |_|\___/|___|___/___\___|
                               `),
		"⭐️ Star to support our work!")
}

func printHelp() {
	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Underline(true).
		Foreground(lipgloss.Color("#00D787"))

	optionNameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5FD7FF"))

	descriptionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AFAFAF"))

	fmt.Print(headerString())
	fmt.Print(lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("#4157ba")).Render("   https://github.com/angga7togk/moesic"))
	fmt.Println()
	fmt.Println()
	fmt.Println(sectionStyle.Render("Usage:"))
	fmt.Println("  " + optionNameStyle.Render(fmt.Sprintf("%-30s", "moesic play")) + descriptionStyle.Render("Play flat moesic"))
	fmt.Println("  " + optionNameStyle.Render(fmt.Sprintf("%-30s", "moesic playlist")) + descriptionStyle.Render("Play playlist"))
	fmt.Println("  " + optionNameStyle.Render(fmt.Sprintf("%-30s", "moesic update")) + descriptionStyle.Render("Update dependencies"))
	fmt.Println("  " + optionNameStyle.Render(fmt.Sprintf("%-30s", "moesic version")) + descriptionStyle.Render("Moesic version"))
	fmt.Println("  " + optionNameStyle.Render(fmt.Sprintf("%-30s", "moesic help")) + descriptionStyle.Render("Moesic Help"))

	fmt.Println()
	fmt.Println(sectionStyle.Render("Options:"))
	fmt.Println("  " + optionNameStyle.Render(fmt.Sprintf("%-30s", "--random, -r")) + descriptionStyle.Render("Random options"))
	fmt.Println("  " + optionNameStyle.Render(fmt.Sprintf("%-30s", "--one, -o")) + descriptionStyle.Render("Just play one moesic not next or skiped"))
	fmt.Println("  " + optionNameStyle.Render(fmt.Sprintf("%-30s", "--test")) + descriptionStyle.Render("Play Megumin's explosion (for Test)"))
	fmt.Print("\n\n")
}

var (
	playlists []data.Playlist = []data.Playlist{}
	flatSongs []data.Moesic   = []data.Moesic{}
	version                   = "1.0.3"
)

func main() {
	InstallDependencies(false)

	playlists = data.GetMoesic()
	flatSongs = data.FlatSongs(playlists)

	args := os.Args

	if len(args) < 2 {
		printHelp()
		return
	}

	switch args[1] {
	case "playlist":
		if len(args) > 2 {
			if argsHas("--random", "-r") {
				p := tea.NewProgram(initialPlayerPlaylistplaylistModel(nil), tea.WithAltScreen())
				if _, err := p.Run(); err != nil {
					fmt.Printf("Alas, there's been an error: %v", err)
					os.Exit(1)
				}
			} else {
				printHelp()
			}
		} else {
			if _, err := tea.NewProgram(initSearchModel()).Run(); err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}
		}
	case "play":
		if len(args) > 2 {
			if argsHas("--random", "-r") {
				p := tea.NewProgram(initialModel(options{
					isPlayOne: argsHas("--one", "-o"),
				}), tea.WithAltScreen())
				if _, err := p.Run(); err != nil {
					fmt.Printf("Alas, there's been an error: %v", err)
					os.Exit(1)
				}
			} else if argsHas("--test") {
				p := tea.NewProgram(initialModel(options{
					isPlayOne: true,
					moesic: &data.Moesic{
						Name:         "Megumin's explosion destroys Chomusuke!",
						PlaylistName: "Konosuba",
						Url:          "https://www.youtube.com/watch?v=IQU49JbStVA",
					},
				}), tea.WithAltScreen())
				if _, err := p.Run(); err != nil {
					fmt.Printf("Alas, there's been an error: %v", err)
					os.Exit(1)
				}
			} else {
				printHelp()
			}
		} else {
			if _, err := tea.NewProgram(initSearchFlatModel()).Run(); err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}
		}
	case "update":
		InstallDependencies(true)
	case "version":
		fmt.Printf("Current version: %s\n", version)
	default:
		printHelp()
	}
}
