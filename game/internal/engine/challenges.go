package engine

import (
	"crypto/rand"
	"embed"
	"math/big"

	"gopkg.in/yaml.v3"
)

//go:embed assets/challenges.yaml
var challengesFS embed.FS

// Challenge represents a single kata challenge.
type Challenge struct {
	ID              string `yaml:"id"`
	Name            string `yaml:"name"`
	Category        string `yaml:"category"`
	Difficulty      int    `yaml:"difficulty"`
	Description     string `yaml:"description"`
	Filetype        string `yaml:"filetype"`
	InitialBuffer   string `yaml:"initial_buffer"`
	ExpectedBuffer  string `yaml:"expected_buffer,omitempty"`
	ValidationType  string `yaml:"validation_type"`
	ExpectedCursor  []int  `yaml:"expected_cursor,omitempty"`
	ExpectedContent string `yaml:"expected_content,omitempty"`
	FunctionName    string `yaml:"function_name,omitempty"`
	CursorStart     []int  `yaml:"cursor_start,omitempty"`
	ParKeystrokes   int    `yaml:"par_keystrokes"`
	GoldBase        int    `yaml:"gold_base"`
	RequiredPlugin  string `yaml:"required_plugin,omitempty"`
}

// ChallengeFile represents the YAML file structure.
type ChallengeFile struct {
	Challenges []Challenge `yaml:"challenges"`
}

// ChallengeManager loads and provides challenges.
type ChallengeManager struct {
	challenges   []Challenge
	byCategory   map[string][]Challenge
	byDifficulty map[int][]Challenge
}

// NewChallengeManager creates a new challenge manager and loads challenges.
func NewChallengeManager() (*ChallengeManager, error) {
	cm := &ChallengeManager{
		challenges:   make([]Challenge, 0),
		byCategory:   make(map[string][]Challenge),
		byDifficulty: make(map[int][]Challenge),
	}

	if err := cm.loadChallenges(); err != nil {
		return nil, err
	}

	return cm, nil
}

func (cm *ChallengeManager) loadChallenges() error {
	data, err := challengesFS.ReadFile("assets/challenges.yaml")
	if err != nil {
		// If embedded file not found, use empty challenges (not an error)
		cm.challenges = nil
		return nil //nolint:nilerr // intentional: missing file means empty challenges
	}

	var file ChallengeFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return err
	}

	cm.challenges = file.Challenges

	// Index by category and difficulty
	for i := range cm.challenges {
		c := &cm.challenges[i]
		cm.byCategory[c.Category] = append(cm.byCategory[c.Category], *c)
		cm.byDifficulty[c.Difficulty] = append(cm.byDifficulty[c.Difficulty], *c)
	}

	return nil
}

// GetChallenge returns a challenge by ID.
func (cm *ChallengeManager) GetChallenge(id string) *Challenge {
	for i := range cm.challenges {
		if cm.challenges[i].ID == id {
			return &cm.challenges[i]
		}
	}
	return nil
}

// GetRandomChallenge returns a random challenge matching the criteria.
func (cm *ChallengeManager) GetRandomChallenge(category string, maxDifficulty int) *Challenge {
	var candidateIndices []int

	for i := range cm.challenges {
		c := &cm.challenges[i]
		if category != "" && c.Category != category {
			continue
		}
		if maxDifficulty > 0 && c.Difficulty > maxDifficulty {
			continue
		}
		candidateIndices = append(candidateIndices, i)
	}

	if len(candidateIndices) == 0 {
		return nil
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(candidateIndices))))
	if err != nil {
		// Fallback to first candidate on error
		return &cm.challenges[candidateIndices[0]]
	}
	idx := candidateIndices[n.Int64()]
	return &cm.challenges[idx]
}

// GetChallengesByCategory returns all challenges in a category.
func (cm *ChallengeManager) GetChallengesByCategory(category string) []Challenge {
	return cm.byCategory[category]
}

// GetCategories returns all available categories.
func (cm *ChallengeManager) GetCategories() []string {
	categories := make([]string, 0, len(cm.byCategory))
	for cat := range cm.byCategory {
		categories = append(categories, cat)
	}
	return categories
}

// Count returns the total number of challenges.
func (cm *ChallengeManager) Count() int {
	return len(cm.challenges)
}

// CountByCategory returns the number of challenges per category.
func (cm *ChallengeManager) CountByCategory() map[string]int {
	counts := make(map[string]int)
	for cat, challenges := range cm.byCategory {
		counts[cat] = len(challenges)
	}
	return counts
}
