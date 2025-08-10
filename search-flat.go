package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type searchFlatModel struct {
	list list.Model
}

func (m searchFlatModel) Init() tea.Cmd {
	return nil
}

func (m searchFlatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {

			if it, ok := m.list.SelectedItem().(item); ok {
				for i, moe := range flatSongs {
					if moe.Name == it.title {
						p := tea.NewProgram(initialModel(options{
							isPlayOne: argsHas("--one", "-o"),
							moesic:    &flatSongs[i],
						}), tea.WithAltScreen())
						if _, err := p.Run(); err != nil {
							fmt.Printf("Alas, there's been an error: %v", err)
							os.Exit(1)
						}
						p.Wait()
						break
					}
				}
			}

			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m searchFlatModel) View() string {
	return docStyle.Render(m.list.View())
}

func initSearchFlatModel() searchFlatModel {

	items := []list.Item{}
	for _, s := range flatSongs {
		items = append(items, item{title: s.Name, desc: s.PlaylistName})
	}
	d := list.NewDefaultDelegate()
	m := searchFlatModel{list: list.New(items, d, 0, 0)}
	m.list.Title = "What moesic would you like to play?"
	m.list.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.NoColor{}).
		Bold(true).
		Foreground(lipgloss.Color("#FF69B4"))
	return m
}
