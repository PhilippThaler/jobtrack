package ui

import (
	"fmt"
	"jobtrack/models"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// ListModel shows all applications as a table.
// On 'a' it tells the App to open the form via openFormMsg.
type ListModel struct {
	list list.Model
}

type appItem struct {
	application models.Application
}

func (i appItem) Title() string {
	star := "  "
	if i.application.IsFavorite {
		star = "⭐"
	}
	return fmt.Sprintf("%s  %s", star, i.application.Company)
}
func (i appItem) Description() string {
	return fmt.Sprintf("%s - %s", i.application.Role, i.application.Status)
}

func (i appItem) FilterValue() string {
	return fmt.Sprintf("%s %s %s", i.application.Company, i.application.Role, i.application.Status)
}

func toItems(applications []models.Application) []list.Item {
	items := make([]list.Item, len(applications))
	for i, app := range applications {
		items[i] = appItem{application: app}
	}
	return items
}

func NewList(apps []models.Application, width, height int) ListModel {
	l := list.New(toItems(apps), list.NewDefaultDelegate(), width, height)
	return ListModel{list: l}
}

func (m ListModel) Init() tea.Cmd { return nil }

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width-appBox.GetHorizontalFrameSize(), msg.Height-appBox.GetVerticalFrameSize())
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "a":
			return m, openFormCmd()
		case "d":
			app := m.list.SelectedItem().(appItem).application
			return m, deleteAppCmd(app)
		case "e":
			app := m.list.SelectedItem().(appItem).application
			return m, openEditFormCmd(app)
		case "f":
			app := m.list.SelectedItem().(appItem).application
			return m, toggleFavoriteCmd(app)
		case "enter":
			app := m.list.SelectedItem().(appItem).application
			return m, openAppDetailsCmd(app)
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ListModel) View() tea.View {
	return tea.NewView(m.list.View())
}
