package nvim

// JSON-RPC 2.0 message types for Neovim communication

// RPCClient is the interface for sending RPC messages to Neovim
// Both Client (stdin/stderr) and SocketServer implement this
type RPCClient interface {
	RequestChallenge(requestID, category string, difficulty int) error
	SendGameState(state string, wave, gold, health, enemies, towers int) error
	SendGameReady() error
	SendGoldUpdate(gold, earned int, source string, speedBonus float64) error
	SendChallengeAvailable(count, nextReward int, nextCategory string) error
	SendGameOver(wave, gold, towers, health int) error
	SendVictory(wave, gold, towers, health int) error
}

// Request represents a JSON-RPC request
type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      int         `json:"id,omitempty"`
}

// Response represents a JSON-RPC response
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      int         `json:"id"`
}

// Notification represents a JSON-RPC notification (no ID, no response expected)
type Notification struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Standard JSON-RPC error codes
const (
	ErrCodeParseError     = -32700
	ErrCodeInvalidRequest = -32600
	ErrCodeMethodNotFound = -32601
	ErrCodeInvalidParams  = -32602
	ErrCodeInternalError  = -32603
)

// Game -> Neovim message types

// ChallengeRequest asks Neovim to present a challenge to the user
type ChallengeRequest struct {
	RequestID  string `json:"request_id"`
	Category   string `json:"category"`
	Difficulty int    `json:"difficulty"`
}

// GameStateUpdate notifies Neovim of game state changes
type GameStateUpdate struct {
	State    string `json:"state"` // playing, paused, game_over, victory
	Wave     int    `json:"wave"`
	Gold     int    `json:"gold"`
	Health   int    `json:"health"`
	Enemies  int    `json:"enemies"`
	Towers   int    `json:"towers"`
}

// GoldUpdate notifies Neovim of gold changes from challenges
type GoldUpdate struct {
	Gold       int     `json:"gold"`        // New gold total
	Earned     int     `json:"earned"`      // Gold earned from this action
	Source     string  `json:"source"`      // "challenge", "mob", "wave_bonus"
	SpeedBonus float64 `json:"speed_bonus"` // Speed bonus multiplier if from challenge
}

// ChallengeAvailable notifies Neovim that challenges are available
type ChallengeAvailable struct {
	Count         int `json:"count"`          // Number of available challenges
	NextReward    int `json:"next_reward"`    // Estimated reward for next challenge
	NextCategory  string `json:"next_category"` // Category of next challenge
}

// Neovim -> Game message types

// ChallengeResult contains the result of a completed challenge
type ChallengeResult struct {
	RequestID      string  `json:"request_id"`
	Success        bool    `json:"success"`
	Skipped        bool    `json:"skipped,omitempty"`
	KeystrokeCount int     `json:"keystroke_count,omitempty"`
	TimeMs         int     `json:"time_ms,omitempty"`
	Efficiency     float64 `json:"efficiency,omitempty"`
	SpeedBonus     float64 `json:"speed_bonus,omitempty"`  // Speed bonus multiplier
	GoldEarned     int     `json:"gold_earned,omitempty"`  // Gold earned from challenge
}

// StartChallengeRequest is sent when user triggers a new challenge
type StartChallengeRequest struct {
	// Empty for now, but can be extended with category preferences
}

// ConfigUpdate contains configuration from Neovim
type ConfigUpdate struct {
	Difficulty     string `json:"difficulty,omitempty"`
	UseNerdFonts   bool   `json:"use_nerd_fonts,omitempty"`
	StartingGold   int    `json:"starting_gold,omitempty"`
	StartingHealth int    `json:"starting_health,omitempty"`
}

// Method names
const (
	// Game -> Neovim
	MethodRequestChallenge   = "request_challenge"
	MethodGameStateUpdate    = "game_state_update"
	MethodGameReady          = "game_ready"
	MethodGoldUpdate         = "gold_update"
	MethodChallengeAvailable = "challenge_available"
	MethodGameOver           = "game_over"
	MethodVictory            = "victory"

	// Neovim -> Game
	MethodChallengeComplete = "challenge_complete"
	MethodStartChallenge    = "start_challenge"
	MethodConfigUpdate      = "config_update"
	MethodPauseGame         = "pause_game"
	MethodResumeGame        = "resume_game"
	MethodRestartGame       = "restart_game"
)

// GameOverParams contains game over notification data
type GameOverParams struct {
	Wave    int `json:"wave"`
	Gold    int `json:"gold"`
	Towers  int `json:"towers"`
	Health  int `json:"health"`
}

// VictoryParams contains victory notification data
type VictoryParams struct {
	Wave    int `json:"wave"`
	Gold    int `json:"gold"`
	Towers  int `json:"towers"`
	Health  int `json:"health"`
}

// NewRequest creates a new JSON-RPC request
func NewRequest(id int, method string, params interface{}) *Request {
	return &Request{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      id,
	}
}

// NewNotification creates a new JSON-RPC notification
func NewNotification(method string, params interface{}) *Notification {
	return &Notification{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
}

// NewResponse creates a new JSON-RPC response
func NewResponse(id int, result interface{}, err *RPCError) *Response {
	return &Response{
		JSONRPC: "2.0",
		Result:  result,
		Error:   err,
		ID:      id,
	}
}

// NewError creates a new RPC error
func NewError(code int, message string) *RPCError {
	return &RPCError{
		Code:    code,
		Message: message,
	}
}
