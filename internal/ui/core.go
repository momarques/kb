package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/teacup/statusbar"
	"github.com/momarques/kibe/internal/kube"
)

type viewState int

const (
	showList viewState = iota
	showTable
	showTab
)

type CoreUI struct {
	viewState
	height int

	client *kube.ClientReady

	listModel list.Model
	*listSelector

	enabledKeys

	tableModel table.Model
	*tableContent
	tableKeyMap

	tabModel tabModel
	tabKeyMap

	headerModel    headerModel
	helpModel      help.Model
	statusbarModel statusbar.Model
	syncBarModel   syncBarModel
	statusLogModel
}

func NewUI() CoreUI {
	selector := newListSelector()

	tableKeyMap := newTableKeyMap()
	tabKeyMap := newTabKeyMap()

	return CoreUI{
		viewState: showList,

		listSelector: selector,
		listModel:    newlistModel(selector),

		enabledKeys: setKeys(tableKeyMap, tabKeyMap),

		tableContent: newTableContent(),
		tableKeyMap:  tableKeyMap,
		tableModel:   newTableModel(),

		tabKeyMap: tabKeyMap,
		tabModel:  newTabModel(),

		headerModel:    headerModel{},
		helpModel:      help.New(),
		statusbarModel: newStatusBarModel(),
		statusLogModel: newStatusLogModel(),
		syncBarModel:   newSyncBarModel(),
	}
}

func (m CoreUI) Init() tea.Cmd {
	return tea.SetWindowTitle("Kibe UI")
}

func (m CoreUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.QuitMsg:
		return m, tea.Quit
	case statusLogMessage:
		return m.updateStatusLog(msg), nil
	}

	switch m.viewState {
	case showList:
		return m.updateListModel(msg)
	case showTable:
		m.enabledKeys = m.setEnabled(m.tableKeyMap.fullHelp()...)
		return m.updateTableModel(msg)
	case showTab:
		switch m.tabModel.tabViewState {
		case contentSelected:
			m.enabledKeys = m.setEnabled(m.tabKeyMap.fullHelp()...)
		case noContentSelected:
			m.enabledKeys = m.setEnabled(m.tabKeyMap.fullHelpWithContentSelected()...)
		}
		return m.updateTabModel(msg)
	}
	return m, nil
}

func (m CoreUI) View() string {
	switch m.viewState {

	case showList:
		return m.listModelView()

	case showTable, showTab:
		return m.composedView()
	}
	return m.View()
}

func (m CoreUI) showHelpLines(helpBindingLines ...[]key.Binding) []string {
	var helpLines []string

	helpStyle := lipgloss.NewStyle().MarginBottom(1)

	for _, line := range helpBindingLines {
		helpLines = append(helpLines, helpStyle.Render(
			m.helpModel.ShortHelpView(line)))
	}
	return helpLines
}

func (m CoreUI) composedView() string {
	var helpBindingLines [][]key.Binding

	switch m.viewState {
	case showTable:
		helpBindingLines = [][]key.Binding{
			m.tableKeyMap.firstHelpLineView(),
			m.tableKeyMap.secondHelpLineView(),
		}

	case showTab:
		switch m.tabModel.tabViewState {
		case noContentSelected:
			helpBindingLines = [][]key.Binding{
				m.tabKeyMap.firstHelpLineView(),
				m.tabKeyMap.secondHelpLineView(),
			}
		case contentSelected:
			helpBindingLines = [][]key.Binding{
				m.tabKeyMap.firstHelpLineViewWithContentSelected(),
				m.tabKeyMap.secondHelpLineView(),
			}
		}
	}

	helpView := lipgloss.JoinVertical(
		lipgloss.Center,
		m.showHelpLines(helpBindingLines...)...)

	leftUtilityPanel := lipgloss.JoinVertical(
		lipgloss.Left,
		m.paginatorModelView(),
		m.syncBarModelView(),
	)

	bottomPanel := lipgloss.JoinVertical(lipgloss.Left,
		m.tabModelView(),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			leftUtilityPanel,
			helpView,
		))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.headerModelView(),
		m.tableModelView(),
		lipgloss.JoinHorizontal(lipgloss.Center,
			bottomPanel,
			m.statusLogModelView()),
		m.statusbarModel.View())
}
