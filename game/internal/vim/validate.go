package vim

import (
	"regexp"
	"strings"
)

// ValidationResult contains the result of validating a challenge
type ValidationResult struct {
	Success    bool
	Efficiency float64
	Message    string
}

// ChallengeSpec contains the challenge validation parameters
type ChallengeSpec struct {
	ValidationType  string
	ExpectedBuffer  string
	ExpectedContent string
	ExpectedCursor  []int
	Pattern         string
	FunctionName    string
	InitialBuffer   string
	ParKeystrokes   int
}

// Validate checks if the editor state matches challenge expectations
func Validate(e *Editor, spec ChallengeSpec) ValidationResult {
	result := ValidationResult{}

	switch spec.ValidationType {
	case "exact_match":
		result.Success = normalizeBuffer(e.Buffer.String()) == normalizeBuffer(spec.ExpectedBuffer)
		if !result.Success {
			result.Message = "Buffer content doesn't match expected"
		}

	case "contains":
		result.Success = strings.Contains(e.Buffer.String(), spec.ExpectedContent)
		if !result.Success {
			result.Message = "Buffer should contain: " + spec.ExpectedContent
		}

	case "cursor_position":
		if len(spec.ExpectedCursor) == 2 {
			result.Success = e.Cursor.Line == spec.ExpectedCursor[0] &&
				e.Cursor.Col == spec.ExpectedCursor[1]
			if !result.Success {
				result.Message = "Cursor not at expected position"
			}
		}

	case "different":
		result.Success = normalizeBuffer(e.Buffer.String()) != normalizeBuffer(spec.InitialBuffer)
		if !result.Success {
			result.Message = "Buffer content unchanged"
		}

	case "pattern":
		matched, _ := regexp.MatchString(spec.Pattern, e.Buffer.String())
		result.Success = matched
		if !result.Success {
			result.Message = "Buffer doesn't match required pattern"
		}

	case "function_exists":
		result.Success = checkFunctionExists(e.Buffer.String(), spec.FunctionName)
		if !result.Success {
			result.Message = "Function " + spec.FunctionName + " not found"
		}

	default:
		// Default to success for unknown validation types
		result.Success = true
	}

	// Calculate efficiency
	if result.Success && e.KeystrokeCount > 0 && spec.ParKeystrokes > 0 {
		result.Efficiency = float64(spec.ParKeystrokes) / float64(e.KeystrokeCount)
		if result.Efficiency > 1.0 {
			result.Efficiency = 1.0 // Cap at 100%
		}
	} else if result.Success {
		result.Efficiency = 1.0
	}

	return result
}

func normalizeBuffer(s string) string {
	// Trim trailing newlines and normalize line endings
	s = strings.TrimRight(s, "\n\r")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return s
}

func checkFunctionExists(content, funcName string) bool {
	patterns := []string{
		`function\s+` + regexp.QuoteMeta(funcName) + `\s*\(`,      // JS/Lua
		`def\s+` + regexp.QuoteMeta(funcName) + `\s*\(`,           // Python
		`func\s+` + regexp.QuoteMeta(funcName) + `\s*\(`,          // Go
		regexp.QuoteMeta(funcName) + `\s*=\s*function`,            // JS function
		`const\s+` + regexp.QuoteMeta(funcName) + `\s*=`,          // JS const
		regexp.QuoteMeta(funcName) + `\s*=\s*\(`,                  // JS arrow
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, content); matched {
			return true
		}
	}
	return false
}
