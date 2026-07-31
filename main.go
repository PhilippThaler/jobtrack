package main

import (
	"fmt"
	"os"
	"path/filepath"

	"jobtrack/store"
	"jobtrack/ui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	var fullPath string
	path, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	fullPath = filepath.Join(path, "jobtrack")
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	fullPath = filepath.Join(fullPath, "jobs.db")

	s, err := store.New(fullPath)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	defer s.Close()

	apps, err := s.ListApplications()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.NewApp(s, apps))
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
