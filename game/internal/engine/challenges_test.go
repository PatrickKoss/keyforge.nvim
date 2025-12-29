package engine

import (
	"testing"
)

func TestDurationTier(t *testing.T) {
	tests := []struct {
		name          string
		parKeystrokes int
		expectedTier  string
	}{
		{"1 keystroke is quick", 1, "quick"},
		{"5 keystrokes is quick", 5, "quick"},
		{"6 keystrokes is standard", 6, "standard"},
		{"15 keystrokes is standard", 15, "standard"},
		{"16 keystrokes is complex", 16, "complex"},
		{"40 keystrokes is complex", 40, "complex"},
		{"41 keystrokes is expert", 41, "expert"},
		{"100 keystrokes is expert", 100, "expert"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Challenge{ParKeystrokes: tt.parKeystrokes}
			if got := c.DurationTier(); got != tt.expectedTier {
				t.Errorf("DurationTier() = %v, want %v", got, tt.expectedTier)
			}
		})
	}
}

func TestParTime(t *testing.T) {
	tests := []struct {
		name          string
		parKeystrokes int
		expectedTime  int
	}{
		{"quick challenges have 5s par time", 3, 5},
		{"standard challenges have 15s par time", 10, 15},
		{"complex challenges have 45s par time", 25, 45},
		{"expert challenges have 90s par time", 50, 90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Challenge{ParKeystrokes: tt.parKeystrokes}
			if got := c.ParTime(); got != tt.expectedTime {
				t.Errorf("ParTime() = %v, want %v", got, tt.expectedTime)
			}
		})
	}
}

func TestChallengeManagerLoadsSuccessfully(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}
	if cm == nil {
		t.Fatal("NewChallengeManager() returned nil")
	}
}

func TestChallengeCount(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	// We should have at least 150 challenges (spec requirement)
	count := cm.Count()
	if count < 150 {
		t.Errorf("Count() = %d, want at least 150", count)
	}
}

func TestChallengeCategories(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	categories := cm.GetCategories()
	categoryMap := make(map[string]bool)
	for _, c := range categories {
		categoryMap[c] = true
	}

	// Expected categories from the spec
	expectedCategories := []string{
		"movement",
		"text-objects",
		"lsp-navigation",
		"search-replace",
		"refactoring",
		"git-operations",
		"window-management",
		"buffer-management",
		"folding",
		"quickfix",
		"diagnostics",
		"telescope",
		"surround",
		"harpoon",
		"formatting",
	}

	for _, expected := range expectedCategories {
		if !categoryMap[expected] {
			t.Errorf("GetCategories() missing category: %s", expected)
		}
	}
}

func TestCategoryChallengeBalance(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	counts := cm.CountByCategory()

	// Each category should have between 5-25 challenges (from spec)
	for cat, count := range counts {
		if count < 5 {
			t.Errorf("Category %s has %d challenges, want at least 5", cat, count)
		}
		if count > 25 {
			t.Errorf("Category %s has %d challenges, want at most 25", cat, count)
		}
	}
}

func TestGetChallengeByID(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	// Test getting a known challenge
	challenge := cm.GetChallenge("movement_end_of_line")
	if challenge == nil {
		t.Skip("Challenge movement_end_of_line not found, skipping")
	}

	if challenge.ID != "movement_end_of_line" {
		t.Errorf("GetChallenge() ID = %v, want movement_end_of_line", challenge.ID)
	}
	if challenge.Category != "movement" {
		t.Errorf("GetChallenge() Category = %v, want movement", challenge.Category)
	}
}

func TestGetChallengeByIDNotFound(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	challenge := cm.GetChallenge("nonexistent_challenge")
	if challenge != nil {
		t.Errorf("GetChallenge() = %v, want nil for nonexistent challenge", challenge)
	}
}

func TestGetRandomChallenge(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	// Get a random movement challenge
	challenge := cm.GetRandomChallenge("movement", 5)
	if challenge == nil {
		t.Skip("No movement challenges found")
	}

	if challenge.Category != "movement" {
		t.Errorf("GetRandomChallenge() Category = %v, want movement", challenge.Category)
	}
	if challenge.Difficulty > 5 {
		t.Errorf("GetRandomChallenge() Difficulty = %v, want <= 5", challenge.Difficulty)
	}
}

func TestGetRandomChallengeNoMatch(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	// Get a challenge from a non-existent category
	challenge := cm.GetRandomChallenge("nonexistent", 5)
	if challenge != nil {
		t.Errorf("GetRandomChallenge() = %v, want nil for nonexistent category", challenge)
	}
}

func TestGetChallengesByCategory(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	challenges := cm.GetChallengesByCategory("movement")
	if len(challenges) == 0 {
		t.Error("GetChallengesByCategory() returned empty slice for 'movement'")
	}

	for _, c := range challenges {
		if c.Category != "movement" {
			t.Errorf("GetChallengesByCategory() returned challenge with Category = %v, want movement", c.Category)
		}
	}
}

func TestDurationTierDistribution(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	tierCounts := make(map[string]int)

	for _, cat := range cm.GetCategories() {
		for _, c := range cm.GetChallengesByCategory(cat) {
			tierCounts[c.DurationTier()]++
		}
	}

	// We should have challenges in at least quick, standard, and complex tiers
	if tierCounts["quick"] == 0 {
		t.Error("No quick tier challenges found")
	}
	if tierCounts["standard"] == 0 {
		t.Error("No standard tier challenges found")
	}
	if tierCounts["complex"] == 0 {
		t.Error("No complex tier challenges found")
	}
}

func TestChallengeHasRequiredFields(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	for _, cat := range cm.GetCategories() {
		for _, c := range cm.GetChallengesByCategory(cat) {
			t.Run(c.ID, func(t *testing.T) {
				if c.ID == "" {
					t.Error("Challenge missing ID")
				}
				if c.Name == "" {
					t.Error("Challenge missing Name")
				}
				if c.Category == "" {
					t.Error("Challenge missing Category")
				}
				if c.Difficulty <= 0 {
					t.Errorf("Challenge Difficulty = %d, want > 0", c.Difficulty)
				}
				if c.Description == "" {
					t.Error("Challenge missing Description")
				}
				if c.ValidationType == "" {
					t.Error("Challenge missing ValidationType")
				}
				if c.ParKeystrokes <= 0 {
					t.Errorf("Challenge ParKeystrokes = %d, want > 0", c.ParKeystrokes)
				}
				if c.GoldBase <= 0 {
					t.Errorf("Challenge GoldBase = %d, want > 0", c.GoldBase)
				}
			})
		}
	}
}

func TestValidationTypes(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	validTypes := map[string]bool{
		"exact_match":     true,
		"contains":        true,
		"cursor_position": true,
		"cursor_on_char":  true,
		"function_exists": true,
		"pattern":         true,
		"different":       true,
	}

	for _, cat := range cm.GetCategories() {
		for _, c := range cm.GetChallengesByCategory(cat) {
			if !validTypes[c.ValidationType] {
				t.Errorf("Challenge %s has invalid validation type: %s", c.ID, c.ValidationType)
			}
		}
	}
}
