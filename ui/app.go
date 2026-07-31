// Package ui contains the Bubble Tea views and the top-level App that
// composes them. Views never talk to each other or to the database:
// they emit messages, and App orchestrates.
package ui

import (
	"jobtrack/models"
	"jobtrack/store"

	tea "charm.land/bubbletea/v2"
)

// ---- Messages views use to talk to the App ----

type showListMsg struct{}

// openFormMsg: the list view, user pressed 'a'.
type openFormMsg struct{}

type openAppDetailsMsg struct {
	app models.Application
}

type backMsg struct{}

// submitFormMsg: the form view, user pressed ctrl+s with valid input.
type submitFormMsg struct {
	app models.Application
}

type deleteAppMsg struct {
	app models.Application
}

type submitDeleteApplicationMsg struct {
	app models.Application
}

// cancelFormMsg: the form view, user pressed esc.
type cancelFormMsg struct{}

func showListCmd() tea.Cmd {
	return func() tea.Msg { return showListMsg{} }
}

func openFormCmd() tea.Cmd {
	return func() tea.Msg { return openFormMsg{} }
}

func openAppDetailsCmd(app models.Application) tea.Cmd {
	return func() tea.Msg { return openAppDetailsMsg{app: app} }
}

func backCmd() tea.Cmd {
	return func() tea.Msg { return backMsg{} }
}

func submitCmd(app models.Application) tea.Cmd {
	return func() tea.Msg { return submitFormMsg{app: app} }
}

func deleteAppCmd(app models.Application) tea.Cmd {
	return func() tea.Msg { return deleteAppMsg{app: app} }
}

func submitDeleteApplicationCmd(app models.Application) tea.Cmd {
	return func() tea.Msg { return submitDeleteApplicationMsg{app: app} }
}

func cancelCmd() tea.Cmd {
	return func() tea.Msg { return cancelFormMsg{} }
}

// App is the root model. It owns the store, the current application
// data, and whichever view is active.
type App struct {
	store *store.Store
	apps  []models.Application
	view  tea.Model
}

func NewApp(s *store.Store, apps []models.Application) App {
	return App{
		store: s,
		apps:  apps,
		view:  NewList(apps),
	}
}

func (m App) Init() tea.Cmd { return nil }

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case showListMsg:
		apps, err := m.store.ListApplications()
		if err != nil {
			return m, nil
		}
		m.apps = apps
		m.view = NewList(apps)
		return m, nil
	case openFormMsg:
		m.view = NewForm()
		return m, nil
	case submitFormMsg:
		if err := m.store.CreateApplication(&msg.app); err != nil {
			// TODO: surface this to the user instead of ignoring it.
			return m, nil
		}
		return m, showListCmd()
	case backMsg:
		m.view = NewList(m.apps)
		return m, nil
	case openAppDetailsMsg:
		m.view = NewApplicationDetail(msg.app)
		return m, nil
	case deleteAppMsg:
		m.view = NewConfirmDelete(msg.app)
		return m, nil
	case submitDeleteApplicationMsg:
		if err := m.store.DeleteApplication(msg.app.ID); err != nil {
			return m, nil
		}
		return m, showListCmd()
	}

	// Everything else goes to whichever view is active.
	view, cmd := m.view.Update(msg)
	m.view = view
	return m, cmd
}

func (m App) View() tea.View {
	return m.view.View()
}
