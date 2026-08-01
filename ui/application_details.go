package ui

import (
	"fmt"
	"jobtrack/models"

	tea "charm.land/bubbletea/v2"
)

type ApplicationDetailModel struct {
	application models.Application
}

func NewApplicationDetail(app models.Application) ApplicationDetailModel {
	return ApplicationDetailModel{application: app}
}

func (m ApplicationDetailModel) Init() tea.Cmd {
	return nil
}

func (m ApplicationDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			return m, backCmd()
		case "e":
			return m, openEditFormCmd(m.application)
		case "d":
			return m, deleteAppCmd(m.application)
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ApplicationDetailModel) View() tea.View {
	s := "Details\n"
	s += "----------------------------------------\n"
	s += fmt.Sprintf("\tCompany: %s\n", m.application.Company)
	s += fmt.Sprintf("\tRole:    %s\n", m.application.Role)
	s += fmt.Sprintf("\tURL:     %s\n", m.application.URL)
	s += fmt.Sprintf("\tNotes:   %s\n", m.application.Notes)
	s += fmt.Sprintf("\tStatus:  %s\n", m.application.Status)
	s += "----------------------------------------\n\n"
	s += "esc: back    e: edit    d: delete    q: quit"
	return tea.NewView(s)
}
