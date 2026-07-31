package ui

import (
	"fmt"
	"jobtrack/models"

	tea "charm.land/bubbletea/v2"
)

type ConfirmDeleteModel struct {
	application models.Application
}

func NewConfirmDelete(application models.Application) ConfirmDeleteModel {
	return ConfirmDeleteModel{application: application}
}

func (m ConfirmDeleteModel) Init() tea.Cmd { return nil }

func (m ConfirmDeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "y":
			return m, submitDeleteApplicationCmd(m.application)
		case "n", "esc":
			return m, openAppDetailsCmd(m.application)
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}
func (m ConfirmDeleteModel) View() tea.View {
	s := "Confirm deletion\n"
	s += "----------------------------------------\n"
	s += fmt.Sprintf("Do you really want to delete your application for %s\n", m.application.Company)
	s += "[y/n]\n"

	return tea.NewView(s)
}
