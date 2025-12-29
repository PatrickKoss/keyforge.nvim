package engine

import (
	"crypto/rand"
	"math/big"
)

const (
	// DefaultMaxRecentChallenges is the number of challenges to remember before allowing repeats.
	DefaultMaxRecentChallenges = 30
	// DefaultMaxRecentCategories is the number of recent categories to track for variety weighting.
	DefaultMaxRecentCategories = 3
)

// ChallengeSelector provides intelligent challenge selection with history tracking.
// It wraps ChallengeManager and maintains per-session state to avoid repetition
// and ensure category variety.
type ChallengeSelector struct {
	manager            *ChallengeManager
	recentChallengeIDs []string // Ring buffer, max DefaultMaxRecentChallenges
	recentCategories   []string // Ring buffer, max DefaultMaxRecentCategories
}

// NewChallengeSelector creates a new selector wrapping the given manager.
func NewChallengeSelector(cm *ChallengeManager) *ChallengeSelector {
	return &ChallengeSelector{
		manager:            cm,
		recentChallengeIDs: make([]string, 0, DefaultMaxRecentChallenges),
		recentCategories:   make([]string, 0, DefaultMaxRecentCategories),
	}
}

// GetChallenge returns a challenge matching criteria while avoiding repetition.
// If category is empty, any category is allowed but recently-used categories are penalized.
func (cs *ChallengeSelector) GetChallenge(category string, maxDifficulty int) *Challenge {
	// Step 1: Get all matching candidates
	candidates := cs.manager.GetCandidates(category, maxDifficulty)
	if len(candidates) == 0 {
		return nil
	}

	// Step 2: Filter out recently shown challenges
	fresh := cs.filterRecent(candidates)

	// Step 3: Handle exhaustion - if all candidates were recently shown
	if len(fresh) == 0 {
		// Fallback: pick least-recently-seen from candidates
		selected := cs.leastRecentFrom(candidates)
		cs.recordShown(selected)
		return selected
	}

	// Step 4: Apply category variety weighting
	weights := cs.computeWeights(fresh)

	// Step 5: Weighted random selection
	selected := cs.weightedRandom(fresh, weights)

	// Step 6: Record and return
	cs.recordShown(selected)
	return selected
}

// Reset clears all history (call when starting a new game).
func (cs *ChallengeSelector) Reset() {
	cs.recentChallengeIDs = cs.recentChallengeIDs[:0]
	cs.recentCategories = cs.recentCategories[:0]
}

// filterRecent removes challenges that appear in recentChallengeIDs.
func (cs *ChallengeSelector) filterRecent(candidates []*Challenge) []*Challenge {
	recentSet := make(map[string]bool, len(cs.recentChallengeIDs))
	for _, id := range cs.recentChallengeIDs {
		recentSet[id] = true
	}

	result := make([]*Challenge, 0, len(candidates))
	for _, c := range candidates {
		if !recentSet[c.ID] {
			result = append(result, c)
		}
	}
	return result
}

// leastRecentFrom returns the candidate that appears earliest in history,
// or the first candidate if none are in history.
func (cs *ChallengeSelector) leastRecentFrom(candidates []*Challenge) *Challenge {
	if len(candidates) == 0 {
		return nil
	}

	// Build position map: ID -> position in recentChallengeIDs
	posMap := make(map[string]int, len(cs.recentChallengeIDs))
	for i, id := range cs.recentChallengeIDs {
		posMap[id] = i
	}

	// Find candidate with lowest position (earliest in history = least recent)
	var bestCandidate *Challenge
	bestPos := len(cs.recentChallengeIDs) + 1 // Beyond end = not in history

	for _, c := range candidates {
		pos, found := posMap[c.ID]
		if !found {
			// Not in history at all - ideal candidate
			return c
		}
		if pos < bestPos {
			bestPos = pos
			bestCandidate = c
		}
	}

	if bestCandidate != nil {
		return bestCandidate
	}
	return candidates[0] // Fallback
}

// computeWeights returns weights for candidates based on category recency.
func (cs *ChallengeSelector) computeWeights(candidates []*Challenge) []float64 {
	weights := make([]float64, len(candidates))
	for i, c := range candidates {
		weights[i] = cs.categoryWeight(c.Category)
	}
	return weights
}

// categoryWeight returns a weight multiplier based on category recency.
// Recently used categories get lower weights.
func (cs *ChallengeSelector) categoryWeight(category string) float64 {
	for i, cat := range cs.recentCategories {
		if cat == category {
			// i=0 is oldest, len-1 is most recent (FIFO queue)
			recency := len(cs.recentCategories) - 1 - i
			switch recency {
			case 0:
				return 0.3 // Most recent: strong penalty
			case 1:
				return 0.6 // Second most recent: medium penalty
			case 2:
				return 0.8 // Third most recent: light penalty
			}
		}
	}
	return 1.0 // Not in recent history: full weight
}

// weightedRandom selects a random item from candidates using weights.
func (cs *ChallengeSelector) weightedRandom(candidates []*Challenge, weights []float64) *Challenge {
	if len(candidates) == 0 {
		return nil
	}
	if len(candidates) == 1 {
		return candidates[0]
	}

	// Compute total weight
	var totalWeight float64
	for _, w := range weights {
		totalWeight += w
	}

	if totalWeight <= 0 {
		// All weights zero - fallback to uniform random
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(candidates))))
		if err != nil {
			return candidates[0]
		}
		return candidates[n.Int64()]
	}

	// Generate random value in [0, totalWeight)
	// Using crypto/rand for consistency with existing code
	// Scale up for integer precision (multiply by 10000 for 4 decimal places)
	const scale = 10000
	maxInt := int64(totalWeight * scale)
	n, err := rand.Int(rand.Reader, big.NewInt(maxInt))
	if err != nil {
		return candidates[0]
	}
	r := float64(n.Int64()) / scale

	// Find the selected candidate
	var cumulative float64
	for i, w := range weights {
		cumulative += w
		if r < cumulative {
			return candidates[i]
		}
	}

	return candidates[len(candidates)-1]
}

// recordShown adds the challenge to history.
func (cs *ChallengeSelector) recordShown(c *Challenge) {
	if c == nil {
		return
	}

	// Add to recent challenges (ring buffer)
	if len(cs.recentChallengeIDs) >= DefaultMaxRecentChallenges {
		// Remove oldest (front)
		cs.recentChallengeIDs = cs.recentChallengeIDs[1:]
	}
	cs.recentChallengeIDs = append(cs.recentChallengeIDs, c.ID)

	// Add to recent categories (ring buffer, avoid duplicates of immediate repeat)
	if len(cs.recentCategories) == 0 || cs.recentCategories[len(cs.recentCategories)-1] != c.Category {
		if len(cs.recentCategories) >= DefaultMaxRecentCategories {
			cs.recentCategories = cs.recentCategories[1:]
		}
		cs.recentCategories = append(cs.recentCategories, c.Category)
	}
}
