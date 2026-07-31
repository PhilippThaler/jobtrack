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
	path := os.Getenv("XDG_CONFIG_HOME")
	if path != "" {
		fullPath = filepath.Join(path, "jobtrack")
		os.Mkdir(fullPath, 0755)
	} else {
		path = os.Getenv("HOME")
		if path != "" {
			fullPath = filepath.Join(path, ".jobtrack")
			os.Mkdir(fullPath, 0755)
		} else {
			fmt.Println("Can't find $XDG_CONFIG_HOME or $HOME")
			fmt.Println("Exiting...")
			os.Exit(1)
		}
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
