package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/momarques/kibe/internal/bindings"
	"github.com/momarques/kibe/internal/kube"
	uistyles "github.com/momarques/kibe/internal/ui/styles"
)

type listModel struct {
	list.Model
	*listSelector
}

func newListModel() listModel {
	selector := newListSelector()

	l := list.New(
		[]list.Item{},
		newItemDelegate(selector), 0, 0)

	l.Styles.Title = uistyles.ViewTitleStyle.Copy()
	l.Styles.HelpStyle = uistyles.HelpStyle.Copy()
	l.Styles.FilterPrompt = uistyles.ListFilterPromptStyle.Copy()
	l.Styles.FilterCursor = uistyles.ListFilterCursorStyle.Copy()
	l.InfiniteScrolling = false
	l.KeyMap.Quit = bindings.New("quit", "q", "ctrl+c")
	return listModel{
		Model: l,

		listSelector: selector,
	}
}

func (m CoreUI) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := uistyles.
			AppStyle.
			Copy().
			GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		m.height = msg.Height
		m.statusBar.SetSize(msg.Width)
		return m, nil

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

	case spinner.TickMsg:
		m.list.spinner, cmd = m.list.spinner.Update(msg)
		return m, cmd

	case headerTitleUpdated:
		m.header.text = msg
		return m, nil

	case *kube.ClientReady:
		m.viewState = showTable
		m.client = msg
		return m, nil

	case statusBarUpdated:
		m.statusBar.SetContent(
			"Resource", m.list.resource,
			fmt.Sprintf("Context: %s", m.list.context),
			fmt.Sprintf("Namespace: %s", m.list.namespace))

		m.statusBar, cmd = m.statusBar.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.list.Model, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m CoreUI) listView() string {
	if m.list.spinnerState == showSpinner {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			fmt.Sprintf("%s%s",
				m.list.spinner.View(),
				m.list.View()),
			m.statusBar.View())
	}
	return lipgloss.JoinVertical(
		lipgloss.Top, m.list.View(),
		m.statusBar.View())
}
