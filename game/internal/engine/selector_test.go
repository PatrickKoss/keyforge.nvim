package engine

import "testing"

func TestSelectorAvoidsRepetition(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	cs := NewChallengeSelector(cm)

	// Get 20 movement challenges (movement has ~20 challenges)
	seen := make(map[string]int)
	for range 20 {
		c := cs.GetChallenge("movement", 3)
		if c == nil {
			t.Fatal("GetChallenge returned nil")
		}
		seen[c.ID]++
		if seen[c.ID] > 1 {
			t.Errorf("Challenge %s repeated within 20 selections (count: %d)", c.ID, seen[c.ID])
		}
	}
}

func TestSelectorCategoryVariety(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	cs := NewChallengeSelector(cm)

	// Get 15 any-category challenges
	categories := make([]string, 0, 15)
	for range 15 {
		c := cs.GetChallenge("", 3)
		if c == nil {
			t.Fatal("GetChallenge returned nil")
		}
		categories = append(categories, c.Category)
	}

	// Check that we don't have 4+ of the same category in a row
	// (3 could happen due to weighting, but 4 should be very rare)
	for i := range len(categories) - 3 {
		if categories[i] == categories[i+1] &&
			categories[i+1] == categories[i+2] &&
			categories[i+2] == categories[i+3] {
			t.Errorf("Same category %s appeared 4 times in a row at index %d", categories[i], i)
		}
	}
}

func TestSelectorReset(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	cs := NewChallengeSelector(cm)

	// Get some challenges
	for range 5 {
		cs.GetChallenge("movement", 3)
	}

	// Reset
	cs.Reset()

	// Verify history is cleared
	if len(cs.recentChallengeIDs) != 0 {
		t.Errorf("recentChallengeIDs not cleared: got %d", len(cs.recentChallengeIDs))
	}
	if len(cs.recentCategories) != 0 {
		t.Errorf("recentCategories not cleared: got %d", len(cs.recentCategories))
	}
}

func TestSelectorExhaustionFallback(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	cs := NewChallengeSelector(cm)

	// Get harpoon challenges (small category: 5 challenges)
	// Should not panic or return nil even after 10 selections
	for i := range 10 {
		c := cs.GetChallenge("harpoon", 3)
		if c == nil {
			t.Fatalf("GetChallenge returned nil on iteration %d", i)
		}
		if c.Category != "harpoon" {
			t.Errorf("Expected harpoon category, got %s", c.Category)
		}
	}
}

func TestSelectorRespectsCategory(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	cs := NewChallengeSelector(cm)

	// Request specific categories
	categories := []string{"movement", "lsp-navigation", "text-objects", "refactoring"}
	for _, cat := range categories {
		c := cs.GetChallenge(cat, 3)
		if c == nil {
			t.Errorf("GetChallenge returned nil for category %s", cat)
			continue
		}
		if c.Category != cat {
			t.Errorf("Expected %s, got %s", cat, c.Category)
		}
	}
}

func TestSelectorRespectsDifficulty(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	cs := NewChallengeSelector(cm)

	// Request max difficulty 1
	for range 20 {
		c := cs.GetChallenge("", 1)
		if c == nil {
			continue // Some iterations might exhaust easy challenges
		}
		if c.Difficulty > 1 {
			t.Errorf("Expected difficulty <= 1, got %d for %s", c.Difficulty, c.ID)
		}
	}
}

func TestSelectorReturnsNilForNoMatch(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	cs := NewChallengeSelector(cm)

	// Request non-existent category
	c := cs.GetChallenge("nonexistent-category", 3)
	if c != nil {
		t.Errorf("Expected nil for nonexistent category, got %v", c)
	}
}

func TestSelectorLeastRecentFallback(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	cs := NewChallengeSelector(cm)

	// Exhaust harpoon (5 challenges) and verify we get the least recent on 6th
	firstFive := make([]string, 5)
	for i := range 5 {
		c := cs.GetChallenge("harpoon", 3)
		if c == nil {
			t.Fatalf("GetChallenge returned nil on iteration %d", i)
		}
		firstFive[i] = c.ID
	}

	// 6th should be the first one (least recent)
	sixth := cs.GetChallenge("harpoon", 3)
	if sixth == nil {
		t.Fatal("GetChallenge returned nil on 6th iteration")
	}
	if sixth.ID != firstFive[0] {
		t.Errorf("Expected least recent challenge %s, got %s", firstFive[0], sixth.ID)
	}
}

func TestCategoryWeightPenalties(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	cs := NewChallengeSelector(cm)

	// Manually set up category history
	cs.recentCategories = []string{"movement", "text-objects", "lsp-navigation"}

	// Test weights - most recent (lsp-navigation) should have lowest weight
	weight1 := cs.categoryWeight("lsp-navigation") // most recent
	weight2 := cs.categoryWeight("text-objects")   // second most recent
	weight3 := cs.categoryWeight("movement")       // third most recent
	weight4 := cs.categoryWeight("refactoring")    // not in history

	if weight1 >= weight2 {
		t.Errorf("Most recent category weight (%f) should be < second most recent (%f)", weight1, weight2)
	}
	if weight2 >= weight3 {
		t.Errorf("Second most recent weight (%f) should be < third most recent (%f)", weight2, weight3)
	}
	if weight3 >= weight4 {
		t.Errorf("Third most recent weight (%f) should be < fresh category (%f)", weight3, weight4)
	}
	if weight4 != 1.0 {
		t.Errorf("Fresh category weight should be 1.0, got %f", weight4)
	}
}

func TestGetCandidates(t *testing.T) {
	cm, err := NewChallengeManager()
	if err != nil {
		t.Fatalf("NewChallengeManager() error = %v", err)
	}

	// Test category filter
	candidates := cm.GetCandidates("movement", 0)
	for _, c := range candidates {
		if c.Category != "movement" {
			t.Errorf("Expected movement category, got %s", c.Category)
		}
	}

	// Test difficulty filter
	candidates = cm.GetCandidates("", 1)
	for _, c := range candidates {
		if c.Difficulty > 1 {
			t.Errorf("Expected difficulty <= 1, got %d for %s", c.Difficulty, c.ID)
		}
	}

	// Test combined filter
	candidates = cm.GetCandidates("movement", 1)
	for _, c := range candidates {
		if c.Category != "movement" {
			t.Errorf("Expected movement category, got %s", c.Category)
		}
		if c.Difficulty > 1 {
			t.Errorf("Expected difficulty <= 1, got %d for %s", c.Difficulty, c.ID)
		}
	}
}
