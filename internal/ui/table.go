package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/momarques/kibe/internal/kube"
	"github.com/momarques/kibe/internal/logging"
	uistyles "github.com/momarques/kibe/internal/ui/styles"
	windowutil "github.com/momarques/kibe/internal/ui/window_util"
)

const tableViewHeightPercentage int = 32

func newTableModel() table.Model {
	t := table.New(
		table.WithFocused(true),
	)

	t.SetStyles(uistyles.NewTableStyle(false))
	t.SetHeight(
		windowutil.ComputeHeightPercentage(tableViewHeightPercentage))
	return t
}

func (m CoreUI) updateTableModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.tableContent.syncState {
	case synced, syncing:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:

			m.tableModel.SetHeight(msg.Height - m.tableModel.Height())
			// m.tableModel.SetWidth(msg.Width - m.tableModel.Width())
			m.tableModel.SetColumns(m.client.ResourceSelected.Columns())
			logging.Log.Infof("window size -> %d x %d", msg.Width, msg.Height)
			logging.Log.Infof("table size -> %d x %d", m.tableModel.Width(), m.tableModel.Height())
			m.helpModel.Width = 20
			m.tableModel, cmd = m.tableModel.Update(msg)
			return m, cmd

		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.tableKeyMap.Quit):
				return m, tea.Quit

			case key.Matches(msg, m.tableKeyMap.Describe):
				selectedResource := m.tableModel.SelectedRow()

				m.tabModel, cmd = m.tabModel.describeResource(m.client, selectedResource[0])
				return m, cmd

			case key.Matches(msg, m.tableKeyMap.PreviousPage, m.tableKeyMap.NextPage):
				m.paginatorModel, _ = m.paginatorModel.Update(msg)
				m.tableModel, cmd = m.applyTableItems(m.tableModel)

				return m, cmd
			}

		case descriptionReady:
			m.viewState = showTab
			m.tabModel.Tabs, m.tabModel.TabContent = msg.tabNames, msg.tabContent
			return m, nil

		case headerItemCountUpdated:
			m.headerModel.itemCount = msg
			return m, nil

		case lastSync:
			m.tableContent.syncState = synced
			m.syncBarModel = m.changeSyncState()

			return m, tea.Batch(tea.Tick(kube.ResquestTimeout, startSyncing))

		case syncState:
			if msg == unsynced {
				m.tableContent.syncState = msg
				m.syncBarModel = m.changeSyncState()
				return m.sync(nil)
			}
		}

	case unsynced:
		return m.sync(msg)
	}

	m.tableModel, cmd = m.tableModel.Update(msg)
	return m, cmd
}

func (m CoreUI) tableModelView() string {
	tableStyle := uistyles.TableStyle

	if m.viewState == showTab {
		tableStyle = uistyles.DimmedTableStyle
		m.tableModel.SetStyles(uistyles.NewTableStyle(true))
		return tableStyle.Render(m.tableModel.View())
	}
	if m.tableContent.columns == nil {
		return lipgloss.NewStyle().
			Height((windowutil.ComputeHeightPercentage(tableViewHeightPercentage) + 3)).
			Render("")
	}
	return tableStyle.Render(m.tableModel.View())
}
