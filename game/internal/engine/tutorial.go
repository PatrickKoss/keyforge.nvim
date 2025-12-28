package engine

// TutorialStep represents a single step in the tutorial
type TutorialStep struct {
	ID          string
	Title       string
	Description string
	Highlight   string // what to highlight: "cursor", "tower", "enemy", "path", "shop"
	WaitFor     string // what to wait for: "move", "place_tower", "kill_enemy", "any_key"
	Position    *struct{ X, Y int }
}

// Tutorial manages the first-time user experience
type Tutorial struct {
	Active      bool
	CurrentStep int
	Steps       []TutorialStep
	Completed   bool
}

// NewTutorial creates a new tutorial
func NewTutorial() *Tutorial {
	return &Tutorial{
		Active:      false,
		CurrentStep: 0,
		Steps:       createTutorialSteps(),
		Completed:   false,
	}
}

func createTutorialSteps() []TutorialStep {
	return []TutorialStep{
		{
			ID:          "welcome",
			Title:       "Welcome to Keyforge!",
			Description: "Defend against waves of bugs by building towers.\nPress any key to continue...",
			WaitFor:     "any_key",
		},
		{
			ID:          "movement",
			Title:       "Movement",
			Description: "Use h/j/k/l (vim keys) or arrow keys to move the cursor.\nTry moving around the grid now.",
			Highlight:   "cursor",
			WaitFor:     "move",
		},
		{
			ID:          "path",
			Title:       "The Path",
			Description: "Enemies follow the highlighted path (░).\nThey enter from the left and exit on the right.\nDon't let them reach the end!",
			Highlight:   "path",
			WaitFor:     "any_key",
		},
		{
			ID:          "towers",
			Title:       "Building Towers",
			Description: "Press 1, 2, or 3 to select a tower type.\nThen press SPACE to place it.\nTowers can't be placed on the path.",
			Highlight:   "shop",
			WaitFor:     "any_key",
		},
		{
			ID:          "place_tower",
			Title:       "Place Your First Tower",
			Description: "Move to an empty cell and press SPACE to place a tower.\nTry placing an Arrow Tower (1) near the path.",
			Highlight:   "cursor",
			WaitFor:     "place_tower",
		},
		{
			ID:          "combat",
			Title:       "Combat",
			Description: "Towers automatically attack enemies in range.\nKilling enemies earns gold (💰).\nWhen enemies reach the end, you lose health (❤️).",
			WaitFor:     "any_key",
		},
		{
			ID:          "upgrades",
			Title:       "Upgrades",
			Description: "Move cursor over a tower and press 'u' to upgrade.\nUpgrades increase damage, range, and attack speed.",
			WaitFor:     "any_key",
		},
		{
			ID:          "waves",
			Title:       "Waves",
			Description: "Survive all 10 waves to win!\nEach wave brings stronger enemies.\nGood luck, vim warrior!",
			WaitFor:     "any_key",
		},
		{
			ID:          "complete",
			Title:       "Tutorial Complete!",
			Description: "You're ready to play!\nPress 'p' to pause anytime.\nPress any key to start the game...",
			WaitFor:     "any_key",
		},
	}
}

// Start begins the tutorial
func (t *Tutorial) Start() {
	t.Active = true
	t.CurrentStep = 0
	t.Completed = false
}

// Skip ends the tutorial early
func (t *Tutorial) Skip() {
	t.Active = false
	t.Completed = true
}

// CurrentStepData returns the current tutorial step
func (t *Tutorial) CurrentStepData() *TutorialStep {
	if t.CurrentStep >= len(t.Steps) {
		return nil
	}
	return &t.Steps[t.CurrentStep]
}

// Advance moves to the next tutorial step
func (t *Tutorial) Advance() bool {
	t.CurrentStep++
	if t.CurrentStep >= len(t.Steps) {
		t.Active = false
		t.Completed = true
		return false
	}
	return true
}

// CheckCondition checks if the current step's condition is met
func (t *Tutorial) CheckCondition(event string) bool {
	step := t.CurrentStepData()
	if step == nil {
		return false
	}

	switch step.WaitFor {
	case "any_key":
		return true
	case "move":
		return event == "move"
	case "place_tower":
		return event == "place_tower"
	case "kill_enemy":
		return event == "kill_enemy"
	default:
		return event == step.WaitFor
	}
}

// HandleEvent processes a game event for tutorial progression
func (t *Tutorial) HandleEvent(event string) bool {
	if !t.Active {
		return false
	}

	if t.CheckCondition(event) {
		return t.Advance()
	}
	return false
}

// IsActive returns whether the tutorial is currently active
func (t *Tutorial) IsActive() bool {
	return t.Active
}

// ShouldShowTutorial checks if tutorial should be offered
func (t *Tutorial) ShouldShowTutorial() bool {
	return !t.Completed
}
