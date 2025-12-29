package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/keyforge/keyforge/internal/engine"
	"github.com/keyforge/keyforge/internal/ui"
)

func main() {
	// Connection flags
	nvimMode := flag.Bool("nvim-mode", false, "Enable Neovim integration mode with RPC")
	rpcSocket := flag.String("rpc-socket", "", "Unix socket path for RPC communication")

	// Game settings flags (defaults from nvim config)
	difficulty := flag.String("difficulty", "normal", "Difficulty level: easy, normal, hard")
	gameSpeed := flag.Float64("game-speed", 1.0, "Game speed multiplier: 0.5, 1.0, 1.5, 2.0")
	startingGold := flag.Int("starting-gold", 200, "Starting gold amount (100-500)")
	startingHealth := flag.Int("starting-health", 100, "Starting health (50-200)")

	flag.Parse()

	// Build settings from flags
	settings := engine.GameSettings{
		Difficulty:     *difficulty,
		GameSpeed:      engine.GameSpeed(*gameSpeed),
		StartingGold:   *startingGold,
		StartingHealth: *startingHealth,
	}
	settings.Validate()

	model := ui.NewModelWithSettings(settings)
	model.NvimMode = *nvimMode

	// In nvim mode, start the RPC server/client
	if *nvimMode {
		if *rpcSocket != "" {
			// Use Unix socket for RPC (preferred)
			model.InitNvimSocket(*rpcSocket)
		} else {
			// Fallback to stdin/stderr (legacy)
			model.InitNvimClient()
		}
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
