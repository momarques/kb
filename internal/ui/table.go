package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func newTableUI() table.Model {
	t := table.New(
		table.WithFocused(true),
	)

	s := table.DefaultStyles()

	s.Cell = s.Cell.Blink(false)
	s.Header = s.Header.Blink(false).Background(lipgloss.Color("#c5636a"))
	s.Selected = s.Selected.Blink(false).Background(lipgloss.Color("#ffb1b5")).Foreground(lipgloss.Color("#322223"))

	// s.Header = s.Header.
	// 	Border(lipgloss.NormalBorder()).
	// 	BorderForeground(lipgloss.Color("99")).
	// 	Bold(true).
	// 	Foreground(lipgloss.Color("99"))

	// s.Cell = s.Cell.
	// 	Border(lipgloss.NormalBorder()).
	// 	BorderForeground(lipgloss.Color("99")).
	// 	Foreground(lipgloss.Color("229"))

	// s.Selected = s.Selected.
	// 	BorderForeground(lipgloss.Color("#ffffff")).
	// 	Foreground(lipgloss.Color("#ffffff"))

	t.SetStyles(s)

	return t
}

func (m CoreUI) updateTableUI(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.tableContent.contentState {
	case loaded:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:

			m.tableUI.SetHeight(msg.Height - m.tableUI.Height())
			m.tableUI.SetWidth(msg.Width - m.tableUI.Width())

			m.tableUI, cmd = m.tableUI.Update(msg)
			return m, cmd

		case tea.KeyMsg:

			switch msg.String() {
			case "esc":
				if m.tableUI.Focused() {
					m.tableUI.Blur()
				} else {
					m.tableUI.Focus()
				}
			case "q", "ctrl+c":
				return m, tea.Quit
			case "enter":
				return m, tea.Batch(
					tea.Printf("Let's go to %s!", m.tableUI.SelectedRow()[1]),
				)
			}
		case headerUpdated:
			m.headerUI.text = msg.text
			m.headerUI.itemCount = msg.itemCount

		default:
			return m, tea.Tick(loadInterval, func(t time.Time) tea.Msg {
				m.tableContent.contentState = notLoaded
				return nil
			})
		}

	case notLoaded:
		m.tableContent.client = m.client

		m.tableUI, cmd = m.tableContent.fetch(m.tableUI)
		return m, cmd
	}

	m.tableUI, cmd = m.tableUI.Update(msg)
	return m, cmd
}

func (m CoreUI) viewTableUI() string {
	tableStyle := lipgloss.NewStyle()
	tableView := lipgloss.Place(5, 5,
		lipgloss.Center,
		lipgloss.Center,
		tableStyle.
			MarginLeft(2).
			Border(lipgloss.DoubleBorder(), true, true, true, true).
			BorderForeground(lipgloss.Color("#ffb8bc")).
			Render(m.tableUI.View()))
	fullTableSize := lipgloss.Width(tableView)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.headerUI.viewHeaderUI(fullTableSize),
		tableView,
		m.statusbarUI.View(),
	)

	// return lipgloss.JoinVertical(lipgloss.Left,
	// 	m.tableUI.View(), m.statusbarUI.View())
	// lipgloss.Place(1, 1, lipgloss.Center, lipgloss.Center, m.tableUI.View()),
	// lipgloss.Place(
	// 	1, 1,
	// 	lipgloss.Center, lipgloss.Bottom,
	// 	m.statusbarUI.View()))
}
