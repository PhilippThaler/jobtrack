package ui

import (
	"fmt"

	"jobtrack/models"

	tea "charm.land/bubbletea/v2"
)

// ListModel shows all applications as a table.
// On 'a' it tells the App to open the form via openFormMsg.
type ListModel struct {
	applications []models.Application
	cursor       int
}

func NewList(apps []models.Application) ListModel {
	return ListModel{applications: apps}
}

func (m ListModel) Init() tea.Cmd { return nil }

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.applications)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.applications) > 0 {
				return m, openAppDetailsCmd(m.applications[m.cursor])
			}
		case "a":
			return m, openFormCmd()
		case "d":
			if len(m.applications) > 0 {
				return m, deleteAppCmd(m.applications[m.cursor])
			}
		case "e":
			if len(m.applications) > 0 {
				return m, openEditFormCmd(m.applications[m.cursor])
			}
		}
	}
	return m, nil
}

func (m ListModel) View() tea.View {
	s := "JobTrack\n"
	s += "----------------------------------------\n"
	for i, app := range m.applications {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %-20s %-25s %s\n", cursor, app.Company, app.Role, app.Status)
	}
	s += "\nenter: details        j/k: move    q: quit\n"
	s += "a: add    e: edit    d: delete\n"

	return tea.NewView(s)
}
