package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type searchModel struct {
	list list.Model
}

func (m searchModel) Init() tea.Cmd {
	return nil
}

func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			if it, ok := m.list.SelectedItem().(item); ok {
				for i, pl := range playlists {
					if pl.Name == it.title {
						p := tea.NewProgram(initialPlayerPlaylistplaylistModel(&playlists[i]), tea.WithAltScreen())
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

func (m searchModel) View() string {
	return docStyle.Render(m.list.View())
}

func initSearchModel() searchModel {

	items := []list.Item{}
	for _, s := range playlists {
		items = append(items, item{title: s.Name})
	}
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	m := searchModel{list: list.New(items, d, 0, 0)}
	m.list.Title = "What playlist would you like to play?"
	return m
}
