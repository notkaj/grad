package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	s "github.com/notkaj/grad/internal/selection"
)

func main() {
	if _, err := tea.NewProgram(s.InitialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
