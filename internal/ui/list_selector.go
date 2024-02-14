package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mistakenelf/teacup/statusbar"
	"github.com/momarques/kibe/internal/bindings"
	"github.com/momarques/kibe/internal/kube"
	uistyles "github.com/momarques/kibe/internal/ui/styles"
	"github.com/samber/lo"
)

type clientState int

const (
	ready clientState = iota
	notReady
)

type selector struct {
	clientState
	spinnerState

	client *kube.ClientReady

	context   string
	namespace string
	resource  string

	chooseKey key.Binding
	// previousKey key.Binding

	spinner   spinner.Model
	statusbar *statusbar.Model

	// previousSelection string
}

func newListSelector(spinner spinner.Model, status *statusbar.Model) *selector {
	return &selector{
		clientState:  notReady,
		spinnerState: hideSpinner,
		chooseKey:    bindings.New("enter", "choose"),
		// previousKey:  bindings.New(",", "previous"),
		spinner:   spinner,
		statusbar: status,
	}
}

func newItemDelegate(s *selector) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = s.update

	d.Styles.SelectedTitle = activeSelectionStyle
	d.Styles.SelectedDesc = activeSelectionStyle

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{s.chooseKey} //, s.previousKey}
	}
	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{{s.chooseKey}} //, s.previousKey}}
	}
	return d
}

func (s *selector) update(msg tea.Msg, m *list.Model) tea.Cmd {
	switch msg := msg.(type) {

	case kube.SelectContext:
		m.Title = "Choose a context to connect"
		s.spinnerState = hideSpinner

		return m.SetItems(msg.Contexts)

	case kube.ContextSelected:
		m.ResetFilter()

		s.client = kube.NewClientReady(msg.C)
		return s.updateStatusBar()

	case kube.SelectNamespace:
		m.Title = "Choose a namespace"
		s.spinnerState = hideSpinner

		// s.setPreviousSelection("context")

		return m.SetItems(msg.Namespaces)

	case kube.NamespaceSelected:
		m.ResetFilter()
		s.client = s.client.WithNamespace(msg.NS)

		return s.updateStatusBar()

	case kube.SelectResource:
		m.Title = "Choose a resource type"
		s.spinnerState = hideSpinner

		// s.setPreviousSelection("namespace")

		return m.SetItems(msg.Resources)

	case kube.ResourceSelected:
		m.ResetFilter()

		s.client = s.client.WithResource(msg.R)
		return s.updateStatusBar()

	case tea.KeyMsg:
		switch {
		// case key.Matches(msg, s.previousKey):
		// 	return s.showPreviousSelection()

		case key.Matches(msg, s.chooseKey):
			switch i := m.SelectedItem().(type) {

			case kube.ContextItem:
				s.context = i.FilterValue()
				s.spinnerState = showSpinner

				return tea.Batch(
					m.NewStatusMessage(uistyles.StatusMessageStyle(
						"Context", s.context, "set")),
					s.contextSelected(s.context),
					s.spinner.Tick)

			case kube.NamespaceItem:
				s.namespace = i.FilterValue()
				s.spinnerState = showSpinner

				return tea.Batch(
					m.NewStatusMessage(uistyles.StatusMessageStyle(
						"Namespace", s.namespace, "selected")),
					s.namespaceSelected(s.namespace),
					s.spinner.Tick)

			case kube.ResourceItem:
				s.resource = i.FilterValue()
				s.spinnerState = showSpinner

				return tea.Batch(
					m.NewStatusMessage(uistyles.StatusMessageStyle(
						"Showing", s.resource+"s")),
					s.resourceSelected(s.resource),
					s.spinner.Tick)
			}
		}

	default:

		if s.context != "" && s.namespace != "" && s.resource != "" {
			s.clientState = ready
		}

		switch s.clientState {
		case ready:
			return s.clientReady()

		case notReady:
			if s.context == "" {
				return kube.NewSelectContext()
			}

			if s.namespace == "" {
				return kube.NewSelectNamespace(s.client)
			}

			if s.resource == "" {
				return kube.NewSelectResource(s.client)
			}
		}
	}
	return s.spinner.Tick
}

func (s *selector) contextSelected(context string) func() tea.Msg {
	return func() tea.Msg {
		return kube.ContextSelected{C: context}
	}
}

func (s *selector) namespaceSelected(namespace string) func() tea.Msg {
	return func() tea.Msg {
		return kube.NamespaceSelected{NS: namespace}
	}
}

func (s *selector) resourceSelected(kind string) func() tea.Msg {
	r, _ := lo.Find(kube.SupportedResources, func(item kube.Resource) bool {
		switch item.Kind() {
		case "Pod", "Service", "Namespace":
			return true
		}
		return false
	})

	return func() tea.Msg {
		return kube.ResourceSelected{R: r}
	}
}

func (s *selector) clientReady() func() tea.Msg {
	return func() tea.Msg {
		return s.client
	}
}

// func (s *selector) showPreviousSelection() func() tea.Msg {
// 	return func() tea.Msg {
// 		switch s.previousSelection {
// 		case "context":
// 			s.context = ""
// 		case "namespace":
// 			s.namespace = ""
// 		case "resource":
// 			s.resource = ""
// 		}
// 		return nil
// 	}
// }

// func (s *selector) setPreviousSelection(previousSelection string) {
// 	s.previousSelection = previousSelection
// }
