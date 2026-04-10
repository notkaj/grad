package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/notkaj/grad/internal/world"
)

func main() {
	if _, err := tea.NewProgram(world.InitialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
