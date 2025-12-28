package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/keyforge/keyforge/internal/ui"
)

func main() {
	nvimMode := flag.Bool("nvim-mode", false, "Enable Neovim integration mode with RPC")
	flag.Parse()

	model := ui.NewModel()
	model.NvimMode = *nvimMode

	// In nvim mode, start the RPC client
	if *nvimMode {
		model.InitNvimClient()
	}

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running keyforge: %v\n", err)
		os.Exit(1)
	}
}
