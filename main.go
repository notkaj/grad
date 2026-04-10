package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/notkaj/grad/internal/chooser"
)

func main() {
	if _, err := tea.NewProgram(chooser.InitialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
