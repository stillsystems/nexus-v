package tui

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	id, title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.id
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

func (m model) View() string {
	if m.choice != "" || m.quitting {
		return ""
	}
	return docStyle.Render(m.list.View())
}

func SelectVariant(variants []string) (string, error) {
	var items []list.Item
	for _, v := range variants {
		desc := "Scaffold using the " + v + " template"
		// Provide better descriptions for built-in variants
		switch v {
		case "command":
			desc = "Basic extension with a registered command"
		case "webview":
			desc = "Extension with a custom HTML webview panel"
		case "language":
			desc = "Language support with syntax highlighting"
		case "theme":
			desc = "Custom color theme extension"
		case "minimal":
			desc = "Bare-bones extension structure"
		case "multi":
			desc = "Advanced template with multiple commands"
		case "web":
			desc = "Web-compatible extension (browser-ready)"
		}
		items = append(items, item{id: v, title: strings.Title(v), desc: desc})
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select Template Variant"

	m := model{list: l}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	res := finalModel.(model).choice
	if res == "" && finalModel.(model).quitting {
		os.Exit(0)
	}

	return res, nil
}
