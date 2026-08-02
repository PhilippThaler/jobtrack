package ui

import (
	"fmt"
	"strings"

	"jobtrack/models"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// Form fields. The first four are text inputs, fieldStatus is special:
// it cycles through statuses instead of accepting text.
const (
	fieldCompany = iota
	fieldRole
	fieldURL
	fieldNotes
	fieldStatus
	fieldFavorite
	fieldCount // always last, gives the total number of fields
)

type FormModel struct {
	inputs    []textinput.Model
	focus     int // which field is active, 0..fieldStatus
	statusIdx int // index into models.AllStatuses
	id        int64
	title     string
	favorite  bool
	errMsg    string
}

func NewForm() FormModel {
	return newForm(nil)
}

func NewEditForm(app models.Application) FormModel {
	return newForm(&app)
}

func newForm(app *models.Application) FormModel {
	company := textinput.New()
	company.Placeholder = "Acme Corp"
	company.Focus()

	role := textinput.New()
	role.Placeholder = "System Administrator"

	url := textinput.New()
	url.Placeholder = "https://example.com/job/123"

	notes := textinput.New()
	notes.Placeholder = "Anything worth remembering"

	f := FormModel{title: "Add Application"}
	if app != nil {
		company.SetValue(app.Company)
		role.SetValue(app.Role)
		url.SetValue(app.URL)
		notes.SetValue(app.Notes)
		f.id = app.ID
		f.statusIdx = statusIndex(app.Status)
		f.title = "Edit Application"
		f.favorite = app.IsFavorite
	}
	f.inputs = []textinput.Model{company, role, url, notes}

	return f
}

func statusIndex(status string) int {
	for i, s := range models.AllStatuses {
		if s == status {
			return i
		}
	}
	return 0
}

func (m FormModel) Init() tea.Cmd { return nil }

// focusField moves the visual cursor to m.focus.
func (m FormModel) focusField() {
	for i := range m.inputs {
		if i == m.focus {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

func (m FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, backCmd()
		case "ctrl+s":
			return m.submit()
		case "tab":
			m.focus = (m.focus + 1) % fieldCount
			m.focusField()
			return m, nil
		case "shift+tab":
			// (focus + fieldCount - 1) % fieldCount wraps backwards
			m.focus = (m.focus + fieldCount - 1) % fieldCount
			m.focusField()
			return m, nil
		case "right", "space":
			if m.focus == fieldStatus {
				m.statusIdx = (m.statusIdx + 1) % len(models.AllStatuses)
				return m, nil
			} else if m.focus == fieldFavorite {
				m.favorite = !m.favorite
				return m, nil
			}
		case "left":
			if m.focus == fieldStatus {
				m.statusIdx = (m.statusIdx + len(models.AllStatuses) - 1) % len(models.AllStatuses)
				return m, nil
			} else if m.focus == fieldFavorite {
				m.favorite = !m.favorite
				return m, nil
			}
		}
	}

	// Any other keypress goes to the focused text input.
	if m.focus < fieldStatus {
		ti, cmd := m.inputs[m.focus].Update(msg)
		m.inputs[m.focus] = ti
		return m, cmd
	}
	return m, nil
}

// submit validates the input and, if valid, sends the data to the App.
func (m FormModel) submit() (tea.Model, tea.Cmd) {
	company := m.inputs[fieldCompany].Value()
	role := m.inputs[fieldRole].Value()
	if company == "" || role == "" {
		m.errMsg = "Company and role are required"
		return m, nil
	}
	return m, submitCmd(models.Application{
		ID:         m.id,
		Company:    company,
		Role:       role,
		URL:        m.inputs[fieldURL].Value(),
		Status:     models.AllStatuses[m.statusIdx],
		Notes:      m.inputs[fieldNotes].Value(),
		IsFavorite: m.favorite,
	})
}

func (m FormModel) View() tea.View {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%s    (esc: cancel, ctrl+s: save)\n", m.title))
	s.WriteString("----------------------------------------\n\n")

	labels := []string{"Company:", "Role:", "URL:", "Notes:", "Status:", "Favorite:"}
	for i, label := range labels {
		marker := " "
		if m.focus == i {
			marker = ">"
		}
		if i == fieldStatus {
			s.WriteString(fmt.Sprintf("%s %-9s %s   (left/right to change)\n", marker, label, models.AllStatuses[m.statusIdx]))
		} else if i == fieldFavorite {
			box := "[ ]"
			if m.favorite {
				box = "[x]"
			}
			s.WriteString(fmt.Sprintf("%s %-9s %s    (space to toggle)\n", marker, label, box))
		} else {
			s.WriteString(fmt.Sprintf("%s %-9s %s\n", marker, label, m.inputs[i].View()))
		}
	}

	if m.errMsg != "" {
		s.WriteString("\n" + m.errMsg + "\n")
	}
	return tea.NewView(s.String())
}
